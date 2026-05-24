// Package ratelimit provides a token-bucket rate limiter for controlling
// the number of log lines emitted per second during streaming output.
package ratelimit

import (
	"sync"
	"time"
)

// Limiter is a token-bucket rate limiter that controls throughput of log lines.
type Limiter struct {
	mu       sync.Mutex
	rate     int           // tokens (lines) per second
	tokens   float64       // current token count
	max      float64       // bucket capacity
	lastTick time.Time
	clock    func() time.Time
}

// Option configures a Limiter.
type Option func(*Limiter)

// WithRate sets the maximum number of lines allowed per second.
// A rate of 0 means unlimited.
func WithRate(linesPerSecond int) Option {
	return func(l *Limiter) {
		l.rate = linesPerSecond
		l.max = float64(linesPerSecond)
	}
}

// withClock overrides the clock used for token refills (for testing).
func withClock(fn func() time.Time) Option {
	return func(l *Limiter) {
		l.clock = fn
	}
}

// New creates a Limiter with the given options.
// By default the rate is unlimited (0).
func New(opts ...Option) *Limiter {
	l := &Limiter{
		clock: time.Now,
	}
	for _, o := range opts {
		o(l)
	}
	if l.rate > 0 {
		l.tokens = float64(l.rate)
		l.max = float64(l.rate)
	}
	l.lastTick = l.clock()
	return l
}

// Allow reports whether a single log line may be emitted right now.
// If the rate is 0 (unlimited), Allow always returns true immediately.
// Otherwise it blocks until a token is available.
func (l *Limiter) Allow() bool {
	if l.rate == 0 {
		return true
	}
	l.mu.Lock()
	defer l.mu.Unlock()

	now := l.clock()
	elapsed := now.Sub(l.lastTick).Seconds()
	l.lastTick = now

	l.tokens += elapsed * float64(l.rate)
	if l.tokens > l.max {
		l.tokens = l.max
	}

	if l.tokens >= 1 {
		l.tokens--
		return true
	}
	return false
}

// Rate returns the configured lines-per-second limit (0 means unlimited).
func (l *Limiter) Rate() int {
	return l.rate
}
