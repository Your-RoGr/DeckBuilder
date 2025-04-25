package testUtils

import (
	"os"
	"path/filepath"
	"testing"
)

func TestNoPanic_NoPanic(t *testing.T) {
	NoPanic(t, func() {
		_ = 1 + 1
	})
}

func TestNoPanic_WithPanic(t *testing.T) {
	tt := &testing.T{}
	NoPanic(tt, func() {
		panic("fail!")
	})
	if !tt.Failed() {
		t.Error("NoPanic: expected test to fail on panic, but it did not")
	}
}

func TestNoPanic_EmptyFunc(t *testing.T) {
	NoPanic(t, func() {})
}

func TestNoPanic_DeferNoPanic(t *testing.T) {
	NoPanic(t, func() {
		defer func() {
			_ = 2 * 2
		}()
	})
}

func TestTempDataDir(t *testing.T) {
	
	tempDir := TempDataDir(t)

	if _, err := os.Stat(tempDir); os.IsNotExist(err) {
		t.Errorf("TempDataDir failed to create directory: %v", err)
	}

	testFile := filepath.Join(tempDir, "testfile.txt")
	content := []byte("test content")
	if err := os.WriteFile(testFile, content, 0644); err != nil {
		t.Errorf("Failed to write to temporary directory: %v", err)
	}

	readContent, err := os.ReadFile(testFile)
	if err != nil {
		t.Errorf("Failed to read from temporary directory: %v", err)
	}

	if string(readContent) != string(content) {
		t.Errorf("File content mismatch. Got %s, want %s", string(readContent), string(content))
	}
}

func TestTempCSVPath(t *testing.T) {
	
	csvPath := TempCSVPath(t)

	dir := filepath.Dir(csvPath)
	filename := filepath.Base(csvPath)

	if _, err := os.Stat(dir); os.IsNotExist(err) {
		t.Errorf("Directory for CSV path doesn't exist: %v", err)
	}

	if filename != "test.csv" {
		t.Errorf("Incorrect CSV filename. Got %s, want test.csv", filename)
	}

	content := []byte("id,name\n1,test\n2,example")
	if err := os.WriteFile(csvPath, content, 0644); err != nil {
		t.Errorf("Failed to write to CSV path: %v", err)
	}

	readContent, err := os.ReadFile(csvPath)
	if err != nil {
		t.Errorf("Failed to read from CSV path: %v", err)
	}

	if string(readContent) != string(content) {
		t.Errorf("CSV content mismatch. Got %s, want %s", string(readContent), string(content))
	}
}

func TestTempDirCleanup(t *testing.T) {
	
	var tempDir string

	t.Run("Create directory", func(t *testing.T) {
		tempDir = TempDataDir(t)
		testFile := filepath.Join(tempDir, "cleanup_test.txt")
		if err := os.WriteFile(testFile, []byte("test"), 0644); err != nil {
			t.Errorf("Failed to write test file: %v", err)
		}
	})

	if _, err := os.Stat(tempDir); !os.IsNotExist(err) {
		t.Errorf("Directory was not cleaned up properly: %v", err)
	}
}
