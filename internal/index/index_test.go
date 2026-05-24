package index

import (
	"strings"
	"testing"
	"time"
)

// sampleLog contains three lines with RFC3339 timestamps.
const sampleLog = `2024-01-10T08:00:00Z INFO  server started
2024-01-10T09:00:00Z INFO  request received
2024-01-10T10:00:00Z ERROR something went wrong
`

func TestBuild_ParsesTimestamps(t *testing.T) {
	r := strings.NewReader(sampleLog)
	idx, err := Build(r)
	if err != nil {
		t.Fatalf("Build: %v", err)
	}
	if idx.Len() != 3 {
		t.Fatalf("expected 3 entries, got %d", idx.Len())
	}
}

func TestBuild_SkipsUnparseable(t *testing.T) {
	log := "no timestamp here\n2024-01-10T08:00:00Z INFO ok\n"
	r := strings.NewReader(log)
	idx, err := Build(r)
	if err != nil {
		t.Fatalf("Build: %v", err)
	}
	if idx.Len() != 1 {
		t.Fatalf("expected 1 entry, got %d", idx.Len())
	}
}

func TestFindStart_ReturnsCorrectOffset(t *testing.T) {
	r := strings.NewReader(sampleLog)
	idx, _ := Build(r)

	// Requesting start at 09:00 should skip the first line.
	from := time.Date(2024, 1, 10, 9, 0, 0, 0, time.UTC)
	offset := idx.FindStart(from)

	// First line is "2024-01-10T08:00:00Z INFO  server started\n" = 42 bytes
	expected := int64(len("2024-01-10T08:00:00Z INFO  server started\n"))
	if offset != expected {
		t.Errorf("FindStart offset = %d, want %d", offset, expected)
	}
}

func TestFindStart_BeforeAllEntries(t *testing.T) {
	r := strings.NewReader(sampleLog)
	idx, _ := Build(r)

	from := time.Date(2024, 1, 1, 0, 0, 0, 0, time.UTC)
	if offset := idx.FindStart(from); offset != 0 {
		t.Errorf("expected offset 0, got %d", offset)
	}
}

func TestFindEnd_ReturnsMinusOneWhenAllInRange(t *testing.T) {
	r := strings.NewReader(sampleLog)
	idx, _ := Build(r)

	to := time.Date(2024, 12, 31, 0, 0, 0, 0, time.UTC)
	if end := idx.FindEnd(to); end != -1 {
		t.Errorf("expected -1, got %d", end)
	}
}

func TestFindEnd_StopsAtBoundary(t *testing.T) {
	r := strings.NewReader(sampleLog)
	idx, _ := Build(r)

	// Only include lines before 10:00; the third line should mark the end.
	to := time.Date(2024, 1, 10, 9, 30, 0, 0, time.UTC)
	end := idx.FindEnd(to)
	if end <= 0 {
		t.Errorf("expected positive end offset, got %d", end)
	}
}

func TestEntries_ReturnsCopy(t *testing.T) {
	r := strings.NewReader(sampleLog)
	idx, _ := Build(r)

	e1 := idx.Entries()
	e1[0].Offset = 9999
	e2 := idx.Entries()
	if e2[0].Offset == 9999 {
		t.Error("Entries should return a copy, not a reference")
	}
}
