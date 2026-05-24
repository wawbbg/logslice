package filter

import (
	"compress/gzip"
	"os"
	"strings"
	"testing"
	"time"
)

func TestFilterFile_PlainText(t *testing.T) {
	f, err := os.CreateTemp("", "logslice-*.log")
	if err != nil {
		t.Fatal(err)
	}
	defer os.Remove(f.Name())

	_, _ = f.WriteString("2024-03-01T08:00:00Z INFO  boot\n")
	_, _ = f.WriteString("2024-03-01T09:00:00Z ERROR crash\n")
	f.Close()

	opts := Options{
		From: time.Date(2024, 3, 1, 7, 0, 0, 0, time.UTC),
		To:   time.Date(2024, 3, 1, 10, 0, 0, 0, time.UTC),
	}

	out := &strings.Builder{}
	res, err := FilterFile(f.Name(), out, opts)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if res.MatchedLines != 2 {
		t.Errorf("expected 2 matched lines, got %d", res.MatchedLines)
	}
}

func TestFilterFile_GzipCompressed(t *testing.T) {
	f, err := os.CreateTemp("", "logslice-*.log.gz")
	if err != nil {
		t.Fatal(err)
	}
	defer os.Remove(f.Name())

	gw := gzip.NewWriter(f)
	_, _ = gw.Write([]byte("2024-03-01T08:00:00Z INFO  compressed log entry\n"))
	gw.Close()
	f.Close()

	opts := Options{
		From: time.Date(2024, 3, 1, 0, 0, 0, 0, time.UTC),
		To:   time.Date(2024, 3, 2, 0, 0, 0, 0, time.UTC),
	}

	out := &strings.Builder{}
	res, err := FilterFile(f.Name(), out, opts)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if res.MatchedLines != 1 {
		t.Errorf("expected 1 matched line, got %d", res.MatchedLines)
	}
	if !strings.Contains(out.String(), "compressed log entry") {
		t.Errorf("expected output to contain 'compressed log entry'")
	}
}

func TestOpenFile_NotFound(t *testing.T) {
	_, err := OpenFile("/nonexistent/path/to/file.log")
	if err == nil {
		t.Error("expected error for missing file, got nil")
	}
}
