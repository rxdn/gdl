package ratelimit

import (
	"fmt"
	"github.com/go-redis/redis"
	"strconv"
	"time"
)

type RedisStore struct {
	*redis.Client
	keyPrefix string
}

func NewRedisStore(client *redis.Client, keyPrefix string) *RedisStore {
	return &RedisStore{
		Client:    client,
		keyPrefix: keyPrefix,
	}
}

func (s *RedisStore) getTTLAndDecrease(endpoint string) (time.Duration, error) {
	key := fmt.Sprintf("%s:%s", s.keyPrefix, endpoint)

	remainingStr, err := s.Get(key).Result()
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

	s.Decr(key) // attempt to decrease, doesn't matter too much if it errors

	if remaining > 0 {
		return 0, nil
	} else { // if we're out of requests, we need to check the TTL of the key
		return s.PTTL(key).Result()
	}
}

func (s *RedisStore) UpdateRateLimit(endpoint string, remaining int, resetAfter time.Duration) {
	key := fmt.Sprintf("%s:%s", s.keyPrefix, endpoint)
	s.Set(key, remaining, resetAfter)
}

func (s *RedisStore) UpdateGlobalRateLimit(resetAfter time.Duration) {
	key := fmt.Sprintf("%s:global", s.keyPrefix)
	s.Set(key, true, resetAfter)
}

func (s *RedisStore) getGlobalTTL() (time.Duration, error) {
	key := fmt.Sprintf("%s:global", s.keyPrefix)
	return s.PTTL(key).Result()
}

func (s *RedisStore) identifyWait(shardId int, largeShardingBuckets int) error {
	key := fmt.Sprintf("%s:identify:%d", s.keyPrefix, shardId % largeShardingBuckets)

	set := false

	for !set {
		var err error
		set, err = s.SetNX(key, 1, IdentifyWait).Result()
		if err != nil {
			return err
		}

		if !set {
			cooldown, err := s.PTTL(key).Result()
			if err != nil && err != redis.Nil { // if err == redis.Nil, cooldown must've expired since running SET
				return err
			}

			<- time.After(cooldown)
		}
	}

	return nil
}
