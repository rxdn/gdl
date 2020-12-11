package interaction

type ApplicationCommand struct {
	Id            uint64                     `json:"id,omitempty"`
	ApplicationId uint64                     `json:"application_id,omitempty"`
	Name          string                     `json:"name"`
	Description   string                     `json:"description"`
	Options       []ApplicationCommandOption `json:"options"`
}
