package main

import "time"

// Check and handle collisions between balls, walls and platforms
func (g *Game) checkCollision(ballBouncesSides bool) {
	w, h := g.ui.ball.Coords()
	d := g.ui.ball.Diam()
	vx, vy := g.ui.ball.Vels()
	sw, sh := g.ui.screen.Size()

	// wall collisions. If ball side bouncing is enabled, bounce it,
	// if it is not... do not bounce it
	if ballBouncesSides {
		if w < 1 && vx < 0 {
			g.ui.ball.SwitchX()
		}
		if w >= sw-2*d-1 && vx > 0 {
			g.ui.ball.SwitchX()
		}
	} else {
		var p *Player
		if w < 1 && vx < 0 {
			p = g.players.P2
		}
		if w >= sw-2*d-1 && vx > 0 {
			p = g.players.P1
		}
		if p != nil {
			p.AddPoint()

			// hardPause cannot be unpaused by player
			g.toggleHardPause() // set to true

			go func() {
				time.Sleep(scoreSleepSecs * time.Second)
				g.toggleHardPause() // set to false
				g.Reset()
			}()
		}
	}
	if h <= 1 && vy < 0 {
		g.ui.ball.SwitchY()
	}
	if h >= sh-d && vy > 0 {
		g.ui.ball.SwitchY()
	}

	// platform collisions
	if !ballBouncesSides {
		p1w, p1h := g.players.P1.Coords()
		p2w, p2h := g.players.P2.Coords()

		// left side
		if h+d > p1h && h < p1h+platformHeight && (p1w+platformWidth)/2 == w/2 {
			g.ui.ball.SwitchX()
		}
		// right side
		if h+d > p2h && h < p2h+platformHeight && p2w/2 == w/2+d {
			g.ui.ball.SwitchX()
		}
	}
}
