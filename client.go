package main

import (
	"github.com/gdamore/tcell"
)

type Client interface {
	Name() string
	RegisterEvent(StateEvent)
	SendUpdate(StateEvent)
}

// Local non-AI player
type LocalClient struct {
	screen     tcell.Screen
	serverChan chan StateEvent
	gameChan   chan StateEvent
}

func NewLocalClient(s tcell.Screen, srv chan StateEvent, gm chan StateEvent) *LocalClient {
	return &LocalClient{s, srv, gm}
}

// Unnecesarry, will probs be removed or replaced
func (c *LocalClient) Name() string {
	return "LocalClient"
}

// Register the event that comes from a server/game controller
func (c *LocalClient) RegisterEvent(evt StateEvent) {
	c.gameChan <- evt
}

// Send the update to server
func (c *LocalClient) SendUpdate(evt StateEvent) {
	c.serverChan <- evt
}

type RemoteClient struct {
}
