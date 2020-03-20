package events

import (
	"github.com/Dot-Rar/gdl/objects"
)

type GuildMemberRemoveEvent func(*GuildMemberRemove)

type GuildMemberRemove struct {
	GuildId string          `json:"guild_id"`
	*objects.User
}

func (cc GuildMemberRemoveEvent) Type() EventType {
	return GUILD_MEMBER_REMOVE
}

func (cc GuildMemberRemoveEvent) New() interface{} {
	return &GuildMemberRemove{}
}

func (cc GuildMemberRemoveEvent) Handle(i interface{}) {
	if t, ok := i.(*GuildMemberRemove); ok {
		cc(t)
	}
}
