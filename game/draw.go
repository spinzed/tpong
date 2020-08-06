package game

import (
	"strconv"
	"strings"

	"github.com/gdamore/tcell"
)

// Draws the gui the updates the terminal.
func (g *Game) drawInTerminal() {
	g.screen.Clear()
	g.drawGUI()
	g.screen.Show()
}

// Draws all gui. Doesn't update the terminal
func (g *Game) drawGUI() {
	if g.theme.IsBgShown() {
		g.drawBackground()
	}

	g.drawOverlay()

	if g.started {
		g.drawScores()
	}

	g.drawPlatforms()

	if g.started {
		g.drawBall()
	}
	if g.paused {
		g.drawPauseText()
	}
	if !g.started {
		g.drawStartText()
	}
}

// Draw background. Doesn't update the terminal.
func (g *Game) drawBackground() {
	w, h := g.screen.Size()

	st := g.theme.GetCurrent().GetBackgroundStyle()

	g.rect(0, w, 0, h, ' ', st)
}

func (g *Game) drawOverlay() {
	w, h := g.screen.Size()

	st := g.theme.GetCurrent().GetOverlayStyle()

	// dashed line in the middle
	for i := 2; i < h; i += 4 {
		g.rect(w/2, w/2+1, i, i+3, ' ', st)
	}
}

func (g *Game) drawStartText() {
	_, h := g.screen.Size()

	text := []string{
		"  Space  - Start Game",
		"    P    - Pause Game",
		"    R    - Restart Round",
		"    T    - Switch Theme",
		"    B    - Toggle Background",
		"    Q    - Quit",
		" = P1 =",
		"    W    - Move Up,     S     - Move Down",
		" = P2 =",
		" ArrowUp - Move Up, ArrowDown - Move Down",
	}

	g.lines(0, h-len(text)-1, text)
}

func (g *Game) drawPauseText() {
	_, h := g.screen.Size()

	text := formatKeys([]string{
		eventTogglePause,
		eventReset,
		eventSwitchTheme,
		eventDestroy,
		eventP1Up,
		eventP1Down,
		eventP2Up,
		eventP2Down,
	})

	finalText := append([]string{" == PAUSE =="}, text...)

	g.lines(0, h-len(finalText)-1, finalText)
}

func formatKeys(eventNames []string) []string {
	var keys []string
	var descs []string

	var maxKeylen int

	alternate := map[string]string{
		"Up":   "ArrowUp",
		"Down": "ArrowDown",
	}

	for _, eventName := range eventNames {
		for key, event := range events {
			if event.Name == eventName {
				var realKey string

				if alternate[key] != "" {
					realKey = alternate[key]
				} else {
					realKey = key
				}
				keys = append(keys, realKey)
				descs = append(descs, event.Description)

				if len(key) > maxKeylen {
					maxKeylen = len(realKey)
				}
			}
		}
		for key, event := range dispEvents {
			if event.Name == eventName {
				var realKey string

				if alternate[key] != "" {
					realKey = alternate[key]
				} else {
					realKey = key
				}
				keys = append(keys, realKey)
				descs = append(descs, event.Description)

				if len(key) > maxKeylen {
					maxKeylen = len(realKey)
				}
			}
		}
	}

	maxKeylen += 2

	var final []string

	for i, key := range keys {
		desc := descs[i]

		thstring := strings.Repeat(" ", (maxKeylen-len(key))/2)
		thstring += key
		thstring += strings.Repeat(" ", maxKeylen-len(key)-((maxKeylen-len(key))/2))
		thstring += "- " + desc

		final = append(final, thstring)
	}

	return final
}

// Draw every player. Doesn't update the terminal.
func (g *Game) drawPlatforms() {
	for _, p := range g.players.GetAll() {
		g.drawPlatform(*p)
	}
}

// Draw specified player. Doesn't update the terminal
func (g *Game) drawPlatform(p Player) {
	pad, _ := p.Coords()

	_, yPos := p.Coords()

	st := g.theme.GetCurrent().GetPlatformStyle()

	g.rect(pad, pad+p.GetWidth(), yPos, yPos+p.GetHeight(), ' ', st)
}

// Draw the ball. Doesn't update the terminal
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
	padY := h/10 - letterHeight/2

	// this makes sure that padding isn't smaller than one
	if padY < 1 {
		padY = 1
	}

	st := g.theme.GetCurrent().GetOverlayStyle()

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
func (g *Game) lines(x int, y int, lines []string) {
	st := g.theme.GetCurrent().GetTextStyle()

	for i, line := range lines {
		g.text(x, y+i, line, st)
	}
}

// Draw text in a straight line. Doesn't update the terminal
func (g *Game) text(x int, y int, chars string, st tcell.Style) {
	for i, char := range chars {
		g.screen.SetContent(x+i, y, char, nil, st)
	}
}
