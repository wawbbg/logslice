// Package checkpoint provides resumable log processing by persisting the last
// successfully consumed byte offset for each log file.
//
// A Store writes checkpoint state to a JSON file on disk. On the next run the
// caller can retrieve the saved offset and seek the log file to that position,
// avoiding re-processing lines that were already handled.
//
// Typical usage:
//
//	s := checkpoint.New("/var/lib/logslice/checkpoints.json")
//
//	rc, startOffset, err := checkpoint.ReaderFrom(s, "/var/log/app.log")
//	if err != nil { ... }
//
//	err = checkpoint.Lines(s, "/var/log/app.log", rc, startOffset, 500,
//	    func(line string) { /* process line */ })
//
// The flushEvery parameter controls how often the offset is written to disk;
// passing 0 disables mid-stream flushing and only saves at EOF.
package checkpoint
