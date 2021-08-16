package interaction

import (
	"github.com/rxdn/gdl/objects"
	"github.com/rxdn/gdl/objects/channel/message"
	"github.com/rxdn/gdl/objects/interaction/component"
	"github.com/rxdn/gdl/objects/member"
	"github.com/rxdn/gdl/objects/user"
)

type Interaction struct {
	Version uint8           `json:"version"`
	Type    InteractionType `json:"type"`
}

type InteractionType uint8

const (
	InteractionTypePing InteractionType = iota + 1
	InteractionTypeApplicationCommand
	InteractionTypeButton
)

type PingInteraction struct {
	Interaction
}

// If GuildId is not null, Member will be present and User will not.
// If GuildId is null, Member will not be present, and User will.
type ApplicationCommandInteraction struct {
	Interaction
	Id            uint64                             `json:"id,string"`
	ApplicationId uint64                             `json:"application_id,string"`
	Data          *ApplicationCommandInteractionData `json:"data"`
	GuildId       objects.NullableSnowflake          `json:"guild_id"`
	ChannelId     uint64                             `json:"channel_id,string"`
	Member        *member.Member                     `json:"member"`
	User          *user.User                         `json:"user"`
	Token         string                             `json:"token"`
}

type ApplicationCommandInteractionData struct {
	Id       uint64                                    `json:"id,string"`
	Name     string                                    `json:"name"`
	Options  []ApplicationCommandInteractionDataOption `json:"options"`
	TargetId uint64                                    `json:"target_id"`
	Type     ApplicationCommandType                    `json:"type"`
}

type ButtonInteraction struct {
	Id            uint64                    `json:"id,string"`
	ApplicationId uint64                    `json:"application_id,string"`
	Data          ButtonInteractionData     `json:"data"`
	GuildId       objects.NullableSnowflake `json:"guild_id"`
	ChannelId     uint64                    `json:"channel_id,string"`
	Message       message.Message           `json:"message"`
	Member        *member.Member            `json:"member"`
	User          *user.User                `json:"user"`
	Token         string                    `json:"token"`
}

type ButtonInteractionData struct {
	CustomId string                  `json:"custom_id"`
	Type     component.ComponentType `json:"component_type"`
}
