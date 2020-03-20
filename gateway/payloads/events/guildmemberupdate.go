package events

import (
	"github.com/Dot-Rar/gdl/objects"
)

type GuildMemberUpdateEvent func(*GuildMemberUpdate)

type GuildMemberUpdate struct {
	GuildId string       `json:"guild_id"`
	Roles   []string     `json:"roles"`
	User    objects.User `json:"user"`
	Nick    string       `json:"nick"`
}

func (cc GuildMemberUpdateEvent) Type() EventType {
	return GUILD_MEMBER_UPDATE
}

func (cc GuildMemberUpdateEvent) New() interface{} {
	return &GuildMemberUpdate{}
}

func (cc GuildMemberUpdateEvent) Handle(i interface{}) {
	if t, ok := i.(*GuildMemberUpdate); ok {
		cc(t)
	}
}
