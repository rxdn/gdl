package cache

import (
	"fmt"
	"github.com/rxdn/gdl/objects/channel"
	"github.com/rxdn/gdl/objects/guild"
	"github.com/rxdn/gdl/objects/guild/emoji"
	"github.com/rxdn/gdl/objects/member"
	"github.com/rxdn/gdl/objects/user"
	"sync"
)

type MemoryCache struct {
	options CacheOptions

	users    map[uint64]user.CachedUser
	userLock sync.RWMutex

	guilds    map[uint64]guild.CachedGuild
	guildLock sync.RWMutex

	// guilds -> users
	members    map[uint64]map[uint64]member.CachedMember
	memberLock sync.RWMutex

	channels    map[uint64]channel.CachedChannel
	channelLock sync.RWMutex

	roles    map[uint64]guild.CachedRole
	roleLock sync.RWMutex

	emojis    map[uint64]emoji.CachedEmoji
	emojiLock sync.RWMutex

	// guilds -> users
	voiceStates    map[uint64]map[uint64]guild.CachedVoiceState
	voiceStateLock sync.RWMutex

	selfLock sync.RWMutex
	self     user.User
}

func NewMemoryCache(cacheOptions CacheOptions) MemoryCache {
	return MemoryCache{
		options:     cacheOptions,
		users:       make(map[uint64]user.CachedUser),
		guilds:      make(map[uint64]guild.CachedGuild),
		members:     make(map[uint64]map[uint64]member.CachedMember),
		channels:    make(map[uint64]channel.CachedChannel),
		roles:       make(map[uint64]guild.CachedRole),
		emojis:      make(map[uint64]emoji.CachedEmoji),
		voiceStates: make(map[uint64]map[uint64]guild.CachedVoiceState),
	}
}

func (c *MemoryCache) Size() string {
	return fmt.Sprintf("%d channels %d members %d roles %d emojis %d voice states %d users", len(c.channels), len(c.members), len(c.roles), len(c.emojis), len(c.voiceStates), len(c.users))
}

func (c *MemoryCache) GetOptions() CacheOptions {
	return c.options
}

func (c *MemoryCache) StoreUser(u user.User) {
	c.StoreUsers([]user.User{u})
}

func (c *MemoryCache) StoreUsers(users []user.User) {
	if c.options.Users {
		c.userLock.Lock()
		defer c.userLock.Unlock()

		for _, user := range users {
			c.users[user.Id] = user.ToCachedUser()
		}
	}
}

func (c *MemoryCache) GetUser(userId uint64) (u user.User, found bool) {
	c.userLock.RLock()
	defer c.userLock.RUnlock()

	var cached user.CachedUser
	cached, found = c.users[userId]
	if found {
		return cached.ToUser(userId), found
	} else {
		return
	}
}

func (c *MemoryCache) StoreGuild(g guild.Guild) {
	c.StoreGuilds([]guild.Guild{g})
}

func (c *MemoryCache) StoreGuilds(guilds []guild.Guild) {
	if c.options.Guilds {
		c.guildLock.Lock()

		for _, guild := range guilds {
			cached := guild.ToCachedGuild()

			// remove unnecessary fields
			if !c.options.Channels {
				cached.Channels = nil
			}
			if !c.options.Roles {
				cached.Roles = nil
			}
			if !c.options.Emojis {
				cached.Emojis = nil
			}

			c.guilds[guild.Id] = cached
		}

		c.guildLock.Unlock()
	}

	for _, guild := range guilds {
		c.StoreChannels(guild.Channels)
		c.StoreMembers(guild.Members, guild.Id)
		c.StoreRoles(guild.Roles, guild.Id)
		c.StoreEmojis(guild.Emojis, guild.Id)
		c.StoreVoiceStates(guild.VoiceStates)
	}
}

func (c *MemoryCache) GetGuild(guildId uint64, withMembers bool) (g guild.Guild, found bool) {
	var cached guild.CachedGuild

	c.guildLock.RLock()
	cached, found = c.guilds[guildId]
	c.guildLock.RUnlock()

	if found {
		g = cached.ToGuild(guildId)

		// re-add fields
		if withMembers {
			g.Members = c.GetGuildMembers(guildId, false)
		}

		g.Channels = c.GetGuildChannels(guildId)
		g.Roles = c.GetGuildRoles(guildId)
		g.Emojis = c.GetGuildEmojis(guildId)
		g.VoiceStates = c.GetGuildVoiceStates(guildId)
	}

	return
}

// You should never use this
func (c *MemoryCache) GetGuilds() (guilds []guild.Guild) {
	c.guildLock.RLock()
	defer c.guildLock.RUnlock()

	guilds = make([]guild.Guild, len(c.guilds))
	i := 0

	for guildId, guild := range c.guilds {
		guilds[i] = guild.ToGuild(guildId)
		i++
	}

	return
}

func (c *MemoryCache) DeleteGuild(guildId uint64) {
	c.guildLock.Lock()
	delete(c.guilds, guildId)
	c.guildLock.Unlock()
}

func (c *MemoryCache) GetGuildCount() (count int) {
	c.guildLock.RLock()
	count = len(c.guilds)
	c.guildLock.RUnlock()
	return
}

func (c *MemoryCache) StoreMember(m member.Member, guildId uint64) {
	c.StoreMembers([]member.Member{m}, guildId)
}

func (c *MemoryCache) StoreMembers(members []member.Member, guildId uint64) {
	if c.options.Members {
		c.memberLock.Lock()
		defer c.memberLock.Unlock()

		if c.members[guildId] == nil {
			c.members[guildId] = make(map[uint64]member.CachedMember)
		}

		for _, member := range members {
			c.members[guildId][member.User.Id] = member.ToCachedMember()
		}
	}
}

func (c *MemoryCache) GetMember(guildId, userId uint64) (m member.Member, found bool) {
	var cached member.CachedMember

	c.memberLock.RLock()

	if c.members[guildId] == nil {
		c.memberLock.RUnlock()
		return
	}

	cached, found = c.members[guildId][userId]
	if found {
		u, userFound := c.GetUser(userId)
		if !userFound {
			u = user.User{Id: userId}
		}

		m = cached.ToMember(u)
	}

	c.memberLock.RUnlock()
	return
}

func (c *MemoryCache) GetGuildMembers(guildId uint64, withUserData bool) (members []member.Member) {
	c.memberLock.RLock()
	defer c.memberLock.RUnlock()

	if c.members[guildId] == nil {
		return
	}

	for userId, cachedMember := range c.members[guildId] {
		// get user
		u := user.User{Id: userId}
		if withUserData {
			if cachedUser, foundUser := c.GetUser(userId); foundUser {
				u = cachedUser
			}
		}

		members = append(members, cachedMember.ToMember(u))
	}

	return
}

func (c *MemoryCache) DeleteMember(userId, guildId uint64) {
	c.memberLock.Lock()
	defer c.memberLock.Unlock()

	if c.members[guildId] != nil {
		delete(c.members[guildId], userId)
	}
}

func (c *MemoryCache) StoreChannel(ch channel.Channel) {
	c.StoreChannels([]channel.Channel{ch})
}

func (c *MemoryCache) StoreChannels(channels []channel.Channel) {
	if c.options.Channels {
		c.channelLock.Lock()
		defer c.channelLock.Unlock()

		for _, channel := range channels {
			c.channels[channel.Id] = channel.ToCachedChannel()

			// Add to guild object
			c.guildLock.Lock()
			if guild, found := c.guilds[channel.GuildId]; found {
				// Check to see if channel already exists
				var channelExists bool
				for _, userId := range guild.Channels {
					if userId == channel.Id {
						channelExists = true
						break
					}
				}

				if !channelExists {
					guild.Channels = append(guild.Channels, channel.Id)
					c.guilds[channel.GuildId] = guild
				}
			}
			c.guildLock.Unlock()
		}
	}
}

func (c *MemoryCache) GetChannel(channelId uint64) (ch channel.Channel, found bool) {
	var cached channel.CachedChannel

	c.channelLock.RLock()
	cached, found = c.channels[channelId]
	c.channelLock.RUnlock()

	ch = cached.ToChannel(channelId, cached.GuildId)
	return
}

func (c *MemoryCache) GetGuildChannels(guildId uint64) (channels []channel.Channel) {
	c.guildLock.RLock()
	guild, found := c.guilds[guildId]
	c.guildLock.RUnlock()
	if !found {
		return
	}

	c.channelLock.RLock()
	for _, channelId := range guild.Channels {
		if cachedChannel, found := c.channels[channelId]; found {
			channels = append(channels, cachedChannel.ToChannel(channelId, guildId))
		}
	}
	c.channelLock.RUnlock()

	return channels
}

func (c *MemoryCache) DeleteChannel(channelId uint64) {
	c.channelLock.Lock()
	cached, found := c.channels[channelId]
	delete(c.channels, channelId)
	c.channelLock.Unlock()

	if found {
		// delete from guild
		c.guildLock.Lock()
		if guild, found := c.guilds[cached.GuildId]; found {
			// iterate channels
			var updated bool
			for i, ch := range guild.Channels {
				if ch == channelId {
					updated = true
					guild.Channels[i] = guild.Channels[len(guild.Channels)-1]
					guild.Channels = guild.Channels[:len(guild.Channels)-1]
					break
				}
			}

			if updated {
				c.guilds[guild.Id] = guild
			}
		}
		c.guildLock.Unlock()
	}
}

func (c *MemoryCache) StoreRole(role guild.Role, guildId uint64) {
	c.StoreRoles([]guild.Role{role}, guildId)
}

func (c *MemoryCache) StoreRoles(roles []guild.Role, guildId uint64) {
	if c.options.Roles {
		c.roleLock.Lock()
		defer c.roleLock.Unlock()

		for _, role := range roles {
			c.roles[role.Id] = role.ToCachedRole(guildId)

			// Add to guild object
			c.guildLock.Lock()
			if guild, found := c.guilds[guildId]; found {
				// Check to see if role already exists
				var roleExists bool
				for _, roleId := range guild.Roles {
					if roleId == role.Id {
						roleExists = true
						break
					}
				}

				if !roleExists {
					guild.Roles = append(guild.Roles, role.Id)
					c.guilds[guildId] = guild
				}
			}
			c.guildLock.Unlock()
		}
	}
}

func (c *MemoryCache) GetRole(roleId uint64) (role guild.Role, found bool) {
	var cachedRole guild.CachedRole

	c.roleLock.RLock()
	defer c.roleLock.RUnlock()

	cachedRole, found = c.roles[roleId]
	role = cachedRole.ToRole(roleId)

	return
}

func (c *MemoryCache) GetGuildRoles(guildId uint64) (roles []guild.Role) {
	// get guild
	c.guildLock.RLock()
	guild, found := c.guilds[guildId]
	c.guildLock.RUnlock()

	if !found {
		return
	}

	c.roleLock.RLock()
	defer c.roleLock.RUnlock()

	for _, roleId := range guild.Roles {
		if role, found := c.roles[roleId]; found {
			roles = append(roles, role.ToRole(roleId))
		}
	}

	return roles
}

func (c *MemoryCache) DeleteRole(roleId uint64) {
	c.roleLock.Lock()
	cached, found := c.roles[roleId]
	delete(c.roles, roleId)
	c.roleLock.Unlock()

	if found {
		// delete from guild
		c.guildLock.Lock()
		if guild, found := c.guilds[cached.GuildId]; found {
			// iterate roles
			var updated bool
			for i, role := range guild.Roles {
				if roleId == role {
					updated = true
					guild.Roles[i] = guild.Roles[len(guild.Roles)-1]
					guild.Roles = guild.Roles[:len(guild.Roles)-1]
					break
				}
			}

			if updated {
				c.guilds[guild.Id] = guild
			}
		}
		c.guildLock.Unlock()
	}
}

func (c *MemoryCache) StoreEmoji(e emoji.Emoji, guildId uint64) {
	c.StoreEmojis([]emoji.Emoji{e}, guildId)
}

func (c *MemoryCache) StoreEmojis(emojis []emoji.Emoji, guildId uint64) {
	c.emojiLock.Lock()

	for _, emoji := range emojis {
		c.emojis[emoji.Id] = emoji.ToCachedEmoji(guildId)

		// Add to guild object
		c.guildLock.Lock()
		if guild, found := c.guilds[guildId]; found {
			// Check to see if emoji already exists
			var emojiExists bool
			for _, emojiId := range guild.Emojis {
				if emojiId == emoji.Id {
					emojiExists = true
					break
				}
			}

			if !emojiExists {
				guild.Emojis = append(guild.Emojis, emoji.Id)
				c.guilds[guildId] = guild
			}
		}
		c.guildLock.Unlock()
	}

	c.emojiLock.Unlock()
}

func (c *MemoryCache) GetEmoji(emojiId uint64) (e emoji.Emoji, found bool) {
	var cached emoji.CachedEmoji

	c.emojiLock.RLock()
	cached, found = c.emojis[emojiId]
	c.emojiLock.RUnlock()

	if !found {
		return
	}

	u, userFound := c.GetUser(cached.User)
	if !userFound {
		u = user.User{Id: cached.User}
	}

	e = cached.ToEmoji(emojiId, u)
	return
}

func (c *MemoryCache) GetGuildEmojis(guildId uint64) (emojis []emoji.Emoji) {
	// get guild
	c.guildLock.RLock()
	guild, found := c.guilds[guildId]
	c.guildLock.RUnlock()

	if !found {
		return
	}

	c.emojiLock.RLock()
	defer c.emojiLock.RUnlock()

	for _, emojiId := range guild.Emojis {
		cached, found := c.emojis[emojiId]
		if !found {
			continue
		}

		u, userFound := c.GetUser(cached.User)
		if !userFound {
			u = user.User{Id: cached.User}
		}

		emojis = append(emojis, cached.ToEmoji(emojiId, u))
	}

	return
}

func (c *MemoryCache) DeleteEmoji(emojiId uint64) {
	c.emojiLock.Lock()
	cached, found := c.emojis[emojiId]
	delete(c.emojis, emojiId)
	c.emojiLock.Unlock()

	if found {
		// delete from guild
		c.guildLock.Lock()
		if guild, found := c.guilds[cached.GuildId]; found {
			// iterate emojis
			var updated bool
			for i, emoji := range guild.Emojis {
				if emoji == emojiId {
					updated = true
					guild.Emojis[i] = guild.Emojis[len(guild.Emojis)-1]
					guild.Emojis = guild.Emojis[:len(guild.Emojis)-1]
					break
				}
			}

			if updated {
				c.guilds[guild.Id] = guild
			}
		}
		c.guildLock.Unlock()
	}
}

func (c *MemoryCache) StoreVoiceState(state guild.VoiceState) {
	c.StoreVoiceStates([]guild.VoiceState{state})
}

func (c *MemoryCache) StoreVoiceStates(states []guild.VoiceState) {
	if c.options.VoiceStates {
		c.voiceStateLock.Lock()
		defer c.voiceStateLock.Unlock()

		for _, state := range states {
			if c.voiceStates[state.GuildId] == nil {
				c.voiceStates[state.GuildId] = make(map[uint64]guild.CachedVoiceState)
			}

			c.voiceStates[state.GuildId][state.UserId] = state.ToCachedVoiceState()
		}
	}
}

func (c *MemoryCache) GetVoiceState(userId, guildId uint64) (state guild.VoiceState, found bool) {
	var cached guild.CachedVoiceState

	c.voiceStateLock.RLock()
	defer c.voiceStateLock.RUnlock()
	if c.voiceStates[guildId] == nil {
		return
	}

	cached, found = c.voiceStates[guildId][userId]
	if found {
		// get member
		m, memberFound := c.GetMember(guildId, userId)
		if !memberFound {
			m = member.Member{
				User: user.User{
					Id: userId,
				},
			}
		}

		state = cached.ToVoiceState(guildId, m)
	}

	return
}

func (c *MemoryCache) GetGuildVoiceStates(guildId uint64) (states []guild.VoiceState) {
	c.voiceStateLock.RLock()
	defer c.voiceStateLock.RUnlock()

	if c.voiceStates[guildId] == nil {
		return
	}

	for userId, cached := range c.voiceStates[guildId] {
		// get member
		m, memberFound := c.GetMember(guildId, userId)
		if !memberFound {
			m = member.Member{
				User: user.User{
					Id: userId,
				},
			}
		}

		states = append(states, cached.ToVoiceState(guildId, m))
	}

	return
}

func (c *MemoryCache) DeleteVoiceState(userId, guildId uint64) {
	c.voiceStateLock.Lock()
	defer c.voiceStateLock.Unlock()

	if c.voiceStates[guildId] == nil {
		return
	}

	delete(c.voiceStates[guildId], userId)
}

func (c *MemoryCache) StoreSelf(self user.User) {
	c.selfLock.Lock()
	c.self = self
	c.selfLock.Unlock()
}

func (c *MemoryCache) GetSelf() (user.User, bool) {
	c.selfLock.RLock()
	self := c.self
	c.selfLock.RUnlock()

	return self, self.Id != 0
}
