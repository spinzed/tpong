package game

// Update the start menu state every tick
func (g *Game) PerformStartMenuTick() {
	g.ball.Move()
	g.checkCollision(true)

	g.drawStartGameMenu()
}

// Update the game state every tick.
func (g *Game) PerformGameTick() {
	// update the screen.
	// this isn't expensive since it just checks for changes on the canvas
	// and if there aren't any, nothing will be updated therefore no bloat.
	// the terminal is updated after everything has been drawn.
	// this should be called after every tick, thats why it's defered
	defer g.drawGameTick()

	if g.paused || g.hardPaused {
		return
	}

	// keys that don't persist when game is paused
	for _, key := range *g.keys {
		switch key {
		case eventP1Up:
			g.movePlayerUp(g.players.P1)
		case eventP1Down:
			g.movePlayerDown(g.players.P1)
		case eventP2Up:
			g.movePlayerUp(g.players.P2)
		case eventP2Down:
			g.movePlayerDown(g.players.P2)
		}
	}

	// move the ball for 1 tick
	g.ball.Move()
	g.checkCollision(false)
}
