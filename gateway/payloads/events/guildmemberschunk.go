package events

import (
	"github.com/Dot-Rar/gdl/objects"
)

type GuildMembersChunkEvent func(*GuildMembersChunk)

type GuildMembersChunk struct {
	GuildId string           `json:"guild_id"`
	Members []objects.Member `json:"members"`
}

func (cc GuildMembersChunkEvent) Type() EventType {
	return GUILD_MEMBERS_CHUNK
}

func (cc GuildMembersChunkEvent) New() interface{} {
	return &GuildMembersChunk{}
}

func (cc GuildMembersChunkEvent) Handle(i interface{}) {
	if t, ok := i.(*GuildMembersChunk); ok {
		cc(t)
	}
}
