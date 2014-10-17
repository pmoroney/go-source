package source

type Router interface {
	Route(CommandMessage) error
	StartStore(EventRecorder)
}

type DefaultRouter struct {
	agents  map[EventSourceID]chan<- CommandMessage
	store   EventRecorder
	persist chan EventMessage
}

func (r *DefaultRouter) StartStore(store EventRecorder) {
	r.store = store
	r.persist = make(chan EventMessage)

	go func() {
		for {
			event := <-r.persist
			r.store.Record(event)
		}
	}()
}

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
	cmd.Err = make(chan error)
	agent <- cmd
	err := <-cmd.Err
	return err
}

func (r *DefaultRouter) startAgent(cmd CommandMessage) (chan<- CommandMessage, error) {
	if r.agents == nil {
		r.agents = make(map[EventSourceID]chan<- CommandMessage, 0)
	}
	events, err := r.store.GetEvents(cmd.ID)
	if err != nil {
		return nil, err
	}
	var ar EventSource
	if events != nil {
		ar = r.loadAgent(cmd, events)
	}
	ar = r.newAgent(cmd)
	ar.Run()
	r.agents[ar.ID] = ar.CommandChan
	return ar.CommandChan, nil
}

func (r *DefaultRouter) newAgent(cmd CommandMessage) EventSource {
	commands := make(chan CommandMessage)
	ar := EventSource{
		ID:          cmd.ID,
		State:       cmd.EmptyState,
		CommandChan: commands,
		PersistChan: r.persist,
	}
	return ar
}

func (r *DefaultRouter) loadAgent(cmd CommandMessage, events []Event) EventSource {
	ar := r.newAgent(cmd)

	for i := range events {
		ar.Apply(events[i])
	}
	return ar
}
