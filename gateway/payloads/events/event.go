package events

type Event interface {
	Type() EventType
	New() interface{}
	Handle(interface{})
}
