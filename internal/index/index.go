// Package index provides byte-offset indexing for large log files,
// enabling fast seeking to time-range boundaries without full scans.
package index

import (
	"bufio"
	"io"
	"time"

	"github.com/yourorg/logslice/internal/parser"
)

// Entry records the byte offset and parsed timestamp of a log line.
type Entry struct {
	Offset    int64
	Timestamp time.Time
}

// Index holds an ordered slice of entries built from a log file.
type Index struct {
	entries []Entry
}

// Build scans r, recording an Entry for every line that contains a
// parseable timestamp. Lines without a recognisable timestamp are skipped.
func Build(r io.ReadSeeker) (*Index, error) {
	if _, err := r.Seek(0, io.SeekStart); err != nil {
		return nil, err
	}

	var entries []Entry
	var offset int64

	scanner := bufio.NewScanner(r)
	for scanner.Scan() {
		line := scanner.Text()
		ts, err := parser.ParseTimestamp(line)
		if err == nil {
			entries = append(entries, Entry{Offset: offset, Timestamp: ts})
		}
		// +1 for the newline byte consumed by the scanner
		offset += int64(len(scanner.Bytes())) + 1
	}
	if err := scanner.Err(); err != nil {
		return nil, err
	}

	return &Index{entries: entries}, nil
}

// Len returns the number of indexed entries.
func (idx *Index) Len() int { return len(idx.entries) }

// FindStart returns the byte offset of the first entry whose timestamp is
// >= from. Returns 0 if no suitable entry exists.
func (idx *Index) FindStart(from time.Time) int64 {
	for _, e := range idx.entries {
		if !e.Timestamp.Before(from) {
			return e.Offset
		}
	}
	return 0
}

// FindEnd returns the byte offset of the first entry whose timestamp is
// strictly after to, so callers can stop reading there. Returns -1 when
// all entries fall within the range.
func (idx *Index) FindEnd(to time.Time) int64 {
	for _, e := range idx.entries {
		if e.Timestamp.After(to) {
			return e.Offset
		}
	}
	return -1
}

// Entries returns a copy of all indexed entries.
func (idx *Index) Entries() []Entry {
	out := make([]Entry, len(idx.entries))
	copy(out, idx.entries)
	return out
}
