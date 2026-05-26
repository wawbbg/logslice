// Package timewindow provides sliding and tumbling time-window aggregation
// over log streams, grouping matched lines into fixed-duration buckets.
package timewindow

import (
	"time"

	"github.com/logslice/logslice/internal/parser"
)

// Bucket holds all lines that fall within a single time window.
type Bucket struct {
	Start time.Time
	End   time.Time
	Lines []string
	Count int
}

// Options configures the windowing behaviour.
type Options struct {
	Size time.Duration // width of each bucket
}

// Windower groups log lines into time buckets of fixed duration.
type Windower struct {
	size    time.Duration
	buckets []*Bucket
	current *Bucket
}

// New returns a Windower with the given options.
// Size defaults to one minute when zero.
func New(opts Options) *Windower {
	if opts.Size <= 0 {
		opts.Size = time.Minute
	}
	return &Windower{size: opts.Size}
}

// Add attempts to parse a timestamp from line and places it in the
// appropriate bucket. Lines whose timestamps cannot be parsed are silently
// dropped. Returns the bucket the line was added to.
func (w *Windower) Add(line string) *Bucket {
	t, ok := parser.ParseTimestamp(line)
	if !ok {
		return nil
	}

	// Align to bucket boundary.
	start := t.Truncate(w.size)
	end := start.Add(w.size)

	if w.current == nil || !w.current.Start.Equal(start) {
		b := &Bucket{Start: start, End: end}
		w.buckets = append(w.buckets, b)
		w.current = b
	}

	w.current.Lines = append(w.current.Lines, line)
	w.current.Count++
	return w.current
}

// Buckets returns all collected buckets in insertion order.
func (w *Windower) Buckets() []*Bucket {
	out := make([]*Bucket, len(w.buckets))
	copy(out, w.buckets)
	return out
}

// Reset clears all accumulated buckets and resets internal state.
func (w *Windower) Reset() {
	w.buckets = w.buckets[:0]
	w.current = nil
}
