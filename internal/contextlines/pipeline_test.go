package contextlines

import (
	"strings"
	"testing"
)

func runExtract(input string, cfg Config, pattern string) string {
	var sb strings.Builder
	r := strings.NewReader(input)
	Extract(r, &sb, cfg, func(line string) bool {
		return strings.Contains(line, pattern)
	})
	return sb.String()
}

func TestExtract_MatchOnly(t *testing.T) {
	input := "alpha\nbeta\ngamma\n"
	out := runExtract(input, Config{}, "beta")
	if !strings.Contains(out, "beta") {
		t.Errorf("expected match in output, got: %q", out)
	}
	if strings.Contains(out, "alpha") || strings.Contains(out, "gamma") {
		t.Errorf("unexpected context lines in output: %q", out)
	}
}

func TestExtract_BeforeContext(t *testing.T) {
	input := "line1\nline2\nMATCH\nline4\n"
	out := runExtract(input, Config{Before: 2, After: 0}, "MATCH")
	if !strings.Contains(out, "line1") {
		t.Errorf("expected line1 in before-context, got: %q", out)
	}
	if !strings.Contains(out, "line2") {
		t.Errorf("expected line2 in before-context, got: %q", out)
	}
	if strings.Contains(out, "line4") {
		t.Errorf("unexpected after line in output: %q", out)
	}
}

func TestExtract_AfterContext(t *testing.T) {
	input := "line1\nMATCH\nline3\nline4\nline5\n"
	out := runExtract(input, Config{Before: 0, After: 2}, "MATCH")
	if !strings.Contains(out, "line3") {
		t.Errorf("expected line3 in after-context, got: %q", out)
	}
	if !strings.Contains(out, "line4") {
		t.Errorf("expected line4 in after-context, got: %q", out)
	}
	if strings.Contains(out, "line5") {
		t.Errorf("unexpected line5 in output: %q", out)
	}
}

func TestExtract_SeparatorBetweenGroups(t *testing.T) {
	input := "a\nMATCH1\nb\nc\nd\ne\nMATCH2\nf\n"
	out := runExtract(input, Config{Before: 0, After: 0}, "MATCH")
	if !strings.Contains(out, "--") {
		t.Errorf("expected separator between groups, got: %q", out)
	}
}

func TestExtract_NoMatch(t *testing.T) {
	input := "foo\nbar\nbaz\n"
	out := runExtract(input, Config{Before: 1, After: 1}, "NOMATCH")
	if out != "" {
		t.Errorf("expected empty output for no match, got: %q", out)
	}
}
