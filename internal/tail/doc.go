// Package tail provides efficient tail-reading of log files.
//
// Unlike reading an entire file and discarding most of it, tail seeks
// backwards from the end of the file in fixed-size chunks, counting
// newlines until the desired number of lines is located. This makes
// it suitable for large log files where only the most recent entries
// are needed.
//
// Example usage:
//
//	lines, err := tail.Lines("/var/log/app.log", 50)
//	if err != nil {
//		log.Fatal(err)
//	}
//	for _, l := range lines {
//		fmt.Println(l)
//	}
package tail
