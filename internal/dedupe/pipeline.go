package dedupe

import (
	"bufio"
	"io"
)

// FilterReader wraps an io.Reader and skips duplicate lines on the fly.
// It is intended for use in streaming pipelines where allocating the full
// input in memory is undesirable.
type FilterReader struct {
	scanner *bufio.Scanner
	dedup   *Deduplicator
	current string
	done    bool
}

// NewFilterReader returns a FilterReader that deduplicates lines from r
// using the provided Deduplicator.
func NewFilterReader(r io.Reader, d *Deduplicator) *FilterReader {
	return &FilterReader{
		scanner: bufio.NewScanner(r),
		dedup:   d,
	}
}

// Next advances to the next unique line. It returns false when the input
// is exhausted or an error occurs.
func (f *FilterReader) Next() bool {
	for f.scanner.Scan() {
		line := f.scanner.Text()
		if !f.dedup.IsDuplicate(line) {
			f.current = line
			return true
		}
	}
	f.done = true
	return false
}

// Line returns the current unique line. Call Next first.
func (f *FilterReader) Line() string {
	return f.current
}

// Err returns any scanner error that occurred during iteration.
func (f *FilterReader) Err() error {
	return f.scanner.Err()
}

// Lines returns all unique lines from r, deduplicating with d.
// It reads the full input into memory and is a convenience wrapper
// around FilterReader for smaller inputs.
func Lines(r io.Reader, d *Deduplicator) ([]string, error) {
	fr := NewFilterReader(r, d)
	var out []string
	for fr.Next() {
		out = append(out, fr.Line())
	}
	if err := fr.Err(); err != nil {
		return nil, err
	}
	return out, nil
}
