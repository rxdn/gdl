package channel

import "time"

type CachedChannel struct {
	GuildId              uint64                `json:"-"`
	Type                 ChannelType           `json:"type"`
	Position             int                   `json:"position"`
	PermissionOverwrites []PermissionOverwrite `json:"permission_overwrites"`
	Name                 string                `json:"name"`
	Topic                string                `json:"topic"`
	Nsfw                 bool                  `json:"nsfw"`
	LastMessageId        uint64                `json:"last_message_id"`
	Bitrate              int                   `json:"bitrate"`
	UserLimit            int                   `json:"user_limit"`
	RateLimitPerUser     int                   `json:"rate_limit_per_user"`
	Icon                 string                `json:"icon"`
	OwnerId              uint64                `json:"owner_id,string"`
	ApplicationId        uint64                `json:"application_id"`
	ParentId             uint64                `json:"parent_id,string,omitempty"`
	LastPinTimestamp     time.Time             `json:"last_pin_timestamp"`
}

func (c *CachedChannel) ToChannel(channelId, guildId uint64) Channel {
	return Channel{
		Id:                   channelId,
		Type:                 c.Type,
		GuildId:              guildId,
		Position:             c.Position,
		PermissionOverwrites: c.PermissionOverwrites,
		Name:                 c.Name,
		Topic:                c.Topic,
		Nsfw:                 c.Nsfw,
		LastMessageId:        c.LastMessageId,
		Bitrate:              c.Bitrate,
		UserLimit:            c.UserLimit,
		RateLimitPerUser:     c.RateLimitPerUser,
		Recipients:           nil,
		Icon:                 c.Icon,
		OwnerId:              c.OwnerId,
		ApplicationId:        c.ApplicationId,
		ParentId:             c.ParentId,
		LastPinTimestamp:     c.LastPinTimestamp,
	}
}
