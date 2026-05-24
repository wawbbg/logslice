package highlight_test

import (
	"strings"
	"testing"

	"github.com/yourorg/logslice/internal/highlight"
)

func TestHighlight_Disabled(t *testing.T) {
	h := highlight.New("", false)
	line := "error: something went wrong"
	got := h.Highlight(line, "error")
	if got != line {
		t.Errorf("disabled: expected unchanged line, got %q", got)
	}
}

func TestHighlight_EmptyPattern(t *testing.T) {
	h := highlight.New("", true)
	line := "no pattern here"
	got := h.Highlight(line, "")
	if got != line {
		t.Errorf("empty pattern: expected unchanged line, got %q", got)
	}
}

func TestHighlight_SingleMatch(t *testing.T) {
	h := highlight.New(highlight.Yellow, true)
	line := "2024-01-01 ERROR failed to connect"
	got := h.Highlight(line, "ERROR")

	if !strings.Contains(got, "ERROR") {
		t.Error("highlighted output should still contain the matched text")
	}
	if !strings.Contains(got, highlight.Yellow) {
		t.Error("highlighted output should contain ANSI color code")
	}
	if !strings.Contains(got, highlight.Reset) {
		t.Error("highlighted output should contain ANSI reset code")
	}
}

func TestHighlight_MultipleMatches(t *testing.T) {
	h := highlight.New(highlight.Red, true)
	line := "warn warn warn"
	got := h.Highlight(line, "warn")

	stripped := highlight.StripANSI(got)
	if stripped != line {
		t.Errorf("stripped output mismatch: got %q, want %q", stripped, line)
	}

	count := strings.Count(got, highlight.Red)
	if count != 3 {
		t.Errorf("expected 3 color codes, got %d", count)
	}
}

func TestHighlight_CaseInsensitive(t *testing.T) {
	h := highlight.New(highlight.Cyan, true)
	line := "Error ERROR error"
	got := h.Highlight(line, "error")

	count := strings.Count(got, highlight.Cyan)
	if count != 3 {
		t.Errorf("case-insensitive: expected 3 matches, got %d", count)
	}
}

func TestStripANSI(t *testing.T) {
	cases := []struct {
		input string
		want  string
	}{
		{"plain text", "plain text"},
		{"\033[31mred\033[0m", "red"},
		{"\033[1m\033[33mbold yellow\033[0m end", "bold yellow end"},
		{"", ""},
	}
	for _, tc := range cases {
		got := highlight.StripANSI(tc.input)
		if got != tc.want {
			t.Errorf("StripANSI(%q) = %q, want %q", tc.input, got, tc.want)
		}
	}
}

func TestNew_DefaultColor(t *testing.T) {
	h := highlight.New("", true)
	line := "match this"
	got := h.Highlight(line, "match")
	if !strings.Contains(got, highlight.Yellow) {
		t.Error("default color should be Yellow")
	}
}
