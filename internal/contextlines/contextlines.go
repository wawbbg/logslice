// Package contextlines provides support for extracting surrounding lines
// around matched log entries, similar to grep's -B (before) and -A (after) flags.
package contextlines

// Config holds the number of lines to include before and after a match.
type Config struct {
	Before int
	After  int
}

// Buffer maintains a sliding window of recent lines and pending after-context.
type Buffer struct {
	cfg      Config
	ring     []string
	head     int
	size     int
	aftLeft  int
}

// New creates a new Buffer with the given context configuration.
func New(cfg Config) *Buffer {
	cap := cfg.Before
	if cap < 1 {
		cap = 1
	}
	return &Buffer{
		cfg:  cfg,
		ring: make([]string, cap),
	}
}

// Feed records a line into the before-context ring buffer.
// Returns true if this line should be emitted due to active after-context.
func (b *Buffer) Feed(line string) bool {
	if b.cfg.Before > 0 {
		b.ring[b.head] = line
		b.head = (b.head + 1) % len(b.ring)
		if b.size < len(b.ring) {
			b.size++
		}
	}
	if b.aftLeft > 0 {
		b.aftLeft--
		return true
	}
	return false
}

// Before returns the buffered lines that precede the current position.
func (b *Buffer) Before() []string {
	if b.size == 0 {
		return nil
	}
	out := make([]string, b.size)
	start := (b.head - b.size + len(b.ring)) % len(b.ring)
	for i := 0; i < b.size; i++ {
		out[i] = b.ring[(start+i)%len(b.ring)]
	}
	return out
}

// OnMatch signals that the current line matched; resets after-context counter.
func (b *Buffer) OnMatch() {
	b.aftLeft = b.cfg.After
}

// AfterPending reports whether after-context lines are still expected.
func (b *Buffer) AfterPending() bool {
	return b.aftLeft > 0
}
