package main

import (
	"log"
	"os"

	"github.com/spinzed/tpong/argparse"
	"github.com/spinzed/tpong/game"
)

func main() {
	options, err := argparse.Parse(os.Args)

	if err != nil {
		log.Fatal(err)
	}

	castedOptns := game.GameSettings(*options)

	g, err := game.Create(&castedOptns)

	if err != nil {
		log.Fatal(err)
	}

	g.Loop()
}
