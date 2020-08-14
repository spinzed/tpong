package game

import "time"

// Check and handle collisions between balls, walls and platforms
func (g *Game) checkCollision(ballBouncesSides bool) {
	w, h := g.ball.Coords()
	d := g.ball.Diam()
	vx, vy := g.ball.Vels()
	sw, sh := g.screen.Size()

	// wall collisions. If ball side bouncing is enabled, bounce it,
	// if it is not... do not bounce it
	if ballBouncesSides {
		if w < 1 && vx < 0 {
			g.ball.SwitchX()
		}
		if w >= sw-2*d-1 && vx > 0 {
			g.ball.SwitchX()
		}
	} else {
		if w < 1 && vx < 0 {
			g.players.P2.AddPoint()

			// hardPause cannot be unpaused by player
			g.hardPaused = true

			go func() {
				time.Sleep(scoreSleepSecs * time.Second)
				g.hardPaused = false
				g.Reset()
			}()
		}
		if w >= sw-2*d-1 && vx > 0 {
			g.players.P1.AddPoint()

			// hardPause cannot be unpaused by player
			g.hardPaused = true

			go func() {
				time.Sleep(scoreSleepSecs * time.Second)
				g.hardPaused = false
				g.Reset()
			}()
		}
	}
	if h <= 1 && vy < 0 {
		g.ball.SwitchY()
	}
	if h >= sh-d && vy > 0 {
		g.ball.SwitchY()
	}

	// platform collisions
	if !ballBouncesSides {
		p1w, p1h := g.players.P1.Coords()
		p2w, p2h := g.players.P2.Coords()

		if h+d > p1h && h < p1h+platformHeight && (p1w+platformWidth)/2+1 == w/2 {
			g.ball.SwitchX()
		}
		if h+d > p2h && h < p2h+platformHeight && p2w/2+1 == w/2+d {
			g.ball.SwitchX()
		}
	}
}
