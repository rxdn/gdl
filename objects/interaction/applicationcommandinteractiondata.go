package interaction

type ApplicationCommandInteractionData struct {
	Id      uint64                                    `json:"id,string"`
	Name    string                                    `json:"name"`
	Options []ApplicationCommandInteractionDataOption `json:"options"`
}
