package channel

import (
	"fmt"
	"github.com/rxdn/gdl/objects"
	"github.com/rxdn/gdl/objects/user"
	"time"
)

type Channel struct {
	Id                   uint64                    `json:"id,string"`
	Type                 ChannelType               `json:"type"`
	GuildId              uint64                    `json:"guild_id,string"`
	Position             int                       `json:"position"`
	PermissionOverwrites []PermissionOverwrite     `json:"permission_overwrites"`
	Name                 string                    `json:"name"`
	Topic                string                    `json:"topic"`
	Nsfw                 bool                      `json:"nsfw"`
	LastMessageId        objects.NullableSnowflake `json:"last_message_id"`
	Bitrate              int                       `json:"bitrate"`
	UserLimit            int                       `json:"user_limit"`
	RateLimitPerUser     int                       `json:"rate_limit_per_user"`
	Recipients           []user.User               `json:"recipients"`
	Icon                 string                    `json:"icon"`
	OwnerId              uint64                    `json:"owner_id,string"` // Owner of a group DM
	ApplicationId        uint64                    `json:"application_id,string"`
	ParentId             objects.NullableSnowflake `json:"parent_id,omitempty"`
	LastPinTimestamp     time.Time                 `json:"last_pin_timestamp"`
	RtcRegion            *string                   `json:"rtc_region"`
	VideoQualityMode     VideoQualityMode          `json:"video_quality_mode"`
	MessageCount         uint64                    `json:"message_count"`
	MemberCount          uint64                    `json:"member_count"`
	ThreadMetadata       *ThreadMetadata           `json:"thread_metadata,omitempty"`
	Member               ThreadMember              `json:"member"`
}

func (c *Channel) Mention() string {
	return fmt.Sprintf("<#%d>", c.Id)
}

func (c *Channel) ToCachedChannel() CachedChannel {
	return CachedChannel{
		GuildId:              c.GuildId,
		Type:                 c.Type,
		Position:             c.Position,
		PermissionOverwrites: c.PermissionOverwrites,
		Name:                 c.Name,
		Topic:                c.Topic,
		Nsfw:                 c.Nsfw,
		LastMessageId:        c.LastMessageId,
		Bitrate:              c.Bitrate,
		UserLimit:            c.UserLimit,
		RateLimitPerUser:     c.RateLimitPerUser,
		Icon:                 c.Icon,
		OwnerId:              c.OwnerId,
		ApplicationId:        c.ApplicationId,
		ParentId:             c.ParentId,
		LastPinTimestamp:     c.LastPinTimestamp,
	}
}
