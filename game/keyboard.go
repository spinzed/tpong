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

// Listen for keyboard events and dispatch them through a channel.
// It will block so it must be called in a separate goroutine
func keyboardListen(k *keylogger.KeyLogger, c chan keyState) {
	kch := k.Read()

	for e := range kch {
		switch e.Type {
		case keylogger.EvKey:
			// I must make separate checks for KeyPress and KeyRelease because
			// it can happen that a key is held down continuously, in that case
			// both methods return false
			if e.KeyPress() {
				c <- keyState{e.KeyString(), true}
			}

			if e.KeyRelease() {
				c <- keyState{e.KeyString(), false}
			}
		}
	}
}
