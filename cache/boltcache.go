package cache

import (
	"encoding/json"
	"fmt"
	"github.com/boltdb/bolt"
	"github.com/rxdn/gdl/objects/channel"
	"github.com/rxdn/gdl/objects/guild"
	"github.com/rxdn/gdl/objects/guild/emoji"
	"github.com/rxdn/gdl/objects/member"
	"github.com/rxdn/gdl/objects/user"
	"os"
	"strconv"
	"strings"
	"sync"
)

type BoltCache struct {
	*bolt.DB
	options CacheOptions

	// TODO: Should we store self in the DB? Seems kinda redundant
	selfLock sync.RWMutex
	self     user.User
}

type BoltOptions struct {
	ClearOnRestart bool
	Path string
	FileMode os.FileMode
	*bolt.Options
}

func NewBoltCache(cacheOptions CacheOptions, boltOptions BoltOptions) BoltCache {
	if boltOptions.ClearOnRestart {
		_ = os.Remove(boltOptions.Path)
	}

	db, err := bolt.Open(boltOptions.Path, boltOptions.FileMode, boltOptions.Options)
	if err != nil {
		panic(err)
	}

	if err := createBuckets(db); err != nil {
		panic(err)
	}

	return BoltCache{
		DB:      db,
		options: cacheOptions,
	}
}

func createBuckets(db *bolt.DB) error {
	return db.Update(func(tx *bolt.Tx) error {
		boltMustCreate(tx, "users")
		boltMustCreate(tx, "guilds")
		boltMustCreate(tx, "members")
		boltMustCreate(tx, "channels")
		boltMustCreate(tx, "roles")
		boltMustCreate(tx, "emojis")
		boltMustCreate(tx, "voice_states")

		return nil
	})
}

func (c *BoltCache) GetOptions() CacheOptions {
	return c.options
}

func (c *BoltCache) StoreUser(u user.User) {
	c.StoreUsers([]user.User{u})
}

func (c *BoltCache) StoreUsers(users []user.User) {
	if c.options.Users {
		_ = c.Batch(func(tx *bolt.Tx) error {
			b := tx.Bucket([]byte("users"))

			for _, u := range users {
				if encoded, err := json.Marshal(u.ToCachedUser()); err == nil {
					if err := b.Put(toBytes(u.Id), encoded); err != nil {
						return err
					}
				} else {
					return err
				}
			}

			return nil
		})
	}
}

func (c *BoltCache) GetUser(userId uint64) (user.User, bool) {
	var u user.CachedUser
	var found bool

	_ = c.View(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte("users"))
		encoded := b.Get(toBytes(userId))

		if encoded == nil {
			return nil
		}

		return json.Unmarshal(encoded, &u)
	})

	return u.ToUser(userId), found
}

func (c *BoltCache) StoreGuild(g guild.Guild) {
	c.StoreGuilds([]guild.Guild{g})
}

func (c *BoltCache) StoreGuilds(guilds []guild.Guild) {
	if c.options.Guilds {
		_ = c.Batch(func(tx *bolt.Tx) error {
			b := tx.Bucket([]byte("guilds"))

			for _, g := range guilds {
				if encoded, err := json.Marshal(g.ToCachedGuild()); err == nil {
					if err := b.Put(toBytes(g.Id), encoded); err != nil {
						return err
					}
				} else {
					return err
				}
			}

			return nil
		})
	}

	for _, guild := range guilds {
		c.StoreChannels(guild.Channels)
		c.StoreMembers(guild.Members, guild.Id)
		c.StoreRoles(guild.Roles, guild.Id)
		c.StoreEmojis(guild.Emojis, guild.Id)
		c.StoreVoiceStates(guild.VoiceStates)
	}
}

func (c *BoltCache) GetGuild(guildId uint64, withMembers bool) (guild.Guild, bool) {
	var cached guild.CachedGuild
	var found bool

	_ = c.View(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte("guilds"))
		encoded := b.Get(toBytes(guildId))

		if encoded == nil {
			return nil
		}

		return json.Unmarshal(encoded, &cached)
	})

	g := cached.ToGuild(guildId)
	g.Channels = c.GetGuildChannels(guildId)
	g.Roles = c.GetGuildRoles(guildId)

	if withMembers {
		g.Members = c.GetGuildMembers(guildId, false)
	}

	g.Emojis = c.GetGuildEmojis(guildId)
	g.VoiceStates = c.GetGuildVoiceStates(guildId)

	return g, found
}

func (c *BoltCache) GetGuilds() []guild.Guild {
	var guilds []guild.Guild

	_ = c.View(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte("guilds"))

		return b.ForEach(func(k, encoded []byte) error {
			guildId, err := strconv.ParseUint(string(k), 10, 64); if err != nil {
				return nil
			}

			var cached guild.CachedGuild
			if err := json.Unmarshal(encoded, &cached); err == nil {
				guilds = append(guilds, cached.ToGuild(guildId))
			}

			return nil
		})
	})

	return guilds
}

func (c *BoltCache) DeleteGuild(guildId uint64) {
	_ = c.Update(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte("guilds"))
		 _ = b.Delete(toBytes(guildId))

		return nil
	})
}

func (c *BoltCache) GetGuildCount() int {
	var count int

	_ = c.View(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte("guilds"))
		cursor := b.Cursor()

		for k, _ := cursor.First(); k != nil; k, _ = cursor.Next() {
			count++
		}

		return nil
	})

	return count
}

func (c *BoltCache) StoreMember(m member.Member, guildId uint64) {
	c.StoreMembers([]member.Member{m}, guildId)
}

func (c *BoltCache) StoreMembers(members []member.Member, guildId uint64) {
	if c.options.Members {
		_ = c.Batch(func(tx *bolt.Tx) error {
			b := tx.Bucket([]byte("members"))

			for _, m := range members {
				if encoded, err := json.Marshal(m.ToCachedMember()); err == nil {
					if err := b.Put(memberToBytes(m.User.Id, guildId), encoded); err != nil {
						return err
					}
				} else {
					return err
				}
			}

			return nil
		})
	}
}

func (c *BoltCache) GetMember(guildId, userId uint64) (member.Member, bool) {
	var cached member.CachedMember
	var found bool

	_ = c.View(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte("members"))
		encoded := b.Get(memberToBytes(userId, guildId))

		if encoded == nil {
			return nil
		}

		return json.Unmarshal(encoded, &cached)
	})

	u, userFound := c.GetUser(userId)
	if !userFound {
		u = user.User{Id:userId}
	}

	return cached.ToMember(u), found
}


func (c *BoltCache) GetGuildMembers(guildId uint64, withUserData bool) []member.Member {
	var members []member.Member

	_ = c.View(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte("members"))

		return b.ForEach(func(k, encoded []byte) error {
			split := strings.Split(string(k), ":")
			if len(split) < 2 {
				return nil
			}

			// Hacky but w/e
			cachedUserId, err := strconv.ParseUint(split[0], 10, 64); if err != nil {
				return nil
			}

			cachedGuildId, err := strconv.ParseUint(split[1], 10, 64); if err != nil {
				return nil
			}

			var cached member.CachedMember
			if err := json.Unmarshal(encoded, &cached); err == nil && cachedGuildId == guildId {
				u := user.User{Id: cachedUserId}

				if withUserData {
					var found bool
					u, found = c.GetUser(cachedUserId)
					if !found {
						u = user.User{Id: cachedUserId}
					}
				}

				members = append(members, cached.ToMember(u))
			}

			return nil
		})
	})

	return members
}

func (c *BoltCache) DeleteMember(userId, guildId uint64) {
	_ = c.Update(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte("members"))
		_ = b.Delete(memberToBytes(userId, guildId))

		return nil
	})
}

type channelWithGuild struct {
	channel.CachedChannel
	guildId uint64
}

func (c *BoltCache) StoreChannel(ch channel.Channel) {
	c.StoreChannels([]channel.Channel{ch})
}

func (c *BoltCache) StoreChannels(channels []channel.Channel) {
	if c.options.Guilds {
		_ = c.Batch(func(tx *bolt.Tx) error {
			b := tx.Bucket([]byte("channels"))

			for _, ch := range channels {
				cwg := channelWithGuild{
					CachedChannel: ch.ToCachedChannel(),
					guildId:       ch.GuildId,
				}

				if encoded, err := json.Marshal(cwg); err == nil {
					if err := b.Put(toBytes(ch.Id), encoded); err != nil {
						return err
					}
				} else {
					return err
				}
			}

			return nil
		})
	}
}

func (c *BoltCache) GetChannel(channelId uint64) (channel.Channel, bool) {
	var cached channelWithGuild
	var found bool

	_ = c.View(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte("channels"))
		encoded := b.Get(toBytes(channelId))

		if encoded == nil {
			return nil
		}

		return json.Unmarshal(encoded, &cached)
	})

	ch := cached.ToChannel(channelId, cached.guildId)
	return ch, found
}

func (c *BoltCache) GetGuildChannels(guildId uint64) []channel.Channel {
	var channels []channel.Channel

	_ = c.View(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte("channels"))

		return b.ForEach(func(k, encoded []byte) error {
			channelId, err := strconv.ParseUint(string(k), 10, 64); if err != nil {
				return nil
			}

			var cached channelWithGuild
			if err := json.Unmarshal(encoded, &cached); err == nil && cached.guildId == guildId {
				channels = append(channels, cached.ToChannel(channelId, cached.guildId))
			}

			return nil
		})
	})

	return channels
}

func (c *BoltCache) DeleteChannel(channelId uint64) {
	_ = c.Update(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte("channels"))
		_ = b.Delete(toBytes(channelId))

		return nil
	})
}

type roleWithGuild struct {
	guild.CachedRole
	guildId uint64
}

func (c *BoltCache) StoreRole(role guild.Role, guildId uint64) {
	c.StoreRoles([]guild.Role{role}, guildId)
}

func (c *BoltCache) StoreRoles(roles []guild.Role, guildId uint64) {
	if c.options.Roles {
		_ = c.Batch(func(tx *bolt.Tx) error {
			b := tx.Bucket([]byte("roles"))

			for _, role := range roles {
				rwg := roleWithGuild{
					CachedRole: role.ToCachedRole(guildId),
					guildId:    guildId,
				}

				if encoded, err := json.Marshal(rwg); err == nil {
					if err := b.Put(toBytes(role.Id), encoded); err != nil {
						return err
					}
				} else {
					return err
				}
			}

			return nil
		})
	}
}

func (c *BoltCache) GetRole(roleId uint64) (guild.Role, bool) {
	var cached roleWithGuild
	var found bool

	_ = c.View(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte("roles"))
		encoded := b.Get(toBytes(roleId))

		if encoded == nil {
			return nil
		}

		return json.Unmarshal(encoded, &cached)
	})

	ch := cached.ToRole(roleId)
	return ch, found
}

func (c *BoltCache) GetGuildRoles(guildId uint64) []guild.Role {
	var roles []guild.Role

	_ = c.View(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte("roles"))

		return b.ForEach(func(k, encoded []byte) error {
			roleId, err := strconv.ParseUint(string(k), 10, 64); if err != nil {
				return nil
			}

			var cached roleWithGuild
			if err := json.Unmarshal(encoded, &cached); err == nil && cached.guildId == guildId {
				roles = append(roles, cached.ToRole(roleId))
			}

			return nil
		})
	})

	return roles
}

func (c *BoltCache) DeleteRole(roleId uint64) {
	_ = c.Update(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte("roles"))
		_ = b.Delete(toBytes(roleId))

		return nil
	})
}

type emojiWithGuild struct {
	emoji.CachedEmoji
	guildId uint64
}

func (c *BoltCache) StoreEmoji(e emoji.Emoji, guildId uint64) {
	c.StoreEmojis([]emoji.Emoji{e}, guildId)
}

func (c *BoltCache) StoreEmojis(emojis []emoji.Emoji, guildId uint64) {
	_ = c.Batch(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte("emojis"))

		for _, emoji := range emojis {
			ewg := emojiWithGuild{
				CachedEmoji: emoji.ToCachedEmoji(guildId),
				guildId:    guildId,
			}

			if encoded, err := json.Marshal(ewg); err == nil {
				if err := b.Put(toBytes(emoji.Id), encoded); err != nil {
					return err
				}
			} else {
				return err
			}
		}

		return nil
	})
}

func (c *BoltCache) GetEmoji(emojiId uint64) (emoji.Emoji, bool) {
	var cached emojiWithGuild
	var found bool

	_ = c.View(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte("emojis"))
		encoded := b.Get(toBytes(emojiId))

		if encoded == nil {
			return nil
		}

		return json.Unmarshal(encoded, &cached)
	})

	u, userFound := c.GetUser(cached.User)
	if !userFound {
		u = user.User{Id: cached.User}
	}

	emoji := cached.ToEmoji(emojiId, u)
	return emoji, found
}

func (c *BoltCache) GetGuildEmojis(guildId uint64) []emoji.Emoji {
	var emojis []emoji.Emoji

	_ = c.View(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte("emojis"))

		return b.ForEach(func(k, encoded []byte) error {
			emojiId, err := strconv.ParseUint(string(k), 10, 64); if err != nil {
				return nil
			}

			var cached emojiWithGuild
			if err := json.Unmarshal(encoded, &cached); err == nil && cached.guildId == guildId {
				u, found := c.GetUser(cached.User)
				if !found {
					u = user.User{Id: cached.User}
				}

				emojis = append(emojis, cached.ToEmoji(emojiId, u))
			}

			return nil
		})
	})

	return emojis
}

func (c *BoltCache) DeleteEmoji(emojiId uint64) {
	_ = c.Update(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte("emojis"))
		_ = b.Delete(toBytes(emojiId))

		return nil
	})
}

func (c *BoltCache) StoreVoiceState(state guild.VoiceState) {
	c.StoreVoiceStates([]guild.VoiceState{state})
}

func (c *BoltCache) StoreVoiceStates(states []guild.VoiceState) {
	_ = c.Batch(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte("voice_states"))

		for _, state := range states {
			if encoded, err := json.Marshal(state.ToCachedVoiceState()); err == nil {
				if err := b.Put(memberToBytes(state.UserId, state.GuildId), encoded); err != nil {
					return err
				}
			} else {
				return err
			}
		}

		return nil
	})
}

func (c *BoltCache) GetVoiceState(userId, guildId uint64) (guild.VoiceState, bool) {
	var cached guild.CachedVoiceState
	var found bool

	_ = c.View(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte("voice_states"))
		encoded := b.Get(memberToBytes(userId, guildId))

		if encoded == nil {
			return nil
		}

		return json.Unmarshal(encoded, &cached)
	})

	m, memberFound := c.GetMember(guildId, userId)
	if !memberFound {
		u, userFound := c.GetUser(userId)
		if !userFound {
			u = user.User{Id: userId}
		}

		m = member.Member{User: u}
	}

	state := cached.ToVoiceState(guildId, m)
	return state, found
}

func (c *BoltCache) GetGuildVoiceStates(guildId uint64) []guild.VoiceState {
	var states []guild.VoiceState

	_ = c.View(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte("voice_states"))

		return b.ForEach(func(k, encoded []byte) error {
			split := strings.Split(string(k), ":")
			if len(split) < 2 {
				return nil
			}

			// Hacky but w/e
			stateUserId, err := strconv.ParseUint(split[0], 10, 64); if err != nil {
				return nil
			}

			stateGuildId, err := strconv.ParseUint(split[1], 10, 64); if err != nil {
				return nil
			}

			var cached guild.CachedVoiceState
			if err := json.Unmarshal(encoded, &cached); err == nil && stateGuildId == guildId {
				m, memberFound := c.GetMember(guildId, stateUserId)
				if !memberFound {
					u, userFound := c.GetUser(stateUserId)
					if !userFound {
						u = user.User{Id: stateUserId}
					}

					m = member.Member{User: u}
				}

				states = append(states, cached.ToVoiceState(guildId, m))
			}

			return nil
		})
	})

	return states
}

func (c *BoltCache) DeleteVoiceState(userId, guildId uint64) {
	_ = c.Update(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte("voice_states"))
		_ = b.Delete(memberToBytes(userId, guildId))

		return nil
	})
}

func (c *BoltCache) StoreSelf(self user.User) {
	c.selfLock.Lock()
	c.self = self
	c.selfLock.Unlock()
}

func (c *BoltCache) GetSelf() (user.User, bool) {
	c.selfLock.RLock()
	self := c.self
	c.selfLock.RUnlock()

	return self, self.Id != 0
}

func boltMustCreate(tx *bolt.Tx, name string) {
	if _, err := tx.CreateBucketIfNotExists([]byte(name)); err != nil {
		panic(err)
	}
}

func memberToBytes(userId, guildId uint64) []byte {
	return []byte(fmt.Sprintf("%d:%d", userId, guildId))
}

func toBytes(i uint64) []byte {
	return []byte(strconv.FormatUint(i, 10))
}
