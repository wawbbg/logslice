package checkpoint

import (
	"bufio"
	"io"
	"os"
)

// ReaderFrom opens the given file at the offset recorded in the Store.
// If no checkpoint exists the file is read from the beginning.
// It returns the reader, the starting offset, and any error.
func ReaderFrom(s *Store, file string) (io.ReadCloser, int64, error) {
	entry, err := s.Load(file)
	if err != nil {
		return nil, 0, err
	}

	f, err := os.Open(file)
	if err != nil {
		return nil, 0, err
	}

	if entry.Offset > 0 {
		if _, err := f.Seek(entry.Offset, io.SeekStart); err != nil {
			f.Close()
			return nil, 0, err
		}
	}
	return f, entry.Offset, nil
}

// Lines reads lines from rc, calling fn for each one.
// It tracks the cumulative byte count (relative to startOffset) and saves a
// checkpoint after every flushEvery lines so progress survives interruptions.
func Lines(s *Store, file string, rc io.ReadCloser, startOffset int64, flushEvery int, fn func(string)) error {
	defer rc.Close()

	scanner := bufio.NewScanner(rc)
	offset := startOffset
	count := 0

	for scanner.Scan() {
		line := scanner.Text()
		offset += int64(len(line)) + 1 // +1 for newline
		count++
		fn(line)

		if flushEvery > 0 && count%flushEvery == 0 {
			if err := s.Save(file, offset); err != nil {
				return err
			}
		}
	}
	if err := scanner.Err(); err != nil {
		return err
	}
	return s.Save(file, offset)
}
