package dataFrame

import (
	"os"
	"os/user"
	"path/filepath"
	"reflect"
	"testing"
)

func TestNewDataFrame(t *testing.T) {

	df := NewDataFrame(',')

	if df == nil {
		t.Fatalf("NewDataFrame returned nil")
	}

	if df.delimiter != ',' {
		t.Errorf("Expected delimiter ',', got %v", df.delimiter)
	}

	if len(df.Columns) != 0 {
		t.Errorf("Expected empty Columns, got %v", df.Columns)
	}

	if len(df.Data) != 0 {
		t.Errorf("Expected empty Data, got %v", df.Data)
	}
}

func TestAddColumn(t *testing.T) {

	df := NewDataFrame(',')
	df.Columns = []string{"A"}
	df.Data = [][]string{
		{"1"},
		{"2"},
	}

	err := df.AddColumn("B", []string{"x", "y"})
	if err != nil {
		t.Fatalf("AddColumn returned error: %v", err)
	}

	if !reflect.DeepEqual(df.Columns, []string{"A", "B"}) {
		t.Errorf("Columns mismatch: %v", df.Columns)
	}

	if !reflect.DeepEqual(df.Data, [][]string{{"1", "x"}, {"2", "y"}}) {
		t.Errorf("Data mismatch: %v", df.Data)
	}

	err = df.AddColumn("C", []string{"z"})
	if err == nil {
		t.Errorf("Expected error for length mismatch, got nil")
	}
}

func TestAddRow(t *testing.T) {

	df := NewDataFrame(';')
	df.Columns = []string{"A", "B"}
	err := df.AddRow([]string{"1", "2"})

	if err != nil {
		t.Fatalf("AddRow error: %v", err)
	}

	if len(df.Data) != 1 || !reflect.DeepEqual(df.Data[0], []string{"1", "2"}) {
		t.Errorf("Data content mismatch: %v", df.Data)
	}

	err = df.AddRow([]string{"3"})
	if err == nil {
		t.Errorf("Expected error for short row, got nil")
	}
}

func TestAddUniqueRow(t *testing.T) {

	df := NewDataFrame(',')
	df.Columns = []string{"A", "B"}
	df.Data = [][]string{{"x", "1"}}
	err := df.AddUniqueRow([]string{"x", "1"})

	if err != nil {
		t.Fatalf("AddUniqueRow returned error: %v", err)
	}

	if len(df.Data) != 1 {
		t.Errorf("Expected row not duplicated, got %v", df.Data)
	}

	err = df.AddUniqueRow([]string{"y", "2"})
	if err != nil {
		t.Fatalf("AddUniqueRow returned error: %v", err)
	}

	if len(df.Data) != 2 {
		t.Errorf("Expected row added uniquely, got %v", df.Data)
	}
}

func TestDeleteRow(t *testing.T) {

	df := NewDataFrame(',')
	df.Columns = []string{"A"}
	df.Data = [][]string{{"one"}, {"two"}, {"three"}}
	err := df.DeleteRow(1)

	if err != nil {
		t.Fatalf("DeleteRow error: %v", err)
	}

	if len(df.Data) != 2 || df.Data[1][0] != "three" {
		t.Errorf("DeleteRow failed, data now %v", df.Data)
	}

	err = df.DeleteRow(999)
	if err == nil {
		t.Errorf("Expected error for out of range index, got nil")
	}
}

func TestGetColumnByName(t *testing.T) {

	df := NewDataFrame(',')
	df.Columns = []string{"a", "b"}
	df.Data = [][]string{{"x", "1"}, {"y", "2"}}
	col, err := df.GetColumnByName("b")

	if err != nil {
		t.Fatalf("GetColumnByName error: %v", err)
	}

	if !reflect.DeepEqual(col, []string{"1", "2"}) {
		t.Errorf("GetColumnByName mismatch: %v", col)
	}

	_, err = df.GetColumnByName("notfound")
	if err == nil {
		t.Errorf("Expected error for not found column")
	}
}

func TestGetAllColumns(t *testing.T) {

	df := NewDataFrame(',')
	df.Columns = []string{"a", "b"}
	df.Data = [][]string{{"x", "1"}, {"y", "2"}}
	cols := df.GetAllColumns()
	expect := [][]string{{"x", "y"}, {"1", "2"}}

	if !reflect.DeepEqual(cols, expect) {
		t.Errorf("GetAllColumns mismatch: got %v, want %v", cols, expect)
	}
}

func TestGetRowsAsStrings(t *testing.T) {

	df := NewDataFrame(',')
	df.Columns = []string{"a", "b"}
	df.Data = [][]string{{"one", "two"}, {"three", "four"}}
	rows := df.GetRowsAsStrings(" ")
	expect := []string{"one two", "three four"}

	if !reflect.DeepEqual(rows, expect) {
		t.Errorf("GetRowsAsStrings got %v, want %v", rows, expect)
	}
}

func TestCSVSaveLoad(t *testing.T) {

	df := NewDataFrame(',')
	df.Columns = []string{"col1", "col2"}
	df.Data = [][]string{{"1", "2"}, {"3", "4"}}

	file := filepath.Join(os.TempDir(), "testdataframe.csv")
	defer os.Remove(file)

	if err := df.SaveCSV(file); err != nil {
		t.Fatalf("SaveCSV error: %v", err)
	}

	df2 := NewDataFrame(',')
	if err := df2.LoadCSV(file); err != nil {
		t.Fatalf("LoadCSV error: %v", err)
	}

	if !reflect.DeepEqual(df2.Columns, df.Columns) {
		t.Errorf("Columns mismatch: %v", df2.Columns)
	}

	if !reflect.DeepEqual(df2.Data, df.Data) {
		t.Errorf("Data mismatch: %v", df2.Data)
	}
}

func TestCreateNewCSV(t *testing.T) {

	file := filepath.Join(os.TempDir(), "testnew.csv")
	defer os.Remove(file)

	err := CreateNewCSV(file, []string{"a", "b"}, ',')
	if err != nil {
		t.Fatalf("CreateNewCSV error: %v", err)
	}

	err = CreateNewCSV(file, []string{"c", "d"}, ',')
	if err != nil {
		t.Fatalf("CreateNewCSV on existing file error: %v", err)
	}

	if _, err := os.Stat(file); err != nil {
		t.Errorf("File does not exist: %v", err)
	}
}

func TestAddRowAndSave(t *testing.T) {

	file := filepath.Join(os.TempDir(), "testaddrow.csv")
	defer os.Remove(file)

	df := NewDataFrame(',')
	df.Columns = []string{"id", "val"}
	df.Data = [][]string{}

	err := df.AddRowAndSave([]string{"1", "x"}, file)
	if err != nil {
		t.Fatalf("AddRowAndSave error: %v", err)
	}

	df2 := NewDataFrame(',')
	if err := df2.LoadCSV(file); err != nil {
		t.Fatalf("LoadCSV error: %v", err)
	}

	if len(df2.Data) != 1 || df2.Data[0][0] != "1" || df2.Data[0][1] != "x" {
		t.Errorf("Saved/Loaded data mismatch: %v", df2.Data)
	}
}

func TestAddUniqueRowAndSave(t *testing.T) {

	file := filepath.Join(os.TempDir(), "testuniquerow.csv")
	defer os.Remove(file)

	df := NewDataFrame(',')
	df.Columns = []string{"name", "num"}
	df.Data = [][]string{{"aaa", "1"}}

	err := df.AddUniqueRowAndSave([]string{"bbb", "2"}, file)
	if err != nil {
		t.Fatalf("AddUniqueRowAndSave error: %v", err)
	}

	err = df.AddUniqueRowAndSave([]string{"bbb", "2"}, file)
	if err != nil {
		t.Fatalf("AddUniqueRowAndSave error: %v", err)
	}

	df2 := NewDataFrame(',')
	if err := df2.LoadCSV(file); err != nil {
		t.Fatalf("LoadCSV error: %v", err)
	}

	if len(df2.Data) != 2 {
		t.Errorf("Expected 2 unique rows, got: %v", df2.Data)
	}
}

func TestDeleteRowAndSave(t *testing.T) {

	file := filepath.Join(os.TempDir(), "testdelrow.csv")
	defer os.Remove(file)

	df := NewDataFrame(',')
	df.Columns = []string{"a"}
	df.Data = [][]string{{"1"}, {"2"}, {"3"}}

	if err := df.DeleteRowAndSave(1, file); err != nil {
		t.Fatalf("DeleteRowAndSave error: %v", err)
	}

	df2 := NewDataFrame(',')
	if err := df2.LoadCSV(file); err != nil {
		t.Fatalf("LoadCSV error: %v", err)
	}

	if len(df2.Data) != 2 || df2.Data[1][0] != "3" {
		t.Errorf("DeleteRowAndSave mismatch: %v", df2.Data)
	}
}

func TestDeleteRowByColumnValue(t *testing.T) {

	df := NewDataFrame(',')
	df.Columns = []string{"id", "name"}
	df.Data = [][]string{{"1", "x"}, {"2", "y"}, {"3", "z"}}
	err := df.DeleteRowByColumnValue("name", "y")

	if err != nil {
		t.Fatalf("DeleteRowByColumnValue error: %v", err)
	}

	if len(df.Data) != 2 {
		t.Errorf("Expected 2 rows after deletion, got %v", len(df.Data))
	}
}

func TestGetTrueFilepath(t *testing.T) {

	usr, _ := user.Current()
	home := usr.HomeDir
	out, err := getTrueFilepath("~/test.csv")

	if err != nil {
		t.Fatalf("getTrueFilepath error: %v", err)
	}

	expect := filepath.Join(home, "test.csv")
	if out != expect {
		t.Errorf("getTrueFilepath: got %v, want %v", out, expect)
	}
}

func TestGetColumnByIndex(t *testing.T) {

	df := NewDataFrame(',')
	df.Columns = []string{"a", "b"}
	df.Data = [][]string{{"1", "2"}, {"3", "4"}}
	col, err := df.getColumnByIndex(1)

	if err != nil {
		t.Fatalf("getColumnByIndex error: %v", err)
	}

	if !reflect.DeepEqual(col, []string{"2", "4"}) {
		t.Errorf("getColumnByIndex mismatch: %v", col)
	}

	_, err = df.getColumnByIndex(5)
	if err == nil {
		t.Errorf("Expected error for out of range column index")
	}
}

func TestDeleteRowByColumnValueAndSave(t *testing.T) {

	file := filepath.Join(os.TempDir(), "testdelbycol.csv")
	defer os.Remove(file)

	df := NewDataFrame(',')
	df.Columns = []string{"id", "name"}
	df.Data = [][]string{{"1", "a"}, {"2", "b"}}
	err := df.DeleteRowByColumnValueAndSave("name", "b", file)

	if err != nil {
		t.Fatalf("DeleteRowByColumnValueAndSave error: %v", err)
	}

	df2 := NewDataFrame(',')
	if err := df2.LoadCSV(file); err != nil {
		t.Fatalf("LoadCSV error: %v", err)
	}

	if len(df2.Data) != 1 || df2.Data[0][1] != "a" {
		t.Errorf("DeleteRowByColumnValueAndSave mismatch: %v", df2.Data)
	}
}
