package ratelimit

import (
	"testing"
	"time"
)

func TestNew_DefaultsToUnlimited(t *testing.T) {
	l := New()
	if l.Rate() != 0 {
		t.Fatalf("expected rate 0 (unlimited), got %d", l.Rate())
	}
}

func TestNew_WithRate(t *testing.T) {
	l := New(WithRate(100))
	if l.Rate() != 100 {
		t.Fatalf("expected rate 100, got %d", l.Rate())
	}
}

func TestAllow_UnlimitedAlwaysTrue(t *testing.T) {
	l := New() // rate 0 = unlimited
	for i := 0; i < 1000; i++ {
		if !l.Allow() {
			t.Fatal("unlimited limiter should always allow")
		}
	}
}

func TestAllow_RespectsBucketCapacity(t *testing.T) {
	// Fix the clock so no time passes — tokens never refill.
	fixed := time.Now()
	l := New(
		WithRate(5),
		withClock(func() time.Time { return fixed }),
	)

	allowed := 0
	for i := 0; i < 10; i++ {
		if l.Allow() {
			allowed++
		}
	}
	// Initial bucket has 5 tokens; no time passes, so only 5 allowed.
	if allowed != 5 {
		t.Fatalf("expected 5 allowed, got %d", allowed)
	}
}

func TestAllow_RefillsOverTime(t *testing.T) {
	current := time.Now()
	l := New(
		WithRate(10),
		withClock(func() time.Time { return current }),
	)

	// Drain the bucket.
	for i := 0; i < 10; i++ {
		l.Allow()
	}

	// Advance clock by 1 second — should refill 10 tokens.
	current = current.Add(time.Second)

	allowed := 0
	for i := 0; i < 15; i++ {
		if l.Allow() {
			allowed++
		}
	}
	if allowed != 10 {
		t.Fatalf("expected 10 allowed after refill, got %d", allowed)
	}
}

func TestAllow_BucketDoesNotExceedMax(t *testing.T) {
	current := time.Now()
	l := New(
		WithRate(5),
		withClock(func() time.Time { return current }),
	)

	// Advance clock by 10 seconds — would add 50 tokens, but max is 5.
	current = current.Add(10 * time.Second)

	allowed := 0
	for i := 0; i < 10; i++ {
		if l.Allow() {
			allowed++
		}
	}
	if allowed != 5 {
		t.Fatalf("expected bucket capped at 5, got %d", allowed)
	}
}
