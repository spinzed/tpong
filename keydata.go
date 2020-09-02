package main

import "strings"

// ScreenKeyData struct contains info about key what events are
// dispatchable on the current screen.
// Keys is not of type KeyEvent because it is needed for the events
// to be quickly accessed via the key.
// AltKeys are used for special cases, per example, when the game is
// paused. Thus, it is optional (nil can be passed)
type ScreenKeyData struct {
	Keys    *map[Key]Event
	AltKeys *map[Key]Event
	Legend  *Legend
}

type Legend struct {
	// type can be either left or middle
	Type     string
	Keys     *[]KeyEvent
	Selected int
}

func newScreenKeyData(k *map[Key]Event, al *map[Key]Event, lk *[]KeyEvent, lt string) *ScreenKeyData {
	legend := &Legend{lt, lk, 0}
	return &ScreenKeyData{k, al, legend}
}

var screens = map[string]*ScreenKeyData{
	screenStartMenu: newScreenKeyData(&keysStart, nil, &legendKeysStart, "middle"),
	screenGame:      newScreenKeyData(&keysGame, &altKeysGame, &legendKeysGame, "left"),
}

// Fetches the keys and event descriptions and formats them.
// Mode dictates should keys be format for left or middle.
// May be moved to another file in the future.
func (s *ScreenKeyData) formatKeys() []string {
	mode := s.Legend.Type
	keys := s.Legend.Keys
	selected := s.Legend.Selected

	if mode != "middle" && mode != "left" {
		panic("Invalid format mode:" + mode)
	}

	var maxKeylen int

	// switch default key names with these
	alternate := map[string]string{
		"Up":   "ArrowUp",
		"Down": "ArrowDown",
	}

	// cycle all event names and get their keys and descs
	for _, keyev := range *keys {
		key := keyev.Key.Name
		realKey := key

		if alternate[key] != "" {
			realKey = alternate[key]
		}

		if len(key) > maxKeylen {
			maxKeylen = len(realKey)
		}
	}

	var final []string

	// format every line and append it to the final line list
	for i, keyevt := range *keys {
		key := keyevt.Key.Name
		var thstring string

		if mode == "left" {
			thstring += strings.Repeat(" ", (maxKeylen-len(key))/2-1)
			if i == selected {
				thstring += ">> "
			} else {
				thstring += "   "
			}
			thstring += key
			thstring += strings.Repeat(" ", maxKeylen-len(key)-((maxKeylen-len(key))/2)+1)
			thstring += "- " + keyevt.Event.Description + " "
		} else {
			thstring += "  "
			if i == selected {
				thstring += ">> "
			}
			thstring += key + " - " + keyevt.Event.Description
			if i == selected {
				thstring += " <<"
			}
			thstring += " "
			if len(thstring)%2 == 1 {
				thstring += " "
			}
		}

		final = append(final, thstring)
	}

	return final
}

// Keys which can be triggered when game is started
var keysStart = map[Key]Event{
	newKey("SPACE", stateClick): newEvent(eventStart, "Start Game"),
	newKey("ENTER", stateClick): newEvent(eventMenuSelect, "Select Menu Action"),
	newKey("Up", stateNormal):   newEvent(eventMenuUp, "Select Action Above"),
	newKey("Down", stateNormal): newEvent(eventMenuDown, "Select Action Below"),
	newKey("Q", stateClick):     newEvent(eventDestroy, "Quit"),
	newKey("P", stateClick):     newEvent(eventTogglePause, "Toggle Pause"),
	newKey("R", stateClick):     newEvent(eventReset, "Reset Round"),
	newKey("T", stateClick):     newEvent(eventSwitchTheme, "Switch Theme"),
	newKey("A", stateClick):     newEvent(eventToggleAI, "Toggle AI"),
}

// Keys which can be triggered when game is started
var keysGame = map[Key]Event{
	newKey("W", stateHold):    newEvent(eventP1Up, "Move Player1 Up"),
	newKey("S", stateHold):    newEvent(eventP1Down, "Move Player1 Down"),
	newKey("Up", stateHold):   newEvent(eventP2Up, "Move Player2 Up"),
	newKey("Down", stateHold): newEvent(eventP2Down, "Move Player2 Down"),
	newKey("Q", stateClick):   newEvent(eventDestroy, "Quit"),
	newKey("P", stateClick):   newEvent(eventTogglePause, "Toggle Pause"),
	newKey("R", stateClick):   newEvent(eventReset, "Reset Round"),
	newKey("T", stateClick):   newEvent(eventSwitchTheme, "Switch Theme"),
}

var altKeysGame = map[Key]Event{
	newKey("W", stateHold):      newEvent(eventP1Up, "Move Player1 Up"),
	newKey("S", stateHold):      newEvent(eventP1Down, "Move Player1 Down"),
	newKey("ENTER", stateClick): newEvent(eventMenuSelect, "Select Menu Action"),
	newKey("Up", stateNormal):   newEvent(eventMenuUp, "Select Action Above"),
	newKey("Down", stateNormal): newEvent(eventMenuDown, "Select Action Below"),
	newKey("Q", stateClick):     newEvent(eventDestroy, "Quit"),
	newKey("P", stateClick):     newEvent(eventTogglePause, "Toggle Pause"),
	newKey("R", stateClick):     newEvent(eventReset, "Reset Round"),
	newKey("T", stateClick):     newEvent(eventSwitchTheme, "Switch Theme"),
}

var legendKeysStart = []KeyEvent{
	newKeyEvent("SPACE", stateClick, eventStart, "Start Game"),
	newKeyEvent("Q", stateClick, eventDestroy, "Quit"),
	newKeyEvent("A", stateClick, eventToggleAI, "Toggle AI"),
	newKeyEvent("P", stateClick, eventTogglePause, "Toggle Pause"),
	newKeyEvent("R", stateClick, eventReset, "Reset Round"),
	newKeyEvent("T", stateClick, eventSwitchTheme, "Switch Theme"),
	newKeyEvent("W", stateHold, eventP1Up, "Move Player1 Up"),
	newKeyEvent("S", stateHold, eventP1Down, "Move Player1 Down"),
	newKeyEvent("Up", stateHold, eventP2Up, "Move Player2 Up"),
	newKeyEvent("Down", stateHold, eventP2Down, "Move Player2 Down"),
}

var legendKeysGame = []KeyEvent{
	newKeyEvent("Q", stateClick, eventDestroy, "Quit"),
	newKeyEvent("P", stateClick, eventTogglePause, "Toggle Pause"),
	newKeyEvent("R", stateClick, eventReset, "Reset Round"),
	newKeyEvent("T", stateClick, eventSwitchTheme, "Switch Theme"),
	newKeyEvent("W", stateHold, eventP1Up, "Move Player1 Up"),
	newKeyEvent("S", stateHold, eventP1Down, "Move Player1 Down"),
	newKeyEvent("Up", stateHold, eventP2Up, "Move Player2 Up"),
	newKeyEvent("Down", stateHold, eventP2Down, "Move Player2 Down"),
}
