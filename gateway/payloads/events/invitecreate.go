package events

import "time"

type InviteCreate struct {
	ChannelId uint64    `json:"channel_id,string"`
	Code      string    `json:"code"`
	CreatedAt time.Time `json:"created_at"`
	GuildId   uint64    `json:"guild_id,string"`
	MaxAge    int       `json:"max_age"` // How long the invite is valid for, in seconds
	MaxUses   int       `json:"max_uses"`
	Temporary bool      `json:"temporary"`
	Uses      int       `json:"uses"` // Will always be 0
}
