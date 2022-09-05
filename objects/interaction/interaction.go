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

type InteractionMetadata struct {
	Interaction
	Id            uint64                    `json:"id,string"`
	ApplicationId uint64                    `json:"application_id,string"`
	GuildId       objects.NullableSnowflake `json:"guild_id"`
	ChannelId     uint64                    `json:"channel_id,string"`
	Member        *member.Member            `json:"member"`
	User          *user.User                `json:"user"`
	Token         string                    `json:"token"`
}

type InteractionType uint8

const (
	InteractionTypePing InteractionType = iota + 1
	InteractionTypeApplicationCommand
	InteractionTypeMessageComponent
	InteractionTypeApplicationCommandAutoComplete
	InteractionTypeModalSubmit
)

type PingInteraction struct {
	Interaction
}

// If GuildId is not null, Member will be present and User will not.
// If GuildId is null, Member will not be present, and User will.
type ApplicationCommandInteraction struct {
	InteractionMetadata
	Data *ApplicationCommandInteractionData `json:"data"`
}

type ApplicationCommandInteractionData struct {
	Id       uint64                                    `json:"id,string"`
	Name     string                                    `json:"name"`
	Type     ApplicationCommandType                    `json:"type"`
	Resolved ResolvedData                              `json:"resolved"`
	Options  []ApplicationCommandInteractionDataOption `json:"options"`
	TargetId uint64                                    `json:"target_id,string"`
}

type MessageComponentInteraction struct {
	InteractionMetadata
	Data    MessageComponentInteractionData `json:"data"`
	Message message.Message                 `json:"message"`
}

type ApplicationCommandAutoCompleteInteraction struct {
	InteractionMetadata
	Data ApplicationCommandAutoCompleteInteractionData `json:"data"`
}

type ApplicationCommandAutoCompleteInteractionData struct {
	Id      uint64                                    `json:"id,string"`
	Name    string                                    `json:"name"`
	Options []ApplicationCommandInteractionDataOption `json:"options"`
	Type    ApplicationCommandType                    `json:"type"`
}

type ModalSubmitInteraction struct {
	InteractionMetadata
	Data ModalSubmitInteractionData `json:"data"`
}

type ModalSubmitInteractionData struct {
	CustomId   string                                `json:"custom_id"`
	Components []ModalSubmitInteractionActionRowData `json:"components"`
}

type ModalSubmitInteractionActionRowData struct {
	Type       component.ComponentType               `json:"type"`
	Components []ModalSubmitInteractionComponentData `json:"components"`
}

type ModalSubmitInteractionComponentData struct {
	Type     component.ComponentType `json:"type"`
	CustomId string                  `json:"custom_id"`
	Value    string                  `json:"value"`
}
