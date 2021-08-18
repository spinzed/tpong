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

type GameState string

// Game states
const (
	gameStateStarting   GameState = "starting"
	gameStateMainMenu   GameState = "mainMenu"
	gameStateInGame     GameState = "inGame"
	gameStatePaused     GameState = "paused"
	gameStateHardPaused GameState = "hardPaused"
	gameStateEnding     GameState = "ending" // to break the game loop
	gameStateEnded      GameState = "ended"  // when game.End is called
)

// Key states
const (
	keyStateNormal  = "stateNormal"
	keyStateClick   = "stateClick"
	keyStateRelease = "stateRelease"
	keyStateHold    = "stateHold"
	keyStateEnd     = "stateHoldEnd"
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
	*NewKeyEvent(NewKey("SPACE", keyStateClick), eventStart),
	*NewKeyEvent(NewKey("ENTER", keyStateClick), eventMenuSelect),
	*NewKeyEvent(NewKey("Up", keyStateNormal), eventMenuUp),
	*NewKeyEvent(NewKey("Down", keyStateNormal), eventMenuDown),
	*NewKeyEvent(NewKey("Q", keyStateClick), eventDestroy),
	*NewKeyEvent(NewKey("P", keyStateClick), eventTogglePause),
	*NewKeyEvent(NewKey("R", keyStateClick), eventReset),
	*NewKeyEvent(NewKey("T", keyStateClick), eventSwitchTheme),
	*NewKeyEvent(NewKey("A", keyStateClick), eventToggleAI),
}

// Keys which can be triggered when game is started
var keysGame = KeyEventMap{
	*NewKeyEvent(NewKey("W", keyStateHold), eventP1Up),
	*NewKeyEvent(NewKey("S", keyStateHold), eventP1Down),
	*NewKeyEvent(NewKey("Up", keyStateHold), eventP2Up),
	*NewKeyEvent(NewKey("Down", keyStateHold), eventP2Down),
	*NewKeyEvent(NewKey("Q", keyStateClick), eventDestroy),
	*NewKeyEvent(NewKey("P", keyStateClick), eventTogglePause),
	*NewKeyEvent(NewKey("R", keyStateClick), eventReset),
	*NewKeyEvent(NewKey("T", keyStateClick), eventSwitchTheme),
}

var keysPause = KeyEventMap{
	*NewKeyEvent(NewKey("W", keyStateHold), eventP1Up),
	*NewKeyEvent(NewKey("S", keyStateHold), eventP1Down),
	*NewKeyEvent(NewKey("ENTER", keyStateClick), eventMenuSelect),
	*NewKeyEvent(NewKey("Up", keyStateNormal), eventMenuUp),
	*NewKeyEvent(NewKey("Down", keyStateNormal), eventMenuDown),
	*NewKeyEvent(NewKey("Q", keyStateClick), eventDestroy),
	*NewKeyEvent(NewKey("P", keyStateClick), eventTogglePause),
	*NewKeyEvent(NewKey("R", keyStateClick), eventReset),
	*NewKeyEvent(NewKey("T", keyStateClick), eventSwitchTheme),
}

var legendKeysStart = KeyEventMap{
	*NewKeyEvent(NewKey("SPACE", keyStateClick), eventStart),
	*NewKeyEvent(NewKey("Q", keyStateClick), eventDestroy),
	*NewKeyEvent(NewKey("P", keyStateClick), eventTogglePause),
	*NewKeyEvent(NewKey("R", keyStateClick), eventReset),
	*NewKeyEvent(NewKey("A", keyStateClick), eventToggleAI),
	*NewKeyEvent(NewKey("T", keyStateClick), eventSwitchTheme),
	*NewKeyEvent(NewKey("W", keyStateHold), eventP1Up),
	*NewKeyEvent(NewKey("S", keyStateHold), eventP1Down),
	*NewKeyEvent(NewKey("Up", keyStateHold), eventP2Up),
	*NewKeyEvent(NewKey("Down", keyStateHold), eventP2Down),
}

var legendKeysPaused = KeyEventMap{
	*NewKeyEvent(NewKey("Q", keyStateClick), eventDestroy),
	*NewKeyEvent(NewKey("P", keyStateClick), eventTogglePause),
	*NewKeyEvent(NewKey("R", keyStateClick), eventReset),
	*NewKeyEvent(NewKey("T", keyStateClick), eventSwitchTheme),
	*NewKeyEvent(NewKey("W", keyStateHold), eventP1Up),
	*NewKeyEvent(NewKey("S", keyStateHold), eventP1Down),
	*NewKeyEvent(NewKey("Up", keyStateHold), eventP2Up),
	*NewKeyEvent(NewKey("Down", keyStateHold), eventP2Down),
}

var screenData = map[GameState]GameInfo{
	gameStateMainMenu: {keysStart, legendKeysStart, "middle"},
	gameStateInGame:   {keysGame, nil, ""},
	gameStatePaused:   {keysPause, legendKeysPaused, "left"},
}
