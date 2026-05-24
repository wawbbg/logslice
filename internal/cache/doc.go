// Package cache implements a lightweight, thread-safe in-memory cache for
// logslice. It stores parsed line offsets and timestamps keyed by file path,
// allowing the filter engine to skip re-scanning portions of a log file that
// have already been indexed during the current process lifetime.
//
// The cache uses a simple eviction strategy: when the number of stored entries
// for a single file exceeds the configured maximum, the oldest half of entries
// are dropped. This keeps memory usage bounded while retaining the most
// recently accessed offsets.
//
// Usage:
//
//	c := cache.New(512)
//	c.Put("app.log", cache.Entry{Offset: 1024, Timestamp: ts, Line: raw})
//	entries, ok := c.Get("app.log")
package cache
