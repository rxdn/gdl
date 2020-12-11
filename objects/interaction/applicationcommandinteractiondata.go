package interaction

type ApplicationCommandInteractionData struct {
	Id      uint64                                    `json:"id"`
	Name    string                                    `json:"name"`
	Options []ApplicationCommandInteractionDataOption `json:"options"`
}
