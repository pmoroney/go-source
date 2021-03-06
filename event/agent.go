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
	router      Router
}

// Apply applies an event to the state by calling state.Apply(Event)
func (a *Agent) Apply(event EventMessage) {
	a.state.Apply(event.Data)
	a.seqID++
}

// Persist persists an event
func (a *Agent) Persist(event Event) error {
	eventMsg := EventMessage{
		Data:      event,
		ID:        a.id,
		SeqID:     a.seqID + 1,
		Timestamp: time.Now(),
		EventType: reflect.TypeOf(event).Name(),
	}

	err := a.router.Record(eventMsg)
	if err != nil {
		return err
	}

	a.Apply(eventMsg)
	return nil
}

// Handle calls the Handle function on the Command. If ErrChan is not nil it sends the error on that channel.
func (a *Agent) Handle(cmd CommandMessage) {
	err := cmd.Cmd.Handle(a)
	if cmd.ErrChan != nil {
		go func() {
			cmd.ErrChan <- err
		}()
	}
}

// Serve starts the agent, it should only be called by the router
func (a *Agent) Serve() {
	//a.timer = *time.NewTimer(a.router.SnapshotInterval())
	// this start a goroutine that will run while the command channel is open
	for {
		select {
		case c, ok := <-a.commandChan:
			if !ok {
				return
			}
			a.Handle(c)

			/*
					// Snapshotting should be scheduled by the router possibly...
				case <-a.timer.C:
					a.takeSnapshot()
					a.timer.Reset(a.router.SnapshotInterval())
			*/
		}
	}
}

func (a *Agent) takeSnapshot() {
	// @TODO implement this :)
}

// Stop stops the agent
func (a *Agent) Stop() {
	close(a.commandChan)
}

// ID returns the unique ID of the Aggregate Root
func (a *Agent) ID() ID {
	return a.id
}

// State returns a copy of the state so that commands can be validated against the state
func (a *Agent) State() State {
	return a.state
}

func NewAgent(state State, router Router) Agent {
	return Agent{
		state:  state,
		router: router,
	}
}
