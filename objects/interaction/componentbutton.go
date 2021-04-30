package interaction

import (
	"encoding/json"
	"github.com/rxdn/gdl/objects/guild/emoji"
)

type Button struct {
	Label    string      `json:"label"`
	CustomId string      `json:"custom_id"`
	Style    ButtonStyle `json:"style"`
	Emoji    emoji.Emoji `json:"emoji"`
	Url      *string     `json:"url,omitempty"`
	Disabled bool        `json:"false"`
}

func (b *Button) Type() ComponentType {
	return ComponentButton
}

type ButtonStyle uint8

const (
	ButtonStylePrimary ButtonStyle = iota + 1
	ButtonStyleSecondary
	ButtonStyleSuccess
	ButtonStyleDanger
	ButtonStyleLink
)

func (b Button) MarshalJSON() ([]byte, error) {
	type WrappedButton Button

	return json.Marshal(struct {
		Type ComponentType `json:"type"`
		WrappedButton
	}{
		Type:   ComponentButton,
		WrappedButton: WrappedButton(b),
	})
}
