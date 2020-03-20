package events

import "github.com/Dot-Rar/gdl/objects"

type GuildRoleCreate struct {
	GuildId uint64        `json:"guild_id,string"`
	Role    *objects.Role ` json:"role"`
}
