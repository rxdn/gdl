package ratelimit

import "time"

type RateLimitStore interface {
	getTTLAndDecrease(endpoint string) (time.Duration, error)
	UpdateRateLimit(endpoint string, remaining int, resetAfter time.Duration)
}
