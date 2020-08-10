package game

// This file defines what action should dispatch on every key and
// events that will be dispatched depending on the state of the game

// Struct event defines an event which can be triggered by a key
type Event struct {
	Name        string
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

// Actions that are dispatched when the player is in start menu
func (g *Game) dispatchStartAction(e keyState) {
	switch e.Name {
	case eventDestroy:
		g.EndLoop()
	case eventStart:
		g.Start()
	case eventReset:
		g.Reset()
	case eventSwitchTheme:
		g.switchTheme()
	}
}

// Actions that are dispatched when game is started
func (g *Game) dispatchGameAction(e keyState) {
	switch e.Name {
	case eventDestroy:
		g.EndLoop()
	case eventStart:
		g.Start()
	case eventTogglePause:
		g.togglePause()
	case eventReset:
		g.Reset()
	case eventSwitchTheme:
		g.switchTheme()
	}
}
