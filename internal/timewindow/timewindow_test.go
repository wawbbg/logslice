package timewindow

import (
	"fmt"
	"testing"
	"time"
)

// logLine produces a log line with an RFC3339 timestamp prefix.
func logLine(t time.Time, msg string) string {
	return fmt.Sprintf("%s %s", t.UTC().Format(time.RFC3339), msg)
}

func TestNew_DefaultsToOneMinute(t *testing.T) {
	w := New(Options{})
	if w.size != time.Minute {
		t.Fatalf("expected 1m default, got %v", w.size)
	}
}

func TestAdd_UnparseableLine_Dropped(t *testing.T) {
	w := New(Options{Size: time.Minute})
	b := w.Add("no timestamp here")
	if b != nil {
		t.Fatal("expected nil bucket for unparseable line")
	}
	if len(w.Buckets()) != 0 {
		t.Fatal("expected zero buckets")
	}
}

func TestAdd_SingleLine_CreatesBucket(t *testing.T) {
	w := New(Options{Size: time.Minute})
	now := time.Date(2024, 1, 15, 10, 30, 45, 0, time.UTC)
	w.Add(logLine(now, "hello"))

	buckets := w.Buckets()
	if len(buckets) != 1 {
		t.Fatalf("expected 1 bucket, got %d", len(buckets))
	}
	if buckets[0].Count != 1 {
		t.Fatalf("expected count 1, got %d", buckets[0].Count)
	}
	expectedStart := now.Truncate(time.Minute)
	if !buckets[0].Start.Equal(expectedStart) {
		t.Fatalf("expected start %v, got %v", expectedStart, buckets[0].Start)
	}
}

func TestAdd_SameWindow_GroupedTogether(t *testing.T) {
	w := New(Options{Size: time.Minute})
	base := time.Date(2024, 1, 15, 10, 30, 0, 0, time.UTC)
	w.Add(logLine(base.Add(5*time.Second), "a"))
	w.Add(logLine(base.Add(30*time.Second), "b"))
	w.Add(logLine(base.Add(59*time.Second), "c"))

	buckets := w.Buckets()
	if len(buckets) != 1 {
		t.Fatalf("expected 1 bucket, got %d", len(buckets))
	}
	if buckets[0].Count != 3 {
		t.Fatalf("expected count 3, got %d", buckets[0].Count)
	}
}

func TestAdd_DifferentWindows_MultipleBuckets(t *testing.T) {
	w := New(Options{Size: time.Minute})
	base := time.Date(2024, 1, 15, 10, 30, 0, 0, time.UTC)
	w.Add(logLine(base, "first window"))
	w.Add(logLine(base.Add(time.Minute), "second window"))
	w.Add(logLine(base.Add(2*time.Minute), "third window"))

	buckets := w.Buckets()
	if len(buckets) != 3 {
		t.Fatalf("expected 3 buckets, got %d", len(buckets))
	}
}

func TestReset_ClearsBuckets(t *testing.T) {
	w := New(Options{Size: time.Minute})
	now := time.Date(2024, 1, 15, 10, 0, 0, 0, time.UTC)
	w.Add(logLine(now, "msg"))
	w.Reset()

	if len(w.Buckets()) != 0 {
		t.Fatal("expected empty buckets after reset")
	}
	// Ensure we can add again after reset.
	w.Add(logLine(now, "after reset"))
	if len(w.Buckets()) != 1 {
		t.Fatal("expected 1 bucket after adding post-reset")
	}
}

func TestBuckets_ReturnsCopy(t *testing.T) {
	w := New(Options{Size: time.Minute})
	now := time.Date(2024, 1, 15, 10, 0, 0, 0, time.UTC)
	w.Add(logLine(now, "msg"))

	b1 := w.Buckets()
	b1[0] = nil // mutate the returned slice
	b2 := w.Buckets()
	if b2[0] == nil {
		t.Fatal("Buckets() should return an independent copy")
	}
}
