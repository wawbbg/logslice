package truncate

import (
	"strings"
	"testing"
)

func TestNew_Defaults(t *testing.T) {
	tr := New()
	if tr.maxBytes != DefaultMaxBytes {
		t.Errorf("expected maxBytes %d, got %d", DefaultMaxBytes, tr.maxBytes)
	}
	if !tr.appendEllipsis {
		t.Error("expected appendEllipsis to be true by default")
	}
}

func TestLine_ShortLineUnchanged(t *testing.T) {
	tr := New(WithMaxBytes(20))
	input := "short line"
	if got := tr.Line(input); got != input {
		t.Errorf("expected %q, got %q", input, got)
	}
}

func TestLine_ExactLengthUnchanged(t *testing.T) {
	tr := New(WithMaxBytes(10), WithEllipsis(false))
	input := "0123456789" // exactly 10 bytes
	if got := tr.Line(input); got != input {
		t.Errorf("expected %q, got %q", input, got)
	}
}

func TestLine_TruncatesWithEllipsis(t *testing.T) {
	tr := New(WithMaxBytes(10))
	input := "0123456789ABCDEF"
	got := tr.Line(input)
	if !strings.HasSuffix(got, ellipsis) {
		t.Errorf("expected ellipsis suffix, got %q", got)
	}
	if len(got) > 10 {
		t.Errorf("expected len <= 10, got %d", len(got))
	}
}

func TestLine_TruncatesWithoutEllipsis(t *testing.T) {
	tr := New(WithMaxBytes(8), WithEllipsis(false))
	input := "ABCDEFGHIJKLMNOP"
	got := tr.Line(input)
	if len(got) != 8 {
		t.Errorf("expected len 8, got %d", len(got))
	}
	if strings.HasSuffix(got, ellipsis) {
		t.Errorf("unexpected ellipsis in %q", got)
	}
}

func TestLine_RespectsUTF8Boundary(t *testing.T) {
	// "é" is 2 bytes (0xC3 0xA9); maxBytes=5 would cut in the middle without boundary check.
	tr := New(WithMaxBytes(5), WithEllipsis(false))
	input := "aébc" // a(1) + é(2) + b(1) + c(1) = 5 bytes total — fits exactly
	got := tr.Line(input)
	if got != input {
		t.Errorf("expected %q unchanged, got %q", input, got)
	}

	// Now force a cut that would land mid-rune.
	input2 := "aébcd" // 6 bytes
	tr2 := New(WithMaxBytes(3), WithEllipsis(false))
	got2 := tr2.Line(input2)
	// Valid UTF-8 must be preserved.
	if !isValidUTF8(got2) {
		t.Errorf("truncated string is not valid UTF-8: %q", got2)
	}
}

func TestLines_AppliestoAll(t *testing.T) {
	tr := New(WithMaxBytes(5), WithEllipsis(false))
	input := []string{"short", "toolongline", "ok"}
	result := tr.Lines(input)
	if len(result[0]) > 5 || len(result[1]) > 5 || len(result[2]) > 5 {
		t.Errorf("one or more lines exceed maxBytes: %v", result)
	}
}

func TestWithMaxBytes_ZeroIgnored(t *testing.T) {
	tr := New(WithMaxBytes(0))
	if tr.maxBytes != DefaultMaxBytes {
		t.Errorf("zero maxBytes should be ignored, got %d", tr.maxBytes)
	}
}

func isValidUTF8(s string) bool {
	for i := 0; i < len(s); {
		r, size := rune(s[i]), 1
		if s[i] >= 0x80 {
			var ok bool
			r, size, ok = decodeRune(s[i:])
			if !ok || r == '\uFFFD' {
				return false
			}
		}
		_ = r
		i += size
	}
	return true
}

func decodeRune(b string) (rune, int, bool) {
	import_rune := []byte(b)
	r, size := rune(0), 0
	switch {
	case import_rune[0]&0xE0 == 0xC0 && len(import_rune) >= 2:
		r = rune(import_rune[0]&0x1F)<<6 | rune(import_rune[1]&0x3F)
		size = 2
	case import_rune[0]&0xF0 == 0xE0 && len(import_rune) >= 3:
		r = rune(import_rune[0]&0x0F)<<12 | rune(import_rune[1]&0x3F)<<6 | rune(import_rune[2]&0x3F)
		size = 3
	default:
		return 0, 1, false
	}
	return r, size, true
}
