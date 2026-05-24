package config

import (
	"testing"
	"time"
)

func TestDefaultConfig(t *testing.T) {
	cfg := DefaultConfig()
	if cfg.Format != "plain" {
		t.Errorf("expected default format 'plain', got %q", cfg.Format)
	}
	if cfg.MaxLines != 0 {
		t.Errorf("expected default MaxLines 0, got %d", cfg.MaxLines)
	}
}

func TestValidate_MissingFile(t *testing.T) {
	cfg := DefaultConfig()
	if err := cfg.Validate(); err == nil {
		t.Fatal("expected error for missing file path")
	}
}

func TestValidate_BadFormat(t *testing.T) {
	cfg := DefaultConfig()
	cfg.FilePath = "some.log"
	cfg.Format = "xml"
	if err := cfg.Validate(); err == nil {
		t.Fatal("expected error for invalid format")
	}
}

func TestValidate_ToBeforeFrom(t *testing.T) {
	cfg := DefaultConfig()
	cfg.FilePath = "some.log"
	cfg.From = time.Now()
	cfg.To = cfg.From.Add(-time.Hour)
	if err := cfg.Validate(); err == nil {
		t.Fatal("expected error when to < from")
	}
}

func TestValidate_NegativeMaxLines(t *testing.T) {
	cfg := DefaultConfig()
	cfg.FilePath = "some.log"
	cfg.MaxLines = -1
	if err := cfg.Validate(); err == nil {
		t.Fatal("expected error for negative max-lines")
	}
}

func TestValidate_Valid(t *testing.T) {
	cfg := DefaultConfig()
	cfg.FilePath = "app.log"
	if err := cfg.Validate(); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
}
