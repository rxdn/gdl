package events

import (
	"github.com/rxdn/gdl/objects/guild"
)

type GuildRoleUpdate struct {
	GuildId uint64     `json:"guild_id,string"`
	Role    guild.Role ` json:"role"`
}
