package gateway

import (
	"github.com/Dot-Rar/gdl/cache"
	"github.com/Dot-Rar/gdl/gateway/payloads/events"
	"github.com/juju/ratelimit"
	"time"
)

type ShardManager struct {
	Token string

	GatewayBucket *ratelimit.Bucket

	TotalShards  int
	MinimumShard int // Inclusive
	MaximumShard int // Inclusive

	Shards map[int]*Shard

	EventBus *events.EventBus

	CacheFactory cache.CacheFactory
}

func NewShardManager(token string, totalShards, minimumShard, maximumShard int, cacheFactory cache.CacheFactory) ShardManager {
	manager := ShardManager{
		Token:         token,
		GatewayBucket: ratelimit.NewBucket(time.Second * 5, 1),

		TotalShards: totalShards,
		MinimumShard: minimumShard,
		MaximumShard: maximumShard,

		EventBus: events.NewEventBus(),
		CacheFactory: cacheFactory,
	}

	shards := make(map[int]*Shard)
	for i := minimumShard; i <= maximumShard; i++ {
		shard := NewShard(&manager, token, i)
		shards[i] = &shard
	}

	manager.Shards = shards

	RegisterCacheListeners(&manager)

	return manager
}

func (sm *ShardManager) Connect() error {
	for _, shard := range sm.Shards {
		go shard.EnsureConnect()
	}

	return nil
}

func (sm *ShardManager) RegisterListeners(listeners ...interface{}) {
	for _, listener := range listeners {
		sm.EventBus.RegisterListener(listener)
	}
}
