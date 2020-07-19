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
	P1 Player
	P2 Player
}

type Game struct {
	screen   tcell.Screen
	players  Players
	event    chan keyState
	ticker   *time.Ticker
	keyboard *keylogger.KeyLogger
	oldstate *terminal.State
}

// Get a pair of players ready and initialised
func newPlayers(initialPos int) Players {
	p1 := newPlayer("p1", initialPos)
	p2 := newPlayer("p2", initialPos)

	return Players{p1, p2}
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
	_, h := g.screen.Size()
	g.players = newPlayers((h - playerHeight) / 2)

	// event channel init
	g.event = make(chan keyState)

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
	g.drawPlayer(g.players.P1)
	g.drawPlayer(g.players.P2)
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
	go keyboardListen(g.keyboard, g.event)

	for {
		select {
		case <-g.ticker.C:
			// update the screen every tick.

			for _, key := range keys {
				switch key {
				case eventDestroy:
					return
				case eventP1Up:
					g.MovePlayerUp(&g.players.P1)
				case eventP1Down:
					g.MovePlayerDown(&g.players.P1)
				case eventP2Up:
					g.MovePlayerUp(&g.players.P2)
				case eventP2Down:
					g.MovePlayerDown(&g.players.P2)
				}
			}

			// these aren't expensive since they just check for changes on the canvas
			// and if there aren't any, nothing will be updated therefore no bloat
			g.screen.Clear()

			g.drawPlayer(g.players.P1)
			g.drawPlayer(g.players.P2)
			g.drawOverlay()

			g.screen.Show()
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
		}
	}
}

// End the game
func (g *Game) End() {
	// clear the screen. The second line connects the command's stdout with
	// the of one of this terminal's session
	clr := exec.Command("clear")
	clr.Stdout = os.Stdout
	clr.Run()

	// bring teminal back from raw mode
	terminal.Restore(0, g.oldstate)
}

// Move player 1 char higher. If player is at the edge, do nothing.
func (g *Game) MovePlayerUp(p *Player) {
	if p.GetYPos() < 1 {
		return
	}

	p.GoUp()
}

// Move player 1 char lower. If player is at the edge, do nothing.
func (g *Game) MovePlayerDown(p *Player) {
	_, h := g.screen.Size()

	if p.GetYPos() > h-p.GetHeight()-1 {
		return
	}

	p.GoDown()
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
	w, _ := g.screen.Size()
	padding := 4

	var leftpad int

	if p.GetTag() == "p1" {
		leftpad = padding
	} else {
		leftpad = w - padding
	}

	yPos := p.GetYPos()

	// temporary copy-pasted color
	st := tcell.StyleDefault.Background(tcell.ColorGray).Foreground(tcell.ColorWhite)

	g.rect(leftpad, leftpad+p.GetWidth(), yPos, yPos+p.GetHeight(), ' ', st)
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
