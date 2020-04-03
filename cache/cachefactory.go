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
			Options:      options,
			Mutex:        &sync.Mutex{},
			users:        make(map[uint64]*user.User),
			usersLock:    &sync.RWMutex{},
			guilds:       make(map[uint64]guild.Guild),
			guildsLock:   &sync.RWMutex{},
			channels:     make(map[uint64]*channel.Channel),
			channelsLock: &sync.RWMutex{},
			roles:        make(map[uint64]*guild.Role),
			rolesLock:    &sync.RWMutex{},
			emojis:       make(map[uint64]*emoji.Emoji),
			emojisLock:   &sync.RWMutex{},
			selfLock:     &sync.RWMutex{},
		}
	}
}
