package component

import "encoding/json"

type ActionRow struct {
	Components []Component
}

func (a *ActionRow) Type() ComponentType {
	return ComponentActionRow
}

func (a ActionRow) MarshalJSON() ([]byte, error) {
	type WrappedActionRow ActionRow

	return json.Marshal(struct {
		Type ComponentType `json:"type"`
		WrappedActionRow
	}{
		Type:             ComponentButton,
		WrappedActionRow: WrappedActionRow(a),
	})
}

