package gateway

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/rxdn/gdl/cache"
	"github.com/rxdn/gdl/gateway/payloads"
	"github.com/rxdn/gdl/gateway/payloads/events"
	"github.com/rxdn/gdl/objects/user"
	"github.com/rxdn/gdl/utils"
	"github.com/sirupsen/logrus"
	"github.com/tatsuworks/czlib"
	"io"
	"log"
	"nhooyr.io/websocket"
	"sync"
	"time"
)

type Shard struct {
	ShardManager *ShardManager
	Token        string
	ShardId      int

	State     State
	StateLock sync.RWMutex

	WebSocket  *websocket.Conn
	Context    context.Context
	ZLibReader io.ReadCloser
	ReadLock   sync.Mutex

	SequenceLock   sync.RWMutex
	SequenceNumber *int

	LastHeartbeat     int64 // Millis
	HeartbeatInterval int
	HasDoneHeartbeat  bool

	LastHeartbeatAcknowledgement int64 // Millis
	HeartbeatMutex               sync.Mutex
	KillHeartbeat                chan struct{}

	SessionId string

	Cache *cache.Cache
}

func NewShard(shardManager *ShardManager, token string, shardId int) Shard {
	cache := shardManager.CacheFactory()

	return Shard{
		ShardManager:                 shardManager,
		Token:                        token,
		ShardId:                      shardId,
		State:                        DEAD,
		Context:                      context.Background(),
		LastHeartbeatAcknowledgement: utils.GetCurrentTimeMillis(),
		Cache:                        &cache,
	}
}

func (s *Shard) EnsureConnect() {
	if err := s.Connect(); err != nil {
		logrus.Warnf("shard %d: Error whilst connecting: %s", s.ShardId, err.Error())
		time.Sleep(500 * time.Millisecond)
		s.EnsureConnect()
	}
}

func (s *Shard) Connect() error {
	logrus.Infof("shard %d: Starting", s.ShardId)

	// Connect to Discord
	s.StateLock.RLock()
	state := s.State
	s.StateLock.RUnlock()
	if state != DEAD {
		return s.Kill()
	}

	s.StateLock.Lock()
	s.State = CONNECTING
	s.StateLock.Unlock()

	// initialise zlib reader
	zlibReader, err := czlib.NewReader(bytes.NewReader(nil))
	if err != nil {
		return err
	}

	s.ZLibReader = zlibReader
	defer zlibReader.Close()

	conn, _, err := websocket.Dial(s.Context, "wss://gateway.discord.gg/?v=6&encoding=json&compress=zlib-stream", &websocket.DialOptions{
		CompressionMode: websocket.CompressionContextTakeover,
	})

	if err != nil {
		s.StateLock.Lock()
		s.State = DEAD
		s.StateLock.Unlock()
		return err
	}

	conn.SetReadLimit(4294967296)

	s.WebSocket = conn

	// Read hello
	if err := s.Read(); err != nil {
		logrus.Warnf("shard %d: Error whilst reading Hello: %s", s.ShardId, err.Error())
		s.Kill()
		return err
	}

	if s.SessionId == "" || s.SequenceNumber == nil {
		s.Identify(s.ShardManager.Presence, s.ShardManager.GuildSubscriptions)
	} else {
		s.Resume()
	}

	logrus.Infof("shard %d: Connected", s.ShardId)

	s.StateLock.Lock()
	s.State = CONNECTED
	s.StateLock.Unlock()

	go func() {
		for {
			// Verify that we are still connected
			s.StateLock.RLock()
			state := s.State
			s.StateLock.RUnlock()
			if state != CONNECTED {
				break
			}

			// Read
			if err := s.Read(); err != nil {
				logrus.Warnf("shard %d: Error whilst reading payload: %s", s.ShardId, err.Error())
				s.Kill()
				s.EnsureConnect()
			}
		}
	}()

	return nil
}

func (s *Shard) Identify(status user.UpdateStatus, guildSubscriptions bool) {
	identify := payloads.NewIdentify(s.ShardId, s.ShardManager.TotalShards, s.Token, status, guildSubscriptions)
	s.ShardManager.GatewayBucket.Wait(1)

	if err := s.Write(identify); err != nil {
		logrus.Warnf("shard %d: Error whilst sending Identify: %s", s.ShardId, err.Error())
		s.Identify(status, guildSubscriptions)
	}
}

func (s *Shard) Resume() {
	s.SequenceLock.RLock()
	resume := payloads.NewResume(s.Token, s.SessionId, *s.SequenceNumber)
	s.SequenceLock.RUnlock()

	logrus.Infof("shard %d: Resuming", s.ShardId)

	if err := s.Write(resume); err != nil {
		logrus.Warnf("shard %d: Error whilst sending Resume: %s", s.ShardId, err.Error())
		s.Identify(s.ShardManager.Presence, s.ShardManager.GuildSubscriptions)
	}
}

func (s *Shard) Read() error {
	defer func() {
		if r := recover(); r != nil {
			logrus.Warnf("Recovered panic while reading: %s", r)
			return
		}
	}()

	s.ReadLock.Lock()
	_, reader, err := s.WebSocket.Reader(s.Context)
	if err != nil {
		s.ReadLock.Unlock()
		return err
	}

	// decompress
	var buffer bytes.Buffer
	s.ZLibReader.(czlib.Resetter).Reset(reader)
	_, err = buffer.ReadFrom(s.ZLibReader)

	s.ReadLock.Unlock()

	if err != nil {
		return err
	}

	data := buffer.Bytes()

	payload, err := payloads.NewPayload(data)
	if err != nil {
		return err
	}

	// Handle new sequence number
	if payload.SequenceNumber != nil {
		s.SequenceLock.Lock()
		s.SequenceNumber = payload.SequenceNumber
		s.SequenceLock.Unlock()
	}

	// Handle payload
	switch payload.Opcode {
	case 0: // Event
		{
			event := events.EventType(payload.EventName)
			s.ExecuteEvent(event, payload.Data)
		}
	case 7: // Reconnect
		{
			s.Kill()
			go s.EnsureConnect()
		}
	case 9: // Invalid session
		{
			s.Kill()
			s.SessionId = ""
			go s.EnsureConnect()
		}
	case 10: // Hello
		{
			hello, err := payloads.NewHello(data)
			if err != nil {
				return err
			}

			s.HeartbeatInterval = hello.EventData.Interval
			s.KillHeartbeat = make(chan struct{})

			ticker := time.NewTicker(time.Duration(int32(s.HeartbeatInterval)) * time.Millisecond)
			go s.CountdownHeartbeat(ticker)
		}
	case 11: // Heartbeat ACK
		{
			_, err := payloads.NewHeartbeackAck(data)
			if err != nil {
				log.Println(err.Error())
				return err
			}

			s.HeartbeatMutex.Lock()
			s.LastHeartbeatAcknowledgement = time.Now().UnixNano() / int64(time.Millisecond)
			s.HeartbeatMutex.Unlock()
		}
	}

	return nil
}

func (s *Shard) Write(payload interface{}) error {
	encoded, err := json.Marshal(payload)
	if err != nil {
		return err
	}

	return s.WriteRaw(encoded)
}

func (s *Shard) WriteRaw(data []byte) error {
	if s.WebSocket == nil {
		msg := fmt.Sprintf("shard %d: WS is closed", s.ShardId)
		log.Println(msg)
		return errors.New(msg)
	}

	err := s.WebSocket.Write(s.Context, websocket.MessageText, data)

	return err
}

func (s *Shard) Kill() error {
	logrus.Infof("killing shard %d", s.ShardId)

	go func() {
		s.KillHeartbeat <- struct{}{}
	}()

	s.StateLock.Lock()
	s.State = DISCONNECTING
	s.StateLock.Unlock()

	var err error
	if s.WebSocket != nil {
		err = s.WebSocket.Close(4000, "unknown")
	}

	s.WebSocket = nil

	s.StateLock.Lock()
	s.State = DEAD
	s.StateLock.Unlock()

	return err
}
