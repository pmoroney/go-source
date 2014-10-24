package event

import (
	"encoding/json"
	. "testing"
)

// Make sure we can turn an EventSourceID into a JSON value
func TestMarshallJSON(t *T) {
	id := NewID()

	i := &struct {
		ID    ID  `json:"id"`
		IDPtr *ID `json:"idPtr"`
	}{
		ID:    id,
		IDPtr: &id,
	}

	b, err := json.Marshal(i)
	if err != nil {
		t.Fatal(err)
	}

	if string(b) != "{\"id\":\""+i.ID.String()+"\",\"idPtr\":\""+i.IDPtr.String()+"\"}" {
		t.Fatal("ID did not marshal correctly")
	}
}

// Make sure we can turn an JSON value into an EventSourceID
func TestUnMarshallJSON(t *T) {
	i := &struct {
		ID    ID  `json:"id"`
		IDPtr *ID `json:"idPtr"`
	}{}

	id := NewID()
	data := []byte("{\"id\":\"" + id.String() + "\",\"idPtr\":\"" + id.String() + "\"}")
	err := json.Unmarshal(data, &i)
	if err != nil {
		t.Fatal(err)
	}

	if i.ID != id {
		t.Error("ID did not unmarshal correctly")
	}

	if *i.IDPtr != id {
		t.Error("ID did not unmarshal correctly")
	}
}
