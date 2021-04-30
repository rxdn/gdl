package component

import (
	"encoding/json"
	"errors"
)

type ComponentType uint8

const (
	ComponentActionRow ComponentType = iota + 1
	ComponentButton
)

type Component struct {
	Type ComponentType `json:"type"`
	ComponentData
}

type ComponentData interface{}

var (
	ErrMissingType = errors.New("component was missing type field")
	ErrUnknownType = errors.New("component had unknown type")
)

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
	default:
		return ErrUnknownType
	}

	if err != nil {
		return err
	}

	c.Type = componentType
	return nil
}
