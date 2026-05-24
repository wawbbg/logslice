// Package highlight provides ANSI terminal color highlighting
// for matched patterns within log lines.
package highlight

import (
	"strings"
)

const (
	// ANSI escape codes
	Reset  = "\033[0m"
	Red    = "\033[31m"
	Green  = "\033[32m"
	Yellow = "\033[33m"
	Cyan   = "\033[36m"
	Bold   = "\033[1m"
)

// Highlighter wraps matched substrings with ANSI color codes.
type Highlighter struct {
	color   string
	enabled bool
}

// New creates a new Highlighter. If enabled is false, Highlight returns
// the input unchanged (useful when output is not a terminal).
func New(color string, enabled bool) *Highlighter {
	if color == "" {
		color = Yellow
	}
	return &Highlighter{color: color, enabled: enabled}
}

// Highlight wraps all case-insensitive occurrences of pattern in line
// with the configured ANSI color code. Returns line unchanged if
// highlighting is disabled or pattern is empty.
func (h *Highlighter) Highlight(line, pattern string) string {
	if !h.enabled || pattern == "" {
		return line
	}

	lower := strings.ToLower(line)
	lowerPat := strings.ToLower(pattern)

	var b strings.Builder
	b.Grow(len(line) + 20)

	offset := 0
	for {
		idx := strings.Index(lower[offset:], lowerPat)
		if idx < 0 {
			b.WriteString(line[offset:])
			break
		}
		abs := offset + idx
		b.WriteString(line[offset:abs])
		b.WriteString(h.color)
		b.WriteString(Bold)
		b.WriteString(line[abs : abs+len(pattern)])
		b.WriteString(Reset)
		offset = abs + len(pattern)
	}

	return b.String()
}

// StripANSI removes ANSI escape sequences from s.
func StripANSI(s string) string {
	var b strings.Builder
	b.Grow(len(s))
	i := 0
	for i < len(s) {
		if s[i] == '\033' && i+1 < len(s) && s[i+1] == '[' {
			i += 2
			for i < len(s) && s[i] != 'm' {
				i++
			}
			i++ // skip 'm'
			continue
		}
		b.WriteByte(s[i])
		i++
	}
	return b.String()
}
