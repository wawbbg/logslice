// Package linenum provides utilities for tracking and formatting line numbers
// during log filtering, allowing consumers to annotate output with the original
// source line positions.
package linenum

import "fmt"

// Tracker keeps a running count of lines seen and matched, and can annotate
// each matched line with its original 1-based source line number.
type Tracker struct {
	current int
	matched int
	offset  int
}

// Option is a functional option for Tracker.
type Option func(*Tracker)

// WithOffset sets the starting line number (default 1).
func WithOffset(n int) Option {
	return func(t *Tracker) {
		if n > 0 {
			t.offset = n - 1
		}
	}
}

// New creates a new Tracker. Line numbers start at 1 unless WithOffset is used.
func New(opts ...Option) *Tracker {
	t := &Tracker{offset: 0}
	for _, o := range opts {
		o(t)
	}
	return t
}

// Advance increments the source-line counter and returns the new line number.
func (t *Tracker) Advance() int {
	t.current++
	return t.current + t.offset
}

// RecordMatch increments the matched-line counter.
func (t *Tracker) RecordMatch() {
	t.matched++
}

// Current returns the current source line number (1-based).
func (t *Tracker) Current() int {
	return t.current + t.offset
}

// Matched returns the total number of matched lines recorded.
func (t *Tracker) Matched() int {
	return t.matched
}

// Format returns a formatted prefix string such as "[42] " for use in output.
func (t *Tracker) Format() string {
	return fmt.Sprintf("[%d] ", t.Current())
}

// Annotate prepends the current line number prefix to the given line.
func (t *Tracker) Annotate(line string) string {
	return t.Format() + line
}
