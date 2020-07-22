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
func newKeyboard() ([]*keylogger.KeyLogger, error) {
	// I implemented this (:b
	kbs := keylogger.FindAllKeyboardDevices()

	if len(kbs) < 1 {
		return nil, newNoKeyboardError("No keyboard detected")
	}

	allKbs := make([]*keylogger.KeyLogger, 0)

	for _, kb := range kbs {
		k, err := keylogger.New(kb)

		if err != nil {
			return nil, err
		}

		allKbs = append(allKbs, k)
	}

	return allKbs, nil
}

type keyState struct {
	Name string
	Down bool
}

func getEvent(m *map[string]string, k string) string {
	return (*m)[k]
}

// fucking hell Go when will you have generics
func getDispEvent(m *map[string]keyState, k string) keyState {
	return (*m)[k]
}

// Listen for keyboard events and dispatch them through a channel.
// It will block so it must be called in a separate goroutine
func keyboardListen(k *keylogger.KeyLogger, c chan keyState, dc chan keyState) {
	kch := k.Read()

	events := map[string]string{
		"W":    eventP1Up,
		"S":    eventP1Down,
		"Up":   eventP2Up,
		"Down": eventP2Down,
	}

	dispEvents := map[string]keyState{
		"Q": {eventDestroy, true},
		"P": {eventTogglePause, true},
		"R": {eventReset, true},
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
				if ev := getDispEvent(&dispEvents, e.KeyString()); ev != (keyState{}) {
					if e.KeyPress() && ev.Down {
						dc <- ev
					}

					if e.KeyRelease() && !ev.Down {
						dc <- ev
					}
				}
			}
		}
	}
}