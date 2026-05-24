package config

import (
	"testing"
)

func TestParseFlags_NoArgs(t *testing.T) {
	_, err := ParseFlags([]string{})
	if err == nil {
		t.Fatal("expected error when no args given")
	}
}

func TestParseFlags_FileOnly(t *testing.T) {
	cfg, err := ParseFlags([]string{"app.log"})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if cfg.FilePath != "app.log" {
		t.Errorf("expected FilePath 'app.log', got %q", cfg.FilePath)
	}
	if cfg.Format != "plain" {
		t.Errorf("expected default format 'plain', got %q", cfg.Format)
	}
}

func TestParseFlags_AllFlags(t *testing.T) {
	args := []string{
		"--from", "2024-01-02T10:00:00Z",
		"--to", "2024-01-02T11:00:00Z",
		"--pattern", "ERROR",
		"--format", "json",
		"--summary",
		"--max-lines", "50",
		"app.log",
	}
	cfg, err := ParseFlags(args)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if cfg.Pattern != "ERROR" {
		t.Errorf("expected pattern 'ERROR', got %q", cfg.Pattern)
	}
	if cfg.Format != "json" {
		t.Errorf("expected format 'json', got %q", cfg.Format)
	}
	if !cfg.ShowSummary {
		t.Error("expected ShowSummary to be true")
	}
	if cfg.MaxLines != 50 {
		t.Errorf("expected MaxLines 50, got %d", cfg.MaxLines)
	}
	if cfg.From.IsZero() {
		t.Error("expected From to be set")
	}
	if cfg.To.IsZero() {
		t.Error("expected To to be set")
	}
}

func TestParseFlags_InvalidFrom(t *testing.T) {
	_, err := ParseFlags([]string{"--from", "not-a-date", "app.log"})
	if err == nil {
		t.Fatal("expected error for invalid --from")
	}
}

func TestParseFlags_InvalidFormat(t *testing.T) {
	_, err := ParseFlags([]string{"--format", "csv", "app.log"})
	if err == nil {
		t.Fatal("expected error for invalid format")
	}
}
