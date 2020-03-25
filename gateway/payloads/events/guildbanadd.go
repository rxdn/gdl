package events

import "github.com/rxdn/gdl/objects"

type GuildBanAdd struct {
	GuildId uint64        `json:"guild_id,string"`
	User    *objects.User `json:"user"`
}
