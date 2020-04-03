package events

import (
	"github.com/rxdn/gdl/objects/guild"
)

type VoiceStateUpdate struct {
	guild.VoiceState
}
