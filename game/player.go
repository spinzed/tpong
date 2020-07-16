package game

type Player struct {
	tag    string
	score  int
	pos    int
	width  int
	height int
}

var playerWidth int = 1
var playerHeight int = 5

// get a new player instance
func newPlayer(tag string, initialPos int) Player {
	p := Player{tag, 0, initialPos, playerWidth, playerHeight}
	return p
}

func (p *Player) GetYPos() int {
	return p.pos
}

func (p *Player) GoUp() {
	p.pos--
}

func (p *Player) GoDown() {
	p.pos++
}

func (p *Player) GetScore() int {
	return p.score
}

func (p *Player) AddPoint() {
	p.score++
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
