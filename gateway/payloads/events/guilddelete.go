package events

import (
	"github.com/Dot-Rar/gdl/objects"
)

type GuildDeleteEvent func(*GuildDelete)

type GuildDelete struct {
	*objects.Guild
}

func (cc GuildDeleteEvent) Type() EventType {
	return GUILD_DELETE
}

func (cc GuildDeleteEvent) New() interface{} {
	return &GuildDelete{}
}

func (cc GuildDeleteEvent) Handle(i interface{}) {
	if t, ok := i.(*GuildDelete); ok {
		cc(t)
	}
}
