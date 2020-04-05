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

func (s *MemoryStore) getTTL(endpoint string) (time.Duration, error) {
	s.Lock()
	defer s.Unlock()

	item, found, _ := s.Cache.GetItem(endpoint)

	if found {
		remaining := item.Data.(int)

		if remaining > 0 {
			return 0, nil
		} else {
			return item.ExpireAt.Sub(time.Now()), nil
		}
	} else { // no bucket is found, obviously not ratelimited yet
		return 0, nil
	}
}

func (s *MemoryStore) UpdateRateLimit(endpoint string, remaining int, resetAfter time.Duration) {
	s.Cache.SetWithTTL(endpoint, remaining, resetAfter)
}
