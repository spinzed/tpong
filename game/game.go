package game

import (
	"os"
	"os/exec"
	"time"

	"github.com/MarinX/keylogger"
	"github.com/gdamore/tcell"
	"golang.org/x/crypto/ssh/terminal"
)

type Players struct {
	P1 *Player
	P2 *Player
}

func (p *Players) Arrayify() []*Player {
	return []*Player{p.P1,p.P2}
}

type Game struct {
	screen     tcell.Screen
	players    *Players
	ball       *Ball
	keys       *[]string
	event      chan keyState
	dispEvent  chan keyState
	ticker     *time.Ticker
	keyboard   []*keylogger.KeyLogger
	oldstate   *terminal.State
	paused     bool
	hardPaused bool
}

// Get a pair of players ready and initialised
func newPlayers(w int, h int, padding int) *Players {
	initialPos := (h - platformHeight) / 2
	p1 := newPlayer(playerP1, initialPos, padding)
	p2 := newPlayer(playerP2, initialPos, w-padding)

	return &Players{p1, p2}
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
func (g *Game) Init() error {
	// This is a workaround since g.screen.Fini is broken. Save terminal old state
	// before making it raw for tcell to use

	// keyboard init
	k, err := newKeyboard()

	if err != nil {
		return err
	}

	g.keyboard = k

	// putting terminal to raw state
	oldState, err := terminal.MakeRaw(0)

	if err != nil {
		return err
	}

	g.oldstate = oldState

	// screen init
	s, err := initScreen()

	if err != nil {
		return err
	}

	g.screen = s

	// players init
	w, h := g.screen.Size()
	g.players = newPlayers(w, h, padding)

	// ball init
	g.ball = newBall((w-ballDiam)/2, 0, ballDiam, 1, 1)

	// keys map init
	keys := make([]string, 0)
	g.keys = &keys

	// event channel init
	g.event = make(chan keyState)

	// dispEvent channel init
	g.dispEvent = make(chan keyState)

	// ticker init
	g.ticker = time.NewTicker(1000000 / framerate * time.Microsecond)

	// initial overlay
	g.drawPlayer(*g.players.P1)
	g.drawPlayer(*g.players.P2)
	g.drawOverlay()

	g.screen.Show()

	// signal that everything is ok
	return nil
}

// Start the game. Must be called after initialization
func (g *Game) Loop() {
	// make sure than end cleanup function is executed
	defer g.End()

	// start the keyboard input listeners
	// Im a bloody genious
	for _, kb := range g.keyboard {
		go keyboardListen(kb, g.event, g.dispEvent)
	}

	for {
		select {
		case <-g.ticker.C:
			// perform 1 game tick
			g.Tick(g.keys)
		case lol := <-g.event:
			// filter the slice from the key if it is in there
			newKeys := make([]string, 0)

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
			case eventTogglePause:
				g.paused = !g.paused
			case eventReset:
				g.Reset()
			}
		}
	}
}

// Reset the game
func (g *Game) Reset() {
	g.ball.Reset()
	g.players.P1.Reset()
	g.players.P2.Reset()
}

// End the game
func (g *Game) End() {
	// there is a problem where keyboard input persists and it is shown in the
	// command line when the game exits. Looking for fix

	// clear the screen. The second line connects the command's stdout with
	// the of one of this terminal's session
	if g.screen != nil {
		g.screen.Clear()
		g.screen.Show()

		clr := exec.Command("clear")
		clr.Stdout = os.Stdout
		clr.Run()
	}

	// bring teminal back from raw mode
	if g.oldstate != nil {
		terminal.Restore(0, g.oldstate)
	}
}

// Move player 1 char higher. If player is at the edge, do nothing.
func (g *Game) movePlayerUp(p *Player) {
	if _, h := p.Coords(); h < 1 {
		return
	}

	p.GoUp()
}

// Move player 1 char lower. If player is at the edge, do nothing.
func (g *Game) movePlayerDown(p *Player) {
	_, sh := g.screen.Size()
	_, ph := p.Coords()

	if ph > sh-p.GetHeight()-1 {
		return
	}

	p.GoDown()
}

func (g *Game) checkCollision() {
	w, h := g.ball.Coords()
	d := g.ball.Diam()
	vx, vy := g.ball.Vels()
	sw, sh := g.screen.Size()

	// wall collisions
	if w < 1 && vx < 0 {
		g.ball.SwitchX()
		g.players.P2.AddPoint()

		// hardPause cannot be unpaused by player
		g.hardPaused = true

		go func() {
			time.Sleep(2 * time.Second)
			g.hardPaused = false
			g.Reset()
		}()
	}
	if w >= sw-2*d-1 && vx > 0 {
		g.ball.SwitchX()
		g.players.P1.AddPoint()

		// hardPause cannot be unpaused by player
		g.hardPaused = true

		go func() {
			time.Sleep(2 * time.Second)
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

	if h+d > p1h && h < p1h+platformHeight && (p1w+platformWidth)/2 == w/2 {
		g.ball.SwitchX()
	}
	if h+d > p2h && h < p2h+platformHeight && p2w/2 == w/2+d {
		g.ball.SwitchX()
	}
}
