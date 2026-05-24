// Package stats provides a lightweight statistics collector for logslice
// processing runs.
//
// A Collector is created at the start of a filtering operation and updated
// as each log line is processed. At the end of a run the caller can use
// WriteSummary to emit a human-readable report, or inspect the exported
// fields directly for programmatic use (e.g. JSON output).
//
// Example usage:
//
//	col := stats.New()
//	for _, line := range lines {
//		matched := filter.Matches(line)
//		col.RecordLine(matched, len(line))
//	}
//	col.WriteSummary(os.Stderr)
package stats
