package cache

import (
	"github.com/rxdn/gdl/objects/channel"
	"github.com/rxdn/gdl/objects/guild"
	"github.com/rxdn/gdl/objects/guild/emoji"
	"github.com/rxdn/gdl/objects/member"
	"github.com/rxdn/gdl/objects/user"
	"sync"
)

type MemoryCache struct {
	Options CacheOptions
	*sync.Mutex

	users map[uint64]*user.User
	usersLock *sync.RWMutex

	guilds map[uint64]guild.Guild
	guildsLock *sync.RWMutex

	channels map[uint64]*channel.Channel
	channelsLock *sync.RWMutex

	roles map[uint64]*guild.Role
	rolesLock *sync.RWMutex

	emojis map[uint64]*emoji.Emoji
	emojisLock *sync.RWMutex

	self user.User
	selfLock *sync.RWMutex
}

func (c *MemoryCache) GetOptions() CacheOptions {
	return c.Options
}

func (c *MemoryCache) StoreUser(user user.User) {
	if c.Options.Users {
		c.usersLock.Lock()
		c.users[user.Id] = &user
		c.usersLock.Unlock()
	}
}

func (c *MemoryCache) GetUser(id uint64) *user.User {
	c.usersLock.RLock()
	user := c.users[id]
	c.usersLock.RUnlock()
	return user
}

func (c *MemoryCache) StoreGuild(g guild.Guild) {
	if c.Options.Guilds {
		if !c.Options.Members {
			g.Members = make([]member.Member, 0)
		}

		if !c.Options.Guilds {
			g.Channels = make([]channel.Channel, 0)
		}

		if !c.Options.Roles {
			g.Roles = make([]guild.Role, 0)
		}

		if !c.Options.Emojis {
			g.Emojis = make([]emoji.Emoji, 0)
		}

		if !c.Options.VoiceStates {
			g.VoiceStates = make([]guild.VoiceState, 0)
		}

		c.Lock()
		c.guildsLock.Lock()
		c.guilds[g.Id] = g

		if c.Options.Users {
			c.usersLock.Lock()
			for i, member := range g.Members {
				c.users[member.User.Id] = &c.guilds[g.Id].Members[i].User
			}
			c.usersLock.Unlock()
		}

		if c.Options.Channels {
			c.channelsLock.Lock()
			for i, channel := range g.Channels {
				c.channels[channel.Id] = &c.guilds[g.Id].Channels[i]
			}
			c.channelsLock.Unlock()
		}

		if c.Options.Roles {
			c.rolesLock.Lock()
			for i, role := range g.Roles {
				c.roles[role.Id] = &c.guilds[g.Id].Roles[i]
			}
			c.rolesLock.Unlock()
		}

		if c.Options.Emojis {
			c.emojisLock.Lock()
			for i, emoji := range g.Emojis {
				c.emojis[emoji.Id] = &c.guilds[g.Id].Emojis[i]
			}
			c.emojisLock.Unlock()
		}

		c.guildsLock.Unlock()
		c.Unlock()
	}
}

func (c *MemoryCache) GetGuild(id uint64) *guild.Guild {
	c.guildsLock.RLock()
	guild, ok := c.guilds[id]
	c.guildsLock.RUnlock()

	if ok {
		return &guild
	} else {
		return nil
	}
}

func (c *MemoryCache) GetGuilds() []guild.Guild {
	c.guildsLock.RLock()
	guilds := make([]guild.Guild, 0)
	for _, guild := range c.guilds {
		guilds = append(guilds, guild)
	}
	c.guildsLock.RUnlock()
	return guilds
}

func (c *MemoryCache) DeleteGuild(id uint64) {
	c.Lock()
	c.guildsLock.Lock()

	// clean up
	g := c.guilds[id]

	if c.Options.Channels {
		c.channelsLock.Lock()
		for _, channel := range g.Channels {
			delete(c.channels, channel.Id)
		}
		c.channelsLock.Unlock()
	}

	if c.Options.Roles {
		c.rolesLock.Lock()
		for _, role := range g.Roles {
			delete(c.roles, role.Id)
		}
		c.rolesLock.Unlock()
	}

	if c.Options.Emojis {
		c.emojisLock.Lock()
		for _, emoji := range g.Emojis {
			delete(c.emojis, emoji.Id)
		}
		c.emojisLock.Unlock()
	}

	delete(c.guilds, id)

	c.guildsLock.Unlock()
	c.Unlock()
}

func (c *MemoryCache) StoreMember(m member.Member, guildId uint64) {
	if c.Options.Members {
		c.Lock()
		c.guildsLock.Lock()
		c.usersLock.Lock()
		c.Unlock()

		g, found := c.guilds[guildId]
		if found {
			// check if we already have the user cached
			index := -1
			var cached member.Member
			for i, guildMember := range g.Members {
				if guildMember.User.Id == m.User.Id {
					index = i
					cached = guildMember
					break
				}
			}

			if index != -1 {
				// handle guild member update event - only a few fields are sent
				cached.Roles = m.Roles
				cached.User = m.User
				cached.Nick = m.Nick
				cached.PremiumSince = m.PremiumSince

				g.Members[index] = cached
			} else {
				g.Members = append(g.Members, m)
				index = len(g.Members) - 1
			}

			if c.Options.Users {
				c.users[m.User.Id] = &g.Members[index].User
			}

			c.guilds[guildId] = g
		}

		c.usersLock.Unlock()
		c.guildsLock.Unlock()
	}
}

func (c *MemoryCache) DeleteMember(userId, guildId uint64) {
	if c.Options.Members {
		c.guildsLock.Lock()

		g, found := c.guilds[guildId]
		if found {
			// check if we already have the user cached
			index := -1
			for i, guildMember := range g.Members {
				if guildMember.User.Id == userId {
					index = i
					break
				}
			}

			if index != -1 {
				g.Members = append(g.Members[:index], g.Members[index+1:]...)
				c.guilds[guildId] = g
			}
		}

		c.guildsLock.Unlock()
	}
}

func (c *MemoryCache) StoreChannel(ch channel.Channel) {
	c.Lock()
	c.guildsLock.Lock()
	c.channelsLock.Lock()
	c.Unlock()

	if c.Options.Channels {
		if ch.Type != channel.ChannelTypeDM && ch.Type != channel.ChannelTypeGroupDM { // store channel on guild
			g, found := c.guilds[ch.GuildId]
			index := -1
			if found {
				// check if the channel already exists
				for i, guildChannel := range g.Channels {
					if guildChannel.Id == ch.Id {
						index = i
						break
					}
				}

				if index != -1 {
					g.Channels[index] = ch
				} else {
					g.Channels = append(g.Channels, ch)
					index = len(g.Channels) - 1
				}

				c.guilds[ch.GuildId] = g

				// store pointer to channel
				c.channels[ch.Id] = &c.guilds[ch.GuildId].Channels[index]
			}
		} else { // store alone
			c.channels[ch.Id] = &ch
		}
	}
	c.guildsLock.Unlock()
	c.channelsLock.Unlock()
}

func (c *MemoryCache) GetChannel(id uint64) *channel.Channel {
	c.channelsLock.RLock()
	channel := c.channels[id]
	c.channelsLock.RUnlock()
	return channel
}

func (c *MemoryCache) DeleteChannel(channelId, guildId uint64) {
	if c.Options.Channels {
		c.Lock()
		c.guildsLock.Lock()
		c.channelsLock.Lock()
		c.Unlock()

		delete(c.channels, channelId)

		// update guild object
		g, found := c.guilds[guildId]
		if found {
			index := -1
			for i, guildChannel := range g.Channels {
				if guildChannel.Id == channelId {
					index = i
					break
				}
			}

			if index != -1 {
				g.Channels = append(g.Channels[:index], g.Channels[index+1:]...)
				c.guilds[guildId] = g
			}
		}

		c.guildsLock.Unlock()
		c.channelsLock.Unlock()
	}
}

func (c *MemoryCache) StoreRole(role guild.Role, guildId uint64) {
	if c.Options.Roles {
		c.Lock()
		c.rolesLock.Lock()
		c.guildsLock.Lock()
		c.Unlock()

		// store on guild object
		g, found := c.guilds[guildId]
		index := -1
		if found {
			// check if the role already exists
			for i, guildRole := range g.Roles {
				if guildRole.Id == role.Id {
					index = i
					break
				}
			}

			if index != -1 {
				g.Roles[index] = role
			} else {
				g.Roles = append(g.Roles, role)
				index = len(g.Roles) - 1
			}

			c.guilds[guildId] = g
		}

		// store pointer to role
		c.roles[role.Id] = &c.guilds[guildId].Roles[index]

		c.guildsLock.Unlock()
		c.rolesLock.Unlock()
	}
}

func (c *MemoryCache) GetRole(id uint64) *guild.Role {
	c.rolesLock.RLock()
	role := c.roles[id]
	c.rolesLock.RUnlock()
	return role
}

func (c *MemoryCache) DeleteRole(roleId, guildId uint64) {
	if c.Options.Roles {
		c.Lock()
		c.guildsLock.Lock()
		c.rolesLock.Lock()
		c.Unlock()

		// update guild object
		g, found := c.guilds[guildId]
		if found {
			index := -1
			for i, guildRole := range g.Roles {
				if guildRole.Id == roleId {
					index = i
					break
				}
			}

			if index != -1 {
				g.Roles = append(g.Roles[:index], g.Roles[index+1:]...)
			}

			c.guilds[guildId] = g
		}

		delete(c.roles, roleId)

		c.guildsLock.Unlock()
		c.rolesLock.Unlock()
	}
}

func (c *MemoryCache) StoreEmoji(emoji emoji.Emoji, guildId uint64) {
	if c.Options.Emojis {
		c.Lock()
		c.emojisLock.Lock()
		c.guildsLock.Lock()
		c.Unlock()

		// store on guild object
		g, found := c.guilds[guildId]
		index := -1
		if found {
			// check if the emoji already exists
			for i, guildEmoji := range g.Emojis {
				if guildEmoji.Id == emoji.Id {
					index = i
					break
				}
			}

			if index != -1 {
				g.Emojis[index] = emoji
			} else {
				g.Emojis = append(g.Emojis, emoji)
				index = len(g.Emojis) - 1
			}

			c.guilds[guildId] = g
		}

		// store pointer to role
		c.emojis[emoji.Id] = &c.guilds[guildId].Emojis[index]

		c.guildsLock.Unlock()
		c.emojisLock.Unlock()
	}
}

func (c *MemoryCache) GetEmoji(id uint64) *emoji.Emoji {
	c.emojisLock.RLock()
	emoji := c.emojis[id]
	c.emojisLock.RUnlock()
	return emoji
}

func (c *MemoryCache) DeleteEmoji(emojiId, guildId uint64) {
	if c.Options.Emojis {
		c.Lock()
		c.guildsLock.Lock()
		c.emojisLock.Lock()
		c.Unlock()

		// update guild object
		g, found := c.guilds[guildId]
		if found {
			index := -1
			for i, guildEmoji := range g.Emojis {
				if guildEmoji.Id == emojiId {
					index = i
					break
				}
			}

			if index != -1 {
				g.Emojis = append(g.Emojis[:index], g.Emojis[index+1:]...)
				c.guilds[guildId] = g
			}
		}

		delete(c.emojis, emojiId)

		c.guildsLock.Unlock()
		c.emojisLock.Unlock()
	}
}

func (c *MemoryCache) StoreVoiceState(state guild.VoiceState) {
	c.guildsLock.Lock()

	g, found := c.guilds[state.GuildId]
	if found {
		index := -1
		for i, guildVoiceState := range g.VoiceStates {
			if guildVoiceState.UserId == state.UserId {
				index = i
				break
			}
		}

		if index != -1 {
			g.VoiceStates[index] = state
		} else {
			g.VoiceStates = append(g.VoiceStates, state)
		}

		c.guilds[state.GuildId] = g
	}

	c.guildsLock.Unlock()
}

func (c *MemoryCache) GetVoiceState(userId, guildId uint64) *guild.VoiceState {
	c.guildsLock.RLock()

	var state *guild.VoiceState
	for _, guildState := range c.guilds[guildId].VoiceStates {
		if guildState.UserId == userId {
			state = &guildState
			break
		}
	}

	c.guildsLock.RUnlock()
	return state
}

func (c *MemoryCache) StoreSelf(self user.User) {
	c.selfLock.Lock()
	c.self = self
	c.selfLock.Unlock()
}

func (c *MemoryCache) GetSelf() *user.User {
	c.selfLock.RLock()
	self := c.self
	c.selfLock.RUnlock()
	return &self
}
