package ratelimit

import (
	"fmt"
	"time"
)

const IdentifyWait = 6 * time.Second

type RateLimitStore interface {
	getTTLAndDecrease(Route) (time.Duration, error)
	getGlobalTTL() (time.Duration, error)
	UpdateRateLimit(route Route, remaining int, resetAfter time.Duration) error
	UpdateGlobalRateLimit(resetAfter time.Duration)
	identifyWait(shardId int, largeShardingBuckets int) error
	getBucket(routeId RouteId) (string, error)
	setBucket(routeId RouteId, bucket string) error
}

func getKey(store RateLimitStore, route Route) (string, error) {
	bucket, err := store.getBucket(route.Id)
	if err != nil {
		return "", err
	}

	var key string
	if bucket == "" {
		key = fmt.Sprintf("%d:%d:%d", route.Type, route.Id, route.Snowflake)
	} else {
		key = bucket
	}

	return key, nil
}
