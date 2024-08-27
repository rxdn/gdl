package interaction

import (
	"github.com/rxdn/gdl/objects/channel/embed"
	"github.com/rxdn/gdl/objects/interaction/component"
)

type Response struct {
	Type ResponseType `json:"type"`
}

// ========================================================
// Pong Response
// ========================================================
type ResponsePong struct {
	Response
}

func NewResponsePong() ResponsePong {
	return ResponsePong{
		Response: Response{
			Type: ResponseTypePong,
		},
	}
}

// ========================================================
// Channel Message Response
// ========================================================
type ResponseChannelMessage struct {
	Response
	Data ApplicationCommandCallbackData `json:"data"`
}

func NewResponseChannelMessage(data ApplicationCommandCallbackData) ResponseChannelMessage {
	return ResponseChannelMessage{
		Response: Response{
			Type: ResponseTypeChannelMessageWithSource,
		},
		Data: data,
	}
}

// ========================================================
// Ack With Source Response
// ========================================================
type ResponseAckWithSource struct {
	Response
	Data struct {
		Flags uint `json:"flags"`
	} `json:"data"`
}

func NewResponseAckWithSource(flags uint) ResponseAckWithSource {
	return ResponseAckWithSource{
		Response: Response{
			Type: ResponseTypeACKWithSource,
		},
		Data: struct {
			Flags uint `json:"flags"`
		}{
			Flags: flags,
		},
	}
}

// ========================================================
// Deferred Message Update Response
// ========================================================
type ResponseDeferredMessageUpdate struct {
	Response
}

func NewResponseDeferredMessageUpdate() ResponseDeferredMessageUpdate {
	return ResponseDeferredMessageUpdate{
		Response: Response{
			Type: ResponseTypeDeferredMessageUpdate,
		},
	}
}

// ========================================================
// Update message Response
// ========================================================
type ResponseUpdateMessage struct {
	Response
	Data ResponseUpdateMessageData `json:"data"`
}

// TODO: Improve this
type ResponseUpdateMessageData struct {
	Content    *string               `json:"content"`
	Embeds     []*embed.Embed        `json:"embeds"`
	Components []component.Component `json:"components"`
}

func NewResponseUpdateMessage(data ResponseUpdateMessageData) ResponseUpdateMessage {
	return ResponseUpdateMessage{
		Response: Response{
			Type: ResponseTypeUpdateMessage,
		},
		Data: data,
	}
}

// ========================================================
// Auto Complete Result Response
// ========================================================
type ApplicationCommandAutoCompleteResultResponse struct {
	Response
	Data struct {
		Choices []ApplicationCommandOptionChoice `json:"choices"`
	} `json:"data"`
}

func NewApplicationCommandAutoCompleteResultResponse(choices []ApplicationCommandOptionChoice) ApplicationCommandAutoCompleteResultResponse {
	return ApplicationCommandAutoCompleteResultResponse{
		Response: Response{
			Type: ResponseTypeApplicationCommandAutoCompleteResult,
		},
		Data: struct {
			Choices []ApplicationCommandOptionChoice `json:"choices"`
		}{
			Choices: choices,
		},
	}
}

// ========================================================
// Modal Response
// ========================================================
type ModalResponse struct {
	Response
	Data ModalResponseData `json:"data"`
}

type ModalResponseData struct {
	CustomId   string                `json:"custom_id"`
	Title      string                `json:"title"`
	Components []component.Component `json:"components"`
}

func NewModalResponse(customId, title string, components []component.Component) ModalResponse {
	return ModalResponse{
		Response: Response{
			Type: ResponseTypeModal,
		},
		Data: ModalResponseData{
			CustomId:   customId,
			Title:      title,
			Components: components,
		},
	}
}

// ========================================================
// Premium Required Response
// ========================================================
type PremiumRequiredResponse struct {
	Response
	Data struct{} `json:"data"`
}

func NewPremiumRequiredResponse() PremiumRequiredResponse {
	return PremiumRequiredResponse{
		Response: Response{
			Type: ResponseTypePremiumRequired,
		},
		Data: struct{}{},
	}
}
