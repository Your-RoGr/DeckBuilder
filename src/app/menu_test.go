package app

import (
	"os"
	"testing"

	"github.com/Your-RoGr/DeckBuilder/src/dataFrame"
	"github.com/Your-RoGr/DeckBuilder/src/testUtils"
)

func TestNewMainMenu_createsMenuAndCSV(t *testing.T) {
	
	path := testUtils.TempCSVPath(t)
	existFilesPath = path

	menu := NewMainMenu()

	if menu == nil {
		t.Fatal("Expected menu to be created")
	}

	if menu.name != "General" {
		t.Errorf("Expected name 'General', got '%s'", menu.name)
	}

	if len(menu.options) != 2 {
		t.Errorf("Expected 2 options, got %d", len(menu.options))
	}

	if _, err := os.Stat(path); err != nil {
		t.Errorf("Expected CSV at %s to be created, err: %v", path, err)
	}
}

func TestNewSubMenu_works(t *testing.T) {

	parent := &Menu{name: "parent"}
	opts := []string{"one", "two"}
	menu := newSubMenu("child", parent, opts)

	if menu.name != "child" {
		t.Error("Bad submenu name")
	}

	if menu.parent != parent {
		t.Error("Parent not set")
	}

	if len(menu.options) != 2 || menu.options[0] != "one" {
		t.Error("Bad submenu options")
	}
}

func TestDataFrame_AddRowAndSave_and_LoadCSV_roundtrip(t *testing.T) {

	path := testUtils.TempCSVPath(t)

	df := dataFrame.NewDataFrame(';')
	df.Columns = []string{"Col"}
	df.Data = [][]string{}

	if err := df.SaveCSV(path); err != nil {
		t.Fatalf("SaveCSV failed: %v", err)
	}

	if err := df.AddRowAndSave([]string{"val1"}, path); err != nil {
		t.Fatalf("AddRowAndSave: %v", err)
	}

	df2 := dataFrame.NewDataFrame(';')

	if err := df2.LoadCSV(path); err != nil {
		t.Fatalf("LoadCSV failed: %v", err)
	}

	if len(df2.Data) != 1 || df2.Data[0][0] != "val1" {
		t.Error("AddRowAndSave/LoadCSV Data mismatch")
	}
}

func TestDataFrame_DeleteRowByColumnValueAndSave(t *testing.T) {

	path := testUtils.TempCSVPath(t)

	df := dataFrame.NewDataFrame(';')
	df.Columns = []string{"Col1"}
	df.Data = [][]string{
		{"foo"},
		{"bar"},
	}
	_ = df.SaveCSV(path)

	if err := df.DeleteRowByColumnValueAndSave("Col1", "foo", path); err != nil {
		t.Fatalf("DeleteRowByColumnValueAndSave: %v", err)
	}

	df2 := dataFrame.NewDataFrame(';')
	_ = df2.LoadCSV(path)

	if len(df2.Data) != 1 || df2.Data[0][0] != "bar" {
		t.Errorf("Expected only [bar] after delete, got: %+v", df2.Data)
	}
}

func TestDataFrame_GetColumnByName(t *testing.T) {

	df := dataFrame.NewDataFrame(';')
	df.Columns = []string{"A", "B"}
	df.Data = [][]string{{"1", "2"}, {"3", "4"}}
	c, err := df.GetColumnByName("B")

	if err != nil {
		t.Fatal(err)
	}

	if len(c) != 2 || c[0] != "2" || c[1] != "4" {
		t.Errorf("GetColumnByName wrong result: %+v", c)
	}
}

func TestDataFrame_GetRowsAsStrings(t *testing.T) {

	df := dataFrame.NewDataFrame(';')
	df.Columns = []string{"A", "B"}
	df.Data = [][]string{{"1", "2"}, {"3", "4"}}
	s := df.GetRowsAsStrings("-")

	if s[0] != "1-2" || s[1] != "3-4" {
		t.Errorf("GetRowsAsStrings wrong result: %+v", s)
	}
}

func TestMenu_selectOption_General_SelectFileFromCatalog(t *testing.T) {

	menu := &Menu{
		name:     "General",
		options:  []string{"Select file from catalog", "Select new file"},
		selected: 0,
		menus:    []*Menu{},
	}

	if err := menu.selectOption(); err != nil {
		t.Fatalf("selectOption failed: %v", err)
	}
	
	if len(menu.menus) == 0 {
		t.Error("submenu was not added")
	}

	if menu.menus[len(menu.menus)-1].name != "Select file from catalog" {
		t.Errorf("Expected submenu name, got %s", menu.menus[len(menu.menus)-1].name)
	}
}

func TestMenu_selectOption_SelectFileFromCatalog_empty(t *testing.T) {

	dfFileOptions = dataFrame.NewDataFrame(';')
	dfFileOptions.Columns = []string{"Option"}
	dfFileOptions.Data = [][]string{}

	menu := &Menu{
		name:     "Select file from catalog",
		options:  []string{"somefile.csv"},
		selected: 0,
	}

	err := menu.selectOption()
	if err == nil || err.Error() != "no file's add new" {
		t.Errorf("Expected error about missing files, got %v", err)
	}
}
