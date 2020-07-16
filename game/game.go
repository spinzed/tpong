package game

import (
	"time"

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
	ticker  *time.Ticker
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
	// screen init
	s, err := initScreen()

	if err != nil {
		return err
	}

	g.screen = s

	// event channel init
	g.event = make(chan string)

	// ticker init
	g.ticker = time.NewTicker(1000000 / framerate * time.Microsecond)

	//players init
	_, h := g.screen.Size()
	g.players = newPlayers((h - playerHeight) / 2)

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
	// start the input loop
	go inputLoop(g.screen, g.event)

	for {
		select {
		case <-g.ticker.C:
			// update the screen every tick.
			// this isn't expensive since this just checks for changes on the canvas
			// and if there aren't any, nothing will be updated therefore no bloat
			g.screen.Clear()

			g.drawPlayer(g.players.P1)
			g.drawPlayer(g.players.P2)
			g.drawOverlay()

			g.screen.Show()
		case lol := <-g.event:
			switch lol {
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
	}

	g.screen.Show()
}

// End the game
func (g *Game) End() {
	// Release the screen resources
	g.screen.Fini()
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
