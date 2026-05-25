// Package merge implements a k-way merge of multiple log streams.
//
// When log data is spread across several files — for example rotated logs or
// logs from different hosts — it is often necessary to interleave them in
// strict timestamp order before further analysis.
//
// Usage:
//
//	f1, _ := os.Open("app.log.1")
//	f2, _ := os.Open("app.log.2")
//	defer f1.Close()
//	defer f2.Close()
//
//	if err := merge.Merge(os.Stdout, []io.Reader{f1, f2}); err != nil {
//		log.Fatal(err)
//	}
//
// Lines that do not contain a recognisable timestamp are written to the output
// immediately, in the order they are encountered, without disrupting the heap.
package merge
