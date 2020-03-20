package events

import (
	"github.com/Dot-Rar/gdl/objects"
)

type GuildRoleCreateEvent func(*GuildRoleCreate)

type GuildRoleCreate struct {
	GuildId uint64       `json:"guild_id,string"`
	Role    objects.Role `json:"role"`
}

func (cc GuildRoleCreateEvent) Type() EventType {
	return GUILD_ROLE_CREATE
}

func (cc GuildRoleCreateEvent) New() interface{} {
	return &GuildRoleCreate{}
}

func (cc GuildRoleCreateEvent) Handle(i interface{}) {
	if t, ok := i.(*GuildRoleCreate); ok {
		cc(t)
	}
}
