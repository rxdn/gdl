package interaction

type ResponseType uint8

const (
	ResponseTypePong ResponseType = iota + 1
	_
	_
	ResponseTypeChannelMessageWithSource
	ResponseTypeACKWithSource
	ResponseTypeDeferredMessageUpdate
	ResponseTypeUpdateMessage
	ResponseTypeApplicationCommandAutoCompleteResult
)
