package objects

import (
	"time"
)

type Channel struct {
	Id                    uint64       `json:"id,string"`
	Type                  int          `json:"type"`
	GuildId               uint64       `json:"guild_id,string"`
	Position              int          `json:"position"`
	PermissionsOverwrites []*Overwrite `json:"permission_overwrites"`
	Name                  string       `json:"name"`
	Topic                 string       `json:"topic"`
	Nsfw                  bool         `json:"nsfw"`
	LastMessageId         uint64       `json:"last_message_id,string"`
	Bitrate               int          `json:"bitrate"`
	UserLimit             int          `json:"user_limit"`
	RateLimitPerUser      int          `json:"rate_limit_per_user"`
	Recipients            []*User      `json:"recipients"`
	Icon                  string       `json:"icon"`
	OwnerId               uint64       `json:"owner_id,string"`
	ApplicationId         uint64       `json:"application_id,string"`
	ParentId              uint64       `json:"parent_id,string"`
	LastPinTimestamp      time.Time    `json:"last_pin_timestamp"`
}
