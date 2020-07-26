package main

import (
	"log"

	"github.com/spinzed/tpong/game"
)

func main() {

	g := game.Game{}

	err := g.Init()

	if err != nil {
		log.Fatal(err)
	}

	g.Loop()
}
