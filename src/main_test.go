package main

import (
	"testing"

	"github.com/Your-RoGr/DeckBuilder/src/testUtils"
	"github.com/nsf/termbox-go"
)

func TestMainLogic(t *testing.T) {
	testUtils.NoPanic(t, func() {

		termbox.Init()
		defer termbox.Close()

		main()
	})
}
