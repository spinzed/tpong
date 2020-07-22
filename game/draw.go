package game

import (
	"github.com/gdamore/tcell"
)

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

// Draws 7-segment-like scores on the top of the screen.
// This looks ugly af. I may refractor it somehow later.
func (g *Game) drawScores() {
	w, h := g.screen.Size()

	sch := h/10
	scw1 := w/4
	scw2 := scw1*3

	score1 := g.players.P1.GetScore()
	score2 := g.players.P2.GetScore()

	num1 := getCellsFromNum(score1)
	num2 := getCellsFromNum(score2)

	st := tcell.StyleDefault.Background(tcell.ColorGray).Foreground(tcell.ColorWhite)

	for _, char := range num1 {
		x := scw1 + int(char[0])
		y := sch + int(char[1])
		g.screen.SetContent(x, y, ' ', nil, st)
	}

	for _, char := range num2 {
		x := scw2 + int(char[0])
		y := sch + int(char[1])
		g.screen.SetContent(x, y, ' ', nil, st)
	}
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
