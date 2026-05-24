// Package output provides formatting and writing utilities for filtered log results.
package output

import (
	"fmt"
	"io"
	"os"
	"strings"
)

// Format defines the output format for log results.
type Format string

const (
	// FormatPlain outputs log lines as-is.
	FormatPlain Format = "plain"
	// FormatJSON wraps each line in a simple JSON envelope.
	FormatJSON Format = "json"
	// FormatNumbered prefixes each line with its line number.
	FormatNumbered Format = "numbered"
)

// Formatter writes log lines to a destination in the configured format.
type Formatter struct {
	format Format
	w      io.Writer
	count  int
}

// NewFormatter creates a Formatter that writes to w using the given format.
// If w is nil, os.Stdout is used.
func NewFormatter(w io.Writer, format Format) *Formatter {
	if w == nil {
		w = os.Stdout
	}
	return &Formatter{format: format, w: w}
}

// WriteLine writes a single log line according to the configured format.
func (f *Formatter) WriteLine(line string) error {
	f.count++
	var out string
	switch f.format {
	case FormatJSON:
		escaped := strings.ReplaceAll(line, `"`, `\"`)
		out = fmt.Sprintf(`{"n":%d,"line":"%s"}\n`, f.count, escaped)
	case FormatNumbered:
		out = fmt.Sprintf("%6d\t%s\n", f.count, line)
	default:
		out = line + "\n"
	}
	_, err := fmt.Fprint(f.w, out)
	return err
}

// WriteLines writes multiple log lines.
func (f *Formatter) WriteLines(lines []string) error {
	for _, l := range lines {
		if err := f.WriteLine(l); err != nil {
			return err
		}
	}
	return nil
}

// Count returns the number of lines written so far.
func (f *Formatter) Count() int {
	return f.count
}

// WriteSummary writes a trailing summary line to w.
func (f *Formatter) WriteSummary(w io.Writer) {
	if w == nil {
		w = os.Stderr
	}
	fmt.Fprintf(w, "# %d line(s) matched\n", f.count)
}
