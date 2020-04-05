package gateway

import (
	"github.com/juju/ratelimit"
	"github.com/rxdn/gdl/gateway/payloads/events"
	restlimiter "github.com/rxdn/gdl/rest/ratelimit"
	"os"
	"os/signal"
	"syscall"
	"time"
)

type ShardManager struct {
	Token string

	GatewayBucket      *ratelimit.Bucket
	ConcurrencyLimiter *restlimiter.Ratelimiter

	ShardOptions ShardOptions
	Shards       map[int]*Shard

	EventBus *events.EventBus
}

func NewShardManager(token string, shardOptions ShardOptions) ShardManager {
	manager := ShardManager{
		Token:              token,
		GatewayBucket:      ratelimit.NewBucket(time.Second*6, 1),
		ConcurrencyLimiter: restlimiter.NewConcurrencyLimiter(shardOptions.RateLimitStore),
		ShardOptions:       shardOptions,
		EventBus:           events.NewEventBus(),
	}

	manager.Shards = make(map[int]*Shard)
	for i := shardOptions.ShardCount.Lowest; i < shardOptions.ShardCount.Highest; i++ {
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
