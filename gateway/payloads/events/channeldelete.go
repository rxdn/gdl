package events

import (
	"github.com/Dot-Rar/gdl/objects"
)

type ChannelDeleteEvent func(*ChannelDelete)

type ChannelDelete struct {
	*objects.Channel
}

func (cc ChannelDeleteEvent) Type() EventType {
	return CHANNEL_DELETE
}

func (cc ChannelDeleteEvent) New() interface{} {
	return &ChannelDelete{}
}

func (cc ChannelDeleteEvent) Handle(i interface{}) {
	if t, ok := i.(*ChannelDelete); ok {
		cc(t)
	}
}
