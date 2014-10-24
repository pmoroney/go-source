package event

import (
	"reflect"
)

// Registry holds the mapping between an event name and it's static type.
// All event type should be registered at bootstrap time so that an event store,
// bus or other services can deserialize messages to concrete types.
var registry map[string]reflect.Type

// RegisterType registers an event type. An existing entry with the same name is overwritten
// if it exists. It will register the type of the element, even if you provide
// a pointer type. For example, *FooBar will be registered as FooBar.
func RegisterType(t reflect.Type) {
	for t.Kind() == reflect.Ptr {
		t = t.Elem()
	}

	if registry == nil {
		registry = make(map[string]reflect.Type, 0)
	}

	n := t.Name()

	registry[n] = t
}

// Register Registers an event type by instance. An existing entry with the same name is
// overwritten if it exists. It will register the type of the element, even if
// you provide a pointer type. For example, *FooBar will be registered as FooBar.
func Register(e Event) {
	t := reflect.TypeOf(e)
	RegisterType(t)
}

// getType gets the static type from an event name. It results `true` for `ok` if
// the type was found; otherwise, `false`.
func getType(n string) (reflect.Type, bool) {
	if registry == nil {
		return nil, false
	}
	r, ok := registry[n]
	return r, ok
}
