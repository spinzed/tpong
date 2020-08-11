package game

import (
	"github.com/MarinX/keylogger"
)

type KeyDispatch struct {
	Name string
	Down bool
}

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

// Listen for keyboard events and dispatch them through a channel.
// It will block so it must be called in a separate goroutine
func keyboardListen(k *keylogger.KeyLogger, c chan KeyDispatch, dc chan KeyDispatch) {
	kch := k.Read()

	for {
		select {
		case e := <-kch:
			switch e.Type {
			case keylogger.EvKey:
				ev := getEvent(&keysGame, e.KeyString())

				switch ev.State {
				case stateClick:
					if e.KeyPress() {
						dc <- KeyDispatch{ev.Name, true}
					}
				case stateHold:
					if e.KeyPress() {
						c <- KeyDispatch{ev.Name, true}
					}
					if e.KeyRelease() {
						c <- KeyDispatch{ev.Name, false}
					}
				}
			}
		}
	}
}
