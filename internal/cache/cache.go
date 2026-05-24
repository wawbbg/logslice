// Package cache provides a simple LRU-style file offset cache for logslice.
// It stores byte offsets of timestamp-indexed lines to speed up repeated
// queries against the same log file.
package cache

import (
	"sync"
	"time"
)

// Entry represents a cached line offset within a log file.
type Entry struct {
	Offset    int64
	Timestamp time.Time
	Line      string
}

// Cache holds a bounded set of file offset entries keyed by file path.
type Cache struct {
	mu      sync.RWMutex
	entries map[string][]Entry
	maxSize int
}

// New creates a new Cache with the given maximum number of entries per file.
func New(maxSize int) *Cache {
	if maxSize <= 0 {
		maxSize = 512
	}
	return &Cache{
		entries: make(map[string][]Entry),
		maxSize: maxSize,
	}
}

// Put stores an Entry for the given file path. If the entry count exceeds
// maxSize the oldest half of entries are evicted.
func (c *Cache) Put(filePath string, e Entry) {
	c.mu.Lock()
	defer c.mu.Unlock()

	bucket := c.entries[filePath]
	bucket = append(bucket, e)
	if len(bucket) > c.maxSize {
		bucket = bucket[c.maxSize/2:]
	}
	c.entries[filePath] = bucket
}

// Get returns all cached entries for the given file path.
func (c *Cache) Get(filePath string) ([]Entry, bool) {
	c.mu.RLock()
	defer c.mu.RUnlock()
	entries, ok := c.entries[filePath]
	return entries, ok
}

// Evict removes all cached entries for the given file path.
func (c *Cache) Evict(filePath string) {
	c.mu.Lock()
	defer c.mu.Unlock()
	delete(c.entries, filePath)
}

// Size returns the total number of cached entries across all files.
func (c *Cache) Size() int {
	c.mu.RLock()
	defer c.mu.RUnlock()
	total := 0
	for _, v := range c.entries {
		total += len(v)
	}
	return total
}
