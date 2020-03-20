package events

type EventBus struct {
	Events []Event
}

func NewEventBus() EventBus {
	return EventBus{}
}

func (e *EventBus) Register(ev Event) {
	e.Events = append(e.Events, ev)
}

func (e *EventBus) GetEventByName(name string) *Event {
	for _, event := range e.Events {
		if string(event.Type()) == name {
			return &event
		}
	}

	return nil
}
