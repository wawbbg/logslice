package linenum

import (
	"bufio"
	"io"
	"strings"
)

// Result holds a single annotated output line together with its source
// line number and whether it was a matched line.
type Result struct {
	LineNo  int
	Line    string
	Matched bool
}

// Annotate reads all lines from r, advances the tracker for each line, and
// returns a Result slice. Lines that contain substr (case-sensitive) are
// flagged as matched and annotated; others are included unannotated when
// includeAll is true.
func Annotate(r io.Reader, tr *Tracker, substr string, includeAll bool) []Result {
	var results []Result
	scanner := bufio.NewScanner(r)
	for scanner.Scan() {
		line := scanner.Text()
		lineNo := tr.Advance()
		matched := substr == "" || strings.Contains(line, substr)
		if matched {
			tr.RecordMatch()
			results = append(results, Result{
				LineNo:  lineNo,
				Line:    tr.Annotate(line),
				Matched: true,
			})
		} else if includeAll {
			results = append(results, Result{
				LineNo:  lineNo,
				Line:    line,
				Matched: false,
			})
		}
	}
	return results
}
