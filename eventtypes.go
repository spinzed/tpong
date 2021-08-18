package main

// This file defines what action should dispatch on every key and
// events that will be dispatched depending on the state of the game

type KeyEventMap []KeyEvent

func (k *KeyEventMap) FindEventByKey(key Key) *Event {
	for _, keyEvent := range *k {
		if key.Name == keyEvent.Key.Name && key.State == keyEvent.Key.State {
			return &keyEvent.Event
		}
	}
	return nil
}

func (k *KeyEventMap) FindKeyEventsByKeyName(keyName string) *[]KeyEvent {
	events := make([]KeyEvent, 0)
	for _, keyEvent := range *k {
		if keyEvent.Key.Name == keyName {
			events = append(events, keyEvent)
		}
	}
	return &events
}

type KeyEvent struct {
	Key
	Event
}

func NewKeyEvent(key Key, event Event) *KeyEvent {
	return &KeyEvent{key, event}
}

type Key struct {
	Name string
	// State describes the key press type the event is triggered, but doesn't define the key.
	// valid states are click, release, hold (and holdEnd) or normal
	// holdEnd state is only used to dismiss the hold state
	State string
}

func NewKey(name string, state string) Key {
	return Key{name, state}
}

type StateEvent struct {
	Event
	State string
}

// Struct event defines an event which can be triggered by a key
type Event struct {
	Name        string
	Description string
}

func NewEvent(name string, desc string) Event {
	// status can be Pulse, Starting or Ending. It is initially not set
	// because it will be set when the event is dispatched
	return Event{name, desc}
}

type GameInfo struct {
	Keys       KeyEventMap
	LegendKeys KeyEventMap
	LegendType string
}
