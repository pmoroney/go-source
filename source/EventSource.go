package source

type EventSource struct {
	ID          EventSourceID
	State       interface{}
	seqID       uint64
	CommandChan chan CommandMessage
	PersistChan chan EventMessage
}

func (source *EventSource) Apply(event Event) {
	event.Apply(source.State)
}

func (source *EventSource) Persist(event Event) {
	source.Apply(event)
	source.seqID += 1
	eventMsg := EventMessage{
		Evt:   event,
		ID:    source.ID,
		SeqID: source.seqID,
	}
	source.PersistChan <- eventMsg
}

func (source *EventSource) Handle(cmd CommandMessage) {
	cmd.Err <- cmd.Cmd.Handle(source)
}

func (source *EventSource) Run() {
	// this start a goroutine that will run forever and just receive and handle commands
	go func() {
		for {
			source.Handle(<-source.CommandChan)
		}
	}()
}
