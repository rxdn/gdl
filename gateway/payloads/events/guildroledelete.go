package events

import (
	"github.com/Dot-Rar/gdl/objects"
)

type GuildRoleDeleteEvent func(*GuildRoleDelete)

type GuildRoleDelete struct {
	GuildId uint64       `json:"guild_id,string"`
	Role    objects.Role `json:"role"`
}

func (cc GuildRoleDeleteEvent) Type() EventType {
	return GUILD_ROLE_DELETE
}

func (cc GuildRoleDeleteEvent) New() interface{} {
	return &GuildRoleDelete{}
}

func (cc GuildRoleDeleteEvent) Handle(i interface{}) {
	if t, ok := i.(*GuildRoleDelete); ok {
		cc(t)
	}
}
