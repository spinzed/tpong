package game

type Ball struct {
	x int
	y int
	initx int
	inity int
	diameter int
	velx int
	vely int
}

func newBall(x int, y int, d int, velx int, vely int) *Ball {
	return &Ball{x, y, x, y, d, velx, vely}
}

// Returns x and y coords of the upper left corner of the ball
func (b *Ball) Coords() (int, int) {
	return b.x, b.y
}

// Returns the diameter
func (b *Ball) Diam() int {
	return b.diameter
}

// Returns the velocities
func (b *Ball) Vels() (int, int) {
	return b.velx, b.vely
}

func (b *Ball) Move() {
	// terminal letter width is 1/2 of height, so 2x is to compensate that
	b.x += 2*b.velx
	b.y += b.vely
}

func (b *Ball) SwitchX() {
	b.velx = -b.velx
}

func (b *Ball) SwitchY() {
	b.vely = -b.vely
}
