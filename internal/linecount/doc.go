// Package linecount counts lines in log files with optional pattern filtering.
//
// It supports plain text and gzip-compressed files. A Result value is returned
// containing the total number of lines, the number that matched the given
// substring pattern, and the number that were skipped.
//
// When the pattern is an empty string, all lines are considered matching.
// Skipped lines are those that could not be read due to encoding issues or
// other per-line errors encountered during scanning.
//
// Example usage:
//
//	res, err := linecount.CountFile("app.log", "error")
//	if err != nil {
//		log.Fatal(err)
//	}
//	fmt.Printf("matched %d of %d lines\n", res.Matched, res.Total)
//
// To count all lines without filtering, pass an empty pattern:
//
//	res, err := linecount.CountFile("app.log", "")
package linecount
