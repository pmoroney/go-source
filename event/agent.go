package event

import (
	"reflect"
	"time"
)

// State is the state of the Aggregate Root.
type State interface {
	// Apply takes an Event and modifies the state accordingly
	Apply(event Event)
}

// Agent is an structure that holds the state of an Aggregate Root, accepts commands, and persits events
type Agent struct {
	id          ID
	state       State
	seqID       uint64
	commandChan chan CommandMessage
	persistChan chan Message
}

// Apply applies an event to the state by calling state.Apply(Event)
func (agent *Agent) Apply(event Event) {
	agent.state.Apply(event)
}

// Persist persists an event
func (agent *Agent) Persist(event Event) {
	agent.Apply(event)
	agent.seqID++
	eventMsg := Message{
		Data:      event,
		ID:        agent.id,
		SeqID:     agent.seqID,
		Timestamp: time.Now(),
		EventType: reflect.TypeOf(event).Name(),
	}
	agent.persistChan <- eventMsg
}

// Handle calls the Handle function on the Command. If ErrChan is not nil it sends the error on that channel.
func (agent *Agent) Handle(cmd CommandMessage) {
	err := cmd.Cmd.Handle(agent)
	if cmd.ErrChan != nil {
		go func() {
			cmd.ErrChan <- err
		}()
	}
}

// Serve starts the agent, it should only be called by the router
func (agent *Agent) Serve() {
	// this start a goroutine that will run while the command channel is open
	for c := range agent.commandChan {
		agent.Handle(c)
	}
}

// Stop stops the agent
func (agent *Agent) Stop() {
	close(agent.commandChan)
}

// ID returns the unique ID of the Aggregate Root
func (agent *Agent) ID() ID {
	return agent.id
}

// State returns a copy of the state so that commands can be validated against the state
func (agent *Agent) State() State {
	return agent.state
}
