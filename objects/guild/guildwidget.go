package guild

import (
	"github.com/rxdn/gdl/objects/channel"
	"github.com/rxdn/gdl/objects/member"
)

type GuildWidget struct {
	Id            uint64            `json:"id,string"`
	Name          string            `json:"name"`
	InstantInvite string            `json:"instant_invite"`
	Channels      []channel.Channel `json:"channels"`
	Members       []member.Member   `json:"members"`
	PresenceCount int               `json:"presence_count"`
}
