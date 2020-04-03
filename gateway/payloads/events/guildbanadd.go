package events

import (
	"github.com/rxdn/gdl/objects/user"
)

type GuildBanAdd struct {
	GuildId uint64    `json:"guild_id,string"`
	User    user.User `json:"user"`
}
