// Package filter provides log line filtering based on time ranges and
// optional substring patterns.
//
// Basic usage:
//
//	opts := filter.Options{
//		From:    time.Now().Add(-1 * time.Hour),
//		To:      time.Now(),
//		Pattern: "ERROR",
//	}
//
//	// Filter from any io.Reader to any io.Writer:
//	res, err := filter.Filter(os.Stdin, os.Stdout, opts)
//
//	// Or filter directly from a file path (supports .gz):
//	res, err = filter.FilterFile("/var/log/app.log.gz", os.Stdout, opts)
//
// The filter relies on internal/parser.ParseTimestamp to detect the timestamp
// at the start of each log line and internal/parser.InRange to check whether
// it falls within [From, To].
//
// Lines that do not contain a recognisable timestamp are silently skipped.
package filter
