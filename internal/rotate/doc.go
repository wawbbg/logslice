// Package rotate handles discovery and sequential reading of rotated log files.
//
// Many log rotation tools produce companion files alongside the active log:
//
//	app.log        — current (live) file
//	app.log.1      — most recently rotated
//	app.log.2      — older
//	app.log.3.gz   — compressed older file
//
// Discover scans the directory of a base log file and returns all rotation
// entries sorted oldest-first so that callers can process them in chronological
// order without manual sorting.
//
// MultiReader wraps the ordered entries into a single io.Reader, transparently
// decompressing gzip files via the filter package. This allows the rest of
// logslice (index building, filtering, tail) to operate on the full rotated
// history as if it were one contiguous stream.
package rotate
