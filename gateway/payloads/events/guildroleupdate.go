package events

import (
	"github.com/Dot-Rar/gdl/objects"
)

type GuildRoleUpdateEvent func(*GuildRoleUpdate)

type GuildRoleUpdate struct {
	GuildId string       `json:"guild_id"`
	Role    objects.Role `json:"role"`
}

func (cc GuildRoleUpdateEvent) Type() EventType {
	return GUILD_ROLE_UPDATE
}

func (cc GuildRoleUpdateEvent) New() interface{} {
	return &GuildRoleUpdate{}
}

func (cc GuildRoleUpdateEvent) Handle(i interface{}) {
	if t, ok := i.(*GuildRoleUpdate); ok {
		cc(t)
	}
}
