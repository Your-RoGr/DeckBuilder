package appUtils

import (
	"strings"

	"github.com/mattn/go-runewidth"
	"github.com/nsf/termbox-go"
)

var fg, bg = termbox.ColorBlack, termbox.ColorMagenta

// DrawHeader displays the application header with the passed title centered on it
func DrawHeader(appName string) {

	width, _ := termbox.Size()
	y := 0

	appName = strings.TrimSpace(appName)
	nameLen := len([]rune(appName))
	start := (width - nameLen) / 2
	if start < 0 {
		start = 0
	}

	// Filling the string with color
	for x := 0; x < width; x++ {
		termbox.SetCell(x, y, ' ', fg, bg)
	}

	// Write centered text
	SetLine(start, y, appName, fg, bg)
}

// DrawVerticalBorders draws vertical lines on the left and right with a thickness of 2 character
func DrawVerticalBorders() {

	width, height := termbox.Size()
	for y := 0; y < height; y++ {

		// Left border (2 columns)
		termbox.SetCell(0, y, rune(0), fg, bg)
		if width > 1 {
			termbox.SetCell(1, y, rune(0), fg, bg)
		}

		// Right border (2 columns)
		if width > 2 {
			termbox.SetCell(width-2, y, rune(0), fg, bg)
		}
		if width > 3 {
			termbox.SetCell(width-1, y, rune(0), fg, bg)
		}
	}
}

func PrintHotkeyBar(msg string, isUp bool) {

	width, height := termbox.Size()
	y := height - 1

	if isUp {
		for x := 0; x < width; x++ {
			termbox.SetCell(x, 1, ' ', fg, bg)
		}
	} else {
		for x := 0; x < width; x++ {
			termbox.SetCell(x, y, ' ', fg, bg)
		}
	}

	if isUp {
		SetLine(2, 1, msg, fg, bg)
	} else {
		SetLine(2, y, msg, fg, bg)
	}
}

// PromptForFilename - Tooltip for entering the file name (in the simple variant - at the bottom of the screen)
func PromptForFilename(prompt string, inputRequire bool) (string, bool) {

	termbox.Clear(termbox.ColorDefault, termbox.ColorDefault)
	PrintHotkeyBar(prompt, true)
	DrawVerticalBorders()
	DrawHeader("DeckBuilder v0.1.0")
	PrintHotkeyBar("Enter - send; Esc - exit.", false)
	termbox.Flush()
	var input []rune

	for {
		switch ev := termbox.PollEvent(); ev.Type {
		case termbox.EventKey:
			if ev.Key == termbox.KeyEnter {
				if len(input) > 0 || !inputRequire {
					return strings.TrimSpace(string(input)), true
				}
			}
			if ev.Key == termbox.KeyEsc {
				return "", false
			}
			if ev.Key == termbox.KeyBackspace || ev.Key == termbox.KeyBackspace2 {
				if len(input) > 0 {
					input = input[:len(input)-1]
				}
			} else if ev.Ch != 0 && len(input) < 128 {
				input = append(input, ev.Ch)
			}
		}

		width, _ := termbox.Size()

		for x := 0; x < width; x++ {
			termbox.SetCell(x+2, 2, ' ', termbox.ColorDefault, termbox.ColorDefault)
		}

		SetLine(2, 2, string(input), termbox.ColorYellow, termbox.ColorDefault)

		DrawVerticalBorders()
		termbox.Flush()
	}
}

func SetLine(x, y int, msg string, fg, bg termbox.Attribute) {
	for _, c := range msg {
		termbox.SetCell(x, y, c, fg, bg)
		x += runewidth.RuneWidth(c)
	}
}
