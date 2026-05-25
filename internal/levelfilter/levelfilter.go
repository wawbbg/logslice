// Package levelfilter provides log-level based filtering for structured log lines.
// It supports standard log levels (DEBUG, INFO, WARN, ERROR, FATAL) and allows
// filtering lines at or above a specified minimum severity level.
package levelfilter

import (
	"strings"
)

// Level represents a log severity level.
type Level int

const (
	LevelDebug Level = iota
	LevelInfo
	LevelWarn
	LevelError
	LevelFatal
	LevelUnknown Level = -1
)

var levelNames = map[string]Level{
	"debug": LevelDebug,
	"info":  LevelInfo,
	"warn":  LevelWarn,
	"warning": LevelWarn,
	"error": LevelError,
	"err":   LevelError,
	"fatal": LevelFatal,
	"crit":  LevelFatal,
	"critical": LevelFatal,
}

// ParseLevel parses a level string into a Level value.
// Returns LevelUnknown if the string is not recognized.
func ParseLevel(s string) Level {
	if l, ok := levelNames[strings.ToLower(strings.TrimSpace(s))]; ok {
		return l
	}
	return LevelUnknown
}

// Filter holds the minimum level threshold for log line filtering.
type Filter struct {
	min Level
}

// New creates a new Filter that passes lines at or above the given minimum level.
func New(min Level) *Filter {
	return &Filter{min: min}
}

// Allow returns true if the line contains a log level at or above the minimum.
// Lines whose level cannot be detected are always allowed through.
func (f *Filter) Allow(line string) bool {
	level := detectLevel(line)
	if level == LevelUnknown {
		return true
	}
	return level >= f.min
}

// detectLevel scans a line for a known level keyword and returns it.
func detectLevel(line string) Level {
	upper := strings.ToUpper(line)
	// Check longer tokens first to avoid partial matches (e.g. WARNING before WARN).
	for _, candidate := range []string{"CRITICAL", "WARNING", "FATAL", "ERROR", "DEBUG", "INFO", "WARN", "ERR", "CRIT"} {
		if strings.Contains(upper, candidate) {
			return levelNames[strings.ToLower(candidate)]
		}
	}
	return LevelUnknown
}
