package game

func (g *Game) Tick(keys *[]string) {
	// update the game state every tick.

	if g.paused || g.hardPaused {
		return
	}

	// keys that don't persist when game is paused
	for _, key := range *keys {
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

	g.drawPlayer(*g.players.P1)
	g.drawPlayer(*g.players.P2)
	g.drawBall()
	g.drawScores()
	g.drawOverlay()

	g.screen.Show()

	// move the ball for 1 tick
	g.checkCollision()
	g.ball.Move()
}
