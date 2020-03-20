package events

import (
	"github.com/Dot-Rar/gdl/objects"
)

type GuildEmojisUpdateEvent func(*GuildEmojisUpdate)

// No need to completely re-hash the entire guild object
type GuildEmojisUpdate struct {
	GuildId string          `json:"guild_id"`
	Emojis  []objects.Emoji `json:"emoji"`
}

func (cc GuildEmojisUpdateEvent) Type() EventType {
	return GUILD_EMOJIS_UPDATE
}

func (cc GuildEmojisUpdateEvent) New() interface{} {
	return &GuildEmojisUpdate{}
}

func (cc GuildEmojisUpdateEvent) Handle(i interface{}) {
	if t, ok := i.(*GuildEmojisUpdate); ok {
		cc(t)
	}
}
