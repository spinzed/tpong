package main

import (
	"time"

	"github.com/MarinX/keylogger"
	"github.com/gdamore/tcell"
)

type Game struct {
	screen       tcell.Screen
	players      *Players
	ball         *Ball
	activeEvents *[]string
	keyDisp      chan KeyEvent
	ticker       *time.Ticker
	keyboard     []*keylogger.KeyLogger
	theme        *ThemeHandler
	keyData      *ScreenKeyData
	started      bool
	paused       bool
	hardPaused   bool
	loopActive   bool
}

type GameSettings struct {
	BgHidden bool
}

// Create a new game instance initialised and ready to go
func CreateGame(options *GameSettings) (*Game, error) {
	g := Game{}

	if err := g.Init(options); err != nil {
		return nil, err
	}

	return &g, nil
}

// Create screen ready to use
func initScreen() (tcell.Screen, error) {
	s, err := tcell.NewScreen()

	if err != nil {
		return s, err
	}

	if err = s.Init(); err != nil {
		return s, err
	}

	return s, nil
}

// Initialize the game. Must be called on new game instance
func (g *Game) Init(optns *GameSettings) error {
	var err error

	defer func() {
		if err != nil {
			g.End()
		}
	}()

	// keyboard init
	k, err := newKeyboard()

	if err != nil {
		return err
	}

	g.keyboard = k

	// screen init
	s, err := initScreen()

	if err != nil {
		return err
	}

	g.screen = s

	// game asset init
	w, h := g.screen.Size()
	// half of the width needs to be divisible with 2 because of ball
	if w/2 % 2 == 1 {
		w += 2
	}
	g.players = newPlayers(w, h, padding)
	g.ball = newBall((w-ballDiam)/2, 0, ballDiam, 1, 1)

	// keyboard channels and key/event state init
	var activeEvents []string
	g.activeEvents = &activeEvents
	g.keyDisp = make(chan KeyEvent)

	// ticker init according to framerate variable
	g.ticker = time.NewTicker(1000000 / framerate * time.Microsecond)

	// initialise default theme
	g.theme = newThemeHandler(!optns.BgHidden)

	// mark that the game loop active - g.Loop
	// should be made to false to end the game
	g.loopActive = true

	// keyData is a struct which stores information about the current screen's,
	// key actions and legend. It doesn't control the screen behavior.
	// all screens are located in the screens map
	// the inital screen is the start menu screen
	g.keyData = screens[screenStartMenu]

	// signal that everything is ok
	return nil
}

// Start the game. Must be called after initialization
func (g *Game) Loop() {
	// make sure that end cleanup function is executed
	defer g.End()

	// polling events from terminal has been substituted for listening for events from
	// /dev/input/event*, but this has to be left dangling like this for g.screen.Fini
	// to function properly
	go func() {
		for {
			g.screen.PollEvent()
		}
	}()

	// start the keyboard input listeners
	// Im a bloody genious
	for _, kb := range g.keyboard {
		go keyboardListen(kb, g.getKeys, g.keyDisp)
	}

	for g.loopActive {
		select {
		case <-g.ticker.C:
			// draw the gui and update the terminal
			g.drawGUI()
			g.updateTerminal()
		case e := <-g.keyDisp:
			g.dispatchEvent(e)
		}
	}
}

// Break the game loop. End cleanup function should invoke when the loop breaks
func (g *Game) EndLoop() {
	// reset the ball from the start menu
	g.loopActive = false
}

// Start the game
func (g *Game) Start() {
	if !g.started {
		// reset the ball from the start menu
		g.ball.Reset()
		g.keyData = screens[screenGame]
		g.started = true
	}
}

// End the game. Must be called when game ends or if g.Init fails for cleanup.
func (g *Game) End() {
	// needs fix: leftover q in terminal when game ends
	// check for nil is required in case the game ends with an error and g.screen is not set,
	// in that case it would panic with nil pointer dereference. If it doesnt exist,
	// then there is no need to clean it up.
	if g.screen != nil {
		g.screen.Fini()
	}
}

// Reset the game
func (g *Game) Reset() {
	g.ball.Reset()

	for _, p := range g.players.GetAll() {
		p.Reset()
	}

	g.drawGUI()
}

// Toggle pause
func (g *Game) togglePause() {
	if !g.started || g.hardPaused {
		return
	}

	g.paused = !g.paused

	if g.paused {
		g.drawLegend()
		g.screen.Show()
	}
}

// Switch next theme and update the terminal
func (g *Game) switchTheme() {
	g.theme.Switch()
	g.drawGUI()
}

// Toggle background on/off and update the terminal
func (g *Game) toggleBackground() {
	g.theme.ToggleBg()
	g.drawGUI()
}

// Gets the current active key map.
func (g *Game) getKeys() *map[Key]Event {
	if g.paused && g.keyData.AltKeys != nil {
		return g.keyData.AltKeys
	}
	return g.keyData.Keys
}

// Move player 1 char higher. If player is at the edge, do nothing.
func (g *Game) movePlayerUp(p *Player) {
	// if the platform is at the top edge, do nothing
	if _, h := p.Coords(); h < 1 {
		return
	}

	p.GoUp()
}

// Move player 1 char lower. If player is at the edge, do nothing.
func (g *Game) movePlayerDown(p *Player) {
	_, sh := g.screen.Size()
	_, ph := p.Coords()

	// if the platform is at the bottom edge, do nothing
	if ph > sh-p.GetHeight()-1 {
		return
	}

	p.GoDown()
}

// Move selected legend item up.
func (g *Game) moveMenuSelectedUp() {
	legend := g.keyData.Legend
	legend.Selected--

	if legend.Selected < 0 {
		legend.Selected = len(*legend.Keys) - 1
	}

	keysarr := *legend.Keys
	ev := keysarr[legend.Selected]

	// if the event is hold up/down, skip it
	if ev.State == stateHold || ev.State == stateHoldEnd {
		g.moveMenuSelectedUp()
	}
}

// Move selected legend item down.
func (g *Game) moveMenuSelectedDown() {
	legend := g.keyData.Legend
	legend.Selected++

	if legend.Selected > len(*legend.Keys)-1 {
		legend.Selected = 0
	}

	keysarr := *legend.Keys
	ev := keysarr[legend.Selected]

	// if the event is hold up/down, skip it
	if ev.State == stateHold || ev.State == stateHoldEnd {
		g.moveMenuSelectedDown()
	}
}

// Dispatches an action according to the selected legend item.
func (g *Game) doSelectedMenuAction() {
	legend := g.keyData.Legend
	keysarr := *legend.Keys

	ev := keysarr[legend.Selected]

	// if it triggers hold or stop hold action, return (since the key isnt actually held)
	if ev.State == stateHold || ev.State == stateHoldEnd {
		return
	}

	g.dispatchEvent(ev)
}

// Triggers an action based on the event an key that have been passed.
func (g *Game) dispatchEvent(e KeyEvent) {
	// filter the key from the slice if it is in there
	switch e.Key.State {
	case stateHold, stateHoldEnd:
		var newEvents []string

		for _, event := range *g.activeEvents {
			if e.Event.Name != event {
				newEvents = append(newEvents, event)
			}
		}

		// update only if there are less than 5 chars and key is pressed down
		if len(*g.activeEvents) < 5 && e.State == stateHold {
			newEvents = append(newEvents, e.Event.Name)
		}

		// update the old array
		g.activeEvents = &newEvents
	case stateClick, stateRelease, stateNormal:
		g.dispatchAction(e.Event)
	}
}

// Calls an action according to the event.
func (g *Game) dispatchAction(e Event) {
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
	case eventMenuUp:
		g.moveMenuSelectedUp()
	case eventMenuDown:
		g.moveMenuSelectedDown()
	case eventMenuSelect:
		g.doSelectedMenuAction()
	}
}
