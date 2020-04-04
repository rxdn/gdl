package guild

import "github.com/rxdn/gdl/objects/user"

type Webhook struct {
	Id        uint64      `json:"id,string"`
	Type      WebhookType `json:"type"`
	GuildId   uint64      `json:"guild_id,string,omitempty"`
	ChannelId uint64      `json:"channel_id,string"`
	User      user.User   `json:"user"`
	Name      string      `json:"name,omitempty"`
	Avatar    string      `json:"avatar,omitempty"`
	Token     string      `json:"token,omitempty"`
}
