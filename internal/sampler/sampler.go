// Package sampler provides log line sampling to reduce output volume
// when processing very large log files. It supports two strategies:
// rate-based sampling (keep every Nth line) and reservoir sampling
// (keep a fixed random subset of lines).
package sampler

import (
	"fmt"
	"math/rand"
)

// Strategy controls how lines are sampled.
type Strategy int

const (
	// RateStrategy keeps every Nth matched line.
	RateStrategy Strategy = iota
	// ReservoirStrategy keeps a random fixed-size subset.
	ReservoirStrategy
)

// Sampler decides which lines to keep.
type Sampler struct {
	strategy  Strategy
	rate      int
	maxLines  int
	counter   int
	reservoir []string
	rng       *rand.Rand
}

// Option configures a Sampler.
type Option func(*Sampler)

// WithRate sets the keep-every-N rate (RateStrategy).
func WithRate(n int) Option {
	return func(s *Sampler) {
		s.rate = n
		s.strategy = RateStrategy
	}
}

// WithReservoir sets the reservoir size (ReservoirStrategy).
func WithReservoir(size int, seed int64) Option {
	return func(s *Sampler) {
		s.maxLines = size
		s.strategy = ReservoirStrategy
		s.reservoir = make([]string, 0, size)
		s.rng = rand.New(rand.NewSource(seed))
	}
}

// New creates a Sampler with the given options.
// Defaults to RateStrategy with rate=1 (keep all lines).
func New(opts ...Option) (*Sampler, error) {
	s := &Sampler{
		strategy: RateStrategy,
		rate:     1,
	}
	for _, o := range opts {
		o(s)
	}
	if s.strategy == RateStrategy && s.rate < 1 {
		return nil, fmt.Errorf("sampler: rate must be >= 1, got %d", s.rate)
	}
	if s.strategy == ReservoirStrategy && s.maxLines < 1 {
		return nil, fmt.Errorf("sampler: reservoir size must be >= 1, got %d", s.maxLines)
	}
	return s, nil
}

// Keep returns true if the line should be included in output.
// For ReservoirStrategy, lines are buffered; call Flush to retrieve them.
func (s *Sampler) Keep(line string) bool {
	s.counter++
	switch s.strategy {
	case RateStrategy:
		return s.counter%s.rate == 0
	case ReservoirStrategy:
		if len(s.reservoir) < s.maxLines {
			s.reservoir = append(s.reservoir, line)
			return false // buffered, not emitted directly
		}
		i := s.rng.Intn(s.counter)
		if i < s.maxLines {
			s.reservoir[i] = line
		}
		return false
	}
	return true
}

// Flush returns all reservoir-sampled lines and resets state.
// For RateStrategy it always returns nil.
func (s *Sampler) Flush() []string {
	if s.strategy != ReservoirStrategy {
		return nil
	}
	out := make([]string, len(s.reservoir))
	copy(out, s.reservoir)
	s.reservoir = s.reservoir[:0]
	s.counter = 0
	return out
}

// Count returns the total number of lines seen so far.
func (s *Sampler) Count() int {
	return s.counter
}
