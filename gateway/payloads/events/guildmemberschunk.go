package events

import (
	"github.com/rxdn/gdl/objects/member"
	"github.com/rxdn/gdl/objects/user"
	"github.com/rxdn/gdl/utils"
)

type GuildMembersChunk struct {
	GuildId   uint64                  `json:"guild_id,string"`
	Members   []*member.Member        `json:"member"`
	NotFound  utils.Uint64StringSlice `json:"not_found,string"`
	Presences []*user.Presence        `json:"presences"`
}
