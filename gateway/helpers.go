package gateway

func (s *Shard) SelfId() uint64 {
	self, _ := s.Cache.GetSelf()
	return self.Id
}

func (s *Shard) SelfAvatar(size int) string {
	self, _ := s.Cache.GetSelf()
	return self.AvatarUrl(size)
}

func (s *Shard) SelfUsername() string {
	self ,_ := s.Cache.GetSelf()
	return self.Username
}

// millis
func (s *Shard) HeartbeatLatency() int64 {
	return s.LastHeartbeatAcknowledgement - s.LastHeartbeat
}

func (s *Shard) GetShardGuildIds() []uint64 {
	s.guildsLock.RLock()
	ids := s.guilds
	s.guildsLock.RUnlock()
	return ids
}
