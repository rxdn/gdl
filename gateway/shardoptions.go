package gateway

import (
	"github.com/rxdn/gdl/cache"
	"github.com/rxdn/gdl/gateway/intents"
	"github.com/rxdn/gdl/objects/user"
	"github.com/rxdn/gdl/rest/ratelimit"
)

type ShardOptions struct {
	ShardCount           ShardCount
	CacheFactory         cache.CacheFactory
	RateLimitStore       ratelimit.RateLimitStore
	GuildSubscriptions   bool
	Presence             user.UpdateStatus
	Hooks                Hooks
	Debug                bool
	Intents              []intents.Intent
	LargeShardingBuckets int // defaults to 1. don't touch unless discord tell you to
}

type ShardCount struct {
	Total   int
	Lowest  int // Inclusive
	Highest int // Exclusive
}
