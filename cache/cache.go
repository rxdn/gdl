package cache

import (
	"github.com/rxdn/gdl/objects/channel"
	"github.com/rxdn/gdl/objects/guild"
	"github.com/rxdn/gdl/objects/guild/emoji"
	"github.com/rxdn/gdl/objects/member"
	"github.com/rxdn/gdl/objects/user"
)

type Cache interface {
	GetOptions() CacheOptions

	StoreUser(user user.User)
	StoreUsers(user []user.User)
	GetUser(id uint64) (user.User, bool)
	GetUsers(ids []uint64) (map[uint64]user.User, error)

	StoreGuild(guild guild.Guild)
	StoreGuilds(guilds []guild.Guild)
	GetGuild(id uint64, withMembers bool) (guild.Guild, bool)
	GetGuilds() []guild.Guild // Note: Guilds will not have Channels, Roles, Members etc to reduce cache lookup time
	DeleteGuild(id uint64)
	GetGuildCount() int
	GetGuildOwner(guildId uint64) (uint64, bool) // Utility function

	StoreMember(member member.Member, guildId uint64)
	StoreMembers(members []member.Member, guildId uint64)
	GetMember(guildId, userId uint64) (member.Member, bool)
	GetGuildMembers(guildId uint64, withUserData bool) []member.Member
	DeleteMember(userId, guildId uint64)

	StoreChannel(channel channel.Channel)
	StoreChannels(channel []channel.Channel)
	GetChannel(id uint64) (channel.Channel, bool)
	GetGuildChannels(guildId uint64) []channel.Channel
	DeleteChannel(channelId uint64)
	DeleteGuildChannels(guildId uint64)

	StoreRole(role guild.Role, guildId uint64)
	StoreRoles(roles []guild.Role, guildId uint64)
	GetRole(id uint64) (guild.Role, bool)
	GetRoles(guildId uint64, ids []uint64) (map[uint64]guild.Role, error)
	GetGuildRoles(guildId uint64) []guild.Role
	DeleteRole(roleId uint64)
	DeleteGuildRoles(guildId uint64)

	StoreEmoji(emoji emoji.Emoji, guildId uint64)
	StoreEmojis(emojis []emoji.Emoji, guildId uint64)
	GetEmoji(id uint64) (emoji.Emoji, bool)
	GetGuildEmojis(id uint64) []emoji.Emoji
	DeleteEmoji(emojiId uint64)

	StoreVoiceState(voiceState guild.VoiceState)
	StoreVoiceStates(voiceStates []guild.VoiceState)
	GetVoiceState(userId, guildId uint64) (guild.VoiceState, bool)
	GetGuildVoiceStates(guildId uint64) []guild.VoiceState
	DeleteVoiceState(userId, guildId uint64)

	StoreSelf(self user.User)
	GetSelf() (user.User, bool)
}
