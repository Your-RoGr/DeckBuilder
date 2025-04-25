package fileUtils

import (
	"fmt"
	"os"
	"path/filepath"
	"sort"

	"github.com/Your-RoGr/DeckBuilder/src/appUtils"
	"github.com/Your-RoGr/DeckBuilder/src/dataFrame"
	"github.com/nsf/termbox-go"
)

// FileChooser implements file/folder selection in the terminal
type FileChooser struct {
	currentDir   string
	entries      []os.DirEntry
	selected     int
	scrollOffset int
}

func (fc *FileChooser) readDir() error {

	entries, err := os.ReadDir(fc.currentDir)
	if err != nil {
		return err
	}

	// Sort: folders first, then files (for convenience)
	sort.SliceStable(entries, func(i, j int) bool {
		if entries[i].IsDir() == entries[j].IsDir() {
			return entries[i].Name() < entries[j].Name()
		}
		return entries[i].IsDir()
	})

	fc.entries = nil
	fc.entries = append(fc.entries, dirEntryUp{
		name: "..",
	})

	fc.entries = append(fc.entries, entries...)
	fc.selected = 0
	return nil
}

// dirEntryUp implements os.DirEntry to exit upwards
type dirEntryUp struct{ name string }

func (d dirEntryUp) Name() string               { return d.name }
func (d dirEntryUp) IsDir() bool                { return true }
func (d dirEntryUp) Type() os.FileMode          { return os.ModeDir }
func (d dirEntryUp) Info() (os.FileInfo, error) { return nil, nil }

// UI method for displaying content
func (fc *FileChooser) redraw() {

	termbox.Clear(termbox.ColorDefault, termbox.ColorDefault)
	_, height := termbox.Size()
	visibleRows := height - 2

	if visibleRows < 1 {
		visibleRows = 1
	}

	if fc.selected < fc.scrollOffset {
		fc.scrollOffset = fc.selected
	}

	if fc.selected >= fc.scrollOffset+visibleRows {
		fc.scrollOffset = fc.selected - visibleRows + 1
	}

	start := fc.scrollOffset
	end := start + visibleRows

	if end > len(fc.entries) {
		end = len(fc.entries)
	}

	for i := start; i < end; i++ {

		entry := fc.entries[i]
		fg, bg := termbox.ColorDefault, termbox.ColorDefault

		if i == fc.selected {
			fg, bg = termbox.ColorBlack, termbox.ColorCyan
		}

		name := entry.Name()

		if entry.IsDir() {
			name = "[" + name + "]"
		}

		appUtils.SetLine(2, i-start+1, name, fg, bg)
	}

	appUtils.DrawVerticalBorders()
	appUtils.DrawHeader("DeckBuilder v0.1.0")
	appUtils.PrintHotkeyBar("  ▲/  ▼- select; A - create file; Enter - open; Esc - exit.", false)
	termbox.Flush()
}

func (fc *FileChooser) Start() (string, error) {

	var err error
	fc.currentDir, err = os.Getwd()

	if err != nil {
		return "", err
	}

	if err := fc.readDir(); err != nil {
		return "", err
	}

	fc.redraw()

	for {
		switch ev := termbox.PollEvent(); ev.Type {
		case termbox.EventKey:
			switch ev.Key {
			case termbox.KeyArrowUp:
				if fc.selected > 0 {
					fc.selected--
				}
			case termbox.KeyArrowDown:
				if fc.selected < len(fc.entries)-1 {
					fc.selected++
				}
			case termbox.KeyEnter:

				entry := fc.entries[fc.selected]

				if entry.IsDir() {

					if entry.Name() == ".." {
						// Go up
						parent := filepath.Dir(fc.currentDir)
						fc.currentDir = parent
					} else {
						// Go inside the folder
						fc.currentDir = filepath.Join(fc.currentDir, entry.Name())
					}

					if err := fc.readDir(); err != nil {
						return "", err
					}
				} else {
					// File selected, return path
					return filepath.Join(fc.currentDir, entry.Name()), nil
				}
			case termbox.KeyEsc, termbox.KeyCtrlC:
				return "", nil
			default:
				if ev.Ch == 'a' || ev.Ch == 'A' {

					name, ok := appUtils.PromptForFilename(
						fmt.Sprintf(
							"%s/... Enter the name of the new file: ",
							fc.currentDir,
						),
						true,
					)

					if ok && name != "" {

						fullpath := filepath.Join(fc.currentDir, name)
						err := dataFrame.CreateNewCSV(fullpath, []string{"Word", "Translation"}, ';')
						
						if err == nil {
							fc.readDir()
						} else {
							appUtils.PromptForFilename("File creation error: "+err.Error(), false)
						}
					}
				}
			}

			fc.redraw()
		case termbox.EventError:
			return "", ev.Err
		}
	}
}
