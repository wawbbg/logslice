package multipattern

import (
	"testing"
)

func TestNew_EmptyPatterns_MatchesAll(t *testing.T) {
	mt, err := New(nil)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !mt.Match("anything") {
		t.Error("expected empty matcher to match all lines")
	}
}

func TestAND_BothPatternsMustMatch(t *testing.T) {
	mt, _ := New([]string{"error", "timeout"})
	if !mt.Match("error: connection timeout") {
		t.Error("expected match when both patterns present")
	}
	if mt.Match("error: unknown") {
		t.Error("expected no match when only one pattern present")
	}
}

func TestOR_EitherPatternSuffices(t *testing.T) {
	mt, _ := New([]string{"warn", "error"}, WithMode(ModeOR))
	if !mt.Match("warn: disk low") {
		t.Error("expected match on 'warn'")
	}
	if !mt.Match("error: crash") {
		t.Error("expected match on 'error'")
	}
	if mt.Match("info: all good") {
		t.Error("expected no match when neither pattern present")
	}
}

func TestNegated_ExcludesMatchingLines(t *testing.T) {
	mt, _ := New([]string{"error", "!timeout"})
	if !mt.Match("error: disk full") {
		t.Error("expected match: has error, no timeout")
	}
	if mt.Match("error: connection timeout") {
		t.Error("expected no match: negated pattern 'timeout' hit")
	}
}

func TestRegex_Pattern(t *testing.T) {
	mt, err := New([]string{"/ERR[0-9]+/"})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !mt.Match("ERR42: something went wrong") {
		t.Error("expected regex match")
	}
	if mt.Match("ERR: no code") {
		t.Error("expected no regex match")
	}
}

func TestRegex_InvalidPattern_ReturnsError(t *testing.T) {
	_, err := New([]string{"/[invalid/"})
	if err == nil {
		t.Error("expected error for invalid regex")
	}
}

func TestCaseInsensitive_PlainPattern(t *testing.T) {
	mt, _ := New([]string{"ERROR"})
	if !mt.Match("error: something") {
		t.Error("expected case-insensitive match")
	}
	if !mt.Match("Error: something") {
		t.Error("expected case-insensitive match for mixed case")
	}
}

func TestNegatedRegex(t *testing.T) {
	mt, err := New([]string{"!/debug/"})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if mt.Match("debug: verbose output") {
		t.Error("expected no match: negated regex hit")
	}
	if !mt.Match("info: startup complete") {
		t.Error("expected match: negated regex not hit")
	}
}
