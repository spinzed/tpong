package main

import (
	"github.com/MarinX/keylogger"
)

// InputHandler handles input and dispatches game events based on
// the KeyEvent array
type InputHandler struct {
	loggers  *[]*keylogger.KeyLogger
	eventMap func() *KeyEventMap
}

// Initialize a new input handler
func NewInputHandler(f func() *KeyEventMap) (*InputHandler, error) {
	kbs := keylogger.FindAllKeyboardDevices()

	if len(kbs) < 1 {
		return nil, NoKeyboardError
	}

	var allKbs []*keylogger.KeyLogger

	for _, kb := range kbs {
		k, err := keylogger.New(kb)

		if err != nil {
			return nil, err
		}
		allKbs = append(allKbs, k)
	}

	return &InputHandler{&allKbs, f}, nil
}

// Listen on all input devices and pass events through one channel.
func (i *InputHandler) Listen(c chan<- StateEvent) {
	for _, l := range *i.loggers {
		inputListen(l, i.eventMap, c)
	}
}

// Listen for keyboard events and dispatch them through a channel.
// It will block so it must be called in a separate goroutine
// Since the keys map is constantly changing, instead of directly passing the keys,
// function that returns keys at execution is passed.
func inputListen(k *keylogger.KeyLogger, f func() *KeyEventMap, c chan<- StateEvent) {
	kch := k.Read()

	for {
		e := <-kch

		if e.Type != keylogger.EvKey {
			continue
		}

		// find every event which is bound to this key
		keyEvents := f().FindKeyEventsByKeyName(e.KeyString())

		// find which event type is bound to that event and dispatch it
		for _, keyEvent := range *keyEvents {
			evt := StateEvent{keyEvent.Event, eventStatePulse}

			if keyEvent.Key.State == stateNormal && !e.KeyRelease() {
				evt.State = eventStatePulse
				c <- evt
			}

			if e.KeyPress() && keyEvent.Key.State == stateClick {
				evt.State = eventStatePulse
				c <- evt
			}

			if e.KeyPress() && keyEvent.Key.State == stateHold {
				evt.State = eventStateStarting
				c <- evt
			}

			if e.KeyRelease() && keyEvent.Key.State == stateHold {
				evt.State = eventStateEnding
				c <- evt
			}
		}
	}
}
