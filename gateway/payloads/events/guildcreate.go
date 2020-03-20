package events

import (
	"github.com/Dot-Rar/gdl/objects"
)

type GuildCreateEvent func(*GuildCreate)

type GuildCreate struct {
	*objects.Guild
}

func (cc GuildCreateEvent) Type() EventType {
	return GUILD_CREATE
}

func (cc GuildCreateEvent) New() interface{} {
	return &GuildCreate{}
}

func (cc GuildCreateEvent) Handle(i interface{}) {
	if t, ok := i.(*GuildCreate); ok {
		cc(t)
	}
}
