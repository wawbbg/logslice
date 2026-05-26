// Package multipattern provides multi-pattern matching with AND/OR/NOT logic
// for filtering log lines against multiple substring or regex patterns.
package multipattern

import (
	"regexp"
	"strings"
)

// Mode controls how multiple patterns are combined.
type Mode int

const (
	// ModeAND requires all patterns to match.
	ModeAND Mode = iota
	// ModeOR requires at least one pattern to match.
	ModeOR
)

// Pattern holds a compiled expression and whether it is negated.
type Pattern struct {
	re      *regexp.Regexp
	raw     string
	negated bool
}

// Matcher evaluates a set of patterns against log lines.
type Matcher struct {
	patterns []Pattern
	mode     Mode
}

// Option configures a Matcher.
type Option func(*Matcher)

// WithMode sets the combination mode (AND / OR).
func WithMode(m Mode) Option {
	return func(mt *Matcher) { mt.mode = m }
}

// New creates a Matcher from the supplied pattern strings.
// A pattern prefixed with '!' is treated as a NOT (negated) pattern.
// Patterns are matched case-insensitively as plain substrings unless they
// are wrapped in '/' delimiters, in which case they are compiled as regex.
func New(patterns []string, opts ...Option) (*Matcher, error) {
	mt := &Matcher{mode: ModeAND}
	for _, o := range opts {
		o(mt)
	}
	for _, p := range patterns {
		var pat Pattern
		if strings.HasPrefix(p, "!") {
			pat.negated = true
			p = p[1:]
		}
		if strings.HasPrefix(p, "/") && strings.HasSuffix(p, "/") && len(p) >= 2 {
			inner := p[1 : len(p)-1]
			re, err := regexp.Compile(inner)
			if err != nil {
				return nil, err
			}
			pat.re = re
		} else {
			pat.raw = strings.ToLower(p)
		}
		mt.patterns = append(mt.patterns, pat)
	}
	return mt, nil
}

// Match reports whether line satisfies the pattern set.
func (m *Matcher) Match(line string) bool {
	if len(m.patterns) == 0 {
		return true
	}
	lower := strings.ToLower(line)
	for _, p := range m.patterns {
		hit := m.eval(p, line, lower)
		switch m.mode {
		case ModeOR:
			if hit {
				return true
			}
		case ModeAND:
			if !hit {
				return false
			}
		}
	}
	return m.mode == ModeAND
}

func (m *Matcher) eval(p Pattern, line, lower string) bool {
	var matched bool
	if p.re != nil {
		matched = p.re.MatchString(line)
	} else {
		matched = strings.Contains(lower, p.raw)
	}
	if p.negated {
		return !matched
	}
	return matched
}
