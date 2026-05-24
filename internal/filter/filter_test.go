package filter

import (
	"strings"
	"testing"
	"time"
)

func mustParse(s string) time.Time {
	t, err := time.Parse(time.RFC3339, s)
	if err != nil {
		panic(err)
	}
	return t
}

var sampleLog = `2024-01-15T10:00:00Z INFO  service started
2024-01-15T10:05:00Z DEBUG request received
2024-01-15T10:10:00Z ERROR connection refused
2024-01-15T10:15:00Z INFO  request completed
2024-01-15T10:20:00Z WARN  high memory usage
`

func TestFilter_TimeRange(t *testing.T) {
	opts := Options{
		From: mustParse("2024-01-15T10:04:00Z"),
		To:   mustParse("2024-01-15T10:12:00Z"),
	}

	out := &strings.Builder{}
	res, err := Filter(strings.NewReader(sampleLog), out, opts)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if res.MatchedLines != 2 {
		t.Errorf("expected 2 matched lines, got %d", res.MatchedLines)
	}
	if res.TotalLines != 5 {
		t.Errorf("expected 5 total lines, got %d", res.TotalLines)
	}
}

func TestFilter_WithPattern(t *testing.T) {
	opts := Options{
		From:    mustParse("2024-01-15T09:00:00Z"),
		To:      mustParse("2024-01-15T11:00:00Z"),
		Pattern: "ERROR",
	}

	out := &strings.Builder{}
	res, err := Filter(strings.NewReader(sampleLog), out, opts)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if res.MatchedLines != 1 {
		t.Errorf("expected 1 matched line, got %d", res.MatchedLines)
	}
	if !strings.Contains(out.String(), "connection refused") {
		t.Errorf("expected output to contain 'connection refused'")
	}
}

func TestFilter_NoMatches(t *testing.T) {
	opts := Options{
		From: mustParse("2025-01-01T00:00:00Z"),
		To:   mustParse("2025-01-02T00:00:00Z"),
	}

	out := &strings.Builder{}
	res, err := Filter(strings.NewReader(sampleLog), out, opts)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if res.MatchedLines != 0 {
		t.Errorf("expected 0 matched lines, got %d", res.MatchedLines)
	}
}

func TestContainsPattern(t *testing.T) {
	cases := []struct {
		line, pattern string
		want          bool
	}{
		{"hello world", "world", true},
		{"hello world", "xyz", false},
		{"ERROR: disk full", "ERROR", true},
		{"", "x", false},
	}
	for _, c := range cases {
		got := containsPattern(c.line, c.pattern)
		if got != c.want {
			t.Errorf("containsPattern(%q, %q) = %v, want %v", c.line, c.pattern, got, c.want)
		}
	}
}
