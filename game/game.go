package game

import (
	"time"

	"github.com/MarinX/keylogger"
	"github.com/gdamore/tcell"
)

type Game struct {
	screen     tcell.Screen
	players    *Players
	ball       *Ball
	keys       *[]string
	event      chan keyState
	dispEvent  chan keyState
	ticker     *time.Ticker
	keyboard   []*keylogger.KeyLogger
	theme      *ThemeHandler
	started    bool
	paused     bool
	hardPaused bool
}

type GameSettings struct {
	BgHidden bool
}

// Create a new game instance initialised and ready to go
func Create(options *GameSettings) (*Game, error) {
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
	g.players = newPlayers(w, h, padding)
	g.ball = newBall((w-ballDiam)/2, 0, ballDiam, 1, 1)

	// keyboard channels and key state init
	var keys []string
	g.keys = &keys
	g.event = make(chan keyState)
	g.dispEvent = make(chan keyState)

	// ticker init according to framerate variable
	g.ticker = time.NewTicker(1000000 / framerate * time.Microsecond)

	// initialise default theme
	g.theme = newThemeHandler(!optns.BgHidden)

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
		go keyboardListen(kb, g.event, g.dispEvent)
	}

	// initial screen overlay
	//g.drawGameTick()

	for {
		select {
		case <-g.ticker.C:
			// draw the gui and update the terminal
			g.drawGUI()
			g.updateTerminal()
		case lol := <-g.event:
			// filter the key from the slice if it is in there
			var newKeys []string

			for _, key := range *g.keys {
				if lol.Name != key {
					newKeys = append(newKeys, key)
				}
			}

			// update only if there are less than 5 chars and key is pressed down
			if len(*g.keys) < 5 && lol.Down {
				newKeys = append(newKeys, lol.Name)
			}

			// update the old array
			g.keys = &newKeys
		case e := <-g.dispEvent:
			switch e.Name {
			case eventDestroy:
				return
			case eventStart:
				g.Start()
			case eventTogglePause:
				g.togglePause()
			case eventReset:
				g.Reset()
			case eventSwitchTheme:
				g.switchTheme()
			case eventToggleBg:
				g.toggleBackground()
			}
		}
	}
}

// Start the game
func (g *Game) Start() {
	g.started = true
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
		g.drawPauseText()
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

// End the game. Must be called when game ends or if g.Init fails for cleanup.
func (g *Game) End() {
	// g.screen.Fini() is fixed, but a leftover "q" is left when terminal is closed.
	// looking for fix
	// check for nil is required in case the game ends with an error and g.screen is not set,
	// in that case it would panic with nil pointer dereference. If it doesnt exist,
	// then there is no need to clean it up.
	if g.screen != nil {
		g.screen.Fini()
	}
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

// Check and handle collisions between balls, walls and platforms
func (g *Game) checkCollision() {
	w, h := g.ball.Coords()
	d := g.ball.Diam()
	vx, vy := g.ball.Vels()
	sw, sh := g.screen.Size()

	// wall collisions
	if w < 1 && vx < 0 {
		g.players.P2.AddPoint()

		// hardPause cannot be unpaused by player
		g.hardPaused = true

		go func() {
			time.Sleep(scoreSleepSecs * time.Second)
			g.hardPaused = false
			g.Reset()
		}()
	}
	if w >= sw-2*d-1 && vx > 0 {
		g.players.P1.AddPoint()

		// hardPause cannot be unpaused by player
		g.hardPaused = true

		go func() {
			time.Sleep(scoreSleepSecs * time.Second)
			g.hardPaused = false
			g.Reset()
		}()
	}
	if h <= 1 && vy < 0 {
		g.ball.SwitchY()
	}
	if h >= sh-d && vy > 0 {
		g.ball.SwitchY()
	}

	// platform collisions
	p1w, p1h := g.players.P1.Coords()
	p2w, p2h := g.players.P2.Coords()

	if h+d > p1h && h < p1h+platformHeight && (p1w+platformWidth)/2+1 == w/2 {
		g.ball.SwitchX()
	}
	if h+d > p2h && h < p2h+platformHeight && p2w/2+1 == w/2+d {
		g.ball.SwitchX()
	}
}
