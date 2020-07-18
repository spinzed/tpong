package game

import (
	"github.com/gdamore/tcell"
)

// Input listener for the screen
// Must be called in separate goroutine since it is blocking
func InputLoopOld(s tcell.Screen, c chan<- string) {
	for {
		e := s.PollEvent()

		switch e := e.(type) {
		case *tcell.EventKey:
			keyMap := map[rune]string{
				'w': eventP1Up,
				's': eventP1Down,
			}

			keyCodeMap := map[tcell.Key]string{
				tcell.KeyUp:    eventP2Up,
				tcell.KeyDown:  eventP2Down,
				tcell.KeyCtrlC: eventDestroy,
			}

			if m := keyMap[e.Rune()]; m != "" {
				c <- m
			} else if k := keyCodeMap[e.Key()]; k != "" {
				c <- k
			}

			// since s.PollEvent is not a channel but a function, I cannot make
			// a select statement with it and a channel which will listen for
			// destroy event. I am still looking for a workaround because ending
			// the input loop this way may have unintended consequences
			if keyCodeMap[e.Key()] == eventDestroy {
				return
			}
		}
	}
}
