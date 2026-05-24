package filter

import (
	"bufio"
	"io"
	"time"

	"github.com/yourorg/logslice/internal/parser"
)

// Options holds configuration for log filtering.
type Options struct {
	From    time.Time
	To      time.Time
	Pattern string // optional substring match
}

// Result holds the outcome of a filter run.
type Result struct {
	MatchedLines int
	TotalLines   int
}

// Filter reads lines from r, writes matching lines to w based on opts.
func Filter(r io.Reader, w io.Writer, opts Options) (Result, error) {
	scanner := bufio.NewScanner(r)
	writer := bufio.NewWriter(w)
	defer writer.Flush()

	var res Result

	for scanner.Scan() {
		line := scanner.Text()
		res.TotalLines++

		ts, err := parser.ParseTimestamp(line)
		if err != nil {
			continue
		}

		if !parser.InRange(ts, opts.From, opts.To) {
			continue
		}

		if opts.Pattern != "" && !containsPattern(line, opts.Pattern) {
			continue
		}

		_, err = writer.WriteString(line + "\n")
		if err != nil {
			return res, err
		}
		res.MatchedLines++
	}

	if err := scanner.Err(); err != nil {
		return res, err
	}

	return res, nil
}

// containsPattern performs a simple case-sensitive substring search.
func containsPattern(line, pattern string) bool {
	return len(line) >= len(pattern) && findSubstring(line, pattern)
}

func findSubstring(s, sub string) bool {
	for i := 0; i <= len(s)-len(sub); i++ {
		if s[i:i+len(sub)] == sub {
			return true
		}
	}
	return false
}
