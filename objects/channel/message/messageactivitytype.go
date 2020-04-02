package message

type MessageActivityType int

const (
	MessageActivityJoin MessageActivityType = iota
	MessageActivitySpectate
	MessageActivityListen
	MessageActivityJoinRequest
)
