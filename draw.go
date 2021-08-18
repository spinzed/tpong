package main

import (
	"strconv"

	"github.com/gdamore/tcell"
)

// Update the terminal.
func (g *Game) updateTerminal() {
	g.ui.Paint()
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
	g.ui.Reset()

	if g.ui.theme.IsBgShown() {
		g.ui.drawBackground()
	}

	g.ui.drawOverlay()
	g.ui.drawScores(*g.players)
	g.ui.drawPlatforms(*g.players)
	g.ui.drawBall()

	if g.state == gameStatePaused {
		g.drawLegend()
	}
}

// Draws the start menu. Doesn't update the terminal.
func (g *Game) drawStartGameMenu() {
	w, h := g.ui.screen.Size()

	g.ui.screen.Clear()

	if g.ui.theme.IsBgShown() {
		g.ui.drawBackground()
	}

	g.ui.drawBall()

	text := "PONG"
	st := g.ui.theme.GetCurrent().GetOverlayStyle()
	g.ui.drawLetters(w/2, h/5*2, startMenuLttrGap, text, st)

	g.drawLegend()
}

// Draws the legend. Doesn't update the terminal.
func (g *Game) drawLegend() {
	g.ui.drawLegend(g.state, g.aiActive)
}

// Draw background. Doesn't update the terminal.
func (ui *UIHandler) drawBackground() {
	w, h := ui.screen.Size()

	st := ui.theme.GetCurrent().GetBackgroundStyle()

	ui.rect(0, w, 0, h, ' ', st)
}

// Draw tghe overlay aka the dashed line. Doesn't update the terminal.
// The score counter may also be added here in the future.
func (ui *UIHandler) drawOverlay() {
	w, h := ui.screen.Size()

	st := ui.theme.GetCurrent().GetOverlayStyle()

	// dashed line in the middle
	for i := 2; i < h; i += 4 {
		ui.rect(w/2, w/2+1, i, i+3, ' ', st)
	}
}

func (ui *UIHandler) drawLegend(state GameState, aiActive bool) {
	w, h := ui.screen.Size()
	legendType := screenData[state].LegendType
	data := screenData[state].LegendKeys

	if legendType == "" || len(data) == 0 {
		return
	}

	x := w / 2
	if legendType != "middle" {
		x = 0
	}

	text := ui.formatKeys(legendType, data, aiActive)
	y := h - len(text) - 1

	if state == gameStatePaused {
		ui.lines(x, y-1, []string{" == PAUSE == "}, legendType)
	}
	ui.lines(x, y, text, legendType)
}

// Draw every player. Doesn't update the terminal.
func (ui *UIHandler) drawPlatforms(players Players) {
	for _, p := range players.GetAll() {
		ui.drawPlatform(*p)
	}
}

// Draw specified player. Doesn't update the terminal.
func (ui *UIHandler) drawPlatform(p Player) {
	pad, _ := p.Coords()
	_, yPos := p.Coords()

	st := ui.theme.GetCurrent().GetPlatformStyle()
	ui.rect(pad, pad+p.GetWidth(), yPos, yPos+p.GetHeight(), ' ', st)
}

// Draw the ball. Doesn't update the terminal.
func (ui *UIHandler) drawBall() {
	x, y := ui.ball.Coords()
	d := ui.ball.Diam()

	st := ui.theme.GetCurrent().GetBallStyle()

	// reason why x2 is x+2*d is because in terminal, height is 2x width so
	// this needed to be done for compensation of width
	ui.rect(x, x+2*d, y, y+d, ' ', st)
}

// Draw 7-segment-like scores on the top of the screen. Doesn't update the terminal.
// This looks ugly af. I may refractor it somehow later.
func (ui *UIHandler) drawScores(players Players) {
	w, h := ui.screen.Size()

	// set the top padding
	padY := h / 10

	// this makes sure that padding isn't smaller than one
	if padY < letterHeight/2 {
		padY = letterHeight / 2
	}

	st := ui.theme.GetCurrent().GetOverlayStyle()

	for _, p := range players.GetAll() {
		// set the middle point of the letter
		padX := w / 4

		if p.GetTag() == playerP2 {
			padX *= 3
		}

		ui.drawLetters(padX, padY, letterGap, strconv.Itoa(p.GetScore()), st)
	}
}

// Draw big letters with the given center.
func (ui *UIHandler) drawLetters(x int, y int, gap int, word string, st tcell.Style) {
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
			ui.screen.SetContent(finalx, finaly, ' ', nil, st)
		}
	}
}

// Draw a rectangle. Doesn't update the terminal.
func (ui *UIHandler) rect(x1 int, x2 int, y1 int, y2 int, mainc rune, style tcell.Style) {
	for i := x1; i < x2; i++ {
		for j := y1; j < y2; j++ {
			// combc is in most cases nil
			ui.screen.SetContent(i, j, mainc, nil, style)
		}
	}
}

// Draw multiple lines of text. Doesn't update the terminal.
// Mode dictates should the lines be aligned to the left or centered
func (ui *UIHandler) lines(x int, y int, lines []string, mode string) {
	st := ui.theme.GetCurrent().GetTextStyle()

	if mode != "left" && mode != "middle" {
		panic("Invalid line draw mode: " + mode)
	}

	for i, line := range lines {
		realx := x

		if mode == "middle" {
			realx -= len(line) / 2
		}
		ui.text(realx, y+i, line, st)
	}
}

// Draw text in a straight line. Doesn't update the terminal.
func (ui *UIHandler) text(x int, y int, chars string, st tcell.Style) {
	for i, char := range chars {
		ui.screen.SetContent(x+i, y, char, nil, st)
	}
}
