package parser

import (
	"testing"
	"time"
)

func TestParseTimestamp_KnownFormats(t *testing.T) {
	cases := []struct {
		input  string
		wantOK bool
	}{
		{"2024-03-15T12:00:00Z", true},
		{"2024-03-15T12:00:00.123456789Z", true},
		{"2024-03-15T12:00:00", true},
		{"2024-03-15 12:00:00", true},
		{"2024/03/15 12:00:00", true},
		{"15/Mar/2024:12:00:00 +0000", true},
		{"not-a-timestamp", false},
		{"", false},
	}

	for _, tc := range cases {
		t.Run(tc.input, func(t *testing.T) {
			got, format, err := ParseTimestamp(tc.input)
			if tc.wantOK {
				if err != nil {
					t.Errorf("expected success, got error: %v", err)
				}
				if got.IsZero() {
					t.Error("expected non-zero time")
				}
				if format == "" {
					t.Error("expected non-empty format string")
				}
			} else {
				if err == nil {
					t.Errorf("expected error for input %q, got time %v", tc.input, got)
				}
			}
		})
	}
}

// TestParseTimestamp_Roundtrip verifies that re-formatting a parsed timestamp
// using the returned format string produces the original input string.
func TestParseTimestamp_Roundtrip(t *testing.T) {
	inputs := []string{
		"2024-03-15T12:00:00Z",
		"2024-03-15 12:00:00",
		"2024/03/15 12:00:00",
		"15/Mar/2024:12:00:00 +0000",
	}

	for _, input := range inputs {
		t.Run(input, func(t *testing.T) {
			parsed, format, err := ParseTimestamp(input)
			if err != nil {
				t.Fatalf("unexpected parse error: %v", err)
			}
			if got := parsed.Format(format); got != input {
				t.Errorf("roundtrip mismatch: Format(%q) = %q, want %q", format, got, input)
			}
		})
	}
}

func TestInRange(t *testing.T) {
	base := time.Date(2024, 3, 15, 12, 0, 0, 0, time.UTC)
	before := base.Add(-time.Hour)
	after := base.Add(time.Hour)

	cases := []struct {
		name  string
		t     time.Time
		start time.Time
		end   time.Time
		want  bool
	}{
		{"within range", base, before, after, true},
		{"on start boundary", before, before, after, true},
		{"on end boundary", after, before, after, true},
		{"before start", before.Add(-time.Minute), base, after, false},
		{"after end", after.Add(time.Minute), before, base, false},
		{"no start bound", before, time.Time{}, base, true},
		{"no end bound", after, base, time.Time{}, true},
		{"unbounded", base, time.Time{}, time.Time{}, true},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			if got := InRange(tc.t, tc.start, tc.end); got != tc.want {
				t.Errorf("InRange() = %v, want %v", got, tc.want)
			}
		})
	}
}
