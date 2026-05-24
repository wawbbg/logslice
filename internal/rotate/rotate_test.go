package rotate_test

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/yourorg/logslice/internal/rotate"
)

func makeFile(t *testing.T, path string) {
	t.Helper()
	if err := os.WriteFile(path, []byte("log line\n"), 0o644); err != nil {
		t.Fatalf("makeFile: %v", err)
	}
}

func TestDiscover_NoRotatedFiles(t *testing.T) {
	dir := t.TempDir()
	base := filepath.Join(dir, "app.log")
	makeFile(t, base)

	entries, err := rotate.Discover(base)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(entries) != 1 {
		t.Fatalf("expected 1 entry, got %d", len(entries))
	}
	if entries[0].Index != 0 {
		t.Errorf("expected index 0 for current file, got %d", entries[0].Index)
	}
}

func TestDiscover_WithRotatedFiles(t *testing.T) {
	dir := t.TempDir()
	base := filepath.Join(dir, "app.log")
	makeFile(t, base)
	makeFile(t, filepath.Join(dir, "app.log.1"))
	makeFile(t, filepath.Join(dir, "app.log.2"))
	makeFile(t, filepath.Join(dir, "app.log.3.gz"))

	entries, err := rotate.Discover(base)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	// Expect: index 3, 2, 1, 0
	if len(entries) != 4 {
		t.Fatalf("expected 4 entries, got %d", len(entries))
	}
	expected := []int{3, 2, 1, 0}
	for i, e := range entries {
		if e.Index != expected[i] {
			t.Errorf("entry[%d]: expected index %d, got %d", i, expected[i], e.Index)
		}
	}
}

func TestDiscover_BaseFileMissing(t *testing.T) {
	_, err := rotate.Discover("/nonexistent/path/app.log")
	if err == nil {
		t.Fatal("expected error for missing base file, got nil")
	}
}

func TestDiscover_IgnoresNonNumericSuffixes(t *testing.T) {
	dir := t.TempDir()
	base := filepath.Join(dir, "app.log")
	makeFile(t, base)
	makeFile(t, filepath.Join(dir, "app.log.bak"))
	makeFile(t, filepath.Join(dir, "app.log.old"))
	makeFile(t, filepath.Join(dir, "app.log.1"))

	entries, err := rotate.Discover(base)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(entries) != 2 {
		t.Fatalf("expected 2 entries (1 rotated + current), got %d", len(entries))
	}
}
