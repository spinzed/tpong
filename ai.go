package main

func (g *Game) aiMove() {
	w, h := g.screen.Size()
	bw, bh := g.ball.Coords()
	pw, ph := g.players.P2.Coords()

	w, bw, pw = w/2, bw/2, pw/2

	if bw < w/2 {
		return
	}

	// if the ball has surpassed half of the screen, calculate the spot
	// where ball will be and move the platform there

	currentW := pw - bw
	currentH := bh
	_, vely := g.ball.Vels()
	if vely > 0 {
		vely = 1
	} else {
		vely = -1
	}

	// in case double bounce occurs, this will shorten the value without any change
	currentW = currentW % (2 * h)

	if currentW > h {
		currentW = currentW % h
		// in this case vertical speed direction of the ball must be switched
		vely *= -1
	}

	currentH = currentH + (currentW * vely)

	if currentH < 0 {
		currentH *= -1
	}

	// now move the platform in the direction where the ball will be
	if currentH > ph+platformHeight/2 {
		g.movePlayerDown(g.players.P2)
	}

	if currentH < ph+platformHeight/2 {
		g.movePlayerUp(g.players.P2)
	}

	//st := g.theme.GetCurrent().GetTextStyle()
	//g.screen.SetContent(10, 10, 'a', []rune(strconv.Itoa(currentW)), st)
	//g.screen.SetContent(10, 11, 'a', []rune(strconv.Itoa(currentH)), st)
	//g.screen.SetContent(10, 12, 'a', []rune(strconv.Itoa(w)), st)
	//g.screen.Show()
}
