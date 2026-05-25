package linenum

import (
	"strings"
	"testing"
)

func TestAnnotate_MatchedLinesAnnotated(t *testing.T) {
	input := "foo bar\nbaz\nfoo qux\n"
	tr := New()
	results := Annotate(strings.NewReader(input), tr, "foo", false)
	if len(results) != 2 {
		t.Fatalf("expected 2 results, got %d", len(results))
	}
	if results[0].LineNo != 1 {
		t.Errorf("expected LineNo=1, got %d", results[0].LineNo)
	}
	if results[0].Line != "[1] foo bar" {
		t.Errorf("unexpected line: %q", results[0].Line)
	}
	if results[1].LineNo != 3 {
		t.Errorf("expected LineNo=3, got %d", results[1].LineNo)
	}
}

func TestAnnotate_IncludeAll_NonMatchNotAnnotated(t *testing.T) {
	input := "match me\nignore\n"
	tr := New()
	results := Annotate(strings.NewReader(input), tr, "match", true)
	if len(results) != 2 {
		t.Fatalf("expected 2 results, got %d", len(results))
	}
	if !results[0].Matched {
		t.Error("first result should be matched")
	}
	if results[1].Matched {
		t.Error("second result should not be matched")
	}
	if results[1].Line != "ignore" {
		t.Errorf("unmatched line should be unannotated, got %q", results[1].Line)
	}
}

func TestAnnotate_EmptySubstr_AllMatch(t *testing.T) {
	input := "line one\nline two\n"
	tr := New()
	results := Annotate(strings.NewReader(input), tr, "", false)
	if len(results) != 2 {
		t.Fatalf("expected 2 results, got %d", len(results))
	}
	for _, r := range results {
		if !r.Matched {
			t.Errorf("expected all lines matched with empty substr")
		}
	}
}

func TestAnnotate_TrackerMatchedCount(t *testing.T) {
	input := "alpha\nbeta\nalpha again\n"
	tr := New()
	Annotate(strings.NewReader(input), tr, "alpha", false)
	if tr.Matched() != 2 {
		t.Fatalf("expected 2 matched, got %d", tr.Matched())
	}
	if tr.Current() != 3 {
		t.Fatalf("expected current=3, got %d", tr.Current())
	}
}

func TestAnnotate_EmptyInput(t *testing.T) {
	tr := New()
	results := Annotate(strings.NewReader(""), tr, "anything", false)
	if len(results) != 0 {
		t.Fatalf("expected 0 results, got %d", len(results))
	}
}
