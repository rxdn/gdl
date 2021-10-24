package interaction

import (
	"encoding/json"
	"github.com/rxdn/gdl/objects/interaction/component"
)

type MessageComponentInteractionData struct {
	MessageComponentInteractionBaseData
	IMessageComponentInteractionData
}

func (d MessageComponentInteractionData) AsButton() ButtonInteractionData  {
	return d.IMessageComponentInteractionData.(ButtonInteractionData)
}

func (d MessageComponentInteractionData) AsSelectMenu() SelectMenuInteractionData  {
	return d.IMessageComponentInteractionData.(SelectMenuInteractionData)
}

type IMessageComponentInteractionData interface {
	Type() component.ComponentType
}

type MessageComponentInteractionBaseData struct {
	ComponentType component.ComponentType `json:"component_type"`
	CustomId      string                  `json:"custom_id"`
}

type ButtonInteractionData struct {
	MessageComponentInteractionBaseData
}

func (d ButtonInteractionData) Type() component.ComponentType {
	return component.ComponentButton
}

type SelectMenuInteractionData struct {
	MessageComponentInteractionBaseData
	Values []string `json:"values"`
}

func (d SelectMenuInteractionData) Type() component.ComponentType {
	return component.ComponentSelectMenu
}

func (d *MessageComponentInteractionData) UnmarshalJSON(data []byte) error {
	var raw map[string]interface{}
	if err := json.Unmarshal(data, &raw); err != nil {
		return err
	}

	var componentType component.ComponentType
	if rawType, ok := raw["component_type"].(float64); ok {
		componentType = component.ComponentType(rawType)
	} else {
		return component.ErrMissingType
	}

	var err error
	switch componentType {
	case component.ComponentActionRow:
		return component.ErrUnknownType
	case component.ComponentButton:
		var parsed ButtonInteractionData
		err = json.Unmarshal(data, &parsed)
		d.IMessageComponentInteractionData = parsed
	case component.ComponentSelectMenu:
		var parsed SelectMenuInteractionData
		err = json.Unmarshal(data, &parsed)
		d.IMessageComponentInteractionData = parsed
	default:
		return component.ErrUnknownType
	}

	if err != nil {
		return err
	}

	d.ComponentType = componentType
	return nil
}
