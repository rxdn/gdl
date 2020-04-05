package ratelimit

import "time"

type RateLimitStore interface {
	getTTLAndDecrease(bucket string) (time.Duration, error)
	UpdateRateLimit(bucket string, remaining int, resetAfter time.Duration)
}
