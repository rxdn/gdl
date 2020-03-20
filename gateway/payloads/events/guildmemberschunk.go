package events

import (
	"github.com/Dot-Rar/gdl/objects"
	"github.com/Dot-Rar/gdl/utils"
)

type GuildMembersChunk struct {
	GuildId   uint64                  `json:"guild_id,string"`
	Members   []*objects.Member       `json:"member"`
	NotFound  utils.Uint64StringSlice `json:"not_found,string"`
	Presences []*objects.Presence     `json:"presences"`
}
