// Package timewindow implements fixed-duration time-window aggregation for
// log streams.
//
// # Overview
//
// A Windower groups log lines into Buckets whose width is determined by the
// Options.Size field (default: 1 minute). Each line is parsed for a leading
// timestamp; lines that carry no recognisable timestamp are silently ignored.
//
// # Usage
//
//	w := timewindow.New(timewindow.Options{Size: 5 * time.Minute})
//	for _, line := range lines {
//		w.Add(line)
//	}
//	for _, b := range w.Buckets() {
//		fmt.Printf("%v – %v : %d lines\n", b.Start, b.End, b.Count)
//	}
//
// Call Reset to reuse the Windower across multiple files or passes.
package timewindow
