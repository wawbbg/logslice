// Package linecount counts lines in log files with optional pattern filtering.
//
// It supports plain text and gzip-compressed files. A Result value is returned
// containing the total number of lines, the number that matched the given
// substring pattern, and the number that were skipped.
//
// Example usage:
//
//	res, err := linecount.CountFile("app.log", "error")
//	if err != nil {
//		log.Fatal(err)
//	}
//	fmt.Printf("matched %d of %d lines\n", res.Matched, res.Total)
package linecount
