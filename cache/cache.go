package cache

import (
	"context"
	"errors"
	"github.com/rxdn/gdl/objects/channel"
	"github.com/rxdn/gdl/objects/guild"
	"github.com/rxdn/gdl/objects/guild/emoji"
	"github.com/rxdn/gdl/objects/member"
	"github.com/rxdn/gdl/objects/user"
)

var ErrNotFound = errors.New("object not found in cache")

type Cache interface {
	Options() CacheOptions

	StoreUser(ctx context.Context, user user.User) error
	StoreUsers(ctx context.Context, user []user.User) error
	GetUser(ctx context.Context, id uint64) (user.User, error)
	GetUsers(ctx context.Context, ids []uint64) (map[uint64]user.User, error)

	StoreGuild(ctx context.Context, guild guild.Guild) error
	StoreGuilds(ctx context.Context, guilds []guild.Guild) error
	GetGuild(ctx context.Context, id uint64) (guild.Guild, error)
	DeleteGuild(ctx context.Context, id uint64) error
	GetGuildCount(ctx context.Context) (int, error)
	GetGuildOwner(ctx context.Context, guildId uint64) (uint64, error)

	StoreMember(ctx context.Context, member member.Member, guildId uint64) error
	StoreMembers(ctx context.Context, members []member.Member, guildId uint64) error
	GetMember(ctx context.Context, guildId, userId uint64) (member.Member, error)
	GetGuildMembers(ctx context.Context, guildId uint64, withUserData bool) ([]member.Member, error)
	DeleteMember(ctx context.Context, userId, guildId uint64) error

	StoreChannel(ctx context.Context, channel channel.Channel) error
	StoreChannels(ctx context.Context, channel []channel.Channel) error
	GetChannel(ctx context.Context, id uint64) (channel.Channel, error)
	GetGuildChannels(ctx context.Context, guildId uint64) ([]channel.Channel, error)
	DeleteChannel(ctx context.Context, channelId uint64) error
	DeleteGuildChannels(ctx context.Context, guildId uint64) error

	StoreRole(ctx context.Context, role guild.Role, guildId uint64) error
	StoreRoles(ctx context.Context, roles []guild.Role, guildId uint64) error
	GetRole(ctx context.Context, id uint64) (guild.Role, error)
	GetRoles(ctx context.Context, guildId uint64, ids []uint64) (map[uint64]guild.Role, error)
	GetGuildRoles(ctx context.Context, guildId uint64) ([]guild.Role, error)
	DeleteRole(ctx context.Context, roleId uint64) error
	DeleteGuildRoles(ctx context.Context, guildId uint64) error

	StoreEmoji(ctx context.Context, emoji emoji.Emoji, guildId uint64) error
	StoreEmojis(ctx context.Context, emojis []emoji.Emoji, guildId uint64) error
	GetEmoji(ctx context.Context, id uint64) (emoji.Emoji, error)
	GetGuildEmojis(ctx context.Context, id uint64) ([]emoji.Emoji, error)
	DeleteEmoji(ctx context.Context, emojiId uint64) error

	StoreVoiceState(ctx context.Context, voiceState guild.VoiceState) error
	StoreVoiceStates(ctx context.Context, voiceStates []guild.VoiceState) error
	GetVoiceState(ctx context.Context, userId, guildId uint64) (guild.VoiceState, error)
	GetGuildVoiceStates(ctx context.Context, guildId uint64) ([]guild.VoiceState, error)
	DeleteVoiceState(ctx context.Context, userId, guildId uint64) error

	StoreSelf(ctx context.Context, self user.User) error
	GetSelf(ctx context.Context) (user.User, error)
}
