package source

import (
	"encoding/json"
	. "testing"
)

// Make sure we can turn an EventSourceID into a JSON value
func TestMarshallJSON(t *T) {
	id := NewEventSourceID()

	i := &struct {
		Id    EventSourceID  `json:"id"`
		IdPtr *EventSourceID `json:"idPtr"`
	}{
		Id:    id,
		IdPtr: &id,
	}

	b, err := json.Marshal(i)
	if err != nil {
		t.Fatal(err)
	}

	if string(b) != "{\"id\":\""+i.Id.String()+"\",\"idPtr\":\""+i.IdPtr.String()+"\"}" {
		t.Fatal("ID did not marshal correctly")
	}
}

// Make sure we can turn an JSON value into an EventSourceID
func TestUnMarshallJSON(t *T) {
	i := &struct {
		Id    EventSourceID  `json:"id"`
		IdPtr *EventSourceID `json:"idPtr"`
	}{}

	id := NewEventSourceID()
	data := []byte("{\"id\":\"" + id.String() + "\",\"idPtr\":\"" + id.String() + "\"}")
	err := json.Unmarshal(data, &i)
	if err != nil {
		t.Fatal(err)
	}

	if i.Id != id {
		t.Error("ID did not unmarshal correctly")
	}

	if *i.IdPtr != id {
		t.Error("ID did not unmarshal correctly")
	}
}
