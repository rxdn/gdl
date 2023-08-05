package cache

import (
	"context"
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

func (c *MemoryCache) Options() CacheOptions {
	return c.options
}

func (c *MemoryCache) StoreUser(ctx context.Context, u user.User) error {
	return c.StoreUsers(ctx, []user.User{u})
}

func (c *MemoryCache) StoreUsers(ctx context.Context, users []user.User) error {
	if c.options.Users {
		c.userLock.Lock()
		defer c.userLock.Unlock()

		for _, user := range users {
			c.users[user.Id] = user.ToCachedUser()
		}
	}

	return nil
}

func (c *MemoryCache) GetUser(ctx context.Context, userId uint64) (user.User, error) {
	c.userLock.RLock()
	defer c.userLock.RUnlock()

	cached, found := c.users[userId]
	if found {
		return cached.ToUser(userId), nil
	} else {
		return user.User{}, ErrNotFound
	}
}

func (c *MemoryCache) GetUsers(ctx context.Context, ids []uint64) (map[uint64]user.User, error) {
	c.userLock.RLock()
	defer c.userLock.RUnlock()

	users := make(map[uint64]user.User)
	for _, id := range ids {
		cached, found := c.users[id]
		if found {
			users[id] = cached.ToUser(id)
		}
	}

	return users, nil
}

func (c *MemoryCache) StoreGuild(ctx context.Context, g guild.Guild) error {
	return c.StoreGuilds(ctx, []guild.Guild{g})
}

func (c *MemoryCache) StoreGuilds(ctx context.Context, guilds []guild.Guild) error {
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
		if err := c.StoreChannels(ctx, guild.Channels); err != nil {
			return err
		}

		if err := c.StoreMembers(ctx, guild.Members, guild.Id); err != nil {
			return err
		}

		if err := c.StoreRoles(ctx, guild.Roles, guild.Id); err != nil {
			return err
		}

		if err := c.StoreEmojis(ctx, guild.Emojis, guild.Id); err != nil {
			return err
		}

		if err := c.StoreVoiceStates(ctx, guild.VoiceStates); err != nil {
			return err
		}
	}

	return nil
}

func (c *MemoryCache) GetGuild(ctx context.Context, guildId uint64) (guild.Guild, error) {
	c.guildLock.RLock()
	cached, found := c.guilds[guildId]
	c.guildLock.RUnlock()

	if found {
		return cached.ToGuild(guildId), nil
	} else {
		return guild.Guild{}, ErrNotFound
	}
}

func (c *MemoryCache) DeleteGuild(ctx context.Context, guildId uint64) error {
	c.guildLock.Lock()
	delete(c.guilds, guildId)
	c.guildLock.Unlock()

	return nil
}

func (c *MemoryCache) GetGuildCount(ctx context.Context) (int, error) {
	c.guildLock.RLock()
	count := len(c.guilds)
	c.guildLock.RUnlock()

	return count, nil
}

func (c *MemoryCache) GetGuildOwner(ctx context.Context, guildId uint64) (uint64, error) {
	c.guildLock.RLock()
	guild, ok := c.guilds[guildId]
	c.guildLock.RUnlock()

	if !ok {
		return 0, ErrNotFound
	}

	return guild.OwnerId, nil
}

func (c *MemoryCache) StoreMember(ctx context.Context, m member.Member, guildId uint64) error {
	return c.StoreMembers(ctx, []member.Member{m}, guildId)
}

func (c *MemoryCache) StoreMembers(ctx context.Context, members []member.Member, guildId uint64) error {
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

	return nil
}

func (c *MemoryCache) GetMember(ctx context.Context, guildId, userId uint64) (member.Member, error) {
	c.memberLock.RLock()
	defer c.memberLock.RUnlock()

	if c.members[guildId] == nil {
		return member.Member{}, ErrNotFound
	}

	cached, found := c.members[guildId][userId]
	if found {
		u, err := c.GetUser(ctx, userId)
		if err == ErrNotFound {
			u = user.User{Id: userId}
		} else if err != nil {
			return member.Member{}, err
		}

		return cached.ToMember(u), nil
	} else {
		return member.Member{}, ErrNotFound
	}
}

func (c *MemoryCache) GetGuildMembers(ctx context.Context, guildId uint64, withUserData bool) ([]member.Member, error) {
	c.memberLock.RLock()
	defer c.memberLock.RUnlock()

	if c.members[guildId] == nil {
		return []member.Member{}, ErrNotFound
	}

	var members []member.Member
	for userId, cachedMember := range c.members[guildId] {
		// Get user
		var u user.User
		if withUserData {
			cachedUser, err := c.GetUser(ctx, userId)
			if err == ErrNotFound {
				u = user.User{Id: userId}
			} else if err != nil {
				return nil, err
			} else {
				u = cachedUser
			}
		} else {
			u = user.User{Id: userId}
		}

		members = append(members, cachedMember.ToMember(u))
	}

	return members, nil
}

func (c *MemoryCache) DeleteMember(ctx context.Context, userId, guildId uint64) error {
	c.memberLock.Lock()
	defer c.memberLock.Unlock()

	if c.members[guildId] != nil {
		delete(c.members[guildId], userId)
	}

	return nil
}

func (c *MemoryCache) StoreChannel(ctx context.Context, ch channel.Channel) error {
	return c.StoreChannels(ctx, []channel.Channel{ch})
}

func (c *MemoryCache) StoreChannels(ctx context.Context, channels []channel.Channel) error {
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

	return nil
}

func (c *MemoryCache) GetChannel(ctx context.Context, channelId uint64) (channel.Channel, error) {
	var cached channel.CachedChannel

	c.channelLock.RLock()
	cached, found := c.channels[channelId]
	c.channelLock.RUnlock()

	if found {
		return cached.ToChannel(channelId, cached.GuildId), nil
	} else {
		return channel.Channel{}, ErrNotFound
	}
}

func (c *MemoryCache) GetGuildChannels(ctx context.Context, guildId uint64) ([]channel.Channel, error) {
	c.guildLock.RLock()
	guild, found := c.guilds[guildId]
	c.guildLock.RUnlock()
	if !found {
		return nil, ErrNotFound
	}

	var channels []channel.Channel

	c.channelLock.RLock()
	for _, channelId := range guild.Channels {
		if cachedChannel, found := c.channels[channelId]; found {
			channels = append(channels, cachedChannel.ToChannel(channelId, guildId))
		}
	}
	c.channelLock.RUnlock()

	return channels, nil
}

func (c *MemoryCache) DeleteChannel(ctx context.Context, channelId uint64) error {
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

	return nil
}

func (c *MemoryCache) DeleteGuildChannels(ctx context.Context, guildId uint64) error {
	c.guildLock.Lock()
	defer c.guildLock.Unlock()

	guild, ok := c.guilds[guildId]
	if !ok {
		return ErrNotFound
	}

	c.channelLock.Lock()
	defer c.channelLock.Unlock()

	for channelId := range c.channels {
		delete(c.channels, channelId)
	}

	guild.Channels = nil
	c.guilds[guildId] = guild

	return nil
}

func (c *MemoryCache) StoreRole(ctx context.Context, role guild.Role, guildId uint64) error {
	return c.StoreRoles(ctx, []guild.Role{role}, guildId)
}

func (c *MemoryCache) StoreRoles(ctx context.Context, roles []guild.Role, guildId uint64) error {
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

	return nil
}

func (c *MemoryCache) GetRole(ctx context.Context, roleId uint64) (guild.Role, error) {
	c.roleLock.RLock()
	defer c.roleLock.RUnlock()

	cachedRole, found := c.roles[roleId]
	if found {
		return cachedRole.ToRole(roleId), nil
	} else {
		return guild.Role{}, ErrNotFound
	}
}

func (c *MemoryCache) GetRoles(ctx context.Context, guildId uint64, ids []uint64) (map[uint64]guild.Role, error) {
	c.roleLock.RLock()
	defer c.roleLock.RUnlock()

	roles := make(map[uint64]guild.Role)
	for _, id := range ids {
		if cachedRole, found := c.roles[id]; found {
			if cachedRole.GuildId == guildId {
				roles[id] = cachedRole.ToRole(id)
			}
		}
	}

	return roles, nil
}

func (c *MemoryCache) GetGuildRoles(ctx context.Context, guildId uint64) ([]guild.Role, error) {
	// get guild
	c.guildLock.RLock()
	g, found := c.guilds[guildId]
	c.guildLock.RUnlock()

	if !found {
		return nil, ErrNotFound
	}

	c.roleLock.RLock()
	defer c.roleLock.RUnlock()

	var roles []guild.Role
	for _, roleId := range g.Roles {
		if role, found := c.roles[roleId]; found {
			roles = append(roles, role.ToRole(roleId))
		}
	}

	return roles, nil
}

func (c *MemoryCache) DeleteRole(ctx context.Context, roleId uint64) error {
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

	return nil
}

func (c *MemoryCache) DeleteGuildRoles(ctx context.Context, guildId uint64) error {
	c.guildLock.Lock()
	defer c.guildLock.Unlock()

	guild, ok := c.guilds[guildId]
	if !ok {
		return ErrNotFound
	}

	c.roleLock.Lock()
	defer c.roleLock.Unlock()

	for roleId := range c.roles {
		delete(c.roles, roleId)
	}

	guild.Roles = nil
	c.guilds[guildId] = guild

	return nil
}

func (c *MemoryCache) StoreEmoji(ctx context.Context, e emoji.Emoji, guildId uint64) error {
	return c.StoreEmojis(ctx, []emoji.Emoji{e}, guildId)
}

func (c *MemoryCache) StoreEmojis(ctx context.Context, emojis []emoji.Emoji, guildId uint64) error {
	c.emojiLock.Lock()

	for _, emoji := range emojis {
		c.emojis[emoji.Id.Value] = emoji.ToCachedEmoji(guildId)

		// Add to guild object
		c.guildLock.Lock()
		if guild, found := c.guilds[guildId]; found {
			// Check to see if emoji already exists
			var emojiExists bool
			for _, emojiId := range guild.Emojis {
				if emojiId == emoji.Id.Value {
					emojiExists = true
					break
				}
			}

			if !emojiExists {
				guild.Emojis = append(guild.Emojis, emoji.Id.Value)
				c.guilds[guildId] = guild
			}
		}
		c.guildLock.Unlock()
	}

	c.emojiLock.Unlock()
	return nil
}

func (c *MemoryCache) GetEmoji(ctx context.Context, emojiId uint64) (emoji.Emoji, error) {
	c.emojiLock.RLock()
	cached, found := c.emojis[emojiId]
	c.emojiLock.RUnlock()

	if !found {
		return emoji.Emoji{}, ErrNotFound
	}

	u, err := c.GetUser(ctx, cached.User)
	if err == ErrNotFound {
		u = user.User{Id: cached.User}
	} else if err != nil {
		return emoji.Emoji{}, err
	}

	return cached.ToEmoji(emojiId, u), nil
}

func (c *MemoryCache) GetGuildEmojis(ctx context.Context, guildId uint64) ([]emoji.Emoji, error) {
	// get guild
	c.guildLock.RLock()
	guild, found := c.guilds[guildId]
	c.guildLock.RUnlock()

	if !found {
		return nil, ErrNotFound
	}

	c.emojiLock.RLock()
	defer c.emojiLock.RUnlock()

	var emojis []emoji.Emoji
	for _, emojiId := range guild.Emojis {
		cached, found := c.emojis[emojiId]
		if !found {
			continue
		}

		u, err := c.GetUser(ctx, cached.User)
		if err == ErrNotFound {
			u = user.User{Id: cached.User}
		} else if err != nil {
			return nil, err
		}

		emojis = append(emojis, cached.ToEmoji(emojiId, u))
	}

	return emojis, nil
}

func (c *MemoryCache) DeleteEmoji(ctx context.Context, emojiId uint64) error {
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

	return nil
}

func (c *MemoryCache) StoreVoiceState(ctx context.Context, state guild.VoiceState) error {
	return c.StoreVoiceStates(ctx, []guild.VoiceState{state})
}

func (c *MemoryCache) StoreVoiceStates(ctx context.Context, states []guild.VoiceState) error {
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

	return nil
}

func (c *MemoryCache) GetVoiceState(ctx context.Context, userId, guildId uint64) (guild.VoiceState, error) {
	var cached guild.CachedVoiceState

	c.voiceStateLock.RLock()
	defer c.voiceStateLock.RUnlock()
	if c.voiceStates[guildId] == nil {
		return guild.VoiceState{}, ErrNotFound
	}

	cached, found := c.voiceStates[guildId][userId]
	if found {
		// get member
		m, err := c.GetMember(ctx, guildId, userId)
		if err == ErrNotFound {
			m = member.Member{
				User: user.User{
					Id: userId,
				},
			}
		} else if err != nil {
			return guild.VoiceState{}, err
		}

		return cached.ToVoiceState(guildId, m), nil
	} else {
		return guild.VoiceState{}, ErrNotFound
	}
}

func (c *MemoryCache) GetGuildVoiceStates(ctx context.Context, guildId uint64) ([]guild.VoiceState, error) {
	c.voiceStateLock.RLock()
	defer c.voiceStateLock.RUnlock()

	if c.voiceStates[guildId] == nil {
		return nil, ErrNotFound
	}

	var states []guild.VoiceState
	for userId, cached := range c.voiceStates[guildId] {
		// get member
		m, err := c.GetMember(ctx, guildId, userId)
		if err == ErrNotFound {
			m = member.Member{
				User: user.User{
					Id: userId,
				},
			}
		} else if err != nil {
			return nil, err
		}

		states = append(states, cached.ToVoiceState(guildId, m))
	}

	return states, nil
}

func (c *MemoryCache) DeleteVoiceState(ctx context.Context, userId, guildId uint64) error {
	c.voiceStateLock.Lock()
	defer c.voiceStateLock.Unlock()

	if c.voiceStates[guildId] == nil {
		return ErrNotFound
	}

	delete(c.voiceStates[guildId], userId)
	return nil
}

func (c *MemoryCache) StoreSelf(ctx context.Context, self user.User) error {
	c.selfLock.Lock()
	c.self = self
	c.selfLock.Unlock()
	return nil
}

func (c *MemoryCache) GetSelf(ctx context.Context) (user.User, error) {
	c.selfLock.RLock()
	self := c.self
	c.selfLock.RUnlock()

	if self.Id == 0 {
		return user.User{}, ErrNotFound
	} else {
		return self, nil
	}
}
