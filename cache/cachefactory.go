package cache

import (
	"github.com/rxdn/gdl/objects/channel"
	"github.com/rxdn/gdl/objects/guild"
	"github.com/rxdn/gdl/objects/guild/emoji"
	"github.com/rxdn/gdl/objects/user"
	"sync"
)

type CacheFactory func() Cache

func MemoryCacheFactory(options CacheOptions) CacheFactory {
	return func() Cache {
		return &MemoryCache{
			Options:     options,
			locks:       make(map[uint64]*sync.RWMutex),
			users:       make(map[uint64]*user.User),
			guilds:      make(map[uint64]*guild.Guild),
			channels:    make(map[uint64]*channel.Channel),
			roles:       make(map[uint64]*guild.Role),
			emojis:      make(map[uint64]*emoji.Emoji),
			voiceStates: make(map[uint64]*guild.VoiceState),
		}
	}
}
