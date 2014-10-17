package source

import "github.com/pjvds/gouuid"

type EventSourceID uuid.UUID

func NewEventSourceID() EventSourceID {
	guid, err := uuid.NewV4()
	if err != nil {
		panic(err)
	}

	return EventSourceID(*guid)
}

func ParseEventSourceID(value string) (id EventSourceID, err error) {
	guid := new(uuid.UUID)
	if guid, err = uuid.ParseHex(value); err == nil {
		id = EventSourceID(*guid)
	}

	return
}

func (id EventSourceID) String() string {
	guid := uuid.UUID(id)
	return guid.String()
}

func (id EventSourceID) MarshalJSON() ([]byte, error) {
	value := uuid.UUID(id)
	return value.MarshalJSON()
}

func (id *EventSourceID) UnmarshalJSON(b []byte) error {
	value := uuid.UUID(*id)
	err := value.UnmarshalJSON(b)
	if err != nil {
		return err
	}

	*id = EventSourceID(value)
	return nil
}
