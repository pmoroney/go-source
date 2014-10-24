package event

import (
	"encoding/json"
	"fmt"
	"log"
	"reflect"
	"time"
)

// Event Represents an event payload that can be fired by event sources.
// An event should:
// * Name is in past tense.
// * Name contains the intent (CustomerMoved vs CustomerAddressCorrected).
// * Contain all the data related to the event.
type Event interface{}

// Message contains all the metadata related to an event
type Message struct {
	ID        ID
	SeqID     uint64
	Timestamp time.Time
	EventType string
	Data      Event
}

// Serialize serializes an event for storage in a Store.
// It currently uses JSON but this can be changed.
// A pluggable serializer would be nice as well
func (e Message) Serialize() ([]byte, error) {
	return json.Marshal(e)
}

// Unserialize unserializes an event from the Store.
// All types that are going to be unserialezed need to be registered by Registrar() or RegistrarType()
func Unserialize(data []byte) (*Message, error) {
	raw := new(struct {
		ID        ID
		SeqID     uint64
		Timestamp time.Time
		EventType string
		Data      json.RawMessage
	})

	err := json.Unmarshal(data, raw)
	if err != nil {
		return nil, err
	}

	e := new(Message)
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
	log.Printf("Event raw: %#v\n", event)
	if err := json.Unmarshal(raw.Data, &event); err != nil {
		return nil, err
	}

	e.Data = reflect.Indirect(reflect.ValueOf(event)).Interface()

	return e, nil
}
