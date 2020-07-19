package game

import (
	"github.com/MarinX/keylogger"
)

type noKeyboardError struct {
	s string
}

func (e *noKeyboardError) Error() string {
	return e.s
}

// Returns pointer to a new error struct instance
func newNoKeyboardError(t string) error {
	return &noKeyboardError{t}
}

// Initialize a new keylogger
func newKeyboard() (*keylogger.KeyLogger, error) {
	kb := keylogger.FindKeyboardDevice()

	if kb == "" {
		return nil, newNoKeyboardError("No keyboard detected")
	}

	k, err := keylogger.New(kb)

	if err != nil {
		return nil, err
	}

	return k, nil
}

type keyState struct {
	Name string
	Down bool
}

func getEvent(m *map[string]string, k string) string {
	return (*m)[k]
}

// Listen for keyboard events and dispatch them through a channel.
// It will block so it must be called in a separate goroutine
func keyboardListen(k *keylogger.KeyLogger, c chan keyState) {
	kch := k.Read()

	events := map[string]string{
		"Q":    eventDestroy,
		"W":    eventP1Up,
		"S":    eventP1Down,
		"Up":   eventP2Up,
		"Down": eventP2Down,
	}

	for {
		select {
		case e := <-kch:
			switch e.Type {
			case keylogger.EvKey:
				// there are separate checks for KeyPress and KeyRelease because
				// it can happen that a key is held down continuously, in that case
				// both methods return false
				if ev := getEvent(&events, e.KeyString()); ev != "" {
					if e.KeyPress() {
						c <- keyState{ev, true}
					}

					if e.KeyRelease() {
						c <- keyState{ev, false}
					}
				}
			}
		}
	}
}
