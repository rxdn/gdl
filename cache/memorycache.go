package cache

import (
	"github.com/rxdn/gdl/objects/channel"
	"github.com/rxdn/gdl/objects/guild"
	"github.com/rxdn/gdl/objects/guild/emoji"
	"github.com/rxdn/gdl/objects/user"
	"sync"
)

type MemoryCache struct {
	Options CacheOptions
	locks map[uint64]*sync.RWMutex

	users map[uint64]*user.User
	guilds map[uint64]*guild.Guild
	channels map[uint64]*channel.Channel
	roles map[uint64]*guild.Role
	emojis map[uint64]*emoji.Emoji
	voiceStates map[uint64]*guild.VoiceState
	self *user.User
}

func (c *MemoryCache) GetOptions() CacheOptions {
	return c.Options
}

func (c *MemoryCache) GetLock(id uint64) *sync.RWMutex {
	lock := c.locks[id]
	if lock == nil {
		lock = &sync.RWMutex{}
		c.locks[id] = lock
	}
	return lock
}

func (c *MemoryCache) StoreUser(user *user.User) {
	c.users[user.Id] = user
}

func (c *MemoryCache) GetUser(id uint64) *user.User {
	lock := c.GetLock(id)
	lock.RLock()
	user := c.users[id]
	lock.RUnlock()
	return user
}

func (c *MemoryCache) StoreGuild(guild *guild.Guild) {
	c.guilds[guild.Id] = guild
}

func (c *MemoryCache) GetGuild(id uint64) *guild.Guild {
	lock := c.GetLock(id)
	lock.RLock()
	guild := c.guilds[id]
	lock.RUnlock()
	return guild
}

func (c *MemoryCache) GetGuilds() []*guild.Guild {
	guilds := make([]*guild.Guild, 0)
	for _, guild := range c.guilds {
		guilds = append(guilds, guild)
	}
	return guilds
}

func (c *MemoryCache) DeleteGuild(id uint64) {
	delete(c.guilds, id)
}

func (c *MemoryCache) StoreChannel(channel *channel.Channel) {
	c.channels[channel.Id] = channel
}

func (c *MemoryCache) GetChannel(id uint64) *channel.Channel {
	lock := c.GetLock(id)
	lock.RLock()
	channel := c.channels[id]
	lock.RUnlock()
	return channel
}

func (c *MemoryCache) DeleteChannel(id uint64) {
	delete(c.channels, id)
}

func (c *MemoryCache) StoreRole(role *guild.Role) {
	c.roles[role.Id] = role
}

func (c *MemoryCache) GetRole(id uint64) *guild.Role {
	lock := c.GetLock(id)
	lock.RLock()
	role := c.roles[id]
	lock.RUnlock()
	return role
}

func (c *MemoryCache) DeleteRole(id uint64) {
	delete(c.roles, id)
}

func (c *MemoryCache) StoreEmoji(emoji *emoji.Emoji) {
	c.emojis[emoji.Id] = emoji
}

func (c *MemoryCache) GetEmoji(id uint64) *emoji.Emoji {
	lock := c.GetLock(id)
	lock.RLock()
	emoji := c.emojis[id]
	lock.RUnlock()
	return emoji
}

func (c *MemoryCache) DeleteEmoji(id uint64) {
	delete(c.emojis, id)
}

func (c *MemoryCache) StoreVoiceState(voiceState *guild.VoiceState) {
	c.voiceStates[voiceState.UserId] = voiceState
}

func (c *MemoryCache) GetVoiceState(user uint64) *guild.VoiceState {
	lock := c.GetLock(user)
	lock.RLock()
	voiceState := c.voiceStates[user]
	lock.RUnlock()
	return voiceState
}

func (c *MemoryCache) DeleteVoiceState(user uint64) {
	delete(c.voiceStates, user)
}

func (c *MemoryCache) StoreSelf(self *user.User) {
	c.self = self
}

func (c *MemoryCache) GetSelf() *user.User {
	return c.self
}
