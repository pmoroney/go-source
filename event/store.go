package event

// Store is an Event Journal
type Store interface {
	// Records the given event.
	Record(Message)

	// Gets the recorded events, or an empty slice if none.
	GetEvents(ID) ([]Event, error)

	//	SubscribeAll(chan EventMessage)
}
