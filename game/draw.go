package game

import (
	"strconv"
	"strings"

	"github.com/gdamore/tcell"
)

func (g *Game) drawOverlay() {
	w, h := g.screen.Size()

	// temporary copy-pasted color
	st := tcell.StyleDefault.Background(tcell.ColorGray).Foreground(tcell.ColorWhite)

	// dashed line in the middle
	for i := 2; i < h; i += 4 {
		g.rect(w/2, w/2+1, i, i+3, ' ', st)
	}
}

func (g *Game) drawStartText() {
	_, h := g.screen.Size()

	// temporary copy-pasted color
	st := tcell.StyleDefault.Background(tcell.ColorGray).Foreground(tcell.ColorWhite)

	text := [][]rune{
		[]rune("  Space  - Start Game"),
		[]rune("    P    - Pause Game"),
		[]rune("    R    - Restart Round"),
		[]rune("    Q    - Quit"),
		[]rune(" = P1 ="),
		[]rune("    W    - Move Up,     S     - Move Down"),
		[]rune(" = P2 ="),
		[]rune(" ArrowUp - Move Up, ArrowDown - Move Down"),
	}

	g.lines(0, h-len(text)-1, text, st)
}

func (g *Game) drawPauseText() {
	_, h := g.screen.Size()

	// temporary copy-pasted color
	st := tcell.StyleDefault.Background(tcell.ColorGray).Foreground(tcell.ColorWhite)

	text := [][]rune{
		[]rune(" == PAUSE =="),
		[]rune("    P    - Unpause Game"),
		[]rune("    R    - Restart Round"),
		[]rune("    Q    - Quit"),
		[]rune(" = P1 ="),
		[]rune("   W     - Move Up,     S     - Move Down"),
		[]rune(" = P2 ="),
		[]rune(" ArrowUp - Move Up, ArrowDown - Move Down"),
	}

	g.lines(0, h-len(text)-1, text, st)
}

func (g *Game) drawPlayers() {
	for _, p := range g.players.GetAll() {
		g.drawPlayer(*p)
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

	// set the top padding
	padY := h/10 - letterHeight/2

	// this makes sure that padding isn't smaller than one
	if padY < 1 {
		padY = 1
	}

	// temporary copy-pasted color
	st := tcell.StyleDefault.Background(tcell.ColorGray).Foreground(tcell.ColorWhite)

	for _, p := range g.players.GetAll() {
		// set the middle point of the letter
		mid := w / 4

		if p.GetTag() == playerP2 {
			mid *= 3
		}

		// get the score and parse it in slice of string digits - 25 -> { "2", "5" }
		score := p.GetScore()
		parsedNums := strings.Split(strconv.Itoa(score), "")

		// get the total length of all letters with spacing
		totalXLen := len(parsedNums)*letterWidth + (len(parsedNums)-1)*letterSpacing

		// starting point - subtract half of the full length of all letters
		start := mid - totalXLen/2

		for i, strchar := range parsedNums {
			// current letter x offset
			xOffset := start + totalXLen/len(parsedNums)*i

			// add between-letter spacing if necessary
			if i != 0 && i < len(parsedNums) {
				xOffset += i * letterSpacing
			}

			// cast the string-char back to int
			parsedNum, err := strconv.Atoi(strchar)

			if err != nil {
				panic(err)
			}

			// get cells for that number
			readyNum := getCellsFromNum(parsedNum)

			// draw each number cell
			for _, char := range readyNum {
				x := xOffset + int(char[0])
				y := padY + int(char[1])
				g.screen.SetContent(x, y, ' ', nil, st)
			}
		}
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

// Draw multiple lines of text. Doesn't update the terminal
func (g *Game) lines(x int, y int, lines [][]rune, st tcell.Style) {
	for i, line := range lines {
		g.text(x, y+i, line, st)
	}
}

// Draw text in a straight line. Doesn't update the terminal
func (g *Game) text(x int, y int, chars []rune, st tcell.Style) {
	for i, char := range chars {
		g.screen.SetContent(x+i, y, char, nil, st)
	}
}
