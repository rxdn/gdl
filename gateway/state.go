package gateway

type State int

const(
	CONNECTED State = iota
	CONNECTING
	DISCONNECTING
	DEAD
)
