package game

// This file defines events that will be dispatched depending on
// the key press and it's type.

// Struct event defines an event which can be triggered by a key
type Event struct {
	Name string
	Description string
}

func newEvent(name string, desc string) Event {
	return Event{name, desc}
}

// These events are fired while the key is held down.
var events = map[string]Event{
	"W":    newEvent(eventP1Up, "Move Player1 Up"),
	"S":    newEvent(eventP1Down, "Move Player1 Down"),
	"Up":   newEvent(eventP2Up, "Move Player2 Up"),
	"Down": newEvent(eventP2Down, "Move Player2 Down"),
}

// For these events, the event is immediately dispatched and
// which results in immediate action call.
var dispEvents = map[string]Event{
	"SPACE": newEvent(eventStart, "Start Game"),
	"Q":     newEvent(eventDestroy, "Quit"),
	"P":     newEvent(eventTogglePause, "Toggle Pause"),
	"R":     newEvent(eventReset, "Reset Round"),
	"T":     newEvent(eventSwitchTheme, "Switch Theme"),
	//		"B":     eventToggleBg,
}

func getEventByName(name string) (*Event, string) {
	for key, event := range dispEvents {
		if event.Name == name {
			return &event, key
		}
	}
	for key, event := range events {
		if event.Name == name {
			return &event, key
		}
	}
	return &Event{}, ""
}
