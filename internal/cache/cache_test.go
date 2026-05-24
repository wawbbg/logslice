package cache

import (
	"fmt"
	"testing"
	"time"
)

func makeEntry(offset int64, line string) Entry {
	return Entry{
		Offset:    offset,
		Timestamp: time.Now(),
		Line:      line,
	}
}

func TestNew_DefaultMaxSize(t *testing.T) {
	c := New(0)
	if c.maxSize != 512 {
		t.Errorf("expected default maxSize 512, got %d", c.maxSize)
	}
}

func TestPutAndGet(t *testing.T) {
	c := New(10)
	e := makeEntry(42, "hello world")
	c.Put("app.log", e)

	entries, ok := c.Get("app.log")
	if !ok {
		t.Fatal("expected entries to exist")
	}
	if len(entries) != 1 {
		t.Fatalf("expected 1 entry, got %d", len(entries))
	}
	if entries[0].Offset != 42 {
		t.Errorf("expected offset 42, got %d", entries[0].Offset)
	}
}

func TestGet_MissingKey(t *testing.T) {
	c := New(10)
	_, ok := c.Get("nonexistent.log")
	if ok {
		t.Error("expected ok=false for missing key")
	}
}

func TestEvict(t *testing.T) {
	c := New(10)
	c.Put("app.log", makeEntry(1, "line1"))
	c.Evict("app.log")
	_, ok := c.Get("app.log")
	if ok {
		t.Error("expected entries to be evicted")
	}
}

func TestSize(t *testing.T) {
	c := New(100)
	for i := 0; i < 5; i++ {
		c.Put("a.log", makeEntry(int64(i), fmt.Sprintf("line%d", i)))
	}
	for i := 0; i < 3; i++ {
		c.Put("b.log", makeEntry(int64(i), fmt.Sprintf("line%d", i)))
	}
	if got := c.Size(); got != 8 {
		t.Errorf("expected size 8, got %d", got)
	}
}

func TestPut_Eviction(t *testing.T) {
	maxSize := 10
	c := New(maxSize)
	for i := 0; i < maxSize+1; i++ {
		c.Put("app.log", makeEntry(int64(i), fmt.Sprintf("line%d", i)))
	}
	entries, ok := c.Get("app.log")
	if !ok {
		t.Fatal("expected entries after eviction")
	}
	if len(entries) > maxSize {
		t.Errorf("expected at most %d entries after eviction, got %d", maxSize, len(entries))
	}
}
