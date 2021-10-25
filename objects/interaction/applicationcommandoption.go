package interaction

import "github.com/rxdn/gdl/objects/channel"

type ApplicationCommandInteractionDataOption struct {
	Name    string                                    `json:"name"`
	Value   interface{}                               `json:"value,omitempty"`
	Options []ApplicationCommandInteractionDataOption `json:"options,omitempty"`
	Focused bool                                      `json:"focused"`
}

type ApplicationCommandOption struct {
	Type         ApplicationCommandOptionType     `json:"type"`
	Name         string                           `json:"name"`
	Description  string                           `json:"description"`
	Default      bool                             `json:"default"`
	Required     bool                             `json:"required"`
	Choices      []ApplicationCommandOptionChoice `json:"choices,omitempty"`
	Autocomplete bool                             `json:"autocomplete"`
	Options      []ApplicationCommandOption       `json:"options,omitempty"`
	ChannelTypes []channel.ChannelType            `json:"channel_types,omitempty"`
}

type ApplicationCommandOptionType uint8

const (
	OptionTypeSubCommand ApplicationCommandOptionType = iota + 1
	OptionTypeSubCommandGroup
	OptionTypeString
	OptionTypeInteger
	OptionTypeBoolean
	OptionTypeUser
	OptionTypeChannel
	OptionTypeRole
)

type ApplicationCommandOptionChoice struct {
	Name  string      `json:"name"`
	Value interface{} `json:"value"` // string or int
}
