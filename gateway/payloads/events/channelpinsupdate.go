package events

import (
	"time"
)

type ChannelPinsUpdate struct {
	GuildId          uint64    `json:"guild_id,string"`
	ChannelId        uint64    `json:"channel_id,string"`
	LastPinTimestamp time.Time `json:"last_pin_timestamp"`
}
