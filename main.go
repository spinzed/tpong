package main

import (
	"github.com/spinzed/tpong/game"
)

func main() {
	g := game.Game{}

	err := g.Init()

	if err != nil {
		panic(err)
	}

	g.Loop()

	g.End()
}
