package event

import "github.com/pjvds/gouuid"

// ID is a UUID with some serilization and parsing functions added
type ID uuid.UUID

// NewID generates a new random V4 UUID
func NewID() ID {
	guid, err := uuid.NewV4()
	if err != nil {
		panic(err)
	}

	return ID(*guid)
}

// ParseID parses a string into an ID value
func ParseID(value string) (id ID, err error) {
	guid := new(uuid.UUID)
	if guid, err = uuid.ParseHex(value); err == nil {
		id = ID(*guid)
	}

	return
}

// String exports the UUID as a string
func (id ID) String() string {
	guid := uuid.UUID(id)
	return guid.String()
}

// MarshalJSON Marshals the value to JSON
func (id ID) MarshalJSON() ([]byte, error) {
	value := uuid.UUID(id)
	return value.MarshalJSON()
}

// UnmarshalJSON unmarshals the ID from JSON
func (id *ID) UnmarshalJSON(b []byte) error {
	value := uuid.UUID(*id)
	err := value.UnmarshalJSON(b)
	if err != nil {
		return err
	}

	*id = ID(value)
	return nil
}
