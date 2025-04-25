package main

import (
	"log"

	"github.com/Your-RoGr/DeckBuilder/src/app"
)

func main() {

	err := app.NewMainMenu().Start()

	if err != nil {
		log.Fatal(err)
	}
}
