// Package redact provides line-level redaction of sensitive patterns
// such as passwords, tokens, IP addresses, and email addresses.
package redact

import (
	"regexp"
	"sync"
)

// Rule describes a single redaction rule: a compiled pattern and its
// replacement string.
type Rule struct {
	pattern     *regexp.Regexp
	replacement string
}

// Redactor applies a set of redaction rules to log lines.
type Redactor struct {
	mu    sync.RWMutex
	rules []Rule
}

// Option is a functional option for configuring a Redactor.
type Option func(*Redactor)

// WithRule adds a redaction rule using a raw regex pattern and a
// replacement string. Returns an error if the pattern is invalid.
func WithRule(pattern, replacement string) (Option, error) {
	re, err := regexp.Compile(pattern)
	if err != nil {
		return nil, err
	}
	return func(r *Redactor) {
		r.rules = append(r.rules, Rule{pattern: re, replacement: replacement})
	}, nil
}

// New creates a Redactor with a default set of built-in rules covering
// common sensitive data patterns (bearer tokens, basic-auth, emails).
func New(opts ...Option) *Redactor {
	r := &Redactor{}
	// Built-in defaults
	defaults := []struct{ pattern, replacement string }{
		{`(?i)(bearer\s+)[A-Za-z0-9\-._~+/]+=*`, `${1}[REDACTED]`},
		{`(?i)(password[=:\s]+)\S+`, `${1}[REDACTED]`},
		{`[a-zA-Z0-9._%+\-]+@[a-zA-Z0-9.\-]+\.[a-zA-Z]{2,}`, `[EMAIL]`},
	}
	for _, d := range defaults {
		re := regexp.MustCompile(d.pattern)
		r.rules = append(r.rules, Rule{pattern: re, replacement: d.replacement})
	}
	for _, opt := range opts {
		opt(r)
	}
	return r
}

// Line applies all redaction rules to the given line and returns the
// sanitised result.
func (r *Redactor) Line(line string) string {
	r.mu.RLock()
	defer r.mu.RUnlock()
	for _, rule := range r.rules {
		line = rule.pattern.ReplaceAllString(line, rule.replacement)
	}
	return line
}

// AddRule compiles and appends a new rule at runtime.
func (r *Redactor) AddRule(pattern, replacement string) error {
	re, err := regexp.Compile(pattern)
	if err != nil {
		return err
	}
	r.mu.Lock()
	r.rules = append(r.rules, Rule{pattern: re, replacement: replacement})
	r.mu.Unlock()
	return nil
}
