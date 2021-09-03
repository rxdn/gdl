package ratelimit

import (
	"context"
	"fmt"
	"github.com/go-redis/redis/v8"
	"strconv"
	"time"
)

type RedisStore struct {
	*redis.Client
	keyPrefix      string
	ContextBuilder func() context.Context
}

func NewRedisStore(client *redis.Client, keyPrefix string) *RedisStore {
	return &RedisStore{
		Client:    client,
		keyPrefix: keyPrefix,
		ContextBuilder: func() context.Context {
			ctx, _ := context.WithTimeout(context.Background(), time.Second*5)
			return ctx
		},
	}
}

func (s *RedisStore) getTTLAndDecrease(route Route) (time.Duration, error) {
	hash, err := getKey(s, route)
	if err != nil {
		return 0, err
	}

	key := fmt.Sprintf("%s:route:%s", s.keyPrefix, hash)

	remainingStr, err := s.Get(s.ContextBuilder(), key).Result()
	if err != nil {
		if err == redis.Nil { // if the key isn't found, then we can't be ratelimited yet
			return 0, nil
		} else { // an actual error occurred
			return 0, err
		}
	}

	remaining, err := strconv.Atoi(remainingStr)
	if err != nil { // some unknown error occurred
		return 0, err
	}

	// If it errors, the key doesn't exist
	s.Decr(s.ContextBuilder(), key)

	if remaining > 0 {
		return 0, nil
	} else { // if we're out of requests, we need to check the TTL of the key
		return s.PTTL(s.ContextBuilder(), key).Result()
	}
}

func (s *RedisStore) UpdateRateLimit(route Route, remaining int, resetAfter time.Duration) error {
	hash, err := getKey(s, route)
	if err != nil {
		return err
	}

	key := fmt.Sprintf("%s:route:%s", s.keyPrefix, hash)
	return s.Set(context.Background(), key, remaining, resetAfter).Err()
}

func (s *RedisStore) UpdateGlobalRateLimit(resetAfter time.Duration) {
	key := fmt.Sprintf("%s:global", s.keyPrefix)
	s.Set(context.Background(), key, true, resetAfter)
}

func (s *RedisStore) getGlobalTTL() (time.Duration, error) {
	key := fmt.Sprintf("%s:global", s.keyPrefix)
	return s.PTTL(s.ContextBuilder(), key).Result()
}

func (s *RedisStore) identifyWait(shardId int, largeShardingBuckets int) error {
	key := fmt.Sprintf("%s:identify:%d", s.keyPrefix, shardId%largeShardingBuckets)

	set := false

	for !set {
		var err error
		set, err = s.SetNX(s.ContextBuilder(), key, 1, IdentifyWait).Result()
		if err != nil {
			return err
		}

		if !set {
			cooldown, err := s.PTTL(s.ContextBuilder(), key).Result()
			if err != nil && err != redis.Nil { // if err == redis.Nil, cooldown must've expired since running SET
				return err
			}

			<-time.After(cooldown)
		}
	}

	return nil
}

func (s *RedisStore) getBucket(routeId RouteId) (string, error) {
	key := fmt.Sprintf("%s:buckets:%d", s.keyPrefix, routeId)

	bucket, err := s.Get(s.ContextBuilder(), key).Result()
	if err == redis.Nil {
		err = nil
	}

	return bucket, err
}

func (s *RedisStore) setBucket(routeId RouteId, bucket string) error {
	key := fmt.Sprintf("%s:buckets:%d", s.keyPrefix, routeId)
	return s.Set(s.ContextBuilder(), key, bucket, 0).Err()
}
