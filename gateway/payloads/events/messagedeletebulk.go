package events

import "github.com/rxdn/gdl/utils"

type MessageDeleteBulk struct {
	Id        utils.Uint64StringSlice `json:"ids"`
	ChannelId uint64                  `json:"channel_id,string"`
	GuildId   uint64                  `json:"guild_id,string"`
}
