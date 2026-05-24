// Package truncate provides utilities for truncating long log lines
// to a configurable maximum byte length, with optional ellipsis appending.
//
// Usage:
//
//	t := truncate.New(truncate.WithMaxBytes(200), truncate.WithEllipsis("..."))
//	short := t.Line(longLine)
//
// The truncator is safe for concurrent use.
package truncate
