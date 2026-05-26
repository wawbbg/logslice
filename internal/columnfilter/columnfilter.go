// Package columnfilter selects or excludes specific fields from structured
// log lines, supporting both JSON and key=value formats.
package columnfilter

import (
	"encoding/json"
	"strings"
)

// Mode controls whether the listed columns are kept or dropped.
type Mode int

const (
	// Include keeps only the named columns.
	Include Mode = iota
	// Exclude drops the named columns and keeps the rest.
	Exclude
)

// Filter selects or drops fields from structured log lines.
type Filter struct {
	columns map[string]struct{}
	mode    Mode
}

// New creates a Filter for the given column names and mode.
func New(columns []string, mode Mode) *Filter {
	cm := make(map[string]struct{}, len(columns))
	for _, c := range columns {
		cm[strings.TrimSpace(c)] = struct{}{}
	}
	return &Filter{columns: cm, mode: mode}
}

// Line applies the column filter to a single log line.
// JSON objects and key=value pairs are handled; plain lines are returned as-is.
func (f *Filter) Line(line string) string {
	if len(f.columns) == 0 {
		return line
	}
	trimmed := strings.TrimSpace(line)
	if strings.HasPrefix(trimmed, "{") {
		return f.filterJSON(trimmed)
	}
	if strings.Contains(trimmed, "=") {
		return f.filterKV(trimmed)
	}
	return line
}

func (f *Filter) filterJSON(line string) string {
	var obj map[string]json.RawMessage
	if err := json.Unmarshal([]byte(line), &obj); err != nil {
		return line
	}
	out := make(map[string]json.RawMessage, len(obj))
	for k, v := range obj {
		_, listed := f.columns[k]
		if (f.mode == Include && listed) || (f.mode == Exclude && !listed) {
			out[k] = v
		}
	}
	b, err := json.Marshal(out)
	if err != nil {
		return line
	}
	return string(b)
}

func (f *Filter) filterKV(line string) string {
	parts := strings.Fields(line)
	var kept []string
	for _, p := range parts {
		idx := strings.IndexByte(p, '=')
		if idx < 0 {
			// not a key=value token — keep unconditionally
			kept = append(kept, p)
			continue
		}
		key := p[:idx]
		_, listed := f.columns[key]
		if (f.mode == Include && listed) || (f.mode == Exclude && !listed) {
			kept = append(kept, p)
		}
	}
	return strings.Join(kept, " ")
}
