package events

import (
	"github.com/Dot-Rar/gdl/objects"
)

type VoiceStateUpdateEvent func(*VoiceStateUpdate)

type VoiceStateUpdate struct {
	*objects.VoiceState
}

func (vs VoiceStateUpdateEvent) Type() EventType {
	return VOICE_STATE_UPDATE
}

func (vs VoiceStateUpdateEvent) New() interface{} {
	return &VoiceStateUpdate{}
}

func (vs VoiceStateUpdateEvent) Handle(i interface{}) {
	if t, ok := i.(*VoiceStateUpdate); ok {
		vs(t)
	}
}
