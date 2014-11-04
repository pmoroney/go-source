package event

// Command is a struct that contains parameters for a command on an Aggregate Root
type Command interface {
	// Handle validates the command and returns an error or nil if the command is successful
	// Handle should also create and persist any events in order to modify the state
	Handle(*Agent) error
}

// CommandMessage has the metadata needed for a command to be routed
type CommandMessage struct {
	Cmd     Command
	ID      ID
	ErrChan chan error
}
