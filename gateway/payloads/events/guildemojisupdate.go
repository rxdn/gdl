package events

import (
	"github.com/rxdn/gdl/objects/guild/emoji"
)

type GuildEmojisUpdate struct {
	GuildId uint64        `json:"guild_id,string"`
	Emojis  []emoji.Emoji `json:"emojis"`
}
