package events

type EventBus struct {
	Listeners []interface{}
}

func NewEventBus() *EventBus {
	return &EventBus{
	}
}

func (e *EventBus) RegisterListener(fn interface{}) {
	e.Listeners = append(e.Listeners, fn)
}
