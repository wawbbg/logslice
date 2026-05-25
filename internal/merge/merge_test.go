package merge

import (
	"bytes"
	"strings"
	"testing"
)

func lines(s string) []string {
	return strings.Split(strings.TrimRight(s, "\n"), "\n")
}

func TestMerge_TwoSortedStreams(t *testing.T) {
	a := strings.NewReader(
		"2024-01-01T10:00:00Z INFO  start\n" +
			"2024-01-01T10:00:02Z INFO  step-a\n" +
			"2024-01-01T10:00:04Z INFO  done-a\n",
	)
	b := strings.NewReader(
		"2024-01-01T10:00:01Z INFO  step-b\n" +
			"2024-01-01T10:00:03Z INFO  step-b2\n" +
			"2024-01-01T10:00:05Z INFO  done-b\n",
	)
	var buf bytes.Buffer
	if err := Merge(&buf, []interface{ Read([]byte) (int, error) }{a, b}); err != nil {
		t.Fatalf("Merge returned error: %v", err)
	}
	got := lines(buf.String())
	want := []string{
		"2024-01-01T10:00:00Z INFO  start",
		"2024-01-01T10:00:01Z INFO  step-b",
		"2024-01-01T10:00:02Z INFO  step-a",
		"2024-01-01T10:00:03Z INFO  step-b2",
		"2024-01-01T10:00:04Z INFO  done-a",
		"2024-01-01T10:00:05Z INFO  done-b",
	}
	if len(got) != len(want) {
		t.Fatalf("expected %d lines, got %d: %v", len(want), len(got), got)
	}
	for i := range want {
		if got[i] != want[i] {
			t.Errorf("line %d: want %q, got %q", i, want[i], got[i])
		}
	}
}

func TestMerge_SingleReader(t *testing.T) {
	a := strings.NewReader(
		"2024-01-01T09:00:00Z DEBUG only\n" +
			"2024-01-01T09:00:01Z DEBUG second\n",
	)
	var buf bytes.Buffer
	if err := Merge(&buf, []interface{ Read([]byte) (int, error) }{a}); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	got := lines(buf.String())
	if len(got) != 2 {
		t.Fatalf("expected 2 lines, got %d", len(got))
	}
}

func TestMerge_EmptyReaders(t *testing.T) {
	var buf bytes.Buffer
	if err := Merge(&buf, []interface{ Read([]byte) (int, error) }{
		strings.NewReader(""),
		strings.NewReader(""),
	}); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if buf.Len() != 0 {
		t.Errorf("expected empty output, got %q", buf.String())
	}
}

func TestMerge_UnparseableLinesPassThrough(t *testing.T) {
	a := strings.NewReader("no timestamp here\n" + "2024-01-01T08:00:00Z INFO ok\n")
	var buf bytes.Buffer
	if err := Merge(&buf, []interface{ Read([]byte) (int, error) }{a}); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !strings.Contains(buf.String(), "no timestamp here") {
		t.Errorf("expected unparseable line to be present in output")
	}
	if !strings.Contains(buf.String(), "INFO ok") {
		t.Errorf("expected parseable line to be present in output")
	}
}
