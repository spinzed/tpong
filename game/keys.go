package game

// This file defines what action should dispatch on every key and
// events that will be dispatched depending on the state of the game

// Struct event defines an event which can be triggered by a key
type Event struct {
	Name        string
	Description string
	// State describes the key press type the event is triggered, but doesn't define the key.
	// valid states are click, release, hold or normal
	State string
}

func newEvent(name string, desc string, state string) Event {
	return Event{name, desc, state}
}

// Keys which can be triggered when game is started
var keysStart = map[string]Event{
	"SPACE": newEvent(eventStart, "Start Game", stateClick),
	"Up":    newEvent(eventP2Up, "Move Player2 Up", stateHold),
	"Down":  newEvent(eventP2Down, "Move Player2 Down", stateHold),
	"Q":     newEvent(eventDestroy, "Quit", stateClick),
	"P":     newEvent(eventTogglePause, "Toggle Pause", stateClick),
	"R":     newEvent(eventReset, "Reset Round", stateClick),
	"T":     newEvent(eventSwitchTheme, "Switch Theme", stateClick),
}

// Keys which can be triggered when game is started
var keysGame = map[string]Event{
	"W":     newEvent(eventP1Up, "Move Player1 Up", stateHold),
	"S":     newEvent(eventP1Down, "Move Player1 Down", stateHold),
	"Up":    newEvent(eventP2Up, "Move Player2 Up", stateHold),
	"Down":  newEvent(eventP2Down, "Move Player2 Down", stateHold),
	"SPACE": newEvent(eventStart, "Start Game", stateClick),
	"Q":     newEvent(eventDestroy, "Quit", stateClick),
	"P":     newEvent(eventTogglePause, "Toggle Pause", stateClick),
	"R":     newEvent(eventReset, "Reset Round", stateClick),
	"T":     newEvent(eventSwitchTheme, "Switch Theme", stateClick),
}

func getEvent(m *map[string]Event, k string) Event {
	// this will be expanded in the future
	return (*m)[k]
}

func getEventByName(m *map[string]Event, name string) (*Event, string) {
	for key, event := range *m {
		if event.Name == name {
			return &event, key
		}
	}
	return &Event{}, ""
}

// Actions that are dispatched when the player is in start menu
func (g *Game) dispatchStartAction(e KeyDispatch) {
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
func (g *Game) dispatchGameAction(e KeyDispatch) {
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
