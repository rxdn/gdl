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
	GetUser(id uint64) (user.User, bool)

	StoreGuild(guild guild.Guild)
	GetGuild(id uint64) (guild.Guild, bool)
	GetGuilds() []guild.Guild
	DeleteGuild(id uint64)

	StoreMember(member member.Member, guildId uint64)
	DeleteMember(userId, guildId uint64)

	StoreChannel(channel channel.Channel)
	GetChannel(id uint64) (channel.Channel, bool)
	DeleteChannel(channelId, guildId uint64)

	StoreRole(role guild.Role, guildId uint64)
	GetRole(id uint64) (guild.Role, bool)
	DeleteRole(roleId, guildId uint64)

	StoreEmoji(emoji emoji.Emoji, guildId uint64)
	GetEmoji(id uint64) (emoji.Emoji, bool)
	DeleteEmoji(emojiId, guildId uint64)

	StoreVoiceState(voiceState guild.VoiceState)
	GetVoiceState(userId, guildId uint64) (guild.VoiceState, bool)

	StoreSelf(self user.User)
	GetSelf() (user.User, bool)
}
