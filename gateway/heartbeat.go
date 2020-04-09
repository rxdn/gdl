package gateway

import (
	"github.com/rxdn/gdl/gateway/payloads"
	"github.com/rxdn/gdl/utils"
	"github.com/sirupsen/logrus"
	"time"
)

func (s *Shard) CountdownHeartbeat(ticker *time.Ticker) {
	loop:
	for {
		select {
		case <-s.killHeartbeat:
			ticker.Stop()
			break loop
		case <-ticker.C:
			s.heartbeatLock.RLock()

			// Check we received an ACK
			timeElapsed := s.lastHeartbeatAcknowledgement - s.lastHeartbeat
			if s.hasDoneHeartbeat && timeElapsed > int64(s.heartbeatInterval) {
				logrus.Warnf("shard %d didn't receive acknowledgement, restarting", s.ShardId)
				s.heartbeatLock.RUnlock()
				s.Kill()
				go s.EnsureConnect()
				return
			}

			s.heartbeatLock.RUnlock()

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

	s.heartbeatLock.Lock()
	s.hasDoneHeartbeat = true
	s.lastHeartbeat = utils.GetCurrentTimeMillis()
	s.heartbeatLock.Unlock()

	return s.write(payload)
}
