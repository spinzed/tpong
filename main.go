package main

import (
	"fmt"
	"github.com/spinzed/tpong/game"
)

func main() {
	g := game.Game{}

	err := g.Init()

	if err != nil {
		panic(err)
	}

	g.Start()

	fmt.Scanln()

	g.Kill()
}
