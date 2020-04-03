package events

import (
	"github.com/rxdn/gdl/objects/member"
)

type GuildMemberAdd struct {
	member.Member
	GuildId uint64 `json:"guild_id,string"`
}
