package component

import (
	"encoding/json"
	"errors"
	"fmt"
)

type ComponentType uint8

const (
	ComponentActionRow ComponentType = iota + 1
	ComponentButton
	ComponentSelectMenu
)

type Component struct {
	Type ComponentType `json:"type"`
	ComponentData
}

type ComponentData interface{
	Type() ComponentType
}

var (
	ErrMissingType  = errors.New("component was missing type field")
	ErrUnknownType  = errors.New("component had unknown type")
	ErrTypeMismatch = errors.New("data did not match component type")
)

func (c Component) MarshalJSON() ([]byte, error) {
	return encode(c.ComponentData)
}

func encode(c ComponentData) (json.RawMessage, error) {
	switch v := c.(type) {
	case ActionRow:
		subComponents := make([]json.RawMessage, len(v.Components))
		for i, sub := range v.Components {
			var err error
			subComponents[i], err = encode(sub.ComponentData)
			if err != nil {
				return nil, err
			}
		}

		data := map[string]interface{}{
			"type":       ComponentActionRow,
			"components": subComponents,
		}
		return json.Marshal(data)
	case Button:
		return json.Marshal(v)
	case SelectMenu:
		return json.Marshal(v)
	default:
		fmt.Println(v)
		return nil, ErrUnknownType
	}
}

func (c *Component) UnmarshalJSON(data []byte) error {
	var raw map[string]interface{}
	if err := json.Unmarshal(data, &raw); err != nil {
		return err
	}

	var componentType ComponentType
	if rawType, ok := raw["type"].(float64); ok {
		componentType = ComponentType(rawType)
	} else {
		return ErrMissingType
	}

	var err error
	switch componentType {
	case ComponentActionRow:
		var parsed ActionRow
		err = json.Unmarshal(data, &parsed)
		c.ComponentData = parsed
	case ComponentButton:
		var parsed Button
		err = json.Unmarshal(data, &parsed)
		c.ComponentData = parsed
	case ComponentSelectMenu:
		var parsed SelectMenu
		err = json.Unmarshal(data, &parsed)
		c.ComponentData = parsed
	default:
		return ErrUnknownType
	}

	if err != nil {
		return err
	}

	c.Type = componentType
	return nil
}
