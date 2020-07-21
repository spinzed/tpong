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

type Game struct {
	screen    tcell.Screen
	players   *Players
	ball      *Ball
	event     chan keyState
	dispEvent chan keyState
	ticker    *time.Ticker
	keyboard  *keylogger.KeyLogger
	oldstate  *terminal.State
	paused    bool
}

// Get a pair of players ready and initialised
func newPlayers(w int, h int, padding int) *Players {
	initialPos := (h - playerHeight) / 2
	p1 := newPlayer("p1", initialPos, padding)
	p2 := newPlayer("p2", initialPos, w-padding)

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
	diam := 1
	g.ball = newBall((w-diam)/2, 0, diam, 1, 1)

	// event channel init
	g.event = make(chan keyState)

	// dispEvent channel init
	g.dispEvent = make(chan keyState)

	// ticker init
	g.ticker = time.NewTicker(1000000 / framerate * time.Microsecond)

	// keyboard init
	k, err := newKeyboard()

	if err != nil {
		// Release the screen on fail. May move this elsewhere
		g.screen.Fini()
		return err
	}

	g.keyboard = k

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

	keys := make([]string, 0)

	// start the keyboard input listeners
	go keyboardListen(g.keyboard, g.event, g.dispEvent)

	for {
		select {
		case <-g.ticker.C:
			// update the game state every tick.

			// event that persist even when game is paused
			for _, key := range keys {
				switch key {
				case eventDestroy:
					return
				}
			}

			if g.paused {
				break
			}

			// keys that don't persist when game is paused
			for _, key := range keys {
				switch key {
				case eventDestroy:
					return
				case eventP1Up:
					g.MovePlayerUp(g.players.P1)
				case eventP1Down:
					g.MovePlayerDown(g.players.P1)
				case eventP2Up:
					g.MovePlayerUp(g.players.P2)
				case eventP2Down:
					g.MovePlayerDown(g.players.P2)
				}
			}

			// update the screen.
			// these aren't expensive since they just check for changes on the canvas
			// and if there aren't any, nothing will be updated therefore no bloat
			g.screen.Clear()

			g.drawPlayer(*g.players.P1)
			g.drawPlayer(*g.players.P2)
			g.drawBall()
			g.drawOverlay()

			g.screen.Show()

			// move the ball for 1 tick
			g.checkCollision()
			g.ball.Move()
		case lol := <-g.event:
			// filter the slice from the key if it is in there
			newKeys := make([]string, 0)

			for _, key := range keys {
				if lol.Name != key {
					newKeys = append(newKeys, key)
				}
			}

			// update only if there are less than 5 chars and key is pressed down
			if len(keys) < 5 && lol.Down {
				newKeys = append(newKeys, lol.Name)
			}

			// update the old array
			keys = newKeys
		case e := <-g.dispEvent:
			switch e.Name {
			case eventTogglePause:
				g.paused = !g.paused
			}
		}
	}
}

// End the game
func (g *Game) End() {
	// clear the screen. The second line connects the command's stdout with
	// the of one of this terminal's session
	g.screen.Clear()
	g.screen.Show()
	clr := exec.Command("clear")
	clr.Stdout = os.Stdout
	clr.Run()

	// bring teminal back from raw mode
	terminal.Restore(0, g.oldstate)
}

// Move player 1 char higher. If player is at the edge, do nothing.
func (g *Game) MovePlayerUp(p *Player) {
	if _, h := p.Coords(); h < 1 {
		return
	}

	p.GoUp()
}

// Move player 1 char lower. If player is at the edge, do nothing.
func (g *Game) MovePlayerDown(p *Player) {
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
	// x portion is temporary, just for testing purposes for now
	if w < 1 && vx < 0 {
		g.ball.SwitchX()
	}
	if w == sw-2*d-1 && vx > 0 {
		g.ball.SwitchX()
	}
	if h == 1 && vy < 0 {
		g.ball.SwitchY()
	}
	if h == sh-d && vy > 0 {
		g.ball.SwitchY()
	}

	// platform collisions
	p1w, p1h := g.players.P1.Coords()
	p2w, p2h := g.players.P2.Coords()

	if h+d > p1h && h < p1h+playerHeight && (p1w+playerWidth)/2 == w/2 {
		g.ball.SwitchX()
	}
	if h+d > p2h && h < p2h+playerHeight && p2w/2 == w/2+d {
		g.ball.SwitchX()
	}
}

// Draw overlay. Doesn't update the terminal
func (g *Game) drawOverlay() {
	w, h := g.screen.Size()

	// temporary copy-pasted color
	st := tcell.StyleDefault.Background(tcell.ColorGray).Foreground(tcell.ColorWhite)

	for i := 2; i < h; i += 4 {
		g.rect(w/2, w/2+1, i, i+3, ' ', st)
	}
}

// Draw specified player. Doesn't update the terminal
func (g *Game) drawPlayer(p Player) {
	pad, _ := p.Coords()

	_, yPos := p.Coords()

	// temporary copy-pasted color
	st := tcell.StyleDefault.Background(tcell.ColorGray).Foreground(tcell.ColorWhite)

	g.rect(pad, pad+p.GetWidth(), yPos, yPos+p.GetHeight(), ' ', st)
}

// Draws the ball. Doesn't update the terminal
func (g *Game) drawBall() {
	x, y := g.ball.Coords()
	d := g.ball.Diam()

	// temporary copy-pasted color
	st := tcell.StyleDefault.Background(tcell.ColorGray).Foreground(tcell.ColorWhite)

	// reason why x2 is x+2*d is because in terminal, height is 2x width so
	// this needed to be done for compensation of width
	g.rect(x, x+2*d, y, y+d, ' ', st)
}

// Draw a rectangle. Doesn't update the terminal
func (g *Game) rect(x1 int, x2 int, y1 int, y2 int, mainc rune, style tcell.Style) {
	for i := x1; i < x2; i++ {
		for j := y1; j < y2; j++ {
			// combc is in most cases nil
			g.screen.SetContent(i, j, mainc, nil, style)
		}
	}
}
