package gateway

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/Dot-Rar/gdl/cache"
	"github.com/Dot-Rar/gdl/gateway/payloads"
	"github.com/Dot-Rar/gdl/gateway/payloads/events"
	"github.com/Dot-Rar/gdl/utils"
	"github.com/gorilla/websocket"
	"github.com/sirupsen/logrus"
	"log"
	"sync"
	"time"
)

type Shard struct {
	ShardManager *ShardManager
	Token        string
	ShardId      int

	State     State
	StateLock sync.Mutex

	WebSocket *websocket.Conn
	ReadLock  sync.Mutex

	SequenceLock   sync.Mutex
	SequenceNumber *int

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
		SequenceNumber:               nil,
		LastHeartbeatAcknowledgement: utils.GetCurrentTimeMillis(),
		Cache:                        &cache,
	}
}

func (s *Shard) EnsureConnect() {
	if err := s.Connect(); err != nil {
		logrus.Warnf("shard %d: Error whilst connecting: %s", s.ShardId, err.Error())
		s.EnsureConnect()
	}
}

func (s *Shard) Connect() error {
	// Connect to Discord
	s.StateLock.Lock()
	if s.State != DEAD {
		return s.Kill()
	}

	s.State = CONNECTING

	conn, _, err := websocket.DefaultDialer.Dial("wss://gateway.discord.gg/?v=6&encoding=json", nil)
	if err != nil {
		s.State = DEAD
		return err
	}
	s.StateLock.Unlock()

	s.WebSocket = conn
	conn.SetCloseHandler(s.OnClose)

	// Read hello
	if err := s.Read(); err != nil {
		logrus.Warnf("shard %d: Error whilst reading Hello: %s", s.ShardId, err.Error())
		s.Kill()
		return err
	}

	if s.SessionId == "" {
		s.Identify()
	} else {
		// TODO: Resume
	}

	logrus.Infof("shard %d: Connected", s.ShardId)

	s.StateLock.Lock()
	s.State = CONNECTED
	s.StateLock.Unlock()

	go func() {
		for {
			// Verify that we are still connected
			s.StateLock.Lock()
			state := s.State
			s.StateLock.Unlock()
			if state != CONNECTED {
				break
			}

			// Read
			if err := s.Read(); err != nil {
				logrus.Warnf("shard %d: Error whilst reading payload: %s", s.ShardId, err.Error())
			}
		}
	}()

	return nil
}

func (s *Shard) Identify() {
	identify := payloads.NewIdentify(s.ShardId, s.ShardManager.TotalShards, s.Token)
	s.ShardManager.GatewayBucket.Wait(1)

	if err := s.Write(identify); err != nil {
		logrus.Warnf("shard %d: Error whilst sending Identify: %s", s.ShardId, err.Error())
		s.Identify()
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
	_, data, err := s.WebSocket.ReadMessage()
	s.ReadLock.Unlock()

	if err != nil {
		return err
	}

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
	case 1: // Heartbeat
		{

		}
	case 7: // Reconnect
		{

		}
	case 9: // Invalid session
		{

		}
	case 10: // Hello
		{
			hello, err := payloads.NewHello(data)
			if err != nil {
				return err
			}

			s.HeartbeatInterval = hello.EventData.Interval
			s.HeartbeatInterval = 10000
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

	err := s.WebSocket.WriteMessage(1, data)

	return err
}

func (s *Shard) OnClose(code int, text string) error {
	logrus.Warnf("shard %d: Discord closed WS", s.ShardId)

	if code == 1000 || code == 1001 || code == 4007 { // Closed gracefully || Invalid seq
		s.SessionId = ""

		s.SequenceLock.Lock()
		s.SequenceNumber = nil
		s.SequenceLock.Unlock()
	}

	s.KillHeartbeat <- struct{}{}

	_ = s.Kill()
	go s.EnsureConnect()
	return nil

}

func (s *Shard) Kill() error {
	s.KillHeartbeat <- struct{}{}

	s.StateLock.Lock()
	s.State = DISCONNECTING
	s.StateLock.Unlock()

	var err error
	if s.WebSocket != nil {
		err = s.WebSocket.Close()
	}

	s.WebSocket = nil

	s.StateLock.Lock()
	s.State = DEAD
	s.StateLock.Unlock()

	return err
}
