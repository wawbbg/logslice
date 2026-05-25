package levelfilter_test

import (
	"testing"

	"github.com/yourorg/logslice/internal/levelfilter"
)

func TestParseLevel_KnownLevels(t *testing.T) {
	cases := []struct {
		input string
		want  levelfilter.Level
	}{
		{"debug", levelfilter.LevelDebug},
		{"DEBUG", levelfilter.LevelDebug},
		{"info", levelfilter.LevelInfo},
		{"INFO", levelfilter.LevelInfo},
		{"warn", levelfilter.LevelWarn},
		{"warning", levelfilter.LevelWarn},
		{"WARNING", levelfilter.LevelWarn},
		{"error", levelfilter.LevelError},
		{"err", levelfilter.LevelError},
		{"fatal", levelfilter.LevelFatal},
		{"crit", levelfilter.LevelFatal},
		{"critical", levelfilter.LevelFatal},
	}
	for _, tc := range cases {
		t.Run(tc.input, func(t *testing.T) {
			got := levelfilter.ParseLevel(tc.input)
			if got != tc.want {
				t.Errorf("ParseLevel(%q) = %d, want %d", tc.input, got, tc.want)
			}
		})
	}
}

func TestParseLevel_Unknown(t *testing.T) {
	if got := levelfilter.ParseLevel("trace"); got != levelfilter.LevelUnknown {
		t.Errorf("expected LevelUnknown, got %d", got)
	}
}

func TestFilter_AllowsAtOrAboveMin(t *testing.T) {
	f := levelfilter.New(levelfilter.LevelWarn)

	cases := []struct {
		line  string
		allow bool
	}{
		{"2024-01-01 DEBUG starting up", false},
		{"2024-01-01 INFO server ready", false},
		{"2024-01-01 WARN disk usage high", true},
		{"2024-01-01 ERROR connection refused", true},
		{"2024-01-01 FATAL out of memory", true},
	}
	for _, tc := range cases {
		t.Run(tc.line, func(t *testing.T) {
			if got := f.Allow(tc.line); got != tc.allow {
				t.Errorf("Allow(%q) = %v, want %v", tc.line, got, tc.allow)
			}
		})
	}
}

func TestFilter_UnknownLevelPassesThrough(t *testing.T) {
	f := levelfilter.New(levelfilter.LevelError)
	line := "2024-01-01 some line with no level keyword"
	if !f.Allow(line) {
		t.Errorf("expected unknown-level line to pass through, but Allow returned false")
	}
}

func TestFilter_DebugMinAllowsAll(t *testing.T) {
	f := levelfilter.New(levelfilter.LevelDebug)
	lines := []string{
		"DEBUG verbose output",
		"INFO startup complete",
		"WARN low memory",
		"ERROR disk full",
		"FATAL core dump",
	}
	for _, l := range lines {
		if !f.Allow(l) {
			t.Errorf("expected Allow(%q) = true with LevelDebug min", l)
		}
	}
}
