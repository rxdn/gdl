package events

import "github.com/rxdn/gdl/objects"

type GuildEmojisUpdate struct {
	GuildId uint64           `json:"guild_id,string"`
	Emojis  []*objects.Emoji `json:"emojis"`
}
