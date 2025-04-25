package testUtils

import (
	"os"
	"path/filepath"
	"testing"
)

func NoPanic(t *testing.T, f func()) {

	defer func() {
		if r := recover(); r != nil {
			t.Errorf("unexpected panic: %v", r)
		}
	}()

	f()
}

func TempDataDir(t *testing.T) string {
	dir := t.TempDir()
	t.Cleanup(func() {
		_ = os.RemoveAll(dir)
	})
	return dir
}

func TempCSVPath(t *testing.T) string {
	return filepath.Join(TempDataDir(t), "test.csv")
}
