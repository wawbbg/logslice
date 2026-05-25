package contextlines

import (
	"testing"
)

func TestNew_InitializesBuffer(t *testing.T) {
	b := New(Config{Before: 3, After: 2})
	if b == nil {
		t.Fatal("expected non-nil buffer")
	}
	if got := b.Before(); len(got) != 0 {
		t.Fatalf("expected empty before-context, got %v", got)
	}
}

func TestBefore_ReturnsUpToNLines(t *testing.T) {
	b := New(Config{Before: 3, After: 0})
	lines := []string{"a", "b", "c", "d"}
	for _, l := range lines {
		b.Feed(l)
	}
	got := b.Before()
	want := []string{"b", "c", "d"}
	if len(got) != len(want) {
		t.Fatalf("expected %v, got %v", want, got)
	}
	for i := range want {
		if got[i] != want[i] {
			t.Errorf("index %d: want %q got %q", i, want[i], got[i])
		}
	}
}

func TestBefore_FewerLinesThanWindow(t *testing.T) {
	b := New(Config{Before: 5, After: 0})
	b.Feed("x")
	b.Feed("y")
	got := b.Before()
	if len(got) != 2 {
		t.Fatalf("expected 2 lines, got %d", len(got))
	}
}

func TestAfterContext_EmittedAfterMatch(t *testing.T) {
	b := New(Config{Before: 0, After: 2})
	b.OnMatch()
	if !b.Feed("post1") {
		t.Error("expected Feed to return true for first after-context line")
	}
	if !b.Feed("post2") {
		t.Error("expected Feed to return true for second after-context line")
	}
	if b.Feed("post3") {
		t.Error("expected Feed to return false after context exhausted")
	}
}

func TestAfterPending_TrueWhileCounterPositive(t *testing.T) {
	b := New(Config{Before: 0, After: 1})
	if b.AfterPending() {
		t.Error("should not be pending before match")
	}
	b.OnMatch()
	if !b.AfterPending() {
		t.Error("should be pending after match")
	}
	b.Feed("line")
	if b.AfterPending() {
		t.Error("should not be pending after context consumed")
	}
}

func TestOnMatch_ResetsCounter(t *testing.T) {
	b := New(Config{Before: 0, After: 2})
	b.OnMatch()
	b.Feed("a")
	// second match resets counter
	b.OnMatch()
	if !b.Feed("b") {
		t.Error("after second match, feed should return true")
	}
	if !b.Feed("c") {
		t.Error("second after-context line should still emit")
	}
	if b.Feed("d") {
		t.Error("third line should not emit")
	}
}
