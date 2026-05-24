package rotate

import (
	"bufio"
	"io"

	"github.com/yourorg/logslice/internal/filter"
)

// MultiReader returns a single io.ReadCloser that concatenates all rotation
// entries in order (oldest first). Gzip files are transparently decompressed
// via filter.OpenFile.
type MultiReader struct {
	entries []Entry
	current io.ReadCloser
	idx     int
}

// NewMultiReader creates a MultiReader for the given entries.
func NewMultiReader(entries []Entry) *MultiReader {
	return &MultiReader{entries: entries}
}

// Read implements io.Reader across all rotation files.
func (m *MultiReader) Read(p []byte) (int, error) {
	for {
		if m.current == nil {
			if m.idx >= len(m.entries) {
				return 0, io.EOF
			}
			rc, err := filter.OpenFile(m.entries[m.idx].Path)
			if err != nil {
				return 0, err
			}
			m.current = rc
			m.idx++
		}
		n, err := m.current.Read(p)
		if err == io.EOF {
			_ = m.current.Close()
			m.current = nil
			if n > 0 {
				return n, nil
			}
			continue
		}
		return n, err
	}
}

// Close closes the currently open file, if any.
func (m *MultiReader) Close() error {
	if m.current != nil {
		err := m.current.Close()
		m.current = nil
		return err
	}
	return nil
}

// Lines returns a line scanner over all rotation entries.
func Lines(entries []Entry) *bufio.Scanner {
	return bufio.NewScanner(NewMultiReader(entries))
}
