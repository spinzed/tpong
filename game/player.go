package game

type Player struct {
	tag    string
	score  int
	xpos   int
	ypos   int
	width  int
	height int
}

// get a new player instance
func newPlayer(tag string, initialPos int, padx int) *Player {
	p := Player{tag, 0, padx, initialPos, playerWidth, playerHeight}
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
