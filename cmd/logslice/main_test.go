package main

import (
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"testing"
)

// buildBinary compiles the logslice binary into a temp dir and returns its path.
func buildBinary(t *testing.T) string {
	t.Helper()
	tmpDir := t.TempDir()
	binPath := filepath.Join(tmpDir, "logslice")
	cmd := exec.Command("go", "build", "-o", binPath, ".")
	cmd.Dir = "."
	out, err := cmd.CombinedOutput()
	if err != nil {
		t.Fatalf("failed to build binary: %v\n%s", err, out)
	}
	return binPath
}

func writeTempLog(t *testing.T, content string) string {
	t.Helper()
	f, err := os.CreateTemp(t.TempDir(), "test-*.log")
	if err != nil {
		t.Fatalf("create temp file: %v", err)
	}
	if _, err := f.WriteString(content); err != nil {
		t.Fatalf("write temp file: %v", err)
	}
	f.Close()
	return f.Name()
}

const sampleLog = `2024-03-01T10:00:00 INFO  service started
2024-03-01T10:01:00 ERROR disk full
2024-03-01T10:02:00 INFO  request handled
2024-03-01T10:03:00 WARN  high memory
`

func TestCLI_NoArgs(t *testing.T) {
	bin := buildBinary(t)
	cmd := exec.Command(bin)
	out, err := cmd.CombinedOutput()
	if err == nil {
		t.Fatal("expected non-zero exit with no args")
	}
	if !strings.Contains(string(out), "Usage") {
		t.Errorf("expected usage message, got: %s", out)
	}
}

func TestCLI_FilterByPattern(t *testing.T) {
	bin := buildBinary(t)
	logFile := writeTempLog(t, sampleLog)
	cmd := exec.Command(bin, "--pattern", "ERROR", logFile)
	out, err := cmd.Output()
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	result := string(out)
	if !strings.Contains(result, "disk full") {
		t.Errorf("expected ERROR line in output, got: %s", result)
	}
	if strings.Contains(result, "service started") {
		t.Errorf("unexpected non-ERROR line in output: %s", result)
	}
}

func TestCLI_SummaryFlag(t *testing.T) {
	bin := buildBinary(t)
	logFile := writeTempLog(t, sampleLog)
	cmd := exec.Command(bin, "--summary", logFile)
	out, err := cmd.Output()
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !strings.Contains(string(out), "lines matched") {
		t.Errorf("expected summary line, got: %s", out)
	}
}

func TestCLI_InvalidStartTime(t *testing.T) {
	bin := buildBinary(t)
	logFile := writeTempLog(t, sampleLog)
	cmd := exec.Command(bin, "--start", "not-a-time", logFile)
	_, err := cmd.Output()
	if err == nil {
		t.Fatal("expected error for invalid --start")
	}
}
