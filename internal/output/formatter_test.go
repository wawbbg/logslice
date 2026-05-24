package output

import (
	"bytes"
	"strings"
	"testing"
)

func TestFormatter_PlainFormat(t *testing.T) {
	var buf bytes.Buffer
	f := NewFormatter(&buf, FormatPlain)

	lines := []string{"2024-01-01 info hello", "2024-01-02 warn world"}
	if err := f.WriteLines(lines); err != nil {
		t.Fatalf("WriteLines error: %v", err)
	}

	got := buf.String()
	for _, l := range lines {
		if !strings.Contains(got, l) {
			t.Errorf("expected output to contain %q", l)
		}
	}
	if f.Count() != 2 {
		t.Errorf("expected count 2, got %d", f.Count())
	}
}

func TestFormatter_NumberedFormat(t *testing.T) {
	var buf bytes.Buffer
	f := NewFormatter(&buf, FormatNumbered)

	if err := f.WriteLine("first line"); err != nil {
		t.Fatalf("WriteLine error: %v", err)
	}
	if err := f.WriteLine("second line"); err != nil {
		t.Fatalf("WriteLine error: %v", err)
	}

	got := buf.String()
	if !strings.Contains(got, "1") || !strings.Contains(got, "first line") {
		t.Errorf("expected numbered output, got: %q", got)
	}
	if !strings.Contains(got, "2") || !strings.Contains(got, "second line") {
		t.Errorf("expected numbered output, got: %q", got)
	}
}

func TestFormatter_JSONFormat(t *testing.T) {
	var buf bytes.Buffer
	f := NewFormatter(&buf, FormatJSON)

	if err := f.WriteLine(`level=info msg="app started"`); err != nil {
		t.Fatalf("WriteLine error: %v", err)
	}

	got := buf.String()
	if !strings.Contains(got, `"n":1`) {
		t.Errorf("expected JSON with n:1, got: %q", got)
	}
	if !strings.Contains(got, `"line":`) {
		t.Errorf("expected JSON with line key, got: %q", got)
	}
}

func TestFormatter_NilWriterDefaultsToStdout(t *testing.T) {
	f := NewFormatter(nil, FormatPlain)
	if f.w == nil {
		t.Error("expected non-nil writer when nil passed")
	}
}

func TestFormatter_WriteSummary(t *testing.T) {
	var out bytes.Buffer
	var summary bytes.Buffer
	f := NewFormatter(&out, FormatPlain)
	_ = f.WriteLines([]string{"a", "b", "c"})
	f.WriteSummary(&summary)

	if !strings.Contains(summary.String(), "3") {
		t.Errorf("expected summary to mention 3, got: %q", summary.String())
	}
}
