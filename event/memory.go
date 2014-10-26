package event

// This is a simple event recorder for testing, it doesnt actually persist the events
type InMemoryEventStore struct {
	events     map[ID][]EventMessage
	eventChans []chan EventMessage
}

func NewInMemoryEventStore() InMemoryEventStore {
	return InMemoryEventStore{
		events: make(map[ID][]EventMessage, 0),
	}
}

func (r InMemoryEventStore) Record(e EventMessage) error {
	_, ok := r.events[e.ID]
	if !ok {
		r.events[e.ID] = make([]EventMessage, 0)
	}
	r.events[e.ID] = append(r.events[e.ID], e)
	return nil
}

func (r InMemoryEventStore) GetEvents(id ID) ([]EventMessage, error) {
	_, ok := r.events[id]
	if !ok {
		return nil, nil
	}
	return r.events[id], nil
}

func (r InMemoryEventStore) SubscribeAll(eventChan chan EventMessage) {
	if r.eventChans == nil {
		r.eventChans = make([]chan EventMessage, 1)
	}
	r.eventChans = append(r.eventChans, eventChan)
}

func (r InMemoryEventStore) Serve() {
}

func (r InMemoryEventStore) Stop() {
}
