package event

// Store is an Event Journal
type Store interface {
	// Records the given event.
	Record(EventMessage) error

	// Gets the recorded events, or an empty slice if none.
	GetEvents(ID) ([]EventMessage, error)

	//	SubscribeAll(chan EventMessage)
}
