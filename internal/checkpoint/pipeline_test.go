package checkpoint

import (
	"io"
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func writeTempLog(t *testing.T, content string) string {
	t.Helper()
	f, err := os.CreateTemp(t.TempDir(), "log-*.log")
	if err != nil {
		t.Fatal(err)
	}
	_, _ = io.WriteString(f, content)
	_ = f.Close()
	return f.Name()
}

func TestLines_ReadsAllLines(t *testing.T) {
	logFile := writeTempLog(t, "alpha\nbeta\ngamma\n")
	s := New(filepath.Join(t.TempDir(), "cp.json"))

	rc, start, err := ReaderFrom(s, logFile)
	if err != nil {
		t.Fatalf("ReaderFrom: %v", err)
	}

	var got []string
	if err := Lines(s, logFile, rc, start, 0, func(l string) { got = append(got, l) }); err != nil {
		t.Fatalf("Lines: %v", err)
	}
	if strings.Join(got, ",") != "alpha,beta,gamma" {
		t.Errorf("got %v", got)
	}
}

func TestLines_CheckpointSavedAfterFlushEvery(t *testing.T) {
	logFile := writeTempLog(t, "line1\nline2\nline3\nline4\n")
	dir := t.TempDir()
	s := New(filepath.Join(dir, "cp.json"))

	rc, start, _ := ReaderFrom(s, logFile)
	_ = Lines(s, logFile, rc, start, 2, func(string) {})

	e, err := s.Load(logFile)
	if err != nil {
		t.Fatalf("Load: %v", err)
	}
	if e.Offset == 0 {
		t.Error("expected non-zero offset after processing")
	}
}

func TestReaderFrom_ResumesFromCheckpoint(t *testing.T) {
	content := "first\nsecond\nthird\n"
	logFile := writeTempLog(t, content)
	dir := t.TempDir()
	s := New(filepath.Join(dir, "cp.json"))

	// Save checkpoint past the first line ("first\n" = 6 bytes).
	_ = s.Save(logFile, 6)

	rc, start, err := ReaderFrom(s, logFile)
	if err != nil {
		t.Fatalf("ReaderFrom: %v", err)
	}

	var got []string
	_ = Lines(s, logFile, rc, start, 0, func(l string) { got = append(got, l) })

	if len(got) != 2 || got[0] != "second" {
		t.Errorf("expected [second third], got %v", got)
	}
}

func TestReaderFrom_MissingFile_ReturnsError(t *testing.T) {
	s := New(filepath.Join(t.TempDir(), "cp.json"))
	_, _, err := ReaderFrom(s, "/nonexistent/path/app.log")
	if err == nil {
		t.Fatal("expected error for missing log file")
	}
}
