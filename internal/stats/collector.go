// Package stats provides log processing statistics collection and reporting.
package stats

import (
	"fmt"
	"io"
	"time"
)

// Collector accumulates statistics during a log filtering run.
type Collector struct {
	StartTime     time.Time
	LinesRead     int
	LinesMatched  int
	LinesSkipped  int
	BytesRead     int64
	ParseErrors   int
}

// New returns an initialized Collector with StartTime set to now.
func New() *Collector {
	return &Collector{
		StartTime: time.Now(),
	}
}

// RecordLine records a processed line, marking whether it matched filters.
func (c *Collector) RecordLine(matched bool, byteLen int) {
	c.LinesRead++
	c.BytesRead += int64(byteLen)
	if matched {
		c.LinesMatched++
	} else {
		c.LinesSkipped++
	}
}

// RecordParseError increments the parse error counter.
func (c *Collector) RecordParseError() {
	c.ParseErrors++
}

// Elapsed returns the duration since the collector was created.
func (c *Collector) Elapsed() time.Duration {
	return time.Since(c.StartTime)
}

// MatchRate returns the fraction of lines that matched, or 0 if no lines were read.
func (c *Collector) MatchRate() float64 {
	if c.LinesRead == 0 {
		return 0
	}
	return float64(c.LinesMatched) / float64(c.LinesRead)
}

// WriteSummary writes a human-readable summary to w.
func (c *Collector) WriteSummary(w io.Writer) {
	fmt.Fprintf(w, "Lines read:    %d\n", c.LinesRead)
	fmt.Fprintf(w, "Lines matched: %d (%.1f%%)\n", c.LinesMatched, c.MatchRate()*100)
	fmt.Fprintf(w, "Lines skipped: %d\n", c.LinesSkipped)
	fmt.Fprintf(w, "Bytes read:    %d\n", c.BytesRead)
	if c.ParseErrors > 0 {
		fmt.Fprintf(w, "Parse errors:  %d\n", c.ParseErrors)
	}
	fmt.Fprintf(w, "Elapsed:       %s\n", c.Elapsed().Round(time.Millisecond))
}
