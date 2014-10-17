package source

// Commands
type Command interface {
	Handle(*EventSource) error
}

type CommandMessage struct {
	Cmd        Command
	ID         EventSourceID
	Err        chan error
	EmptyState interface{}
}
