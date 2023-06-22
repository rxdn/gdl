package channel

import (
	"github.com/rxdn/gdl/objects"
	"time"
)

type PartialChannel struct {
	Id               uint64                    `json:"id,string"`
	Type             ChannelType               `json:"type"`
	GuildId          uint64                    `json:"guild_id,string"`
	Position         int                       `json:"position"`
	Name             string                    `json:"name"`
	Topic            string                    `json:"topic"`
	Nsfw             bool                      `json:"nsfw"`
	LastMessageId    objects.NullableSnowflake `json:"last_message_id"`
	ParentId         objects.NullableSnowflake `json:"parent_id,omitempty"`
	LastPinTimestamp time.Time                 `json:"last_pin_timestamp"`
}
