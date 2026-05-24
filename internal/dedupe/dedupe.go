// Package dedupe provides line-level deduplication for log streams.
// It tracks recently seen lines using a fixed-size ring buffer and
// suppresses consecutive or near-consecutive duplicate entries.
package dedupe

import "hash/fnv"

// Option configures a Deduplicator.
type Option func(*Deduplicator)

// WithWindowSize sets the number of recent line hashes to remember.
// Defaults to 256.
func WithWindowSize(n int) Option {
	return func(d *Deduplicator) {
		if n > 0 {
			d.windowSize = n
		}
	}
}

// WithCaseSensitive controls whether comparison is case-sensitive.
// Defaults to true.
func WithCaseSensitive(cs bool) Option {
	return func(d *Deduplicator) {
		d.caseSensitive = cs
	}
}

// Deduplicator filters duplicate log lines within a sliding window.
type Deduplicator struct {
	windowSize    int
	caseSensitive bool
	ring          []uint64
	pos           int
	seen          map[uint64]struct{}
}

// New creates a new Deduplicator with the given options.
func New(opts ...Option) *Deduplicator {
	d := &Deduplicator{
		windowSize:    256,
		caseSensitive: true,
	}
	for _, o := range opts {
		o(d)
	}
	d.ring = make([]uint64, d.windowSize)
	d.seen = make(map[uint64]struct{}, d.windowSize)
	return d
}

// IsDuplicate reports whether line has been seen within the current window.
// If not a duplicate, the line is recorded.
func (d *Deduplicator) IsDuplicate(line string) bool {
	h := d.hash(line)
	if _, ok := d.seen[h]; ok {
		return true
	}
	// Evict the oldest entry from the ring.
	old := d.ring[d.pos]
	if old != 0 {
		delete(d.seen, old)
	}
	d.ring[d.pos] = h
	d.seen[h] = struct{}{}
	d.pos = (d.pos + 1) % d.windowSize
	return false
}

// Reset clears all recorded state.
func (d *Deduplicator) Reset() {
	d.ring = make([]uint64, d.windowSize)
	d.seen = make(map[uint64]struct{}, d.windowSize)
	d.pos = 0
}

func (d *Deduplicator) hash(s string) uint64 {
	if !d.caseSensitive {
		s = toLower(s)
	}
	h := fnv.New64a()
	_, _ = h.Write([]byte(s))
	return h.Sum64()
}

// toLower is a simple ASCII lower-case to avoid importing strings.
func toLower(s string) string {
	b := []byte(s)
	for i, c := range b {
		if c >= 'A' && c <= 'Z' {
			b[i] = c + 32
		}
	}
	return string(b)
}
