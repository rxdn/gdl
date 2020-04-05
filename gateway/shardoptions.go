package gateway

import (
	"github.com/rxdn/gdl/cache"
	"github.com/rxdn/gdl/objects/user"
	"github.com/rxdn/gdl/rest/ratelimit"
)

type ShardOptions struct {
	ShardCount         ShardCount
	CacheFactory       cache.CacheFactory
	RateLimitStore     ratelimit.RateLimitStore
	GuildSubscriptions bool
	Presence           user.UpdateStatus
}

type ShardCount struct {
	Total   int
	Lowest  int // Inclusive
	Highest int // Exclusive
}
