package gateway

import (
	"github.com/Dot-Rar/gdl/gateway/payloads/events"
)

func RegisterCacheListeners(sm *ShardManager) {
	sm.RegisterListeners(GuildCreateListener)
}

func GuildCreateListener(s *Shard, e *events.GuildCreate) {
	(*s.Cache).StoreGuild(e.Guild)

	for _, member := range e.Members {
		(*s.Cache).StoreUser(member.User)
	}

	for _, channel := range e.Channels {
		(*s.Cache).StoreChannel(channel)
	}

	for _, role := range e.Roles {
		(*s.Cache).StoreRole(role)
	}

	for _, emoji := range e.Emojis {
		(*s.Cache).StoreEmoji(emoji)
	}
}
