package stats_test

import (
	"bytes"
	"strings"
	"testing"
	"time"

	"github.com/yourorg/logslice/internal/stats"
)

func TestNew_InitializesStartTime(t *testing.T) {
	before := time.Now()
	c := stats.New()
	after := time.Now()

	if c.StartTime.Before(before) || c.StartTime.After(after) {
		t.Errorf("StartTime %v not in expected range [%v, %v]", c.StartTime, before, after)
	}
}

func TestRecordLine_Matched(t *testing.T) {
	c := stats.New()
	c.RecordLine(true, 42)
	c.RecordLine(true, 10)

	if c.LinesRead != 2 {
		t.Errorf("LinesRead = %d, want 2", c.LinesRead)
	}
	if c.LinesMatched != 2 {
		t.Errorf("LinesMatched = %d, want 2", c.LinesMatched)
	}
	if c.LinesSkipped != 0 {
		t.Errorf("LinesSkipped = %d, want 0", c.LinesSkipped)
	}
	if c.BytesRead != 52 {
		t.Errorf("BytesRead = %d, want 52", c.BytesRead)
	}
}

func TestRecordLine_NotMatched(t *testing.T) {
	c := stats.New()
	c.RecordLine(false, 20)

	if c.LinesSkipped != 1 {
		t.Errorf("LinesSkipped = %d, want 1", c.LinesSkipped)
	}
	if c.LinesMatched != 0 {
		t.Errorf("LinesMatched = %d, want 0", c.LinesMatched)
	}
}

func TestMatchRate(t *testing.T) {
	c := stats.New()
	if c.MatchRate() != 0 {
		t.Errorf("MatchRate on empty collector = %f, want 0", c.MatchRate())
	}

	c.RecordLine(true, 1)
	c.RecordLine(false, 1)

	if got := c.MatchRate(); got != 0.5 {
		t.Errorf("MatchRate = %f, want 0.5", got)
	}
}

func TestRecordParseError(t *testing.T) {
	c := stats.New()
	c.RecordParseError()
	c.RecordParseError()

	if c.ParseErrors != 2 {
		t.Errorf("ParseErrors = %d, want 2", c.ParseErrors)
	}
}

func TestWriteSummary_ContainsKeyFields(t *testing.T) {
	c := stats.New()
	c.RecordLine(true, 100)
	c.RecordLine(false, 50)
	c.RecordParseError()

	var buf bytes.Buffer
	c.WriteSummary(&buf)
	out := buf.String()

	for _, want := range []string{"Lines read", "Lines matched", "Lines skipped", "Bytes read", "Parse errors", "Elapsed"} {
		if !strings.Contains(out, want) {
			t.Errorf("WriteSummary output missing %q\nGot:\n%s", want, out)
		}
	}
}
