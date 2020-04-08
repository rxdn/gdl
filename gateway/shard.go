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
	"github.com/rxdn/gdl/utils"
	"github.com/sirupsen/logrus"
	"github.com/tatsuworks/czlib"
	"io"
	"log"
	"nhooyr.io/websocket"
	"runtime/debug"
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

	Cache cache.Cache

	guilds     []uint64
	guildsLock *sync.RWMutex
}

func NewShard(shardManager *ShardManager, token string, shardId int) Shard {
	cache := shardManager.ShardOptions.CacheFactory()

	return Shard{
		ShardManager:                 shardManager,
		Token:                        token,
		ShardId:                      shardId,
		State:                        DEAD,
		Context:                      context.Background(),
		LastHeartbeatAcknowledgement: utils.GetCurrentTimeMillis(),
		Cache:                        cache,
		guildsLock:                   &sync.RWMutex{},
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
	s.StateLock.Lock() // Can't RLock - potential state issue
	state := s.State

	if state != DEAD {
		s.StateLock.Unlock()
		if err := s.Kill(); err != nil {
			return err
		}
		s.StateLock.Lock()
	}

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
	if err := s.read(); err != nil {
		logrus.Warnf("shard %d: Error whilst reading Hello: %s", s.ShardId, err.Error())
		s.Kill()
		return err
	}

	if s.SessionId == "" || s.SequenceNumber == nil {
		s.identify()
	} else {
		s.resume()
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
			if err := s.read(); err != nil {
				logrus.Warnf("shard %d: Error whilst reading payload: %s", s.ShardId, err.Error())

				s.StateLock.Lock()
				state := s.State
				s.StateLock.Unlock()

				if state == CONNECTED {
					s.Kill()
					s.EnsureConnect()
				}
			}
		}
	}()

	return nil
}

func (s *Shard) identify() {
	// call hook
	if s.ShardManager.ShardOptions.Hooks.IdentifyHook != nil {
		s.ShardManager.ShardOptions.Hooks.IdentifyHook(s)
	}

	// build payload
	identify := payloads.NewIdentify(s.ShardId, s.ShardManager.ShardOptions.ShardCount.Total, s.Token, s.ShardManager.ShardOptions.Presence, s.ShardManager.ShardOptions.GuildSubscriptions)

	// wait for ratelimit
	if err := s.ShardManager.RateLimiter.IdentifyWait(); err != nil {
		logrus.Warnf("shard %d: Error whilst waiting on identify ratelimit: %s", s.ShardId, err.Error())
	}

	if err := s.write(identify); err != nil {
		logrus.Warnf("shard %d: Error whilst sending Identify: %s", s.ShardId, err.Error())
		s.identify()
	}
}

func (s *Shard) resume() {
	s.SequenceLock.RLock()
	resume := payloads.NewResume(s.Token, s.SessionId, *s.SequenceNumber)
	s.SequenceLock.RUnlock()

	logrus.Infof("shard %d: Resuming", s.ShardId)

	if err := s.write(resume); err != nil {
		logrus.Warnf("shard %d: Error whilst sending Resume: %s", s.ShardId, err.Error())
		s.identify()
	}
}

func (s *Shard) read() error {
	defer func() {
		if r := recover(); r != nil {
			logrus.Warnf("Recovered panic while reading: %s", r)
			s.Kill()
			go s.EnsureConnect()
		}
	}()

	buffer, err := s.readData()

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
			go s.ExecuteEvent(event, payload.Data)
		}
	case 7: // Reconnect
		{
			logrus.Infof("shard %d: received reconnect payload from discord", s.ShardId)

			if s.ShardManager.ShardOptions.Hooks.ReconnectHook != nil {
				s.ShardManager.ShardOptions.Hooks.ReconnectHook(s)
			}

			s.Kill()
			go s.EnsureConnect()
		}
	case 9: // Invalid session
		{
			logrus.Infof("shard %d: received invalid session payload from discord", s.ShardId)
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

func (s *Shard) readData() (bytes.Buffer, error) {
	var buffer bytes.Buffer

	s.ReadLock.Lock()
	defer s.ReadLock.Unlock()

	if s.WebSocket == nil {
		return buffer, errors.New("websocket is nil")
	}

	_, reader, err := s.WebSocket.Reader(s.Context)
	if err != nil {
		return buffer, err
	}

	// decompress
	s.ZLibReader.(czlib.Resetter).Reset(reader)
	_, err = buffer.ReadFrom(s.ZLibReader)
	return buffer, err
}

func (s *Shard) write(payload interface{}) error {
	encoded, err := json.Marshal(payload)
	if err != nil {
		return err
	}

	return s.writeRaw(encoded)
}

func (s *Shard) writeRaw(data []byte) error {
	if s.WebSocket == nil {
		msg := fmt.Sprintf("shard %d: WS is closed", s.ShardId)
		log.Println(msg)
		return errors.New(msg)
	}

	err := s.WebSocket.Write(s.Context, websocket.MessageText, data)

	return err
}

func (s *Shard) Kill() error {
	if s.ShardManager.ShardOptions.Debug {
		debug.PrintStack()
	}

	logrus.Infof("killing shard %d", s.ShardId)

	go func() {
		s.KillHeartbeat <- struct{}{}
	}()

	s.ZLibReader.Close()

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

	logrus.Infof("killed shard %d", s.ShardId)

	return err
}
