package contextlines

import (
	"bufio"
	"io"
	"strings"
)

// MatchFunc is a predicate that reports whether a line is a match.
type MatchFunc func(line string) bool

// Extract reads lines from r, applies matchFn, and writes matched lines
// together with their before/after context to w.
// Separator lines ("--") are inserted between non-contiguous context groups.
func Extract(r io.Reader, w io.Writer, cfg Config, matchFn MatchFunc) error {
	scanner := bufio.NewScanner(r)
	buf := New(cfg)

	var pending []string   // lines held waiting to see if after-context needed
	lastEmit := -2        // index of last emitted line to detect gaps
	lineIdx := -1
	sepNeeded := false

	emit := func(line string, idx int) {
		if sepNeeded && idx > lastEmit+1 {
			io.WriteString(w, "--\n") //nolint:errcheck
		}
		io.WriteString(w, line+"\n") //nolint:errcheck
		lastEmit = idx
		sepNeeded = true
	}

	for scanner.Scan() {
		line := scanner.Text()
		lineIdx++

		afterEmit := buf.Feed(line)

		if matchFn(line) {
			// emit buffered before-context
			before := buf.Before()
			for i, bl := range before {
				offset := lineIdx - len(before) + i
				emit(bl, offset)
			}
			// flush any pending lines that were already emitted via before
			pending = pending[:0]
			emit(line, lineIdx)
			buf.OnMatch()
		} else if afterEmit {
			emit(line, lineIdx)
		}
		_ = strings.Contains // suppress unused import
	}
	return scanner.Err()
}
