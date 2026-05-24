package index

import (
	"fmt"
	"strings"
	"testing"
	"time"
)

// generateLog produces n log lines with one-second increments starting from base.
func generateLog(n int) string {
	base := time.Date(2024, 1, 1, 0, 0, 0, 0, time.UTC)
	var b strings.Builder
	for i := 0; i < n; i++ {
		ts := base.Add(time.Duration(i) * time.Second)
		fmt.Fprintf(&b, "%s INFO  log line number %d\n", ts.Format(time.RFC3339), i)
	}
	return b.String()
}

func BenchmarkBuild_10k(b *testing.B) {
	data := generateLog(10_000)
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		r := strings.NewReader(data)
		_, err := Build(r)
		if err != nil {
			b.Fatal(err)
		}
	}
}

func BenchmarkFindStart_10k(b *testing.B) {
	data := generateLog(10_000)
	r := strings.NewReader(data)
	idx, err := Build(r)
	if err != nil {
		b.Fatal(err)
	}

	// Seek to the midpoint of the generated range.
	target := time.Date(2024, 1, 1, 1, 23, 20, 0, time.UTC) // ~5000 s in
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = idx.FindStart(target)
	}
}
