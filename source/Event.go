package source

import (
	"encoding/json"
	"fmt"
	"reflect"
	"time"
)

// Represents an event payload that can be fired by event sources.
// An event should:
// * Name is in past tense.
// * Name contains the intent (CustomerMoved vs CustomerAddressCorrected).
// * Contain all the data related to the event.
type Event interface{}

type EventMessage struct {
	ID        EventSourceID
	SeqID     uint64
	Timestamp time.Time
	EventType string
	Data      Event
}

func (e EventMessage) Serialize() ([]byte, error) {
	return json.Marshal(e)
}

func Unserialize(data []byte) (*EventMessage, error) {
	raw := new(struct {
		ID        EventSourceID
		SeqID     uint64
		EventType string
		Timestamp time.Time
		Data      json.RawMessage
	})

	err := json.Unmarshal(data, raw)
	if err != nil {
		return nil, err
	}

	e := new(EventMessage)
	e.ID = raw.ID
	e.SeqID = raw.SeqID
	e.Timestamp = raw.Timestamp
	e.EventType = raw.EventType

	eventType, ok := getType(e.EventType)
	if !ok {
		return nil, fmt.Errorf("No known type for '%v', register it first", e.EventType)
	}

	eventValue := reflect.New(eventType)
	event := eventValue.Interface()
	if err := json.Unmarshal(raw.Data, &event); err != nil {
		return nil, err
	}

	e.Data = reflect.Indirect(reflect.ValueOf(event)).Interface()

	return e, nil
}
