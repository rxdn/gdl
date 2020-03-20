package cache

import "github.com/Dot-Rar/gdl/objects"

type Cache interface {
	StoreUser(user *objects.User)
	GetUser(id uint64) *objects.User

	StoreGuild(guild *objects.Guild)
	GetGuild(id uint64) *objects.Guild

	StoreChannel(channel *objects.Channel)
	GetChannel(id uint64) *objects.Channel

	StoreRole(role *objects.Role)
	GetRole(id uint64) *objects.Role

	StoreEmoji(emoji *objects.Emoji)
	GetEmoji(id uint64) *objects.Emoji
}
