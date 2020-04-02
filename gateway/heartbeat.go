package gateway

import (
	"github.com/rxdn/gdl/gateway/payloads"
	"github.com/rxdn/gdl/utils"
	"github.com/sirupsen/logrus"
	"time"
)

func (s *Shard) CountdownHeartbeat(ticker *time.Ticker) {
	for {
		select {
		case <-s.KillHeartbeat:
			ticker.Stop()
			break
		case <-ticker.C:
			s.HeartbeatMutex.Lock()

			// Check we received an ACK
			timeElapsed := utils.GetCurrentTimeMillis() - s.LastHeartbeatAcknowledgement
			if s.HasDoneHeartbeat && timeElapsed > int64(s.HeartbeatInterval) {
				logrus.Warnf("shard %d didn't receive acknowledgement, restarting", s.ShardId)
				s.Kill()
				go s.EnsureConnect()
			}
			s.HeartbeatMutex.Unlock()

			if err := s.Heartbeat(); err != nil {
				logrus.Warnf("shard %d didn't heartbeat failed, restarting: %s", s.ShardId, err.Error())
				s.Kill()
				go s.EnsureConnect()
			}
		}
	}
}

func (s *Shard) Heartbeat() error {
	s.SequenceLock.RLock()
	payload := payloads.NewHeartbeat(s.SequenceNumber)
	s.SequenceLock.RUnlock()

	s.HasDoneHeartbeat = true
	s.LastHeartbeat = utils.GetCurrentTimeMillis()

	return s.Write(payload)
}
