package interaction

type InteractionType uint8

const (
	InteractionTypePing InteractionType = iota + 1
	InteractionTypeApplicationCommand
)
