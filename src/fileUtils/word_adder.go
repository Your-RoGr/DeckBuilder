package fileUtils

import (
	"bufio"
	"errors"
	"fmt"
	"os"
	"strings"

	"github.com/Your-RoGr/DeckBuilder/src/appUtils"
	"github.com/Your-RoGr/DeckBuilder/src/dataFrame"
)

// WordAdder — structure for selecting a file and adding words to it.
type WordAdder struct {
	mode     string
	filePath string // Path to the selected file
	df       *dataFrame.DataFrame
}

// Start — allows selecting a file (txt/csv) and adding new words separated by ;.
func (wa *WordAdder) Start(mode string, path string) error {

	wa.mode = mode
	wa.filePath = path
	wa.df = dataFrame.NewDataFrame(';')
	wa.df.LoadCSV(path)

	// Loop for entering words
	switch wa.mode {
	case "Word":
		for {
			err := wordMode(wa)

			if err != nil {
				if err.Error() == "break" {
					return nil
				}
				return err
			}
		}
	case "Word-Translate":
		for {
			err := wordTranslateMode(wa)

			if err != nil {
				if err.Error() == "break" {
					return nil
				}
				return err
			}
		}
	default:
		return nil
	}
}

func wordMode(wa *WordAdder) error {

	word, ok := appUtils.GetInput("Enter a word to add: ", true)
	if !ok || strings.TrimSpace(word) == "" {
		return errors.New("break")
	}

	word = strings.TrimSpace(word)

	// Check if the word already exists in the file
	exists, err := wordExistsInFile(wa.filePath, word)
	if err != nil {
		appUtils.GetInput("Error: "+err.Error(), false)
		return err
	}

	if exists {
		appUtils.GetInput(fmt.Sprintf("'%s' already exists!", word), false)
		return nil
	}

	err = wa.df.AddUniqueRowAndSave([]string{word, ""}, wa.filePath)
	if err != nil {
		appUtils.GetInput("Error: "+err.Error(), false)
		return err
	}

	appUtils.GetInput(fmt.Sprintf("%s added!", word), false)
	return nil
}

func wordTranslateMode(wa *WordAdder) error {

	word, ok := appUtils.GetInput("Enter a word to add: ", true)
	if !ok || strings.TrimSpace(word) == "" {
		return errors.New("break")
	}

	word = strings.TrimSpace(word)

	// Check if the word already exists in the file
	exists, err := wordExistsInFile(wa.filePath, word)
	if err != nil {
		appUtils.GetInput("Error: "+err.Error(), false)
		return err
	}

	if exists {
		appUtils.GetInput(fmt.Sprintf("'%s' already exists!", word), false)
		return nil
	}

	translate, ok := appUtils.GetInput("Enter a translate to add: ", true)
	if !ok || strings.TrimSpace(word) == "" {
		return errors.New("strings.TrimSpace(word) == \"\"")
	}

	translate = strings.TrimSpace(translate)

	err = wa.df.AddUniqueRowAndSave([]string{word, translate}, wa.filePath)
	if err != nil {
		appUtils.GetInput("Error: "+err.Error(), false)
		return err
	}

	appUtils.GetInput(fmt.Sprintf("%s - %s added!", word, translate), false)
	return nil
}

// wordExistsInFile — checks if a word already exists in the file (case-insensitive)
func wordExistsInFile(filePath, targetWord string) (bool, error) {

	f, err := os.Open(filePath)
	if err != nil {
		return false, err
	}
	defer f.Close()

	targetWordLower := strings.ToLower(strings.TrimSpace(targetWord))

	scanner := bufio.NewScanner(f)
	for scanner.Scan() {
		// Split line by semicolon
		words := strings.Split(scanner.Text(), ";")
		for _, w := range words {
			w = strings.ToLower(strings.TrimSpace(w))
			if w == targetWordLower {
				return true, nil
			}
		}
	}
	return false, scanner.Err()
}
