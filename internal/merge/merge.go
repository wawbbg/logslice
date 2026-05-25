// Package merge provides utilities for merging multiple sorted log streams
// into a single chronologically ordered output.
package merge

import (
	"bufio"
	"container/heap"
	"io"
	"time"

	"github.com/user/logslice/internal/parser"
)

// entry holds a single log line along with its parsed timestamp and the
// index of the source reader it came from.
type entry struct {
	line      string
	ts        time.Time
	sourceIdx int
}

// minHeap implements heap.Interface for entry values ordered by timestamp.
type minHeap []entry

func (h minHeap) Len() int            { return len(h) }
func (h minHeap) Less(i, j int) bool  { return h[i].ts.Before(h[j].ts) }
func (h minHeap) Swap(i, j int)       { h[i], h[j] = h[j], h[i] }
func (h *minHeap) Push(x interface{}) { *h = append(*h, x.(entry)) }
func (h *minHeap) Pop() interface{} {
	old := *h
	n := len(old)
	x := old[n-1]
	*h = old[:n-1]
	return x
}

// Merge reads from each reader in rs, parses the timestamp of every line, and
// writes lines to w in ascending timestamp order. Lines whose timestamps cannot
// be parsed are emitted immediately in source order without affecting the heap.
func Merge(w io.Writer, rs []io.Reader) error {
	scanners := make([]*bufio.Scanner, len(rs))
	for i, r := range rs {
		scanners[i] = bufio.NewScanner(r)
	}

	h := &minHeap{}
	heap.Init(h)

	// Seed the heap with the first line from each scanner.
	for i, sc := range scanners {
		if sc.Scan() {
			line := sc.Text()
			ts, ok := parser.ParseTimestamp(line)
			if !ok {
				if _, err := io.WriteString(w, line+"\n"); err != nil {
					return err
				}
				continue
			}
			heap.Push(h, entry{line: line, ts: ts, sourceIdx: i})
		}
	}

	for h.Len() > 0 {
		e := heap.Pop(h).(entry)
		if _, err := io.WriteString(w, e.line+"\n"); err != nil {
			return err
		}
		sc := scanners[e.sourceIdx]
		for sc.Scan() {
			line := sc.Text()
			ts, ok := parser.ParseTimestamp(line)
			if !ok {
				if _, err := io.WriteString(w, line+"\n"); err != nil {
					return err
				}
				continue
			}
			heap.Push(h, entry{line: line, ts: ts, sourceIdx: e.sourceIdx})
			break
		}
	}
	return nil
}
