package rotate_test

import (
	"bufio"
	"path/filepath"
	"strings"
	"testing"

	"github.com/yourorg/logslice/internal/rotate"
)

func TestMultiReader_ConcatenatesFiles(t *testing.T) {
	dir := t.TempDir()

	writeLog := func(name, content string) string {
		p := filepath.Join(dir, name)
		makeFile(t, p)
		// overwrite with actual content
		if err := writeContent(p, content); err != nil {
			t.Fatalf("writeLog: %v", err)
		}
		return p
	}

	p1 := writeLog("app.log.2", "line-old1\nline-old2\n")
	p2 := writeLog("app.log.1", "line-mid\n")
	p3 := writeLog("app.log", "line-new\n")

	entries := []rotate.Entry{
		{Path: p1, Index: 2},
		{Path: p2, Index: 1},
		{Path: p3, Index: 0},
	}

	mr := rotate.NewMultiReader(entries)
	defer mr.Close()

	scanner := bufio.NewScanner(mr)
	var lines []string
	for scanner.Scan() {
		lines = append(lines, scanner.Text())
	}
	if err := scanner.Err(); err != nil {
		t.Fatalf("scanner error: %v", err)
	}

	expected := []string{"line-old1", "line-old2", "line-mid", "line-new"}
	if len(lines) != len(expected) {
		t.Fatalf("expected %d lines, got %d: %v", len(expected), len(lines), lines)
	}
	for i, l := range lines {
		if l != expected[i] {
			t.Errorf("line[%d]: expected %q, got %q", i, expected[i], l)
		}
	}
}

func TestLines_Helper(t *testing.T) {
	dir := t.TempDir()
	p := filepath.Join(dir, "app.log")
	_ = writeContent(p, "alpha\nbeta\ngamma\n")

	entries := []rotate.Entry{{Path: p, Index: 0}}
	scanner := rotate.Lines(entries)
	var got []string
	for scanner.Scan() {
		got = append(got, scanner.Text())
	}
	if strings.Join(got, ",") != "alpha,beta,gamma" {
		t.Errorf("unexpected lines: %v", got)
	}
}

import "os"

func writeContent(path, content string) error {
	return os.WriteFile(path, []byte(content), 0o644)
}
