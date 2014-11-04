package event

import (
	"time"

	"github.com/thejerf/suture"
)

type Router interface {
	Route(CommandMessage) error
	Record(EventMessage) error
}

// Router creates agents and routes commands to them
// It also retrieves Events that the Agents have persisted and passes them on to the Store
type DefaultRouter struct {
	agents           map[ID]chan<- CommandMessage
	store            Store
	supervisor       *suture.Supervisor
	snapshotInterval time.Duration
	state            State
}

// SetStore creates the connection to the Event Store of choice
func (r *DefaultRouter) SetStore(store Store) {
	r.store = store
}

func (r *DefaultRouter) SetState(state State) {
	r.state = state
}

// Serve starts the router and persists all incoming events to the store
func (r *DefaultRouter) Serve() {
	r.snapshotInterval = 2 * time.Minute
	r.supervisor = suture.NewSimple("Router")
	r.supervisor.Serve()
}

func (r *DefaultRouter) Record(event EventMessage) error {
	return r.store.Record(event)
}

// Stop stops the store
func (r *DefaultRouter) Stop() {
	r.supervisor.Stop()
}

func (r *DefaultRouter) SendCommand(cmd CommandMessage, fail *bool) error {
	cmd.ErrChan = make(chan error)
	err := r.Route(cmd)
	if err != nil {
		*fail = true
		return err
	}

	return <-cmd.ErrChan
}

// Route takes a command and routes it to the Agent if it is already running.
// If the Agent is not running it creates an Agent and hydrates it with the events from the Store
func (r *DefaultRouter) Route(cmd CommandMessage) error {
	var agent chan<- CommandMessage
	var ok bool
	agent, ok = r.agents[cmd.ID]
	if !ok {
		// make a new agent or load from the store
		var err error
		agent, err = r.startAgent(cmd)
		if err != nil {
			return err
		}
	}
	agent <- cmd
	return nil
}

func (r *DefaultRouter) startAgent(cmd CommandMessage) (chan<- CommandMessage, error) {
	if r.agents == nil {
		r.agents = make(map[ID]chan<- CommandMessage, 0)
	}

	events, err := r.store.GetEvents(cmd.ID)
	if err != nil {
		return nil, err
	}

	var ar Agent
	if events != nil {
		ar = r.loadAgent(cmd, events)
	}

	ar = r.newAgent(cmd)

	r.supervisor.Add(&ar)

	r.agents[ar.id] = ar.commandChan
	return ar.commandChan, nil
}

func (r *DefaultRouter) newAgent(cmd CommandMessage) Agent {
	commands := make(chan CommandMessage)
	ar := Agent{
		id:          cmd.ID,
		state:       r.state,
		commandChan: commands,
		router:      r,
	}
	return ar
}

func (r *DefaultRouter) loadAgent(cmd CommandMessage, events []EventMessage) Agent {
	ar := r.newAgent(cmd)

	for i := range events {
		ar.Apply(events[i])
	}
	return ar
}

func (r DefaultRouter) SnapshotInterval() time.Duration {
	return r.snapshotInterval
}
