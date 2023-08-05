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

	gatewayBucketLock sync.RWMutex
	gatewayBuckets    []*ratelimit.Bucket

	globalLock       sync.RWMutex
	globalResetAfter time.Time

	bucketLock sync.RWMutex
	buckets    map[RouteId]string
}

func NewMemoryStore() *MemoryStore {
	cache := ttlcache.NewCache()
	cache.SkipTtlExtensionOnHit(true)

	return &MemoryStore{
		Cache: cache,
	}
}

func (s *MemoryStore) getTTLAndDecrease(route Route) (time.Duration, error) {
	key, err := getKey(s, route)
	if err != nil {
		return 0, err
	}

	s.Lock()
	defer s.Unlock()

	item, found, _ := s.Cache.GetItem(key)

	if found {
		remaining := item.Data.(int)
		ttl := item.ExpireAt.Sub(time.Now())

		if remaining > 0 {
			s.Cache.SetWithTTL(key, remaining-1, ttl)
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

func (s *MemoryStore) UpdateRateLimit(route Route, remaining int, resetAfter time.Duration) error {
	key, err := getKey(s, route)
	if err != nil {
		return err
	}

	s.Lock()
	defer s.Unlock()

	s.Cache.SetWithTTL(key, remaining, resetAfter)
	return nil
}

func (s *MemoryStore) UpdateGlobalRateLimit(resetAfter time.Duration) {
	s.globalLock.Lock()
	s.globalResetAfter = time.Now().Add(resetAfter)
	s.globalLock.Unlock()
}

func (s *MemoryStore) identifyWait(shardId int, largeShardingBuckets int) error {
	s.gatewayBucketLock.Lock()

	if len(s.gatewayBuckets) < largeShardingBuckets {
		for i := len(s.gatewayBuckets); i < largeShardingBuckets; i++ {
			s.gatewayBuckets = append(s.gatewayBuckets, ratelimit.NewBucket(IdentifyWait, 1))
		}
	}

	bucket := s.gatewayBuckets[shardId%largeShardingBuckets]
	s.gatewayBucketLock.Unlock()

	bucket.Wait(1)
	return nil
}

func (s *MemoryStore) getBucket(routeId RouteId) (string, error) {
	s.bucketLock.RLock()
	defer s.bucketLock.RUnlock()

	if bucket, ok := s.buckets[routeId]; ok {
		return bucket, nil
	} else {
		return "", nil
	}
}

func (s *MemoryStore) setBucket(routeId RouteId, bucket string) error {
	s.bucketLock.Lock()
	defer s.bucketLock.Unlock()

	s.buckets[routeId] = bucket
	return nil
}
