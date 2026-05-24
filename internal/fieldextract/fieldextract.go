// Package fieldextract provides utilities for extracting named fields
// from structured log lines (JSON or key=value formats).
package fieldextract

import (
	"encoding/json"
	"strings"
)

// Format represents the detected or configured log line format.
type Format int

const (
	FormatAuto    Format = iota // detect automatically
	FormatJSON                 // JSON object per line
	FormatKeyValue             // key=value pairs
)

// Extractor pulls named fields from a log line.
type Extractor struct {
	format Format
}

// New returns an Extractor configured with the given format.
func New(format Format) *Extractor {
	return &Extractor{format: format}
}

// Extract returns a map of field names to values parsed from line.
// Returns nil if the line cannot be parsed.
func (e *Extractor) Extract(line string) map[string]string {
	line = strings.TrimSpace(line)
	if line == "" {
		return nil
	}
	switch e.format {
	case FormatJSON:
		return extractJSON(line)
	case FormatKeyValue:
		return extractKeyValue(line)
	default: // FormatAuto
		if strings.HasPrefix(line, "{") {
			return extractJSON(line)
		}
		return extractKeyValue(line)
	}
}

// Field returns a single field value from line, or empty string if not found.
func (e *Extractor) Field(line, key string) string {
	fields := e.Extract(line)
	if fields == nil {
		return ""
	}
	return fields[key]
}

func extractJSON(line string) map[string]string {
	var raw map[string]interface{}
	if err := json.Unmarshal([]byte(line), &raw); err != nil {
		return nil
	}
	out := make(map[string]string, len(raw))
	for k, v := range raw {
		switch val := v.(type) {
		case string:
			out[k] = val
		case nil:
			out[k] = ""
		default:
			b, _ := json.Marshal(val)
			out[k] = string(b)
		}
	}
	return out
}

func extractKeyValue(line string) map[string]string {
	out := make(map[string]string)
	for _, token := range tokenize(line) {
		idx := strings.IndexByte(token, '=')
		if idx <= 0 {
			continue
		}
		key := token[:idx]
		val := strings.Trim(token[idx+1:], `"`)
		out[key] = val
	}
	if len(out) == 0 {
		return nil
	}
	return out
}

// tokenize splits a key=value line respecting quoted values.
func tokenize(line string) []string {
	var tokens []string
	var cur strings.Builder
	inQuote := false
	for _, ch := range line {
		switch {
		case ch == '"':
			inQuote = !inQuote
			cur.WriteRune(ch)
		case ch == ' ' && !inQuote:
			if cur.Len() > 0 {
				tokens = append(tokens, cur.String())
				cur.Reset()
			}
		default:
			cur.WriteRune(ch)
		}
	}
	if cur.Len() > 0 {
		tokens = append(tokens, cur.String())
	}
	return tokens
}
