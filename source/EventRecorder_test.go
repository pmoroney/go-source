package source

import (
	. "testing"
)

type TestEvent struct {
	Foo string
}

// Make sure we can record events
func TestRecord(t *T) {
	Register(TestEvent{})
	event := EventMessage{
		Data: TestEvent{
			Foo: "bar",
		},
		ID:    NewEventSourceID(),
		SeqID: 1,
	}

	recoder := NewInMemoryEventRecorder()
	recoder.Record(event)

	events, err := recoder.GetEvents(event.ID)
	if err != nil {
		t.Fatal(err)
	}

	if len(events) != 1 {
		t.Fatal("Events Length != 1")
	}
}

// Make sure we can record events, even if they are the same
func TestRecordSameEventTwice(t *T) {
	Register(TestEvent{})
	event := EventMessage{
		Data: TestEvent{
			Foo: "bar",
		},
		ID:    NewEventSourceID(),
		SeqID: 1,
	}

	recoder := NewInMemoryEventRecorder()
	recoder.Record(event)
	event.SeqID = 2
	recoder.Record(event)

	events, err := recoder.GetEvents(event.ID)
	if err != nil {
		t.Fatal(err)
	}

	if len(events) != 2 {
		t.Fatal("Events Length != 2")
	}
}
