package gateway

import (
	"github.com/Dot-Rar/gdl/gateway/payloads/events"
	"github.com/Dot-Rar/gdl/objects"
	"github.com/sirupsen/logrus"
)

func RegisterCacheListeners(sm *ShardManager) {
	sm.RegisterListeners(
		readyListener,
		channelCreateListener,
		channelUpdateListener,
		channelDeleteListener,
		channelPinsUpdateListener,
		guildCreateListener,
		guildUpdateListener,
		guildDeleteListener,
		guildEmojisUpdateListeners,
		guildMemberAddListener,
		guildMemberRemoveListener,
		guildMemberUpdateListener,
		guildMembersChunkListener,
		guildRoleCreateListener,
		guildRoleUpdateListener,
		guildRoleDeleteListener,
		userUpdateListener,
		voiceStateUpdateListener,
	)
}

func readyListener(s *Shard, e *events.Ready) {
	logrus.Infof("shard %d: received ready", s.ShardId)

	s.SessionId = e.SessionId

	if (*s.Cache).GetOptions().Guilds {
		for _, guild := range e.Guilds {
			(*s.Cache).GetLock(guild.Id).Lock()
			(*s.Cache).StoreGuild(guild)
			(*s.Cache).GetLock(guild.Id).Unlock()
		}
	}
}

func channelCreateListener(s *Shard, e *events.ChannelCreate) {
	if (*s.Cache).GetOptions().Channels {
		(*s.Cache).GetLock(e.Channel.Id).Lock()
		(*s.Cache).StoreChannel(e.Channel)
		(*s.Cache).GetLock(e.Channel.Id).Unlock()
	}
}

func channelUpdateListener(s *Shard, e *events.ChannelUpdate) {
	if (*s.Cache).GetOptions().Channels {
		(*s.Cache).GetLock(e.Channel.Id).Lock()
		(*s.Cache).StoreChannel(e.Channel)
		(*s.Cache).GetLock(e.Channel.Id).Unlock()
	}
}

func channelDeleteListener(s *Shard, e *events.ChannelDelete) {
	if (*s.Cache).GetOptions().Channels {
		(*s.Cache).GetLock(e.Channel.Id).Lock()
		(*s.Cache).DeleteChannel(e.Channel.Id)
		(*s.Cache).GetLock(e.Channel.Id).Unlock()
	}
}

func channelPinsUpdateListener(s *Shard, e *events.ChannelPinsUpdate) {
	if (*s.Cache).GetOptions().Channels {
		(*s.Cache).GetLock(e.ChannelId).Lock()

		channel := (*s.Cache).GetChannel(e.ChannelId)
		channel.LastPinTimestamp = e.LastPinTimestamp
		(*s.Cache).StoreChannel(channel)

		(*s.Cache).GetLock(e.ChannelId).Unlock()
	}
}

func guildCreateListener(s *Shard, e *events.GuildCreate) {
	if (*s.Cache).GetOptions().Guilds {
		modified := *e.Guild
		if !(*s.Cache).GetOptions().Users {
			modified.Members = make([]*objects.Member, 0)
		}
		if !(*s.Cache).GetOptions().Channels {
			modified.Channels = make([]*objects.Channel, 0)
		}
		if !(*s.Cache).GetOptions().Roles {
			modified.Roles = make([]*objects.Role, 0)
		}
		if !(*s.Cache).GetOptions().Emojis {
			modified.Emojis = make([]*objects.Emoji, 0)
		}

		(*s.Cache).GetLock(e.Guild.Id).Lock()
		(*s.Cache).StoreGuild(&modified)
		(*s.Cache).GetLock(e.Guild.Id).Unlock()
	}

	if (*s.Cache).GetOptions().Users {
		for _, member := range e.Members {
			(*s.Cache).GetLock(member.User.Id).Lock()
			(*s.Cache).StoreUser(member.User)
			(*s.Cache).GetLock(member.User.Id).Unlock()
		}
	}

	if (*s.Cache).GetOptions().Channels {
		for _, channel := range e.Channels {
			(*s.Cache).GetLock(channel.Id).Lock()
			(*s.Cache).StoreChannel(channel)
			(*s.Cache).GetLock(channel.Id).Unlock()
		}
	}

	if (*s.Cache).GetOptions().Roles {
		for _, role := range e.Roles {
			(*s.Cache).GetLock(role.Id).Lock()
			(*s.Cache).StoreRole(role)
			(*s.Cache).GetLock(role.Id).Unlock()
		}
	}

	if (*s.Cache).GetOptions().Emojis {
		for _, emoji := range e.Emojis {
			(*s.Cache).GetLock(emoji.Id).Lock()
			(*s.Cache).StoreEmoji(emoji)
			(*s.Cache).GetLock(emoji.Id).Unlock()
		}
	}
}

func guildUpdateListener(s *Shard, e *events.GuildUpdate) {
	if (*s.Cache).GetOptions().Guilds {
		modified := *e.Guild
		if !(*s.Cache).GetOptions().Users {
			modified.Members = make([]*objects.Member, 0)
		}
		if !(*s.Cache).GetOptions().Channels {
			modified.Channels = make([]*objects.Channel, 0)
		}
		if !(*s.Cache).GetOptions().Roles {
			modified.Roles = make([]*objects.Role, 0)
		}
		if !(*s.Cache).GetOptions().Emojis {
			modified.Emojis = make([]*objects.Emoji, 0)
		}

		(*s.Cache).GetLock(e.Guild.Id).Lock()
		(*s.Cache).StoreGuild(&modified)
		(*s.Cache).GetLock(e.Guild.Id).Unlock()
	}

	if (*s.Cache).GetOptions().Users {
		for _, member := range e.Members {
			(*s.Cache).GetLock(member.User.Id).Lock()
			(*s.Cache).StoreUser(member.User)
			(*s.Cache).GetLock(member.User.Id).Unlock()
		}
	}

	if (*s.Cache).GetOptions().Channels {
		for _, channel := range e.Channels {
			(*s.Cache).GetLock(channel.Id).Lock()
			(*s.Cache).StoreChannel(channel)
			(*s.Cache).GetLock(channel.Id).Unlock()
		}
	}

	if (*s.Cache).GetOptions().Roles {
		for _, role := range e.Roles {
			(*s.Cache).GetLock(role.Id).Lock()
			(*s.Cache).StoreRole(role)
			(*s.Cache).GetLock(role.Id).Unlock()
		}
	}

	if (*s.Cache).GetOptions().Emojis {
		for _, emoji := range e.Emojis {
			(*s.Cache).GetLock(emoji.Id).Lock()
			(*s.Cache).StoreEmoji(emoji)
			(*s.Cache).GetLock(emoji.Id).Unlock()
		}
	}
}

func guildDeleteListener(s *Shard, e *events.GuildDelete) {
	if (*s.Cache).GetOptions().Guilds && e.Unavailable == nil { // If the unavailable field is not set, the user was removed from the guild
		(*s.Cache).GetLock(e.Id).Lock()
		(*s.Cache).DeleteGuild(e.Id)
		(*s.Cache).GetLock(e.Id).Unlock()
	}
}

func guildEmojisUpdateListeners(s *Shard, e *events.GuildEmojisUpdate) {
	if (*s.Cache).GetOptions().Emojis {
		for _, emoji := range e.Emojis {
			(*s.Cache).GetLock(emoji.Id).Lock()
			(*s.Cache).StoreEmoji(emoji)
			(*s.Cache).GetLock(emoji.Id).Unlock()
		}

		if (*s.Cache).GetOptions().Guilds {
			(*s.Cache).GetLock(e.GuildId).Lock()

			guild := (*s.Cache).GetGuild(e.GuildId)
			guild.Emojis = e.Emojis
			(*s.Cache).StoreGuild(guild)

			(*s.Cache).GetLock(e.GuildId).Unlock()
		}
	}
}

func guildMemberAddListener(s *Shard, e *events.GuildMemberAdd) {
	if (*s.Cache).GetOptions().Users {
		(*s.Cache).GetLock(e.User.Id).Lock()
		(*s.Cache).StoreUser(e.User)
		(*s.Cache).GetLock(e.User.Id).Unlock()

		if (*s.Cache).GetOptions().Guilds {
			(*s.Cache).GetLock(e.GuildId).Lock()

			guild := (*s.Cache).GetGuild(e.GuildId)
			guild.Members = append(guild.Members, e.Member)
			(*s.Cache).StoreGuild(guild)

			(*s.Cache).GetLock(e.GuildId).Unlock()
		}
	}
}

func guildMemberRemoveListener(s *Shard, e *events.GuildMemberRemove) {
	(*s.Cache).GetLock(e.GuildId).Lock()
	guild := (*s.Cache).GetGuild(e.GuildId)

	index := -1
	for i, member := range guild.Members {
		if member != nil && member.User != nil && member.User.Id == e.User.Id {
			index = i
		}
	}

	if index != -1 {
		guild.Members = append(guild.Members[:index], guild.Members[index+1:]...)
	}

	(*s.Cache).StoreGuild(guild)
	(*s.Cache).GetLock(e.GuildId).Unlock()
}

func guildMemberUpdateListener(s *Shard, e *events.GuildMemberUpdate) {
	if (*s.Cache).GetOptions().Users {
		(*s.Cache).GetLock(e.User.Id).Lock()
		(*s.Cache).StoreUser(e.User)
		(*s.Cache).GetLock(e.User.Id).Unlock()

		// Update guild
		if (*s.Cache).GetOptions().Guilds {
			(*s.Cache).GetLock(e.GuildId).Lock()

			guild := (*s.Cache).GetGuild(e.GuildId)
			for _, member := range guild.Members {
				if member != nil && member.User != nil && member.User.Id == e.User.Id {
					member.Roles = e.Roles
					member.User = e.User
					member.Nick = e.Nick
					break
				}
			}

			(*s.Cache).GetLock(e.GuildId).Unlock()
		}
	}
}

func guildMembersChunkListener(s *Shard, e *events.GuildMembersChunk) {
	if (*s.Cache).GetOptions().Users {
		for _, member := range e.Members {
			(*s.Cache).GetLock(member.User.Id)
			(*s.Cache).StoreUser(member.User)
			(*s.Cache).GetLock(member.User.Id).Unlock()
		}

		if (*s.Cache).GetOptions().Guilds {
			(*s.Cache).GetLock(e.GuildId).Lock()

			guild := (*s.Cache).GetGuild(e.GuildId)

			// Create new member slice
			members := (*e).Members
			for _, existingMember := range guild.Members {
				found := false
				internal:
				for _, newMember := range members {
					if existingMember.User.Id == newMember.User.Id {
						found = true
						break internal
					}
				}

				if !found {
					members = append(members, existingMember)
				}
			}

			guild.Members = members
			(*s.Cache).StoreGuild(guild)

			(*s.Cache).GetLock(e.GuildId).Unlock()
		}
	}
}

func guildRoleCreateListener(s *Shard, e *events.GuildRoleCreate) {
	if (*s.Cache).GetOptions().Roles {
		(*s.Cache).GetLock(e.Role.Id).Lock()
		(*s.Cache).StoreRole(e.Role)
		(*s.Cache).GetLock(e.Role.Id).Unlock()

		// Update guild
		if (*s.Cache).GetOptions().Guilds {
			(*s.Cache).GetLock(e.GuildId).Lock()

			guild := (*s.Cache).GetGuild(e.GuildId)
			guild.Roles = append(guild.Roles, e.Role)
			(*s.Cache).StoreGuild(guild)

			(*s.Cache).GetLock(e.GuildId).Unlock()
		}
	}
}

func guildRoleUpdateListener(s *Shard, e *events.GuildRoleUpdate) {
	if (*s.Cache).GetOptions().Roles {
		(*s.Cache).GetLock(e.Role.Id).Lock()
		(*s.Cache).StoreRole(e.Role)
		(*s.Cache).GetLock(e.Role.Id).Unlock()

		// Update guild
		if (*s.Cache).GetOptions().Guilds {
			(*s.Cache).GetLock(e.GuildId).Lock()

			guild := (*s.Cache).GetGuild(e.GuildId)
			guild.Roles = append(guild.Roles, e.Role)
			(*s.Cache).StoreGuild(guild)

			(*s.Cache).GetLock(e.GuildId).Unlock()
		}
	}
}

func guildRoleDeleteListener(s *Shard, e *events.GuildRoleDelete) {
	if (*s.Cache).GetOptions().Roles {
		(*s.Cache).GetLock(e.RoleId).Lock()
		(*s.Cache).DeleteRole(e.RoleId)
		(*s.Cache).GetLock(e.RoleId).Unlock()

		// Update guild
		if (*s.Cache).GetOptions().Guilds {
			(*s.Cache).GetLock(e.GuildId).Lock()

			guild := (*s.Cache).GetGuild(e.GuildId)

			index := -1
			for i, role := range guild.Roles {
				if role != nil && role.Id == e.RoleId {
					index = i
				}
			}

			if index != -1 {
				guild.Roles = append(guild.Roles[:index], guild.Roles[index+1:]...)
			}

			(*s.Cache).StoreGuild(guild)

			(*s.Cache).GetLock(e.GuildId).Unlock()
		}
	}
}

func userUpdateListener(s *Shard, e *events.UserUpdate) {
	if (*s.Cache).GetOptions().Users {
		(*s.Cache).GetLock(e.Id).Lock()
		(*s.Cache).StoreUser(e.User)
		(*s.Cache).GetLock(e.Id).Unlock()
	}
}

func voiceStateUpdateListener(s *Shard, e *events.VoiceStateUpdate) {
	if (*s.Cache).GetOptions().VoiceStates {

		(*s.Cache).GetLock(e.UserId).Lock()
		(*s.Cache).StoreVoiceState(e.VoiceState)
		(*s.Cache).GetLock(e.UserId).Unlock()

		// Update guild
		if (*s.Cache).GetOptions().Guilds {
			(*s.Cache).GetLock(e.GuildId).Lock()

			guild := (*s.Cache).GetGuild(e.GuildId)
			guild.VoiceStates = append(guild.VoiceStates, e.VoiceState)
			(*s.Cache).StoreGuild(guild)

			(*s.Cache).GetLock(e.GuildId).Unlock()
		}
	}
}
