package main

import "errors"

// General game constants
const (
	framerate        = 30
	ballDiam         = 1
	platformWidth    = 1
	platformHeight   = 6
	padding          = 10
	letterWidth      = 8
	letterHeight     = 7
	letterGap        = 1
	initialDelaySecs = 2
	scoreSleepSecs   = 2
	startMenuLttrGap = 8
	ballRandomness   = 20
	playerP1         = "p1"
	playerP2         = "p2"
)

// Key states
const (
	stateClick   = "stateClick"
	stateRelease = "stateRelease"
	stateHold    = "stateHold"
	stateHoldEnd = "stateHoldEnd"
	stateNormal  = "stateHeld"
)

// StatedEvent states
const (
	eventStatePulse    = "pulse"
	eventStateStarting = "starting"
	eventStateEnding   = "ending"
)

// Screens
const (
	screenStartMenu = "screenMainMenu"
	screenGame      = "screenGame"
)

// Errors
var NoKeyboardError = errors.New("no keyboard detected")

// Game events
var (
	eventP1Up        = NewEvent("eventP1Up", "Move Player1 Up")
	eventP1Down      = NewEvent("eventP1Down", "Move Player1 Down")
	eventP2Up        = NewEvent("eventP2Up", "Move Player2 Up")
	eventP2Down      = NewEvent("eventP2Down", "Move Player2 Down")
	eventMenuUp      = NewEvent("eventMenuUp", "Select Action Above")
	eventMenuDown    = NewEvent("eventMenuDown", "Select Action Below")
	eventMenuSelect  = NewEvent("eventMenuSelect", "Select Menu Action")
	eventTogglePause = NewEvent("eventTogglePause", "Toggle Pause")
	eventStart       = NewEvent("eventStart", "Start Game")
	eventDestroy     = NewEvent("eventDestroy", "Quit")
	eventReset       = NewEvent("eventReset", "Reset Round")
	eventSwitchTheme = NewEvent("eventSwitchTheme", "Switch Theme")
	eventToggleBg    = NewEvent("eventToggleBg", "Toggle Background Visibility")
	eventToggleAI    = NewEvent("eventToggleAI", "Toggle AI")
)

// Keys which can be triggered when game is started
var keysStart = KeyEventMap{
	*NewKeyEvent(NewKey("SPACE", stateClick), eventStart),
	*NewKeyEvent(NewKey("ENTER", stateClick), eventMenuSelect),
	*NewKeyEvent(NewKey("Up", stateNormal), eventMenuUp),
	*NewKeyEvent(NewKey("Down", stateNormal), eventMenuDown),
	*NewKeyEvent(NewKey("Q", stateClick), eventDestroy),
	*NewKeyEvent(NewKey("P", stateClick), eventTogglePause),
	*NewKeyEvent(NewKey("R", stateClick), eventReset),
	*NewKeyEvent(NewKey("T", stateClick), eventSwitchTheme),
	*NewKeyEvent(NewKey("A", stateClick), eventToggleAI),
}

// Keys which can be triggered when game is started
var keysGame = KeyEventMap{
	*NewKeyEvent(NewKey("W", stateHold), eventP1Up),
	*NewKeyEvent(NewKey("S", stateHold), eventP1Down),
	*NewKeyEvent(NewKey("Up", stateHold), eventP2Up),
	*NewKeyEvent(NewKey("Down", stateHold), eventP2Down),
	*NewKeyEvent(NewKey("Q", stateClick), eventDestroy),
	*NewKeyEvent(NewKey("P", stateClick), eventTogglePause),
	*NewKeyEvent(NewKey("R", stateClick), eventReset),
	*NewKeyEvent(NewKey("T", stateClick), eventSwitchTheme),
}

var altKeysGame = KeyEventMap{
	*NewKeyEvent(NewKey("W", stateHold), eventP1Up),
	*NewKeyEvent(NewKey("S", stateHold), eventP1Down),
	*NewKeyEvent(NewKey("ENTER", stateClick), eventMenuSelect),
	*NewKeyEvent(NewKey("Up", stateNormal), eventMenuUp),
	*NewKeyEvent(NewKey("Down", stateNormal), eventMenuDown),
	*NewKeyEvent(NewKey("Q", stateClick), eventDestroy),
	*NewKeyEvent(NewKey("P", stateClick), eventTogglePause),
	*NewKeyEvent(NewKey("R", stateClick), eventReset),
	*NewKeyEvent(NewKey("T", stateClick), eventSwitchTheme),
}

var legendKeysStart = KeyEventMap{
	*NewKeyEvent(NewKey("SPACE", stateClick), eventStart),
	*NewKeyEvent(NewKey("Q", stateClick), eventDestroy),
	*NewKeyEvent(NewKey("P", stateClick), eventTogglePause),
	*NewKeyEvent(NewKey("R", stateClick), eventReset),
	*NewKeyEvent(NewKey("A", stateClick), eventToggleAI),
	*NewKeyEvent(NewKey("T", stateClick), eventSwitchTheme),
	*NewKeyEvent(NewKey("W", stateHold), eventP1Up),
	*NewKeyEvent(NewKey("S", stateHold), eventP1Down),
	*NewKeyEvent(NewKey("Up", stateHold), eventP2Up),
	*NewKeyEvent(NewKey("Down", stateHold), eventP2Down),
}

var legendKeysGame = KeyEventMap{
	*NewKeyEvent(NewKey("Q", stateClick), eventDestroy),
	*NewKeyEvent(NewKey("P", stateClick), eventTogglePause),
	*NewKeyEvent(NewKey("R", stateClick), eventReset),
	*NewKeyEvent(NewKey("T", stateClick), eventSwitchTheme),
	*NewKeyEvent(NewKey("W", stateHold), eventP1Up),
	*NewKeyEvent(NewKey("S", stateHold), eventP1Down),
	*NewKeyEvent(NewKey("Up", stateHold), eventP2Up),
	*NewKeyEvent(NewKey("Down", stateHold), eventP2Down),
}
