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

func newPlayers(initialPos int) Players {
	p1 := newPlayer("p1", initialPos)
	p2 := newPlayer("p2", initialPos)

	return Players{p1, p2}
}

// Creates screen ready to use.
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

func (g *Game) Start() {
	//defer g.screen.Fini()

	g.screen.Show()
}

func (g *Game) Kill() {
	g.screen.Fini()
}

func (g *Game) drawOverlay() {
	w, h := g.screen.Size()

	st := tcell.StyleDefault.Background(tcell.ColorGray).Foreground(tcell.ColorWhite)

	for i := 2; i < h; i += 4 {
		g.rect(w/2, w/2+1, i, i+3, ' ', st)
	}
}

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

	st := tcell.StyleDefault.Background(tcell.ColorGray).Foreground(tcell.ColorWhite)

	g.rect(leftpad, leftpad+p.GetWidth(), yPos, yPos+p.GetHeight(), ' ', st)
}

// draws a rectangle
func (g *Game) rect(x1 int, x2 int, y1 int, y2 int, mainc rune, style tcell.Style) {
	for i := x1; i < x2; i++ {
		for j := y1; j < y2; j++ {
			// combc is in most cases nil
			g.screen.SetContent(i, j, mainc, nil, style)
		}
	}
}
