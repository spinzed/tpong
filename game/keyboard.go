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

func getEvent(m *map[string]Event, k string) Event {
	// this will be expanded in the future
	return (*m)[k]
}

// Listen for keyboard events and dispatch them through a channel.
// It will block so it must be called in a separate goroutine
func keyboardListen(k *keylogger.KeyLogger, c chan keyState, dc chan keyState) {
	kch := k.Read()

	for {
		select {
		case e := <-kch:
			switch e.Type {
			case keylogger.EvKey:
				// there are separate checks for KeyPress and KeyRelease because
				// it can happen that a key is held down continuously, in that case
				// both methods return false
				if ev := getEvent(&events, e.KeyString()); ev.Name != "" {
					if e.KeyPress() {
						c <- keyState{ev.Name, true}
					}

					if e.KeyRelease() {
						c <- keyState{ev.Name, false}
					}
				}
				if ev := getEvent(&dispEvents, e.KeyString()); ev.Name != "" {
					if e.KeyPress() {
						dc <- keyState{ev.Name, true}
					}
				}
			}
		}
	}
}
