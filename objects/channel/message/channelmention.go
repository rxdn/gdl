package message

import "github.com/rxdn/gdl/objects/channel"

type ChannelMention struct {
	Id      uint64              `json:"id,string"`
	GuildId uint64              `json:"guild_id,string"`
	Type    channel.ChannelType `json:"type"`
	Name    string              `json:"name"` // channel name
}
