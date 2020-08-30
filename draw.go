package main

import (
	"strconv"

	"github.com/gdamore/tcell"
)

// Update the terinal.
func (g *Game) updateTerminal() {
	g.screen.Show()
}

// Draw the GUI aka the current screen. Doesn't update the terminal.
func (g *Game) drawGUI() {
	if g.started {
		g.PerformGameTick()
	} else {
		g.PerformStartMenuTick()
	}
}

// Draw the GUI for the current game tick. Doesn't update the terminal.
func (g *Game) drawGameTick() {
	g.screen.Clear()

	if g.theme.IsBgShown() {
		g.drawBackground()
	}

	g.drawOverlay()
	g.drawScores()
	g.drawPlatforms()
	g.drawBall()

	if g.paused {
		g.drawLegend()
	}
}

// Draws the start menu. Doesn't update the terminal.
func (g *Game) drawStartGameMenu() {
	w, h := g.screen.Size()

	g.screen.Clear()

	if g.theme.IsBgShown() {
		g.drawBackground()
	}

	g.drawBall()

	text := "PONG"
	st := g.theme.GetCurrent().GetOverlayStyle()
	g.drawLetters(w/2, h/5*2, startMenuLttrGap, text, st)

	g.drawLegend()
}

// Draw background. Doesn't update the terminal.
func (g *Game) drawBackground() {
	w, h := g.screen.Size()

	st := g.theme.GetCurrent().GetBackgroundStyle()

	g.rect(0, w, 0, h, ' ', st)
}

// Draw tghe overlay aka the dashed line. Doesn't update the terminal.
// The score counter may also be added here in the future.
func (g *Game) drawOverlay() {
	w, h := g.screen.Size()

	st := g.theme.GetCurrent().GetOverlayStyle()

	// dashed line in the middle
	for i := 2; i < h; i += 4 {
		g.rect(w/2, w/2+1, i, i+3, ' ', st)
	}
}

func (g *Game) drawLegend() {
	w, h := g.screen.Size()
	legend := g.keyData.Legend

	x := w / 2
	if legend.Type != "middle" {
		x = 0
	}

	text := g.keyData.formatKeys()
	y := h - len(text) - 1

	g.lines(x, y, text, legend.Type)
}

// Draw every player. Doesn't update the terminal.
func (g *Game) drawPlatforms() {
	for _, p := range g.players.GetAll() {
		g.drawPlatform(*p)
	}
}

// Draw specified player. Doesn't update the terminal.
func (g *Game) drawPlatform(p Player) {
	pad, _ := p.Coords()
	_, yPos := p.Coords()

	st := g.theme.GetCurrent().GetPlatformStyle()
	g.rect(pad, pad+p.GetWidth(), yPos, yPos+p.GetHeight(), ' ', st)
}

// Draw the ball. Doesn't update the terminal.
func (g *Game) drawBall() {
	x, y := g.ball.Coords()
	d := g.ball.Diam()

	st := g.theme.GetCurrent().GetBallStyle()

	// reason why x2 is x+2*d is because in terminal, height is 2x width so
	// this needed to be done for compensation of width
	g.rect(x, x+2*d, y, y+d, ' ', st)
}

// Draw 7-segment-like scores on the top of the screen. Doesn't update the terminal.
// This looks ugly af. I may refractor it somehow later.
func (g *Game) drawScores() {
	w, h := g.screen.Size()

	// set the top padding
	padY := h / 10

	// this makes sure that padding isn't smaller than one
	if padY < letterHeight/2 {
		padY = letterHeight / 2
	}

	st := g.theme.GetCurrent().GetOverlayStyle()

	for _, p := range g.players.GetAll() {
		// set the middle point of the letter
		padX := w / 4

		if p.GetTag() == playerP2 {
			padX *= 3
		}

		g.drawLetters(padX, padY, letterGap, strconv.Itoa(p.GetScore()), st)
	}
}

// Draw big letters with the given center.
func (g *Game) drawLetters(x int, y int, gap int, word string, st tcell.Style) {
	letterNum := len(word)
	totalLen := letterNum*letterWidth + (letterNum-1)*gap
	startx := x - totalLen/2
	offsety := y - letterHeight/2

	for i, letter := range word {
		letterCells := getCellsFromChar(string(letter))
		offsetx := startx + totalLen/letterNum*i

		offsetx += (i - 1) * gap

		for _, cell := range letterCells {
			finalx := offsetx + int(cell[0])
			finaly := offsety + int(cell[1])
			g.screen.SetContent(finalx, finaly, ' ', nil, st)
		}
	}
}

// Draw a rectangle. Doesn't update the terminal.
func (g *Game) rect(x1 int, x2 int, y1 int, y2 int, mainc rune, style tcell.Style) {
	for i := x1; i < x2; i++ {
		for j := y1; j < y2; j++ {
			// combc is in most cases nil
			g.screen.SetContent(i, j, mainc, nil, style)
		}
	}
}

// Draw multiple lines of text. Doesn't update the terminal.
// Mode dictates should the lines be aligned to the left or centered
func (g *Game) lines(x int, y int, lines []string, mode string) {
	st := g.theme.GetCurrent().GetTextStyle()

	if mode != "left" && mode != "middle" {
		panic("Invalid line draw mode: " + mode)
	}

	for i, line := range lines {
		realx := x

		if mode == "middle" {
			realx -= len(line) / 2
		}
		g.text(realx, y+i, line, st)
	}
}

// Draw text in a straight line. Doesn't update the terminal.
func (g *Game) text(x int, y int, chars string, st tcell.Style) {
	for i, char := range chars {
		g.screen.SetContent(x+i, y, char, nil, st)
	}
}
