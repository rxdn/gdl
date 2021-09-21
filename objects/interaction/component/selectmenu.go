package component

import (
	"encoding/json"
	"github.com/rxdn/gdl/objects/guild/emoji"
)

type SelectMenu struct {
	CustomId    string         `json:"custom_id"`
	Options     []SelectOption `json:"options"`
	Placeholder string         `json:"placeholder,omitempty"`
	MinValues   *int           `json:"min_values,omitempty"`
	MaxValues   *int           `json:"max_values,omitempty"`
	Disabled    bool           `json:"disabled"`
}

type SelectOption struct {
	Label       string       `json:"label"`
	Value       string       `json:"value"`
	Description string       `json:"description,omitempty"`
	Emoji       *emoji.Emoji `json:"emoji,omitempty"`
	Default     bool         `json:"default"`
}

func (s SelectMenu) Type() ComponentType {
	return ComponentSelectMenu
}

func (s SelectMenu) MarshalJSON() ([]byte, error) {
	type WrappedSelectMenu SelectMenu

	return json.Marshal(struct {
		Type ComponentType `json:"type"`
		WrappedSelectMenu
	}{
		Type:              ComponentSelectMenu,
		WrappedSelectMenu: WrappedSelectMenu(s),
	})
}

func BuildSelectMenu(data SelectMenu) Component {
	return Component{
		Type:          ComponentSelectMenu,
		ComponentData: data,
	}
}
