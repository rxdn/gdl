package cache

import (
	"github.com/Dot-Rar/gdl/objects"
)

type CacheFactory func() Cache

func MemoryCacheFactory(options CacheOptions) CacheFactory {
	return func() Cache {
		return &MemoryCache{
			Options:  options,
			users:    make(map[uint64]*objects.User),
			guilds:   make(map[uint64]*objects.Guild),
			channels: make(map[uint64]*objects.Channel),
			roles:    make(map[uint64]*objects.Role),
			emojis:   make(map[uint64]*objects.Emoji),
		}
	}
}
