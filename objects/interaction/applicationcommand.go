package interaction

type ApplicationCommand struct {
	Id            uint64                     `json:"id"`
	ApplicationId uint64                     `json:"application_id"`
	Name          string                     `json:"name"`
	Description   string                     `json:"description"`
	Options       []ApplicationCommandOption `json:"options"`
}
