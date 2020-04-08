package channel

import "time"

type CachedChannel struct {
	Type                 ChannelType           `db:"type"`
	Position             int                   `db:"position"`
	PermissionOverwrites []PermissionOverwrite `db:"permission_overwrites"`
	Name                 string                `db:"name"`
	Topic                string                `db:"topic"`
	Nsfw                 bool                  `db:"nsfw"`
	LastMessageId        uint64                `db:"last_message_id"`
	Bitrate              int                   `db:"bitrate"`
	UserLimit            int                   `db:"user_limit"`
	RateLimitPerUser     int                   `db:"rate_limit_per_user"`
	Icon                 string                `db:"icon"`
	OwnerId              uint64                `db:"owner_id,string"`
	ApplicationId        uint64                `db:"application_id"`
	ParentId             uint64                `db:"parent_id,string,omitempty"`
	LastPinTimestamp     time.Time             `db:"last_pin_timestamp"`
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
