package gateway

import (
	"github.com/Dot-Rar/gdl/gateway/payloads"
	"github.com/Dot-Rar/gdl/utils"
	"github.com/sirupsen/logrus"
	"time"
)

func (s *Shard) CountdownHeartbeat(ticker *time.Ticker) {
	for {
		select {
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
		case <-s.KillHeartbeat:
			ticker.Stop()
			break
		}
	}
}

func (s *Shard) Heartbeat() error {
	s.SequenceLock.RLock()
	payload := payloads.NewHeartbeat(s.SequenceNumber)
	s.SequenceLock.RUnlock()

	s.HasDoneHeartbeat = true

	return s.Write(payload)
}
