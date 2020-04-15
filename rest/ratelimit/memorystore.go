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

	bucketLock     sync.Mutex
	gatewayBuckets []*ratelimit.Bucket
}

func NewMemoryStore() *MemoryStore {
	buckets := make([]*ratelimit.Bucket, 16)
	for i := 0; i < 16; i++ {
		buckets[i] = ratelimit.NewBucket(IdentifyWait, 1)
	}

	cache := ttlcache.NewCache()
	return &MemoryStore{
		Cache:          cache,
		gatewayBuckets: buckets,
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

func (s *MemoryStore) UpdateRateLimit(endpoint string, remaining int, resetAfter time.Duration) {
	s.Lock()
	s.Cache.SetWithTTL(endpoint, remaining, resetAfter)
	s.Unlock()
}

func (s *MemoryStore) identifyWait(shardId int, largeSharding bool) error {
	var bucket *ratelimit.Bucket
	if largeSharding {
		bucket = s.gatewayBuckets[shardId % 16]
	} else {
		bucket = s.gatewayBuckets[0]
	}

	bucket.Wait(1)
	return nil
}
