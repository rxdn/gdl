package cache

import (
	"github.com/rxdn/gdl/objects/channel"
	"github.com/rxdn/gdl/objects/guild"
	"github.com/rxdn/gdl/objects/guild/emoji"
	"github.com/rxdn/gdl/objects/user"
	"sync"
)

type Cache interface {
	GetOptions() CacheOptions
	GetLock(uint64) *sync.RWMutex

	StoreUser(user *user.User)
	GetUser(id uint64) *user.User

	StoreGuild(guild *guild.Guild)
	GetGuild(id uint64) *guild.Guild
	GetGuilds() []*guild.Guild
	DeleteGuild(id uint64)

	StoreChannel(channel *channel.Channel)
	GetChannel(id uint64) *channel.Channel
	DeleteChannel(id uint64)

	StoreRole(role *guild.Role)
	GetRole(id uint64) *guild.Role
	DeleteRole(id uint64)

	StoreEmoji(emoji *emoji.Emoji)
	GetEmoji(id uint64) *emoji.Emoji
	DeleteEmoji(id uint64)

	StoreVoiceState(voiceState *guild.VoiceState)
	GetVoiceState(user uint64) *guild.VoiceState
	DeleteVoiceState(user uint64)

	StoreSelf(self *user.User)
	GetSelf() *user.User
}
