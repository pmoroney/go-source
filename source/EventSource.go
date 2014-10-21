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
