package main

import (
	"math/rand"
	"time"
)

type Ball struct {
	x        int
	y        int
	initx    int
	inity    int
	diameter int
	velx     int
	vely     int
}

func newBall(x int, y int, d int, velx int, vely int) *Ball {
	ball := &Ball{x, y, x, y, d, velx, vely}

	// add some randomness to the initial ball spawn
	ball.Reset()

	return ball
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
	b.x += 2 * b.velx
	b.y += b.vely
}

func (b *Ball) SwitchX() {
	b.velx = -b.velx
}

func (b *Ball) SwitchY() {
	b.vely = -b.vely
}

func (b *Ball) Reset() {
	// add some randomness to the start position
	// must be even otherwise the ball may noclip
	seed := time.Now().UnixNano()
	rand.Seed(seed)
	rnd := (rand.Intn(2*ballRandomness+1) - ballRandomness) * 2

	b.x = b.initx + rnd
	b.y = b.inity

	// if ball is more to the left, it will go right and the other way
	if rnd*b.velx > 0 {
		b.SwitchX()
	}

	// if the ball was going up, change its direction
	if b.vely < 0 {
		b.SwitchY()
	}
}
