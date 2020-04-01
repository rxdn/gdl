package events

import (
	"github.com/rxdn/gdl/objects/guild"
)

type GuildDelete struct {
	*guild.Guild
}
