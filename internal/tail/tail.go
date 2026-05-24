// Package tail provides functionality for reading the last N lines
// of a log file efficiently without loading the entire file into memory.
package tail

import (
	"bufio"
	"fmt"
	"io"
	"os"
)

const defaultChunkSize = 4096

// Lines reads the last n lines from the file at the given path.
// It seeks from the end of the file in chunks to avoid reading the
// entire file when only a small tail is needed.
func Lines(path string, n int) ([]string, error) {
	if n <= 0 {
		return nil, fmt.Errorf("tail: n must be positive, got %d", n)
	}

	f, err := os.Open(path)
	if err != nil {
		return nil, fmt.Errorf("tail: open %s: %w", path, err)
	}
	defer f.Close()

	size, err := fileSize(f)
	if err != nil {
		return nil, err
	}

	offset, err := findOffset(f, size, n)
	if err != nil {
		return nil, err
	}

	if _, err := f.Seek(offset, io.SeekStart); err != nil {
		return nil, fmt.Errorf("tail: seek: %w", err)
	}

	return readLines(f, n)
}

func fileSize(f *os.File) (int64, error) {
	info, err := f.Stat()
	if err != nil {
		return 0, fmt.Errorf("tail: stat: %w", err)
	}
	return info.Size(), nil
}

// findOffset scans backwards through the file to find the byte offset
// at which the (size - n)th newline occurs.
func findOffset(f *os.File, size int64, n int) (int64, error) {
	if size == 0 {
		return 0, nil
	}

	newlines := 0
	chunk := make([]byte, defaultChunkSize)
	pos := size

	for pos > 0 {
		read := int64(defaultChunkSize)
		if pos < read {
			read = pos
		}
		pos -= read

		if _, err := f.Seek(pos, io.SeekStart); err != nil {
			return 0, fmt.Errorf("tail: seek: %w", err)
		}

		buf := chunk[:read]
		if _, err := io.ReadFull(f, buf); err != nil {
			return 0, fmt.Errorf("tail: read chunk: %w", err)
		}

		for i := int(read) - 1; i >= 0; i-- {
			if buf[i] == '\n' {
				newlines++
				if newlines > n {
					return pos + int64(i) + 1, nil
				}
			}
		}
	}

	return 0, nil
}

func readLines(r io.Reader, max int) ([]string, error) {
	var lines []string
	scanner := bufio.NewScanner(r)
	for scanner.Scan() {
		lines = append(lines, scanner.Text())
		if len(lines) >= max {
			break
		}
	}
	if err := scanner.Err(); err != nil {
		return nil, fmt.Errorf("tail: scan: %w", err)
	}
	return lines, nil
}
