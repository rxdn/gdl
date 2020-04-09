package gateway

import (
	"fmt"
	"github.com/rxdn/gdl/gateway/payloads"
	"github.com/rxdn/gdl/utils"
	"github.com/sirupsen/logrus"
	"time"
)

func (s *Shard) CountdownHeartbeat(ticker *time.Ticker) {
	loop:
	for {
		select {
		case <-s.KillHeartbeat:
			ticker.Stop()
			break loop
		case <-ticker.C:
			s.HeartbeatMutex.Lock()

			// Check we received an ACK
			timeElapsed := s.LastHeartbeatAcknowledgement - s.LastHeartbeat
			fmt.Printf("current time millis: %d; last ack: %d; elapsed: %d; interval: %d\n", utils.GetCurrentTimeMillis(), s.LastHeartbeatAcknowledgement, timeElapsed, s.HeartbeatInterval)
			if s.HasDoneHeartbeat && timeElapsed > int64(s.HeartbeatInterval) {
				logrus.Warnf("shard %d didn't receive acknowledgement, restarting", s.ShardId)
				s.HeartbeatMutex.Unlock()
				s.Kill()
				go s.EnsureConnect()
				return
			}

			s.HeartbeatMutex.Unlock()

			if err := s.Heartbeat(); err != nil {
				logrus.Warnf("shard %d heartbeat failed, restarting: %s", s.ShardId, err.Error())
				s.Kill()
				go s.EnsureConnect()
			}
		}
	}
}

func (s *Shard) Heartbeat() error {
	s.sequenceLock.RLock()
	payload := payloads.NewHeartbeat(s.sequenceNumber)
	s.sequenceLock.RUnlock()

	s.HasDoneHeartbeat = true
	s.LastHeartbeat = utils.GetCurrentTimeMillis()

	return s.write(payload)
}
