package dataFrame

import (
	"encoding/csv"
	"errors"
	"os"
	"os/user"
	"path/filepath"
	"strings"
)

// DataFrame is a structure for storing tabular data
type DataFrame struct {
	Columns   []string   // column names
	Data      [][]string // rows of data (each row is a slice of strings)
	delimiter rune       // field delimiter (e.g. ',', ';', '\t')
}

// NewDataFrame creates an empty DataFrame with a specified delimiter
func NewDataFrame(delimiter rune) *DataFrame {
	return &DataFrame{
		delimiter: delimiter,
	}
}

// getTrueFilepath return filepath with /home/{user}/... if used ~/...
func getTrueFilepath(filePath string) (string, error) {

	if strings.HasPrefix(filePath, "~") {
		usr, err := user.Current()
		if err != nil {
			return "", err
		}
		filePath = filepath.Join(usr.HomeDir, filePath[1:])
	}

	return filePath, nil
}

// CreateNewCSV creates a new csv file with the specified columns and saves it at filePath.
// If the file already exists, it not will be overwritten.
func CreateNewCSV(filePath string, columns []string, delimiter rune) error {

	filePath, err := getTrueFilepath(filePath)

	if err != nil {
		return err
	}

	dir := filepath.Dir(filePath)

	// Create all necessary intermediate directories.
	if err := os.MkdirAll(dir, 0755); err != nil {
		return err
	}

	if _, err := os.Stat(filePath); os.IsNotExist(err) {

		df := NewDataFrame(delimiter)
		df.Columns = columns
		df.Data = [][]string{}

		return df.SaveCSV(filePath)
	} else {
		return nil
	}
}

// LoadCSV loads data from a CSV file into the DataFrame
func (df *DataFrame) LoadCSV(filePath string) error {

	filePath, err := getTrueFilepath(filePath)

	if err != nil {
		return err
	}

	file, err := os.Open(filePath)
	if err != nil {
		return err
	}
	defer file.Close()

	reader := csv.NewReader(file)
	reader.Comma = df.delimiter
	records, err := reader.ReadAll()
	if err != nil {
		return err
	}

	if len(records) == 0 {
		return errors.New("csv file is empty")
	}

	df.Columns = records[0]
	df.Data = records[1:]

	return nil
}

// SaveCSV saves the DataFrame to a CSV file
func (df *DataFrame) SaveCSV(filePath string) error {

	filePath, err := getTrueFilepath(filePath)

	if err != nil {
		return err
	}

	file, err := os.Create(filePath)
	if err != nil {
		return err
	}
	defer file.Close()

	writer := csv.NewWriter(file)
	writer.Comma = df.delimiter
	defer writer.Flush()

	// Write column names
	if err := writer.Write(df.Columns); err != nil {
		return err
	}

	// Write data rows
	for _, row := range df.Data {
		if err := writer.Write(row); err != nil {
			return err
		}
	}

	return nil
}

// AddColumn adds a column to df
func (df *DataFrame) AddColumn(name string, values []string) error {

	if len(values) != len(df.Data) {
		return errors.New("values length does not match number of rows")
	}

	df.Columns = append(df.Columns, name)

	for i, row := range df.Data {
		df.Data[i] = append(row, values[i])
	}

	return nil
}

// AddRow adds a row to df.Data
func (df *DataFrame) AddRow(row []string) error {

	if len(row) != len(df.Columns) {
		return errors.New("row length does not match number of columns")
	}

	df.Data = append(df.Data, row)

	return nil
}

// AddRowAndSave adds a row and saves the DataFrame to a CSV file
func (df *DataFrame) AddRowAndSave(row []string, filepath string) error {

	if err := df.AddRow(row); err != nil {
		return err
	}

	return df.SaveCSV(filepath)
}

// AddUniqueRow adds a row only if it is not already in the DataFrame
func (df *DataFrame) AddUniqueRow(row []string) error {

	if len(row) != len(df.Columns) {
		return errors.New("row length does not match number of columns")
	}

	for _, existing := range df.Data {

		duplicate := true

		for i := range row {
			if existing[i] != row[i] {
				duplicate = false
				break
			}
		}

		if duplicate {
			return nil
		}
	}

	df.Data = append(df.Data, row)

	return nil
}

// AddUniqueRowAndSave adds a row only if it is not already in the DataFrame
// and saves the DataFrame to a CSV file
func (df *DataFrame) AddUniqueRowAndSave(row []string, filepath string) error {

	if err := df.AddUniqueRow(row); err != nil {
		return err
	}

	return df.SaveCSV(filepath)
}

// DeleteRow deletes a row by its index
func (df *DataFrame) DeleteRow(index int) error {

	if index < 0 || index >= len(df.Data) {
		return errors.New("index out of range")
	}

	df.Data = append(df.Data[:index], df.Data[index+1:]...)

	return nil
}

// DeleteRowAndSave deletes the row by index and saves the DataFrame to a CSV file
func (df *DataFrame) DeleteRowAndSave(index int, filePath string) error {

	if err := df.DeleteRow(index); err != nil {
		return err
	}

	return df.SaveCSV(filePath)
}

// DeleteRowByColumnValue deletes the first row by value in the specified column
func (df *DataFrame) DeleteRowByColumnValue(columnName, value string) error {

	colIdx := -1

	for i, name := range df.Columns {
		if name == columnName {
			colIdx = i
			break
		}
	}

	if colIdx == -1 {
		return errors.New("column name not found")
	}

	err := df.DeleteRow(colIdx)

	return err
}

// DeleteRowByColumnValueAndSave deletes the row by index and saves the DataFrame to a CSV file
// and saves the DataFrame to a CSV file
func (df *DataFrame) DeleteRowByColumnValueAndSave(columnName, value, filePath string) error {

	if err := df.DeleteRowByColumnValue(columnName, value); err != nil {
		return err
	}

	return df.SaveCSV(filePath)
}

// getColumnByIndex returns all column values by index
func (df *DataFrame) getColumnByIndex(idx int) ([]string, error) {

	if idx < 0 || idx >= len(df.Columns) {
		return nil, errors.New("column index out of range")
	}

	column := make([]string, len(df.Data))

	for i, row := range df.Data {
		if idx < len(row) {
			column[i] = row[idx]
		} else {
			column[i] = ""
		}
	}

	return column, nil
}

// GetColumnByName returns all column values by name
func (df *DataFrame) GetColumnByName(name string) ([]string, error) {

	var idx = -1

	for i, n := range df.Columns {
		if n == name {
			idx = i
			break
		}
	}

	if idx == -1 {
		return nil, errors.New("column name not found")
	}

	return df.getColumnByIndex(idx)
}

// GetAllColumns returns a two-dimensional array: each inner array is a separate column
func (df *DataFrame) GetAllColumns() [][]string {

	result := make([][]string, len(df.Columns))

	for colIdx := range df.Columns {
		result[colIdx], _ = df.getColumnByIndex(colIdx)
	}

	return result
}

// GetRowsAsStrings returns a slice of rows, where each row is a
// these are the combined values of a single DataFrame data string using a delimiter
func (df *DataFrame) GetRowsAsStrings(delimeter string) []string {

	rows := make([]string, len(df.Data))

	for i, row := range df.Data {
		rows[i] = strings.Join(row, delimeter)
	}

	return rows
}
