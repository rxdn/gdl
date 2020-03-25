package events

import (
	"github.com/rxdn/gdl/objects"
)

type GuildMemberAdd struct {
	*objects.Member
	GuildId uint64 `json:"guild_id,string"`
}
