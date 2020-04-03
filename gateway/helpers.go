package gateway

func (s *Shard) SelfId() uint64 {
	self := s.Cache.GetSelf()
	if self != nil {
		return self.Id
	}
	return 0
}

func (s *Shard) SelfAvatar(size int) string {
	self := s.Cache.GetSelf()
	if self != nil {
		return self.AvatarUrl(size)
	}
	return ""
}

func (s *Shard) SelfUsername() string {
	self := s.Cache.GetSelf()
	if self != nil {
		return self.Username
	}
	return ""
}

// millis
func (s *Shard) HeartbeatLatency() int64 {
	return s.LastHeartbeatAcknowledgement - s.LastHeartbeat
}
