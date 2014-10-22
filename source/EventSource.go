package source

import (
	"reflect"
	"time"
)

type eventState interface {
	Apply(event Event)
}

type EventSource struct {
	ID          EventSourceID
	State       eventState
	seqID       uint64
	CommandChan chan CommandMessage
	PersistChan chan EventMessage
}

func (source *EventSource) Apply(event Event) {
	source.State.Apply(event)
}

func (source *EventSource) Persist(event Event) {
	source.Apply(event)
	source.seqID += 1
	eventMsg := EventMessage{
		Data:      event,
		ID:        source.ID,
		SeqID:     source.seqID,
		Timestamp: time.Now(),
		EventType: reflect.TypeOf(event).Name(),
	}
	source.PersistChan <- eventMsg
}

func (source *EventSource) Handle(cmd CommandMessage) {
	err := cmd.Cmd.Handle(source)
	if cmd.Err != nil {
		cmd.Err <- err
	}
}

func (source *EventSource) Serve() {
	// this start a goroutine that will run while the command channel is open
	for c := range source.CommandChan {
		source.Handle(c)
	}
}

func (source *EventSource) Stop() {
	close(source.CommandChan)
}
