package memory

import (
	"log"

	"github.com/pmoroney/go-source/event"
)

// This is a simple event recorder for testing, it doesnt actually persist the events
type InMemoryEventStore struct {
	events     map[event.ID][]event.Event
	eventChans []chan event.Message
}

func NewInMemoryEventStore() InMemoryEventStore {
	return InMemoryEventStore{
		events: make(map[event.ID][]event.Event, 0),
	}
}

func (r InMemoryEventStore) Record(e event.Message) {
	_, ok := r.events[e.ID]
	if !ok {
		r.events[e.ID] = make([]event.Event, 0)
	}
	r.events[e.ID] = append(r.events[e.ID], e.Data)
	log.Printf("Recorded Event: %+v\n", e)
}

func (r InMemoryEventStore) GetEvents(id event.ID) ([]event.Event, error) {
	_, ok := r.events[id]
	if !ok {
		return nil, nil
	}
	return r.events[id], nil
}

func (r InMemoryEventStore) SubscribeAll(eventChan chan event.Message) {
	if r.eventChans == nil {
		r.eventChans = make([]chan event.Message, 1)
	}
	r.eventChans = append(r.eventChans, eventChan)
}

func (r InMemoryEventStore) Serve() {
}

func (r InMemoryEventStore) Stop() {
}
