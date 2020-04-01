package events

import (
	"github.com/rxdn/gdl/objects/member"
)

type TypingStart struct {
	ChannelId uint64         `json:"channel_id,string"`
	GuildId   uint64         `json:"guild_id,string"`
	UserId    uint64         `json:"user_id,string"`
	Timestamp uint64         `json:"timestamp"` // Unix timestamp
	Member    *member.Member `json:"member"`
}
