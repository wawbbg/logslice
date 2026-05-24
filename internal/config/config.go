// Package config handles parsing and validation of logslice configuration,
// including CLI flags and optional config file support.
package config

import (
	"errors"
	"time"
)

// Config holds all runtime configuration for a logslice run.
type Config struct {
	// Input
	FilePath string

	// Time range filters
	From time.Time
	To   time.Time

	// Pattern filter (substring or regex)
	Pattern string

	// Output options
	Format      string // "plain", "json", "numbered"
	ShowSummary bool
	OutputFile  string

	// Behaviour
	MaxLines int
}

// DefaultConfig returns a Config with sensible defaults.
func DefaultConfig() *Config {
	return &Config{
		Format:   "plain",
		MaxLines: 0, // 0 means unlimited
	}
}

// Validate checks that the Config is internally consistent.
func (c *Config) Validate() error {
	if c.FilePath == "" {
		return errors.New("file path must not be empty")
	}

	validFormats := map[string]bool{"plain": true, "json": true, "numbered": true}
	if !validFormats[c.Format] {
		return errors.New("format must be one of: plain, json, numbered")
	}

	if !c.From.IsZero() && !c.To.IsZero() && c.To.Before(c.From) {
		return errors.New("--to must not be before --from")
	}

	if c.MaxLines < 0 {
		return errors.New("max-lines must be >= 0")
	}

	return nil
}
