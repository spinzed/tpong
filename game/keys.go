package game

// This file defines what action should dispatch on every key and
// events that will be dispatched depending on the state of the game

type KeyEvent struct {
	Key
	Event
}

func newKeyEvent(kName string, kState string, eName string, eDesc string) KeyEvent {
	return KeyEvent{newKey(kName, kState), newEvent(eName, eDesc)}
}

type Key struct {
	Name string
	// State describes the key press type the event is triggered, but doesn't define the key.
	// valid states are click, release, hold (and holdEnd) or normal
	// holdEnd state is only used to dismiss the hold state
	State string
}

func newKey(name string, state string) Key {
	return Key{name, state}
}

// Struct event defines an event which can be triggered by a key
type Event struct {
	Name        string
	Description string
}

func newEvent(name string, desc string) Event {
	return Event{name, desc}
}

func getEvent(m *map[Key]Event, k Key) Event {
	// this will be expanded in the future
	return (*m)[k]
}
