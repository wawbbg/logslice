package sampler

import (
	"testing"
)

func TestNew_DefaultsToRateOne(t *testing.T) {
	s, err := New()
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if s.rate != 1 {
		t.Errorf("expected rate 1, got %d", s.rate)
	}
	if s.strategy != RateStrategy {
		t.Errorf("expected RateStrategy")
	}
}

func TestNew_InvalidRate(t *testing.T) {
	_, err := New(WithRate(0))
	if err == nil {
		t.Fatal("expected error for rate=0")
	}
}

func TestNew_InvalidReservoir(t *testing.T) {
	_, err := New(WithReservoir(0, 42))
	if err == nil {
		t.Fatal("expected error for reservoir size=0")
	}
}

func TestRateStrategy_KeepsEveryN(t *testing.T) {
	s, _ := New(WithRate(3))
	results := make([]bool, 9)
	for i := range results {
		results[i] = s.Keep("line")
	}
	// counter goes 1..9; keep when counter%3==0 => indices 2,5,8
	expected := []bool{false, false, true, false, false, true, false, false, true}
	for i, got := range results {
		if got != expected[i] {
			t.Errorf("index %d: expected %v got %v", i, expected[i], got)
		}
	}
}

func TestRateStrategy_RateOne_KeepsAll(t *testing.T) {
	s, _ := New(WithRate(1))
	for i := 0; i < 10; i++ {
		if !s.Keep("line") {
			t.Errorf("rate=1 should keep all lines, failed at i=%d", i)
		}
	}
}

func TestReservoirStrategy_BuffersLines(t *testing.T) {
	s, _ := New(WithReservoir(3, 1))
	lines := []string{"a", "b", "c", "d", "e"}
	for _, l := range lines {
		if s.Keep(l) {
			t.Error("ReservoirStrategy.Keep should always return false")
		}
	}
	if s.Count() != 5 {
		t.Errorf("expected count 5, got %d", s.Count())
	}
}

func TestReservoirStrategy_FlushReturnsLines(t *testing.T) {
	s, _ := New(WithReservoir(3, 99))
	for _, l := range []string{"x", "y", "z"} {
		s.Keep(l)
	}
	out := s.Flush()
	if len(out) != 3 {
		t.Fatalf("expected 3 lines, got %d", len(out))
	}
}

func TestReservoirStrategy_FlushResetsState(t *testing.T) {
	s, _ := New(WithReservoir(5, 7))
	for i := 0; i < 10; i++ {
		s.Keep("line")
	}
	s.Flush()
	if s.Count() != 0 {
		t.Errorf("expected count 0 after flush, got %d", s.Count())
	}
	if len(s.reservoir) != 0 {
		t.Errorf("expected empty reservoir after flush")
	}
}

func TestRateStrategy_FlushReturnsNil(t *testing.T) {
	s, _ := New(WithRate(2))
	s.Keep("a")
	if s.Flush() != nil {
		t.Error("RateStrategy Flush should return nil")
	}
}
