package fileUtils

import (
	"testing"

	"github.com/Your-RoGr/DeckBuilder/src/dataFrame"
	"github.com/Your-RoGr/DeckBuilder/src/testUtils"
	"github.com/nsf/termbox-go"
)

func TestWAStartWorkLogic(t *testing.T) {
	testUtils.NoPanic(t, func() {
		
		wordAdder := &WordAdder{}
		filePath := testUtils.TempCSVPath(t)

		termbox.Init()
		defer termbox.Close()

		_ = wordAdder.Start("Word", filePath)
	})
}

func TestWAStartWorkTranslateLogic(t *testing.T) {
	testUtils.NoPanic(t, func() {

		wordAdder := &WordAdder{}
		filePath := testUtils.TempCSVPath(t)

		termbox.Init()
		defer termbox.Close()

		_ = wordAdder.Start("Word-Translate", filePath)
	})
}

func TestWordModeLogic(t *testing.T) {
	testUtils.NoPanic(t, func() {

		wa := &WordAdder{}
		filePath := testUtils.TempCSVPath(t)

		wa.mode = "Word"
		wa.filePath = filePath
		wa.df = dataFrame.NewDataFrame(';')
		wa.df.LoadCSV(filePath)

		termbox.Init()
		defer termbox.Close()

		_ = wordMode(wa)
	})
}

func TestWordTranslateModeLogic(t *testing.T) {
	testUtils.NoPanic(t, func() {

		wa := &WordAdder{}
		filePath := testUtils.TempCSVPath(t)

		wa.mode = "Word-Translate"
		wa.filePath = filePath
		wa.df = dataFrame.NewDataFrame(';')
		wa.df.LoadCSV(filePath)

		termbox.Init()
		defer termbox.Close()

		_ = wordTranslateMode(wa)
	})
}

func TestWordExistsInFileLogic(t *testing.T) {
	testUtils.NoPanic(t, func() {

		termbox.Init()
		defer termbox.Close()

		filePath := testUtils.TempCSVPath(t)
		_, _ = wordExistsInFile(filePath, "aaa")
	})
}
