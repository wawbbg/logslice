// Package highlight provides terminal-aware ANSI color highlighting
// for matched substrings within log lines.
//
// Usage:
//
//	h := highlight.New(highlight.Yellow, isTerminal)
//	formatted := h.Highlight(line, pattern)
//
// When the output destination is not a terminal (e.g. a file or pipe),
// pass enabled=false to New so that raw log content is preserved without
// embedded escape sequences.
//
// StripANSI can be used to remove any pre-existing ANSI codes from
// input lines before further processing.
package highlight
