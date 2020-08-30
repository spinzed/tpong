package main

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

	var allKbs []*keylogger.KeyLogger

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
// Since the keys map is constantly changing, instead of directly passing the keys,
// function that returns keys at execution is passed.
func keyboardListen(k *keylogger.KeyLogger, ky func() *map[Key]Event, c chan KeyEvent) {
	kch := k.Read()

	for {
		select {
		case e := <-kch:
			switch e.Type {
			case keylogger.EvKey:
				key := newKey(e.KeyString(), stateNormal)

				if ev := getEvent(ky(), key); !e.KeyRelease() && ev.Name != "" {
					c <- KeyEvent{key, ev}
				}
				if e.KeyPress() {
					key.State = stateClick
					if ev := getEvent(ky(), key); ev.Name != "" {
						c <- KeyEvent{key, ev}
					}

					key.State = stateHold
					if ev := getEvent(ky(), key); ev.Name != "" {
						c <- KeyEvent{key, ev}
					}
				}
				if e.KeyRelease() {
					key.State = stateRelease
					if ev := getEvent(ky(), key); ev.Name != "" {
						c <- KeyEvent{key, ev}
					}

					key.State = stateHold
					if ev := getEvent(ky(), key); ev.Name != "" {
						// stateHoldEnd is only used when signal needs to be dispatched
						// that stateHold action must be dismissed
						key.State = stateHoldEnd
						c <- KeyEvent{key, ev}
					}
				}
			}
		}
	}
}
