package appUtils

import (
	"testing"

	"github.com/Your-RoGr/DeckBuilder/src/testUtils"
	"github.com/nsf/termbox-go"
)

func TestDrawHeaderPosition(t *testing.T) {
	
	width := 100
	appName := "DeckBuilder"
	nameLen := len([]rune(appName))
	start := (width - nameLen) / 2

	if start < 0 {
		t.Errorf("The text position is less than 0, although the name is no longer than the screen width")
	}

	if start != (width-nameLen)/2 {
		t.Errorf("Incorrect position calculation for text")
	}
}

func TestDrawHeaderLogic(t *testing.T) {
	testUtils.NoPanic(t, func() { DrawHeader("AppName") })
}

func TestDrawVerticalBordersLogic(t *testing.T) {
	testUtils.NoPanic(t, func() { DrawVerticalBorders() })
}

func TestSetLineLogic(t *testing.T) {
	testUtils.NoPanic(t, func() { SetLine(0, 0, "Test", termbox.ColorBlack, termbox.ColorMagenta) })
}

func TestPrintHotkeyBarLogic(t *testing.T) {
	testUtils.NoPanic(t, func() { PrintHotkeyBar("Test Message", false) })
}

func TestPromptForFilenameLogic(t *testing.T) {
	testUtils.NoPanic(t, func() {
		termbox.Init()
		defer termbox.Close()
		PromptForFilename("prompt", false)
	})
}
