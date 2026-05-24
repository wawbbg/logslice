package parser

import (
	"fmt"
	"time"
)

// Common log timestamp formats to attempt parsing
var knownFormats = []string{
	time.RFC3339Nano,
	time.RFC3339,
	"2006-01-02T15:04:05.999999999",
	"2006-01-02T15:04:05",
	"2006-01-02 15:04:05.999999999",
	"2006-01-02 15:04:05",
	"2006/01/02 15:04:05",
	"02/Jan/2006:15:04:05 -0700",
}

// ParseTimestamp attempts to parse a timestamp string using known formats.
// Returns the parsed time and the matched format string, or an error if none match.
func ParseTimestamp(raw string) (time.Time, string, error) {
	for _, format := range knownFormats {
		t, err := time.Parse(format, raw)
		if err == nil {
			return t, format, nil
		}
	}
	return time.Time{}, "", fmt.Errorf("parser: unrecognized timestamp format: %q", raw)
}

// InRange reports whether t falls within [start, end] (inclusive).
// A zero start or end is treated as unbounded on that side.
func InRange(t, start, end time.Time) bool {
	if !start.IsZero() && t.Before(start) {
		return false
	}
	if !end.IsZero() && t.After(end) {
		return false
	}
	return true
}
