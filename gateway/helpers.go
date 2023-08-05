package gateway

import "context"

func (s *Shard) SelfId() uint64 {
	self, _ := s.Cache.GetSelf(context.Background())
	return self.Id
}

func (s *Shard) SelfAvatar(size int) string {
	self, _ := s.Cache.GetSelf(context.Background())
	return self.AvatarUrl(size)
}

func (s *Shard) SelfUsername() string {
	self, _ := s.Cache.GetSelf(context.Background())
	return self.Username
}

// millis
func (s *Shard) HeartbeatLatency() int64 {
	return s.lastHeartbeatAcknowledgement - s.lastHeartbeat
}
