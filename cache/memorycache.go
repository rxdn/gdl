package cache

import (
	"github.com/Dot-Rar/gdl/objects"
)

type MemoryCache struct {
	Options CacheOptions

	users map[uint64]*objects.User
	guilds map[uint64]*objects.Guild
	channels map[uint64]*objects.Channel
	roles map[uint64]*objects.Role
	emojis map[uint64]*objects.Emoji
}

func (c *MemoryCache) StoreUser(user *objects.User) {
	if c.Options.Users {
		c.users[user.Id] = user
	}
}

func (c *MemoryCache) GetUser(id uint64) *objects.User {
	return c.users[id]
}

func (c *MemoryCache) StoreGuild(guild *objects.Guild) {
	if c.Options.Guilds && guild != nil {
		modified := *guild
		if !c.Options.Channels {
			modified.Channels = make([]*objects.Channel, 0)
		}
		if !c.Options.Users {
			modified.Members = make([]*objects.Member, 0)
		}
		if !c.Options.Emojis {
			modified.Emojis = make([]*objects.Emoji, 0)
		}
		if !c.Options.Roles {
			modified.Roles = make([]*objects.Role, 0)
		}

		c.guilds[guild.Id] = &modified
	}
}

func (c *MemoryCache) GetGuild(id uint64) *objects.Guild {
	return c.guilds[id]
}

func (c *MemoryCache) StoreChannel(channel *objects.Channel) {
	if c.Options.Channels {
		c.channels[channel.Id] = channel
	}
}

func (c *MemoryCache) GetChannel(id uint64) *objects.Channel {
	return c.channels[id]
}

func (c *MemoryCache) StoreRole(role *objects.Role) {
	if c.Options.Roles {
		c.roles[role.Id] = role
	}
}

func (c *MemoryCache) GetRole(id uint64) *objects.Role {
	return c.roles[id]
}

func (c *MemoryCache) StoreEmoji(emoji *objects.Emoji) {
	if c.Options.Channels {
		c.emojis[emoji.Id] = emoji
	}
}

func (c *MemoryCache) GetEmoji(id uint64) *objects.Emoji {
	return c.emojis[id]
}