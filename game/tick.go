package game

func (g *Game) Tick() {
	// update the game state every tick.

	if !g.started || g.paused || g.hardPaused {
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

	// update the screen.
	// these aren't expensive since they just check for changes on the canvas
	// and if there aren't any, nothing will be updated therefore no bloat
	g.screen.Clear()

	g.drawScores()
	g.drawOverlay()
	g.drawPlayers()

	// first start is for initial delay
	// move the ball for 1 tick
	g.ball.Move()
	g.checkCollision()
	g.drawBall()

	// at last, update the terminal
	g.screen.Show()
}
