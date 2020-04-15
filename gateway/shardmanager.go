package gateway

import (
	"github.com/rxdn/gdl/gateway/payloads/events"
	"github.com/rxdn/gdl/rest/ratelimit"
	"github.com/rxdn/gdl/rest/request"
	"os"
	"os/signal"
	"syscall"
)

type ShardManager struct {
	Token string

	RateLimiter *ratelimit.Ratelimiter

	ShardOptions ShardOptions
	Shards       map[int]*Shard

	EventBus *events.EventBus
}

func NewShardManager(token string, shardOptions ShardOptions) *ShardManager {
	manager := &ShardManager{
		Token:        token,
		RateLimiter:  ratelimit.NewRateLimiter(shardOptions.RateLimitStore, shardOptions.LargeSharding),
		ShardOptions: shardOptions,
		EventBus:     events.NewEventBus(),
	}

	manager.Shards = make(map[int]*Shard)
	for i := shardOptions.ShardCount.Lowest; i < shardOptions.ShardCount.Highest; i++ {
		shard := NewShard(manager, token, i)
		manager.Shards[i] = &shard
	}

	request.Hook = shardOptions.Hooks.RestHook

	RegisterCacheListeners(manager)

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
		shard.guildsLock.RLock()
		for _, shardGuild := range shard.guilds {
			if shardGuild == guildId {
				shard.guildsLock.RUnlock()
				return shard
			}
		}
		shard.guildsLock.RUnlock()
	}

	return nil
}

func (sm *ShardManager) WaitForInterrupt() {
	ch := make(chan os.Signal, 1)
	signal.Notify(ch, syscall.SIGINT, syscall.SIGTERM, os.Interrupt, os.Kill)
	<-ch
}
