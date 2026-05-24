// Package truncate provides utilities for truncating long log lines
// to a configurable maximum byte length, preserving valid UTF-8 boundaries.
package truncate

import "unicode/utf8"

const (
	// DefaultMaxBytes is the default maximum line length in bytes.
	DefaultMaxBytes = 4096
	// ellipsis is appended to truncated lines.
	ellipsis = "..."
)

// Truncator truncates lines that exceed a maximum byte length.
type Truncator struct {
	maxBytes int
	appendEllipsis bool
}

// Option configures a Truncator.
type Option func(*Truncator)

// WithMaxBytes sets the maximum number of bytes per line.
func WithMaxBytes(n int) Option {
	return func(t *Truncator) {
		if n > 0 {
			t.maxBytes = n
		}
	}
}

// WithEllipsis controls whether truncated lines get a trailing "...".
func WithEllipsis(enabled bool) Option {
	return func(t *Truncator) {
		t.appendEllipsis = enabled
	}
}

// New returns a Truncator with the given options.
// By default it uses DefaultMaxBytes and appends an ellipsis.
func New(opts ...Option) *Truncator {
	t := &Truncator{
		maxBytes:       DefaultMaxBytes,
		appendEllipsis: true,
	}
	for _, o := range opts {
		o(t)
	}
	return t
}

// Line returns the line unchanged if it fits within maxBytes, otherwise
// it clips the line at a valid UTF-8 rune boundary and optionally appends
// an ellipsis.
func (t *Truncator) Line(line string) string {
	if len(line) <= t.maxBytes {
		return line
	}
	cutoff := t.maxBytes
	if t.appendEllipsis && cutoff > len(ellipsis) {
		cutoff -= len(ellipsis)
	}
	// Walk back to a valid UTF-8 boundary.
	for cutoff > 0 && !utf8.RuneStart(line[cutoff]) {
		cutoff--
	}
	truncated := line[:cutoff]
	if t.appendEllipsis {
		truncated += ellipsis
	}
	return truncated
}

// Lines applies Line to every element of the slice in-place and returns it.
func (t *Truncator) Lines(lines []string) []string {
	for i, l := range lines {
		lines[i] = t.Line(l)
	}
	return lines
}
