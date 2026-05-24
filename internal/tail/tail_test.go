package tail_test

import (
	"fmt"
	"os"
	"strings"
	"testing"

	"github.com/yourorg/logslice/internal/tail"
)

func writeTempLog(t *testing.T, lines []string) string {
	t.Helper()
	f, err := os.CreateTemp(t.TempDir(), "tail-*.log")
	if err != nil {
		t.Fatalf("create temp file: %v", err)
	}
	defer f.Close()
	_, err = f.WriteString(strings.Join(lines, "\n") + "\n")
	if err != nil {
		t.Fatalf("write temp file: %v", err)
	}
	return f.Name()
}

func makeLines(n int) []string {
	lines := make([]string, n)
	for i := range lines {
		lines[i] = fmt.Sprintf("2024-01-01T00:00:%02d INFO line %d", i%60, i)
	}
	return lines
}

func TestLines_LastN(t *testing.T) {
	all := makeLines(100)
	path := writeTempLog(t, all)

	got, err := tail.Lines(path, 10)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(got) != 10 {
		t.Fatalf("expected 10 lines, got %d", len(got))
	}
	for i, line := range got {
		want := all[90+i]
		if line != want {
			t.Errorf("line[%d]: got %q, want %q", i, line, want)
		}
	}
}

func TestLines_FewerThanN(t *testing.T) {
	all := makeLines(5)
	path := writeTempLog(t, all)

	got, err := tail.Lines(path, 20)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(got) != 5 {
		t.Fatalf("expected 5 lines, got %d", len(got))
	}
}

func TestLines_EmptyFile(t *testing.T) {
	f, err := os.CreateTemp(t.TempDir(), "empty-*.log")
	if err != nil {
		t.Fatal(err)
	}
	f.Close()

	got, err := tail.Lines(f.Name(), 5)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(got) != 0 {
		t.Fatalf("expected 0 lines, got %d", len(got))
	}
}

func TestLines_InvalidN(t *testing.T) {
	path := writeTempLog(t, []string{"line1"})
	_, err := tail.Lines(path, 0)
	if err == nil {
		t.Fatal("expected error for n=0, got nil")
	}
}

func TestLines_FileNotFound(t *testing.T) {
	_, err := tail.Lines("/nonexistent/path/file.log", 5)
	if err == nil {
		t.Fatal("expected error for missing file, got nil")
	}
}

func TestLines_ExactlyN(t *testing.T) {
	all := makeLines(10)
	path := writeTempLog(t, all)

	got, err := tail.Lines(path, 10)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(got) != 10 {
		t.Fatalf("expected 10 lines, got %d", len(got))
	}
}
