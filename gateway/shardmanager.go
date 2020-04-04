package gateway

import (
	"github.com/juju/ratelimit"
	"github.com/rxdn/gdl/cache"
	"github.com/rxdn/gdl/gateway/payloads/events"
	"github.com/rxdn/gdl/objects/user"
	"os"
	"os/signal"
	"syscall"
	"time"
)

type ShardManager struct {
	Token string

	GatewayBucket *ratelimit.Bucket

	ShardOptions ShardOptions
	Shards       map[int]*Shard

	EventBus *events.EventBus

	CacheFactory cache.CacheFactory

	Presence           user.UpdateStatus
	GuildSubscriptions bool
}

func NewShardManager(token string, shardOptions ShardOptions, cacheFactory cache.CacheFactory) ShardManager {
	manager := ShardManager{
		Token:         token,
		GatewayBucket: ratelimit.NewBucket(time.Second*6, 1),
		ShardOptions:  shardOptions,
		EventBus:      events.NewEventBus(),
		CacheFactory:  cacheFactory,
	}

	manager.Shards = make(map[int]*Shard)
	for i := shardOptions.Lowest; i < shardOptions.Highest; i++ {
		shard := NewShard(&manager, token, i)
		manager.Shards[i] = &shard
	}

	RegisterCacheListeners(&manager)

	return manager
}

func (sm *ShardManager) Connect() {
	for _, shard := range sm.Shards {
		go shard.EnsureConnect()
	}
}

func (sm *ShardManager) RegisterListeners(listeners ...interface{}) {
	for _, listener := range listeners {
		sm.EventBus.RegisterListener(listener)
	}
}

func (sm *ShardManager) ShardForGuild(guildId uint64) *Shard {
	for _, shard := range sm.Shards {
		if shard.Cache.GetGuild(guildId) != nil {
			return shard
		}
	}

	return nil
}

func (sm *ShardManager) WaitForInterrupt() {
	ch := make(chan os.Signal, 1)
	signal.Notify(ch, syscall.SIGINT, syscall.SIGTERM, os.Interrupt, os.Kill)
	<-ch
}
