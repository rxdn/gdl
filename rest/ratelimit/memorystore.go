package ratelimit

import (
	"github.com/TicketsBot/ttlcache"
	"sync"
	"time"
)

type MemoryStore struct {
	sync.Mutex
	Cache *ttlcache.Cache // handles mutex for us
}

func NewMemoryStore() *MemoryStore {
	cache := ttlcache.NewCache()
	return &MemoryStore{
		Cache: cache,
	}
}

func (s *MemoryStore) getTTLAndDecrease(endpoint string) (time.Duration, error) {
	s.Lock()
	defer s.Unlock()

	item, found, _ := s.Cache.GetItem(endpoint)

	if found {
		remaining := item.Data.(int)
		ttl := item.ExpireAt.Sub(time.Now())

		s.Cache.SetWithTTL(endpoint, remaining - 1, ttl)

		if remaining > 0 {
			return 0, nil
		} else {
			return ttl, nil
		}
	} else { // no bucket is found, obviously not ratelimited yet
		return 0, nil
	}
}

func (s *MemoryStore) UpdateRateLimit(endpoint string, remaining int, resetAfter time.Duration) {
	s.Lock()
	s.Cache.SetWithTTL(endpoint, remaining, resetAfter)
	s.Unlock()
}
