package game

import (
	"github.com/gdamore/tcell"
)

type Players struct {
	P1 Player
	P2 Player
}

type Game struct {
	screen  tcell.Screen
	players Players
	event   chan string
}

// Get a pair of players rady and initialised
func newPlayers(initialPos int) Players {
	p1 := newPlayer("p1", initialPos)
	p2 := newPlayer("p2", initialPos)

	return Players{p1, p2}
}

// Creates screen ready to use
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

// Initialise the game. Must be called on new game instance
func (g *Game) Init() error {
	s, err := initScreen()

	if err != nil {
		return err
	}

	g.screen = s

	_, h := g.screen.Size()
	g.players = newPlayers((h - playerHeight) / 2)

	g.drawPlayer(g.players.P1)
	g.drawPlayer(g.players.P2)
	g.drawOverlay()

	g.screen.Show()

	// signal that everything is ok
	return nil
}

// Starts the game. Must be called after initialisation
func (g *Game) Loop() {
	//defer g.screen.Fini()

	// this starts the input loop. Until the event loop is implemented, it is for now
	// synchronous, but it the future it should be in its own goroutine.
	inputLoop(g.screen, g.event)

	g.screen.Show()
}

// End the game
func (g *Game) End() {
	// Release the screen resources
	g.screen.Fini()
}


// Draws overlay. Doesn't update the terminal
func (g *Game) drawOverlay() {
	w, h := g.screen.Size()

	// temporary copy-pasted color
	st := tcell.StyleDefault.Background(tcell.ColorGray).Foreground(tcell.ColorWhite)

	for i := 2; i < h; i += 4 {
		g.rect(w/2, w/2+1, i, i+3, ' ', st)
	}
}


// Draws specified player. Doesn't update the terminal
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

// Draws a rectangle. Doesn't update the terminal
func (g *Game) rect(x1 int, x2 int, y1 int, y2 int, mainc rune, style tcell.Style) {
	for i := x1; i < x2; i++ {
		for j := y1; j < y2; j++ {
			// combc is in most cases nil
			g.screen.SetContent(i, j, mainc, nil, style)
		}
	}
}
