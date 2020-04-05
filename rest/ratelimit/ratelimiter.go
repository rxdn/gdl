package ratelimit

import (
	"sync"
	"time"
)

// Big thanks to https://github.com/spencersharkey for sharing his ratelimiter with me

type Ratelimiter struct {
	sync.Mutex
	Store           RateLimitStore
}

func NewConcurrencyLimiter(store RateLimitStore) *Ratelimiter {
	return &Ratelimiter{
		Store:           store,
	}
}

func (l *Ratelimiter) ExecuteCall(endpoint string, ch chan error) {
	ttl, err := l.Store.getTTL(endpoint)
	if err != nil { // if an error occurred, we should cancel the request
		ch <- err
		return
	}

	if ttl > 0 {
		<-time.After(ttl)
		l.ExecuteCall(endpoint, ch)
	} else {
		ch <- nil
	}
}
