package objects

import "github.com/artemis-org/cache/discord/objects"

type WebhookType int

const (
	WebhookTypeIncoming        WebhookType = 1
	WebhookTypeChannelFollower WebhookType = 2
)

type Webhook struct {
	Id        uint64        `json:"id,string"`
	Type      WebhookType   `json:"type"`
	GuildId   uint64        `json:"guild_id,string,omitempty"`
	ChannelId uint64        `json:"channel_id,string"`
	User      *objects.User `json:"user"`
	Name      string        `json:"name,omitempty"`
	Avatar    string        `json:"avatar,omitempty"`
	Token     string        `json:"token,omitempty"`
}
