package memory

import (
	. "testing"

	"github.com/pmoroney/go-source/event"
)

type TestEvent struct {
	Foo string
}

// Make sure we can record events
func TestRecord(t *T) {
	event.Register(TestEvent{})
	event := event.EventMessage{
		Data: TestEvent{
			Foo: "bar",
		},
		ID:    event.NewID(),
		SeqID: 1,
	}

	recoder := NewInMemoryEventStore()
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
	event.Register(TestEvent{})
	event := event.EventMessage{
		Data: TestEvent{
			Foo: "bar",
		},
		ID:    event.NewID(),
		SeqID: 1,
	}

	recoder := NewInMemoryEventStore()
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
