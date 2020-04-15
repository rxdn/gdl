package ratelimit

import "time"

const IdentifyWait = 6 * time.Second

type RateLimitStore interface {
	getTTLAndDecrease(bucket string) (time.Duration, error)
	UpdateRateLimit(bucket string, remaining int, resetAfter time.Duration)
	identifyWait(shardId int, largeShardingBuckets int) error
}
