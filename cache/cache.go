package cache

import (
	"github.com/Dot-Rar/gdl/objects"
	"sync"
)

type Cache interface {
	GetOptions() CacheOptions
	GetLock(uint64) *sync.RWMutex

	StoreUser(user *objects.User)
	GetUser(id uint64) *objects.User

	StoreGuild(guild *objects.Guild)
	GetGuild(id uint64) *objects.Guild
	DeleteGuild(id uint64)

	StoreChannel(channel *objects.Channel)
	GetChannel(id uint64) *objects.Channel
	DeleteChannel(id uint64)

	StoreRole(role *objects.Role)
	GetRole(id uint64) *objects.Role
	DeleteRole(id uint64)

	StoreEmoji(emoji *objects.Emoji)
	GetEmoji(id uint64) *objects.Emoji
	DeleteEmoji(id uint64)

	StoreVoiceState(voiceState *objects.VoiceState)
	GetVoiceState(user uint64) *objects.VoiceState
	DeleteVoiceState(user uint64)

	StoreSelf(self *objects.User)
	GetSelf() *objects.User
}
