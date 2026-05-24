package dedupe

import (
	"fmt"
	"testing"
)

func TestNew_Defaults(t *testing.T) {
	d := New()
	if d.windowSize != 256 {
		t.Fatalf("expected windowSize 256, got %d", d.windowSize)
	}
	if !d.caseSensitive {
		t.Fatal("expected caseSensitive true by default")
	}
}

func TestIsDuplicate_FirstOccurrence(t *testing.T) {
	d := New()
	if d.IsDuplicate("hello world") {
		t.Fatal("first occurrence should not be a duplicate")
	}
}

func TestIsDuplicate_SecondOccurrence(t *testing.T) {
	d := New()
	d.IsDuplicate("hello world")
	if !d.IsDuplicate("hello world") {
		t.Fatal("second occurrence should be a duplicate")
	}
}

func TestIsDuplicate_DifferentLines(t *testing.T) {
	d := New()
	d.IsDuplicate("line one")
	if d.IsDuplicate("line two") {
		t.Fatal("different line should not be a duplicate")
	}
}

func TestIsDuplicate_CaseSensitive(t *testing.T) {
	d := New(WithCaseSensitive(true))
	d.IsDuplicate("ERROR: disk full")
	if d.IsDuplicate("error: disk full") {
		t.Fatal("case-sensitive mode: different case should not be duplicate")
	}
}

func TestIsDuplicate_CaseInsensitive(t *testing.T) {
	d := New(WithCaseSensitive(false))
	d.IsDuplicate("ERROR: disk full")
	if !d.IsDuplicate("error: disk full") {
		t.Fatal("case-insensitive mode: different case should be duplicate")
	}
}

func TestWindowEviction(t *testing.T) {
	const window = 4
	d := New(WithWindowSize(window))

	// Fill the window with unique lines.
	for i := 0; i < window; i++ {
		d.IsDuplicate(fmt.Sprintf("line-%d", i))
	}
	// "line-0" should have been evicted; adding it again is not a duplicate.
	if d.IsDuplicate("line-0") {
		t.Fatal("line-0 should have been evicted from the window")
	}
}

func TestReset_ClearsState(t *testing.T) {
	d := New()
	d.IsDuplicate("persistent line")
	d.Reset()
	if d.IsDuplicate("persistent line") {
		t.Fatal("after Reset, line should not be considered duplicate")
	}
}

func TestWithWindowSize_InvalidIgnored(t *testing.T) {
	d := New(WithWindowSize(-1))
	if d.windowSize != 256 {
		t.Fatalf("invalid window size should be ignored, got %d", d.windowSize)
	}
}

func BenchmarkIsDuplicate(b *testing.B) {
	d := New()
	lines := make([]string, 1000)
	for i := range lines {
		lines[i] = fmt.Sprintf("2024-01-01T00:00:00Z INFO request id=%d status=200", i)
	}
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		d.IsDuplicate(lines[i%len(lines)])
	}
}
