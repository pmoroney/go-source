package source

// Represents an event payload that can be fired by event sources.
// An event should:
// * Name is in past tense.
// * Name contains the intent (CustomerMoved vs CustomerAddressCorrected).
// * Contain all the data related to the event.
type Event interface {
	Apply(interface{})
}

type EventMessage struct {
	Evt   Event
	ID    EventSourceID
	SeqID uint64
}
