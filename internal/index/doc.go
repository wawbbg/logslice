// Package index builds a lightweight byte-offset index over a log file so
// that logslice can seek directly to the start of a requested time range
// instead of scanning every line from the beginning.
//
// # Usage
//
//	file, _ := os.Open("app.log")
//	idx, err := index.Build(file)
//	if err != nil {
//		log.Fatal(err)
//	}
//
//	startOffset := idx.FindStart(from)
//	_, _ = file.Seek(startOffset, io.SeekStart)
//
// The index stores one Entry per parseable line, keeping memory usage
// proportional to the number of timestamped lines rather than total bytes.
package index
