package game

type Player struct {
	tag    string
	score  int
	xpos   int
	ypos   int
	width  int
	height int
	initx  int
	inity  int
}

// get a new player instance
func newPlayer(tag string, initialPos int, padx int) *Player {
	p := Player{tag, 0, padx, initialPos, platformWidth, platformHeight, padx, initialPos}
	return &p
}

func (p *Player) Coords() (int, int) {
	return p.xpos, p.ypos
}

func (p *Player) GoUp() {
	p.ypos--
}

func (p *Player) GoDown() {
	p.ypos++
}

func (p *Player) AddPoint() {
	p.score++
}

func (p *Player) GetScore() int {
	return p.score
}

func (p *Player) GetWidth() int {
	return p.width
}

func (p *Player) GetHeight() int {
	return p.height
}

func (p *Player) GetTag() string {
	return p.tag
}

func (p *Player) Reset() {
	p.ypos = p.initx
	p.ypos = p.inity
}

// Players struct includes both players individually and a method that returns them both
type Players struct {
	P1 *Player
	P2 *Player
}

func (p *Players) GetAll() []*Player {
	return []*Player{p.P1, p.P2}
}

// Get a pair of players ready and initialised
func newPlayers(w int, h int, padding int) *Players {
	initialPos := (h - platformHeight) / 2
	p1 := newPlayer(playerP1, initialPos, padding)
	p2 := newPlayer(playerP2, initialPos, w-padding)

	return &Players{p1, p2}
}
