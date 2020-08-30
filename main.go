package main

import (
	"log"
	"os"
)

func main() {
	options, err := ArgParse(os.Args)

	if err != nil {
		log.Fatal(err)
	}

	castedOptns := GameSettings(*options)

	g, err := CreateGame(&castedOptns)

	if err != nil {
		log.Fatal(err)
	}

	g.Loop()
}
