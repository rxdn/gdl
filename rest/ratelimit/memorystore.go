package ratelimit

import (
	"github.com/TicketsBot/ttlcache"
	"github.com/juju/ratelimit"
	"sync"
	"time"
)

type MemoryStore struct {
	sync.Mutex
	Cache *ttlcache.Cache // handles mutex for us

	bucketLock     sync.RWMutex
	gatewayBuckets []*ratelimit.Bucket

	globalLock       sync.RWMutex
	globalResetAfter time.Time
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

		s.Cache.SetWithTTL(endpoint, remaining-1, ttl)

		if remaining > 0 {
			return 0, nil
		} else {
			return ttl, nil
		}
	} else { // no bucket is found, obviously not ratelimited yet
		return 0, nil
	}
}

func (s *MemoryStore) getGlobalTTL() (time.Duration, error) {
	s.globalLock.RLock()
	defer s.globalLock.RUnlock()
	return s.globalResetAfter.Sub(time.Now()), nil
}

func (s *MemoryStore) UpdateRateLimit(endpoint string, remaining int, resetAfter time.Duration) {
	s.Lock()
	s.Cache.SetWithTTL(endpoint, remaining, resetAfter)
	s.Unlock()
}

func (s *MemoryStore) UpdateGlobalRateLimit(resetAfter time.Duration) {
	s.globalLock.Lock()
	s.globalResetAfter = time.Now().Add(resetAfter)
	s.globalLock.Unlock()
}

func (s *MemoryStore) identifyWait(shardId int, largeShardingBuckets int) error {
	s.bucketLock.Lock()

	if len(s.gatewayBuckets) < largeShardingBuckets {
		for i := len(s.gatewayBuckets); i < largeShardingBuckets; i++ {
			s.gatewayBuckets = append(s.gatewayBuckets, ratelimit.NewBucket(IdentifyWait, 1))
		}
	}

	bucket := s.gatewayBuckets[shardId%largeShardingBuckets]
	s.bucketLock.Unlock()

	bucket.Wait(1)
	return nil
}
