package event

import "github.com/thejerf/suture"

// Router creates agents and routes commands to them
// It also retrieves Events that the Agents have persisted and passes them on to the Store
type Router struct {
	agents     map[ID]chan<- CommandMessage
	store      Store
	persist    chan Message
	supervisor *suture.Supervisor
}

// SetStore creates the connection to the Event Store of choice
func (r *Router) SetStore(store Store) {
	r.store = store
	r.persist = make(chan Message)
}

// Serve starts the router and persists all incoming events to the store
func (r *Router) Serve() {
	r.supervisor = suture.NewSimple("Router")
	r.supervisor.ServeBackground()
	for event := range r.persist {
		r.store.Record(event)
	}
}

// Stop stops the store
func (r *Router) Stop() {
	close(r.persist)
}

// Route takes a command and routes it to the Agent if it is already running.
// If the Agent is not running it creates an Agent and hydrates it with the events from the Store
func (r *Router) Route(cmd CommandMessage) error {
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

func (r *Router) startAgent(cmd CommandMessage) (chan<- CommandMessage, error) {
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

func (r *Router) newAgent(cmd CommandMessage) Agent {
	commands := make(chan CommandMessage)
	ar := Agent{
		id:          cmd.ID,
		state:       cmd.ZeroState,
		commandChan: commands,
		persistChan: r.persist,
	}
	return ar
}

func (r *Router) loadAgent(cmd CommandMessage, events []Event) Agent {
	ar := r.newAgent(cmd)

	for i := range events {
		ar.Apply(events[i])
	}
	return ar
}
