package cache

import (
	"github.com/rxdn/gdl/objects"
	"sync"
)

type MemoryCache struct {
	Options CacheOptions
	locks map[uint64]*sync.RWMutex

	users map[uint64]*objects.User
	guilds map[uint64]*objects.Guild
	channels map[uint64]*objects.Channel
	roles map[uint64]*objects.Role
	emojis map[uint64]*objects.Emoji
	voiceStates map[uint64]*objects.VoiceState
	self *objects.User
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

func (c *MemoryCache) StoreUser(user *objects.User) {
	c.users[user.Id] = user
}

func (c *MemoryCache) GetUser(id uint64) *objects.User {
	lock := c.GetLock(id)
	lock.RLock()
	user := c.users[id]
	lock.RUnlock()
	return user
}

func (c *MemoryCache) StoreGuild(guild *objects.Guild) {
	c.guilds[guild.Id] = guild
}

func (c *MemoryCache) GetGuild(id uint64) *objects.Guild {
	lock := c.GetLock(id)
	lock.RLock()
	guild := c.guilds[id]
	lock.RUnlock()
	return guild
}

func (c *MemoryCache) DeleteGuild(id uint64) {
	delete(c.guilds, id)
}

func (c *MemoryCache) StoreChannel(channel *objects.Channel) {
	c.channels[channel.Id] = channel
}

func (c *MemoryCache) GetChannel(id uint64) *objects.Channel {
	lock := c.GetLock(id)
	lock.RLock()
	channel := c.channels[id]
	lock.RUnlock()
	return channel
}

func (c *MemoryCache) DeleteChannel(id uint64) {
	delete(c.channels, id)
}

func (c *MemoryCache) StoreRole(role *objects.Role) {
	c.roles[role.Id] = role
}

func (c *MemoryCache) GetRole(id uint64) *objects.Role {
	lock := c.GetLock(id)
	lock.RLock()
	role := c.roles[id]
	lock.RUnlock()
	return role
}

func (c *MemoryCache) DeleteRole(id uint64) {
	delete(c.roles, id)
}

func (c *MemoryCache) StoreEmoji(emoji *objects.Emoji) {
	c.emojis[emoji.Id] = emoji
}

func (c *MemoryCache) GetEmoji(id uint64) *objects.Emoji {
	lock := c.GetLock(id)
	lock.RLock()
	emoji := c.emojis[id]
	lock.RUnlock()
	return emoji
}

func (c *MemoryCache) DeleteEmoji(id uint64) {
	delete(c.emojis, id)
}

func (c *MemoryCache) StoreVoiceState(voiceState *objects.VoiceState) {
	c.voiceStates[voiceState.UserId] = voiceState
}

func (c *MemoryCache) GetVoiceState(user uint64) *objects.VoiceState {
	lock := c.GetLock(user)
	lock.RLock()
	voiceState := c.voiceStates[user]
	lock.RUnlock()
	return voiceState
}

func (c *MemoryCache) DeleteVoiceState(user uint64) {
	delete(c.voiceStates, user)
}

func (c *MemoryCache) StoreSelf(self *objects.User) {
	c.self = self
}

func (c *MemoryCache) GetSelf() *objects.User {
	return c.self
}
