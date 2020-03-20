package events

import (
	"github.com/Dot-Rar/gdl/objects"
)

type ChannelCreateEvent func(*ChannelCreate)

type ChannelCreate struct {
	*objects.Channel
}

func (cc ChannelCreateEvent) Type() EventType {
	return CHANNEL_CREATE
}

func (cc ChannelCreateEvent) New() interface{} {
	return &ChannelCreate{}
}

func (cc ChannelCreateEvent) Handle(i interface{}) {
	if t, ok := i.(*ChannelCreate); ok {
		cc(t)
	}
}
