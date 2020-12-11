package interaction

import "github.com/rxdn/gdl/objects/member"

type Interaction struct {
	Id        uint64                             `json:"id"`
	Type      InteractionType                    `json:"type"`
	Data      *ApplicationCommandInteractionData `json:"data"`
	GuildId   uint64                             `json:"guild_id"`
	ChannelId uint64                             `json:"channel_id"`
	Member    member.Member                      `json:"member"`
	Token     string                             `json:"token"`
}
