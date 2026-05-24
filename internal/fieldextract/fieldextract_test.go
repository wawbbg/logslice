package fieldextract

import (
	"testing"
)

func TestExtract_JSON(t *testing.T) {
	e := New(FormatJSON)
	fields := e.Extract(`{"level":"info","msg":"started","port":8080}`)
	if fields == nil {
		t.Fatal("expected non-nil fields")
	}
	if fields["level"] != "info" {
		t.Errorf("level: got %q, want %q", fields["level"], "info")
	}
	if fields["msg"] != "started" {
		t.Errorf("msg: got %q, want %q", fields["msg"], "started")
	}
	// numeric field should be marshalled back to string
	if fields["port"] != "8080" {
		t.Errorf("port: got %q, want %q", fields["port"], "8080")
	}
}

func TestExtract_JSON_Invalid(t *testing.T) {
	e := New(FormatJSON)
	if got := e.Extract("not json"); got != nil {
		t.Errorf("expected nil for invalid JSON, got %v", got)
	}
}

func TestExtract_KeyValue(t *testing.T) {
	e := New(FormatKeyValue)
	fields := e.Extract(`time=2024-01-02T15:04:05Z level=warn msg="disk full"`)
	if fields == nil {
		t.Fatal("expected non-nil fields")
	}
	if fields["level"] != "warn" {
		t.Errorf("level: got %q, want %q", fields["level"], "warn")
	}
	if fields["msg"] != "disk full" {
		t.Errorf("msg: got %q, want %q", fields["msg"], "disk full")
	}
}

func TestExtract_KeyValue_NoEquals(t *testing.T) {
	e := New(FormatKeyValue)
	if got := e.Extract("plain log line with no pairs"); got != nil {
		t.Errorf("expected nil for no key=value pairs, got %v", got)
	}
}

func TestExtract_Auto_DetectsJSON(t *testing.T) {
	e := New(FormatAuto)
	fields := e.Extract(`{"env":"prod"}`)
	if fields == nil || fields["env"] != "prod" {
		t.Errorf("auto JSON detection failed: %v", fields)
	}
}

func TestExtract_Auto_DetectsKeyValue(t *testing.T) {
	e := New(FormatAuto)
	fields := e.Extract("service=api status=200")
	if fields == nil || fields["service"] != "api" || fields["status"] != "200" {
		t.Errorf("auto kv detection failed: %v", fields)
	}
}

func TestExtract_EmptyLine(t *testing.T) {
	e := New(FormatAuto)
	if got := e.Extract(""); got != nil {
		t.Errorf("expected nil for empty line, got %v", got)
	}
}

func TestField_ReturnsValue(t *testing.T) {
	e := New(FormatKeyValue)
	val := e.Field("level=error msg=oops", "level")
	if val != "error" {
		t.Errorf("Field: got %q, want %q", val, "error")
	}
}

func TestField_MissingKey(t *testing.T) {
	e := New(FormatKeyValue)
	val := e.Field("level=error", "missing")
	if val != "" {
		t.Errorf("Field missing key: got %q, want empty", val)
	}
}

func TestField_UnparsableLine(t *testing.T) {
	e := New(FormatKeyValue)
	val := e.Field("no pairs here", "level")
	if val != "" {
		t.Errorf("Field unparsable: got %q, want empty", val)
	}
}
