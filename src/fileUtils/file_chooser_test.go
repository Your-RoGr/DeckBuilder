package fileUtils

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/Your-RoGr/DeckBuilder/src/testUtils"
	"github.com/nsf/termbox-go"
)

func createTestDir(t *testing.T) string {

	dir, err := os.MkdirTemp("", "filechooser_test")
	if err != nil {
		t.Fatalf("failed to create temp dir: %v", err)
	}

	os.Mkdir(filepath.Join(dir, "folder1"), 0755)
	os.WriteFile(filepath.Join(dir, "file1.txt"), []byte("abc"), 0644)

	return dir
}

func removeTestDir(dir string, t *testing.T) {
	err := os.RemoveAll(dir)
	if err != nil {
		t.Errorf("failed to cleanup temp dir: %v", err)
	}
}

func TestFileChooser_ReadDir(t *testing.T) {

	dir := createTestDir(t)
	defer removeTestDir(dir, t)

	fc := &FileChooser{currentDir: dir}

	err := fc.readDir()
	if err != nil {
		t.Fatalf("readDir failed: %v", err)
	}

	var names []string
	for _, entry := range fc.entries {
		names = append(names, entry.Name())
	}

	want := map[string]bool{"..": true, "folder1": true, "file1.txt": true}

	for _, n := range names {
		if !want[n] {
			t.Errorf("unexpected entry: %s", n)
		}
		delete(want, n)
	}

	for n := range want {
		t.Errorf("missing entry: %s", n)
	}
}

func TestFileChooser_FoldersFirst(t *testing.T) {

	dir := createTestDir(t)
	defer removeTestDir(dir, t)

	fc := &FileChooser{currentDir: dir}
	if err := fc.readDir(); err != nil {
		t.Fatalf("readDir failed: %v", err)
	}

	if len(fc.entries) < 3 {
		t.Fatalf("entries length too short: %d", len(fc.entries))
	}

	if fc.entries[0].Name() != ".." || !fc.entries[0].IsDir() {
		t.Error("first entry must be .. directory")
	}

	if fc.entries[1].Name() != "folder1" || !fc.entries[1].IsDir() {
		t.Error("second entry must be folder1 directory")
	}

	if fc.entries[2].Name() != "file1.txt" || fc.entries[2].IsDir() {
		t.Error("third entry must be file1.txt file")
	}
}

func TestDirEntryUp(t *testing.T) {

	up := dirEntryUp{name: ".."}

	if up.Name() != ".." {
		t.Error("dirEntryUp.Name() should return ..")
	}

	if !up.IsDir() {
		t.Error("dirEntryUp must be directory")
	}

	if up.Type() != os.ModeDir {
		t.Error("dirEntryUp.Type() should be os.ModeDir")
	}

	info, err := up.Info()
	if info != nil || err != nil {
		t.Error("dirEntryUp.Info() should return nil, nil")
	}
}

func TestRedrawLogic(t *testing.T) {
	testUtils.NoPanic(t, func() {

		dir := createTestDir(t)
		defer removeTestDir(dir, t)

		termbox.Init()
		defer termbox.Close()

		fc := &FileChooser{currentDir: dir}
		fc.redraw()
	})
}

func TestStartLogic(t *testing.T) {
	testUtils.NoPanic(t, func() {

		dir := createTestDir(t)
		defer removeTestDir(dir, t)

		termbox.Init()
		defer termbox.Close()

		fc := &FileChooser{currentDir: dir}
		fc.Start()
	})
}
