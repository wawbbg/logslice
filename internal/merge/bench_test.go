package merge

import (
	"bytes"
	"fmt"
	"io"
	"strings"
	"testing"
	"time"
)

// generateStream returns a reader containing n log lines whose timestamps
// start at base and advance by step per line.
func generateStream(n int, base time.Time, step time.Duration) io.Reader {
	var sb strings.Builder
	for i := 0; i < n; i++ {
		ts := base.Add(time.Duration(i) * step)
		fmt.Fprintf(&sb, "%s INFO  bench line %d\n", ts.UTC().Format(time.RFC3339), i)
	}
	return strings.NewReader(sb.String())
}

func BenchmarkMerge_TwoStreams_10k(b *testing.B) {
	base := time.Date(2024, 1, 1, 0, 0, 0, 0, time.UTC)
	for i := 0; i < b.N; i++ {
		b.StopTimer()
		r1 := generateStream(5000, base, 2*time.Second)
		r2 := generateStream(5000, base.Add(time.Second), 2*time.Second)
		b.StartTimer()
		if err := Merge(io.Discard, []interface{ Read([]byte) (int, error) }{r1, r2}); err != nil {
			b.Fatal(err)
		}
	}
}

func BenchmarkMerge_FourStreams_10k(b *testing.B) {
	base := time.Date(2024, 1, 1, 0, 0, 0, 0, time.UTC)
	for i := 0; i < b.N; i++ {
		b.StopTimer()
		readers := make([]interface{ Read([]byte) (int, error) }, 4)
		for j := 0; j < 4; j++ {
			readers[j] = generateStream(2500, base.Add(time.Duration(j)*time.Second), 4*time.Second)
		}
		var buf bytes.Buffer
		b.StartTimer()
		if err := Merge(&buf, readers); err != nil {
			b.Fatal(err)
		}
	}
}
