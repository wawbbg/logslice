package columnfilter

import (
	"testing"
)

func TestNew_EmptyColumns_PassThrough(t *testing.T) {
	f := New(nil, Include)
	line := `{"level":"info","msg":"hello"}`
	if got := f.Line(line); got != line {
		t.Fatalf("expected passthrough, got %q", got)
	}
}

func TestLine_JSON_Include(t *testing.T) {
	f := New([]string{"level", "msg"}, Include)
	line := `{"level":"info","msg":"hello","ts":"2024-01-01"}`
	got := f.Line(line)
	if strings.Contains(got, "ts") {
		t.Fatalf("excluded field 'ts' present in output: %q", got)
	}
	if !strings.Contains(got, "level") || !strings.Contains(got, "msg") {
		t.Fatalf("included fields missing from output: %q", got)
	}
}

func TestLine_JSON_Exclude(t *testing.T) {
	f := New([]string{"ts"}, Exclude)
	line := `{"level":"info","msg":"hello","ts":"2024-01-01"}`
	got := f.Line(line)
	if strings.Contains(got, "ts") {
		t.Fatalf("excluded field 'ts' still present: %q", got)
	}
	if !strings.Contains(got, "level") {
		t.Fatalf("non-excluded field 'level' missing: %q", got)
	}
}

func TestLine_KV_Include(t *testing.T) {
	f := New([]string{"level", "msg"}, Include)
	line := "level=info msg=hello ts=2024-01-01"
	got := f.Line(line)
	if strings.Contains(got, "ts=") {
		t.Fatalf("excluded field present: %q", got)
	}
	if !strings.Contains(got, "level=info") || !strings.Contains(got, "msg=hello") {
		t.Fatalf("included fields missing: %q", got)
	}
}

func TestLine_KV_Exclude(t *testing.T) {
	f := New([]string{"ts"}, Exclude)
	line := "level=info msg=hello ts=2024-01-01"
	got := f.Line(line)
	if strings.Contains(got, "ts=") {
		t.Fatalf("excluded field present: %q", got)
	}
	if !strings.Contains(got, "level=info") {
		t.Fatalf("non-excluded field missing: %q", got)
	}
}

func TestLine_PlainText_Passthrough(t *testing.T) {
	f := New([]string{"level"}, Include)
	line := "this is a plain log line with no structure"
	if got := f.Line(line); got != line {
		t.Fatalf("plain line should pass through unchanged, got %q", got)
	}
}

func TestLine_InvalidJSON_Passthrough(t *testing.T) {
	f := New([]string{"level"}, Include)
	line := "{not valid json"
	if got := f.Line(line); got != line {
		t.Fatalf("invalid JSON should pass through, got %q", got)
	}
}

func TestLine_KV_NonKVTokensKept(t *testing.T) {
	f := New([]string{"msg"}, Include)
	line := "2024-01-01 msg=hello level=info"
	got := f.Line(line)
	if !strings.Contains(got, "2024-01-01") {
		t.Fatalf("non-kv prefix should be kept: %q", got)
	}
	if !strings.Contains(got, "msg=hello") {
		t.Fatalf("included kv field missing: %q", got)
	}
}

func init() {
	// imported via blank identifier to satisfy test file dependency
	_ = strings.Contains
}

import "strings"
