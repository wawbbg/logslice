// Package dedupe implements log-line deduplication for logslice.
//
// It provides a sliding-window Deduplicator that hashes each line with
// FNV-64a and suppresses re-occurrences within a configurable window.
// The window size trades memory for recall: a larger window catches
// duplicates that are further apart in the stream.
//
// Basic usage:
//
//	d := dedupe.New(dedupe.WithWindowSize(512))
//	for _, line := range logLines {
//		if !d.IsDuplicate(line) {
//			fmt.Println(line)
//		}
//	}
//
// For streaming pipelines use FilterReader or the Lines convenience helper:
//
//	lines, err := dedupe.Lines(reader, dedupe.New())
package dedupe
