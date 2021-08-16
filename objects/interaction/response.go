package interaction

import (
	"github.com/rxdn/gdl/objects/channel/embed"
	"github.com/rxdn/gdl/objects/interaction/component"
)

type Response struct {
	Type ResponseType `json:"type"`
}

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

type ResponseUpdateMessage struct {
	Response
	Data ResponseUpdateMessageData `json:"data"`
}

// TODO: Improve this
type ResponseUpdateMessageData struct {
	Content    *string               `json:"content,omitempty"`
	Embeds     []*embed.Embed        `json:"embeds,omitempty"`
	Components []component.Component `json:"components,omitempty"`
}

func NewResponseUpdateMessage(data ResponseUpdateMessageData) ResponseUpdateMessage {
	return ResponseUpdateMessage{
		Response: Response{
			Type: ResponseTypeUpdateMessage,
		},
		Data: data,
	}
}
