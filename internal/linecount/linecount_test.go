package linecount_test

import (
	"bytes"
	"compress/gzip"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/yourorg/logslice/internal/linecount"
)

func writePlain(t *testing.T, lines []string) string {
	t.Helper()
	f, err := os.CreateTemp(t.TempDir(), "log-*.log")
	if err != nil {
		t.Fatalf("create temp: %v", err)
	}
	defer f.Close()
	for _, l := range lines {
		f.WriteString(l + "\n")
	}
	return f.Name()
}

func writeGzip(t *testing.T, lines []string) string {
	t.Helper()
	path := filepath.Join(t.TempDir(), "log.gz")
	f, err := os.Create(path)
	if err != nil {
		t.Fatalf("create gzip: %v", err)
	}
	defer f.Close()
	gw := gzip.NewWriter(f)
	for _, l := range lines {
		gw.Write([]byte(l + "\n"))
	}
	gw.Close()
	return path
}

func TestCountFile_AllLines(t *testing.T) {
	lines := []string{"alpha", "beta", "gamma"}
	path := writePlain(t, lines)
	res, err := linecount.CountFile(path, "")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if res.Total != 3 || res.Matched != 3 || res.Skipped != 0 {
		t.Errorf("got %+v, want Total=3 Matched=3 Skipped=0", res)
	}
}

func TestCountFile_WithPattern(t *testing.T) {
	lines := []string{"error: disk full", "info: all good", "error: oom"}
	path := writePlain(t, lines)
	res, err := linecount.CountFile(path, "error")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if res.Total != 3 || res.Matched != 2 || res.Skipped != 1 {
		t.Errorf("got %+v, want Total=3 Matched=2 Skipped=1", res)
	}
}

func TestCountFile_GzipCompressed(t *testing.T) {
	lines := []string{"line1", "line2", "line3", "line4"}
	path := writeGzip(t, lines)
	res, err := linecount.CountFile(path, "")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if res.Total != 4 {
		t.Errorf("got Total=%d, want 4", res.Total)
	}
}

func TestCountFile_NotFound(t *testing.T) {
	_, err := linecount.CountFile("/nonexistent/path.log", "")
	if err == nil {
		t.Fatal("expected error for missing file")
	}
}

func TestCountReader_EmptyInput(t *testing.T) {
	res, err := linecount.CountReader(strings.NewReader(""), "")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if res.Total != 0 {
		t.Errorf("got Total=%d, want 0", res.Total)
	}
}

func TestCountReader_NoPattern(t *testing.T) {
	input := bytes.NewBufferString("a\nb\nc\n")
	res, err := linecount.CountReader(input, "")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if res.Total != 3 || res.Matched != 3 {
		t.Errorf("got %+v", res)
	}
}
