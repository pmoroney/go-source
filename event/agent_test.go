package event

import (
	"testing"
)

// Mock State
type agentTestState struct {
	applied bool
}

func (s *agentTestState) Apply(e Event) {
	switch m := e.(type) {
	case string:
		if m == "fake event" {
			s.applied = true
		}
	}
}

func TestAgentApply(t *testing.T) {
	state := agentTestState{}
	agent := Agent{
		seqID: 123,
		state: &state,
	}
	eventMsg := EventMessage{
		Data:  "fake event",
		SeqID: 124,
	}
	agent.Apply(eventMsg)

	if !state.applied {
		t.Fatal("Event did not modify original state")
	}

	if agent.seqID == 123 {
		t.Error("Event did not change the SeqID")
	}

	if agent.seqID != 124 {
		t.Error("Event did not change the SeqID to the right value")
	}
}

func TestAgentPersist(t *testing.T) {
	id, err := ParseID("f68884bf-4961-40b0-6b1d-4a1fe4e73269")
	if err != nil {
		t.Fatal(err)
	}
	state := agentTestState{}
	store := NewInMemoryEventStore()
	router := Router{
		store: store,
	}
	agent := Agent{
		id:     id,
		seqID:  123,
		state:  &state,
		router: &router,
	}
	err = agent.Persist("fake event")
	if err != nil {
		t.Fatal(err)
	}

	if !state.applied {
		t.Fatal("Event did not modify original state")
	}

	events, err := store.GetEvents(id)
	if err != nil {
		t.Fatal(err)
	}

	if len(events) != 1 {
		t.Fatal("Persist did not store the event in the store")
	}

	event := events[0]

	if event.ID != id {
		t.Error("EventMessage does not have the correct ID embeded into it")
	}

	if event.SeqID != 124 {
		t.Error("EventMessage does not have the correct SeqID embeded into it")
	}

	if event.EventType != "string" {
		t.Error("EventMessage does not have the correct EventType embeded into it")
	}

	if event.Data != "fake event" {
		t.Error("EventMessage does not have the correct Data embeded into it")
	}
}
