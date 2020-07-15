package game

import (
	"fmt"

	"github.com/gdamore/tcell"
)

// event listener for the screen
// must be called in separate goroutine since it is blocking
func inputLoop(s tcell.Screen, c chan string) {
	for {
		e := s.PollEvent()

		switch e := e.(type) {
		case *tcell.EventKey:
			keyMap := map[rune]string{
				'w': eventP1Up,
				's': eventP1Down,
			}

			keyCodeMap := map[tcell.Key]string{
				tcell.KeyUp:   eventP2Up,
				tcell.KeyDown: eventP2Down,
				tcell.KeyCtrlC: eventDestroy,
			}

			// end the loop if CtrlC is pressed.
			// this is a temporary solution until I make the event loop.
			// event loop must end the input loop, not the other way round
			if keyCodeMap[e.Key()] == eventDestroy {
				return
			}

			if keyMap[e.Rune()] != "" {
				fmt.Println("test")
			} else if keyCodeMap[e.Key()] != "" {
				fmt.Println("heyyy")
			}

			fmt.Println(e.Key())
			fmt.Println(e.Rune())
		}
	}
}
