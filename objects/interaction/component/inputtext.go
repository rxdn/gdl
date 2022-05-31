package component

import (
	"encoding/json"
)

type InputText struct {
	Style       TextStyleTypes `json:"style"`
	CustomId    string         `json:"custom_id"`
	Label       string         `json:"label"`
	Placeholder *string        `json:"placeholder,omitempty"`
	MinLength   *uint32        `json:"min_length,omitempty"`
	MaxLength   *uint32        `json:"max_length,omitempty"`
	Required    *bool          `json:"required,omitempty"`
	Value       *string        `json:"value,omitempty"`
}

type TextStyleTypes uint8

const (
	TextStyleShort TextStyleTypes = iota + 1
	TextStyleParagraph
)

func (i InputText) Type() ComponentType {
	return ComponentInputText
}

func (i InputText) MarshalJSON() ([]byte, error) {
	type WrappedInputText InputText

	return json.Marshal(struct {
		Type ComponentType `json:"type"`
		WrappedInputText
	}{
		Type:             ComponentInputText,
		WrappedInputText: WrappedInputText(i),
	})
}

func BuildInputText(data InputText) Component {
	return Component{
		Type:          ComponentInputText,
		ComponentData: data,
	}
}
