package interaction

type InteractionResponse struct {
	Type InteractionResponseType                   `json:"type"`
	Data InteractionApplicationCommandCallbackData `json:"data,omitempty"`
}
