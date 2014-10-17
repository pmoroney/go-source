package source

type EventRecorder interface {
	// Records the given event.
	Record(EventMessage)

	// Gets the recorded events, or an empty slice if none.
	GetEvents(EventSourceID) ([]Event, error)
}

// This is a simple event recorder for testing, it doesnt actually persist the events
type InMemoryEventRecorder struct {
	events map[EventSourceID][]Event
}

func NewInMemoryEventRecorder() InMemoryEventRecorder {
	return InMemoryEventRecorder{
		events: make(map[EventSourceID][]Event, 0),
	}
}

func (r InMemoryEventRecorder) Record(e EventMessage) {
	_, ok := r.events[e.ID]
	if !ok {
		r.events[e.ID] = make([]Event, 0)
	}
	r.events[e.ID] = append(r.events[e.ID], e.Evt)
}

func (r InMemoryEventRecorder) GetEvents(id EventSourceID) ([]Event, error) {
	_, ok := r.events[id]
	if !ok {
		return nil, nil
	}
	return r.events[id], nil
}
