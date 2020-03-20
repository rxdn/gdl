package events

import (
	"github.com/Dot-Rar/gdl/objects"
)

type GuildMemberAddEvent func(*GuildMemberAdd)

type GuildMemberAdd struct {
	GuildId string          `json:"guild_id"`
	*objects.Member
}

func (cc GuildMemberAddEvent) Type() EventType {
	return GUILD_MEMBER_ADD
}

func (cc GuildMemberAddEvent) New() interface{} {
	return &GuildMemberAdd{}
}

func (cc GuildMemberAddEvent) Handle(i interface{}) {
	if t, ok := i.(*GuildMemberAdd); ok {
		cc(t)
	}
}
