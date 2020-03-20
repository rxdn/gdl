package events

import (
	"github.com/Dot-Rar/gdl/objects"
)

type GuildMemberAdd struct {
	*objects.Member
	GuildId uint64 `json:"guild_id,string"`
}
