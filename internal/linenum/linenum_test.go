package linenum

import (
	"fmt"
	"testing"
)

func TestNew_DefaultsToLineOne(t *testing.T) {
	tr := New()
	if tr.Current() != 0 {
		t.Fatalf("expected initial current=0, got %d", tr.Current())
	}
}

func TestAdvance_IncrementsCounter(t *testing.T) {
	tr := New()
	for i := 1; i <= 5; i++ {
		n := tr.Advance()
		if n != i {
			t.Fatalf("step %d: expected %d, got %d", i, i, n)
		}
	}
}

func TestWithOffset_ShiftsLineNumbers(t *testing.T) {
	tr := New(WithOffset(100))
	n := tr.Advance()
	if n != 100 {
		t.Fatalf("expected 100, got %d", n)
	}
	n = tr.Advance()
	if n != 101 {
		t.Fatalf("expected 101, got %d", n)
	}
}

func TestWithOffset_ZeroIgnored(t *testing.T) {
	tr := New(WithOffset(0))
	n := tr.Advance()
	if n != 1 {
		t.Fatalf("expected 1, got %d", n)
	}
}

func TestRecordMatch_CountsMatches(t *testing.T) {
	tr := New()
	tr.RecordMatch()
	tr.RecordMatch()
	if tr.Matched() != 2 {
		t.Fatalf("expected 2 matches, got %d", tr.Matched())
	}
}

func TestFormat_ReturnsExpectedPrefix(t *testing.T) {
	tr := New()
	tr.Advance()
	got := tr.Format()
	want := "[1] "
	if got != want {
		t.Fatalf("expected %q, got %q", want, got)
	}
}

func TestAnnotate_PrependsLineNumber(t *testing.T) {
	tr := New()
	tr.Advance()
	got := tr.Annotate("hello world")
	want := "[1] hello world"
	if got != want {
		t.Fatalf("expected %q, got %q", want, got)
	}
}

func TestAnnotate_MultipleLines(t *testing.T) {
	tr := New()
	lines := []string{"alpha", "beta", "gamma"}
	for i, line := range lines {
		tr.Advance()
		got := tr.Annotate(line)
		want := fmt.Sprintf("[%d] %s", i+1, line)
		if got != want {
			t.Fatalf("line %d: expected %q, got %q", i+1, want, got)
		}
	}
}
