// Package levelfilter implements severity-level based filtering for log lines.
//
// It recognises the standard log levels DEBUG, INFO, WARN/WARNING, ERROR/ERR,
// and FATAL/CRIT/CRITICAL and allows callers to discard lines below a chosen
// minimum severity.
//
// Lines that contain no recognisable level keyword are always passed through so
// that non-standard or plain-text log entries are never silently dropped.
//
// Usage:
//
//	f := levelfilter.New(levelfilter.LevelWarn)
//	for _, line := range lines {
//	    if f.Allow(line) {
//	        fmt.Println(line)
//	    }
//	}
package levelfilter
