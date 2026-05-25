// Package linecount provides utilities for counting lines in log files,
// including support for plain text and gzip-compressed files.
package linecount

import (
	"bufio"
	"compress/gzip"
	"fmt"
	"io"
	"os"
	"strings"
)

// Result holds the outcome of a line count operation.
type Result struct {
	Total   int64
	Matched int64
	Skipped int64
}

// CountFile counts lines in the given file, optionally filtering by a
// substring pattern. Pass an empty pattern to count all lines.
func CountFile(path, pattern string) (Result, error) {
	f, err := os.Open(path)
	if err != nil {
		return Result{}, fmt.Errorf("linecount: open %q: %w", path, err)
	}
	defer f.Close()

	r, err := newReader(path, f)
	if err != nil {
		return Result{}, fmt.Errorf("linecount: create reader: %w", err)
	}

	return count(r, pattern)
}

// CountReader counts lines from an arbitrary io.Reader.
func CountReader(r io.Reader, pattern string) (Result, error) {
	return count(r, pattern)
}

// newReader wraps the file in a gzip reader when the path ends in ".gz".
func newReader(path string, f io.Reader) (io.Reader, error) {
	if strings.HasSuffix(path, ".gz") {
		gr, err := gzip.NewReader(f)
		if err != nil {
			return nil, err
		}
		return gr, nil
	}
	return f, nil
}

func count(r io.Reader, pattern string) (Result, error) {
	var res Result
	scanner := bufio.NewScanner(r)
	for scanner.Scan() {
		line := scanner.Text()
		res.Total++
		if pattern == "" || strings.Contains(line, pattern) {
			res.Matched++
		} else {
			res.Skipped++
		}
	}
	if err := scanner.Err(); err != nil {
		return res, fmt.Errorf("linecount: scan: %w", err)
	}
	return res, nil
}
