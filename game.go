package main

import (
	"time"
)

type Game struct {
	ui           *UIHandler
	input        *InputHandler
	players      *Players
	activeEvents *[]Event
	keyDisp      chan StateEvent
	ticker       *time.Ticker
	state        GameState
	lastState    GameState
	server       *Server
	localClient  Client
	started      bool
	aiActive     bool
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

// Initialize the game. Must be called on new game instance
func (g *Game) Init(optns *GameSettings) error {
	var err error

	defer func() {
		if err != nil {
			g.End()
		}
	}()

	// screen init
	if g.ui, err = NewUIHandler(optns.BgHidden); err != nil {
		return err
	}
	// input handler init
	if g.input, err = NewInputHandler(g.getKeys); err != nil {
		return err
	}

	// players init
	w, h := g.ui.screen.Size()
	g.players = newPlayers(w, h, padding)

	// input channels and key/event state init
	g.activeEvents = new([]Event)
	g.keyDisp = make(chan StateEvent)

	// ticker init according to framerate variable
	g.ticker = time.NewTicker(1000000 / framerate * time.Microsecond)

	// mark the inital game state
	g.state = gameStateStarting

	// signal that everything is ok
	return nil
}

// Start the game. Must be called after initialization
func (g *Game) Loop() {
	// make sure that end cleanup function is executed
	defer g.End()

	// start the keyboard input listeners
	g.input.Listen(g.keyDisp)

	// set the main menu game state
	g.SetState(gameStateMainMenu)

	// if the game is not ending or ended
	for g.state != gameStateEnding && g.state != gameStateEnded {
		select {
		case <-g.ticker.C:
			// draw the gui and update the terminal.
			// this MUST be the only place where these methods are called.
			g.drawGUI()
			g.updateTerminal()
		case e := <-g.keyDisp:
			g.registerEvent(e)
		}
	}
}

// Break the game loop. End cleanup function should invoke when the loop breaks
func (g *Game) EndLoop() {
	g.SetState(gameStateEnding)
}

// Start the game
func (g *Game) Start() {
	if g.started {
		return
	}

	// setup the server and register the clients
	g.server = NewServer()

	recv := make(chan StateEvent)

	// register the local client
	g.localClient = NewLocalClient(g.ui.screen, g.server.EventListener, recv)
	g.server.RegisterClient(g.localClient)

	// register every event that the client recieved from the server
	go func() {
		for {
			evt := <-recv
			g.dispatchEvent(evt)
		}
	}()

	// reset the ball from the start menu
	g.ui.ball.Reset()
	g.SetState(gameStateInGame)
	g.started = true
}

// End the game. Must be called when game ends or if g.Init fails for cleanup.
func (g *Game) End() {
	// needs fix: leftover q in terminal when game ends
	// check for nil is required in case the game ends with an error and g.ui.screen is not set,
	// in that case it would panic with nil pointer dereference. If it doesnt exist,
	// then there is no need to clean it up.
	g.ui.Destroy()
}

// Reset the game
func (g *Game) Reset() {
	g.ui.ball.Reset()

	for _, p := range g.players.GetAll() {
		p.Reset()
	}
}

//func (g *Game) RegisterClient(p Player) {
//	c := NewLocalClient(p, g.ui.screen)
//	g.localClients = append(g.localClients, c)
//	g.server.RegisterClient(c)
//}

func (g *Game) SetState(s GameState) {
	g.lastState = g.state
	g.state = s
}

// Universal pause toggler
func (g *Game) __tpause(state GameState) {
	if !g.started && state == gameStatePaused {
		return
	}

	if g.state == state {
		g.SetState(g.lastState)
	} else {
		g.SetState(state)
	}
}

// Toggle pause
func (g *Game) togglePause() {
	g.__tpause(gameStatePaused)
}

// Toggle hard pause
func (g *Game) toggleHardPause() {
	g.__tpause(gameStateHardPaused)
}

// Switch next theme and update the terminal
func (g *Game) switchTheme() {
	g.ui.theme.Switch()
}

// Toggle background on/off and update the terminal
func (g *Game) toggleBackground() {
	g.ui.theme.ToggleBg()
}

func (g *Game) toggleAI() {
	g.aiActive = !g.aiActive
}

// Gets the current active key map.
func (g *Game) getKeys() *KeyEventMap {
	stuff := screenData[g.state].Keys

	// if the game is hard paused, preserve the event key map,
	// but block their effects elsewhere
	if g.state == gameStateHardPaused {
		stuff = screenData[g.lastState].Keys
	}
	return &stuff
}

// Gets the current active legend key map.
func (g *Game) getLegendKeys() *KeyEventMap {
	stuff := screenData[g.state].Keys
	return &stuff
}

// Gets the current active legend key map.
func (g *Game) getLegendType() string {
	return screenData[g.state].LegendType
}

// Move player 1 char higher. If player is at the edge, do nothing.
func (g *Game) movePlayerUp(p *Player) {
	// if the platform is at the top edge, do nothing
	if _, h := p.Coords(); h < 1 {
		return
	}

	p.GoUp()
}

// Move player 1 char lower. If the player is at the edge, do nothing.
func (g *Game) movePlayerDown(p *Player) {
	_, sh := g.ui.ScreenSize()
	_, ph := p.Coords()
	// if the platform is at the bottom edge, do nothing
	if ph > sh-p.GetHeight()-1 {
		return
	}

	p.GoDown()
}

// Move selected legend item up.
func (g *Game) moveMenuSelectedUp() {
	g.ui.MoveMenuSelectedUp(g.state)
}

// Move selected legend item down.
func (g *Game) moveMenuSelectedDown() {
	g.ui.MoveMenuSelectedDown(g.state)
}

func (g *Game) doSelectedMenuAction() {
	evt := g.ui.GetSelectedMenuAction(g.state)
	evtState := StateEvent{evt, eventStatePulse}

	g.registerEvent(evtState)
}

func containsEvent(slice []Event, elm Event) bool {
	for _, e := range slice {
		if e == elm {
			return true
		}
	}
	return false
}

func filterEvent(e Event, old []Event) []Event {
	var newEvents []Event

	for _, event := range old {
		if e.Name != event.Name {
			newEvents = append(newEvents, event)
		}
	}
	return newEvents
}

// Decides the type of a event and depending on it, it either triggers an event
// or passes to the server so that it can broadcast it to every node.
// This is usually called when a new input has been passed by the input handler.
func (g *Game) registerEvent(e StateEvent) {
	if containsEvent(localEvents, e.Event) {
		g.dispatchEvent(e)
		return
	}
	g.localClient.SendUpdate(e)
}

// Triggers an action based on the event an key that have been passed.
func (g *Game) dispatchEvent(e StateEvent) {
	// filter the key from the slice if it is in there
	switch e.State {
	case eventStateStarting:
		// in case the event existed in the array, remove it
		newEvents := filterEvent(e.Event, *g.activeEvents)

		// update only if there are less than 5 chars and key is pressed down
		if len(*g.activeEvents) < 5 {
			newEvents = append(newEvents, e.Event)
		}

		// update the old array
		g.activeEvents = &newEvents
	case eventStateEnding:
		*g.activeEvents = filterEvent(e.Event, *g.activeEvents)
	case eventStatePulse:
		g.__rawDispatchAction(e.Event)
	}
}

// Calls an action according to the event.
func (g *Game) __rawDispatchAction(e Event) {
	switch e {
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
	case eventToggleAI:
		g.toggleAI()
	}
}
