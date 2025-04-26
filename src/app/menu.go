package app

import (
	"errors"
	"fmt"

	"github.com/Your-RoGr/DeckBuilder/src/appUtils"
	"github.com/Your-RoGr/DeckBuilder/src/dataFrame"
	"github.com/Your-RoGr/DeckBuilder/src/fileUtils"
	"github.com/nsf/termbox-go"
)

var existFilesPath = "~/.local/share/DeckBuilder/data/existFiles.csv"
var dfFileOptions *dataFrame.DataFrame

type Menu struct {
	name         string
	options      []string
	selected     int
	running      bool
	menus        []*Menu
	parent       *Menu
	scrollOffset int
}

// NewMenu создает новое меню с переданными опциями
func NewMainMenu() *Menu {

	err := dataFrame.CreateNewCSV(existFilesPath, []string{"Option"}, ';')

	if err != nil {
		panic(err)
	}

	dfFileOptions = dataFrame.NewDataFrame(';')
	dfFileOptions.LoadCSV(existFilesPath)

	options := []string{
		"Select file from catalog",
		"Select new file",
	}

	menus := make([]*Menu, 2)

	return &Menu{
		name:         "General",
		options:      options,
		selected:     0,
		running:      true,
		menus:        menus,
		parent:       nil,
		scrollOffset: 0,
	}
}

func newSubMenu(name string, parent *Menu, options []string) *Menu {

	menus := make([]*Menu, 0)

	return &Menu{
		name:         name,
		options:      options,
		selected:     0,
		running:      true,
		menus:        menus,
		parent:       parent,
		scrollOffset: 0,
	}
}

// Start запускает цикл отображения меню и обработки клавиш
func (m *Menu) Start() error {

	m.running = true

	if err := termbox.Init(); err != nil {
		return err
	}
	defer termbox.Close()

	for m.running {
		m.draw()
		switch ev := termbox.PollEvent(); ev.Type {
		case termbox.EventKey:
			switch ev.Key {
			case termbox.KeyArrowUp:
				if m.selected > 0 {
					m.selected--
				}
			case termbox.KeyArrowDown:
				if m.selected < len(m.options)-1 {
					m.selected++
				}
			case termbox.KeyEnter:
				err := m.selectOption()

				if err != nil {
					appUtils.GetInput(err.Error(), false)
				}
			case termbox.KeyEsc:
				m.running = false

				if m.parent != nil {
					m.parent.Start()
				}
			default:
				if ev.Ch == 'd' || ev.Ch == 'D' {
					if m.parent != nil {
						if m.parent.name == "Select file from catalog" {

							input, ok := appUtils.GetInput(
								fmt.Sprintf("Delete file %s? (y)", m.options[m.selected]),
								true,
							)

							if ok && input == "y" {
								err := dfFileOptions.DeleteRowByColumnValueAndSave(
									dfFileOptions.Columns[0],
									m.options[m.selected],
									existFilesPath,
								)

								m.options = append(m.options[:m.selected], m.options[m.selected+1:]...)

								if len(m.options) < 1 {
									m.running = false
									m.parent.Start()
								}

								if err != nil {
									appUtils.PrintHotkeyBar(err.Error(), true)
								}
							}
						}
					}
				}
			}
		case termbox.EventInterrupt, termbox.EventError:
			m.running = false
		}
	}

	return nil
}

func (m *Menu) draw() {

	termbox.Clear(termbox.ColorDefault, termbox.ColorDefault)
	_, height := termbox.Size()
	visibleRows := height - 2

	if visibleRows < 1 {
		visibleRows = 1
	}

	if m.selected < m.scrollOffset {
		m.scrollOffset = m.selected
	}
	if m.selected >= m.scrollOffset+visibleRows {
		m.scrollOffset = m.selected - visibleRows + 1
	}

	start := m.scrollOffset
	end := start + visibleRows
	if end > len(m.options) {
		end = len(m.options)
	}

	for i := start; i < end; i++ {

		fg := termbox.ColorWhite
		bg := termbox.ColorDefault

		if i == m.selected {
			fg = termbox.ColorBlack
			bg = termbox.ColorCyan
		}

		appUtils.SetLine(2, i-start+1, m.options[i], fg, bg)
	}

	appUtils.DrawVerticalBorders()
	appUtils.DrawHeader("DeckBuilder v0.1.2")

	if m.parent != nil && m.parent.name == "Select file from catalog" {
		appUtils.PrintHotkeyBar("  ▲/  ▼- select; D - delete; Enter - select; Esc - exit.", false)
	} else {
		appUtils.PrintHotkeyBar("  ▲/  ▼- select; Enter - select; Esc - exit.", false)
	}

	termbox.Flush()
}

func (m *Menu) selectOption() error {

	switch m.name {
	case "General":
		switch m.selected {
		// "Select file from catalog"
		case 0:

			options := []string{
				"Word",
				"Word-Translate",
				"Show",
			}

			menu := newSubMenu(m.options[m.selected], m, options)
			m.menus = append(m.menus, menu)
			menu.Start()
		//  "Select new file"
		case 1:

			chooser := &fileUtils.FileChooser{}
			path, err := chooser.Start()

			if err != nil {
				return err
			}

			if path != "" {

				err := dfFileOptions.AddRowAndSave([]string{path}, existFilesPath)

				if err != nil {
					appUtils.PrintHotkeyBar(fmt.Sprintf("Error: %s", err.Error()), true)
				} else {
					appUtils.PrintHotkeyBar(fmt.Sprintf("%s - successfully added", path), true)
				}
			}
		}
	case "Select file from catalog":

		options, err := dfFileOptions.GetColumnByName(dfFileOptions.Columns[0])
		if err != nil {
			return err
		}

		if len(options) > 0 {
			menu := newSubMenu(m.options[m.selected], m, options)
			m.menus = append(m.menus, menu)
			menu.Start()
		} else {
			return errors.New("no file's add new")
		}
	case "Select new file":
		return nil
	case "Word", "Word-Translate":

		wordAdder := &fileUtils.WordAdder{}
		err := wordAdder.Start(m.name, m.options[m.selected])

		if err != nil {
			return err
		}
	case "Show":

		df := dataFrame.NewDataFrame(';')
		err := df.LoadCSV(m.options[m.selected])
		if err != nil {
			return err
		}

		options := df.GetRowsAsStrings(" - ")

		if len(options) > 0 {
			menu := newSubMenu(m.options[m.selected], m, options)
			m.menus = append(m.menus, menu)
			menu.Start()
		} else {
			return errors.New("no word's add new")
		}
	default:
		return nil
	}
	return nil
}
