package checkpoint

import (
	"os"
	"path/filepath"
	"testing"
	"time"
)

func tempStore(t *testing.T) *Store {
	t.Helper()
	dir := t.TempDir()
	s := New(filepath.Join(dir, "checkpoints.json"))
	fixed := time.Date(2024, 1, 15, 12, 0, 0, 0, time.UTC)
	s.now = func() time.Time { return fixed }
	return s
}

func TestLoad_MissingFile_ReturnsZero(t *testing.T) {
	s := tempStore(t)
	e, err := s.Load("/var/log/app.log")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if e.Offset != 0 {
		t.Errorf("expected zero offset, got %d", e.Offset)
	}
}

func TestSaveAndLoad_RoundTrip(t *testing.T) {
	s := tempStore(t)
	const file = "/var/log/app.log"
	const want int64 = 4096

	if err := s.Save(file, want); err != nil {
		t.Fatalf("Save: %v", err)
	}
	e, err := s.Load(file)
	if err != nil {
		t.Fatalf("Load: %v", err)
	}
	if e.Offset != want {
		t.Errorf("offset: got %d, want %d", e.Offset, want)
	}
	if e.File != file {
		t.Errorf("file: got %q, want %q", e.File, file)
	}
}

func TestSave_UpdatesExistingEntry(t *testing.T) {
	s := tempStore(t)
	const file = "/var/log/app.log"

	_ = s.Save(file, 100)
	_ = s.Save(file, 9999)

	e, err := s.Load(file)
	if err != nil {
		t.Fatalf("Load: %v", err)
	}
	if e.Offset != 9999 {
		t.Errorf("expected 9999, got %d", e.Offset)
	}
}

func TestDelete_RemovesEntry(t *testing.T) {
	s := tempStore(t)
	const file = "/var/log/app.log"

	_ = s.Save(file, 512)
	if err := s.Delete(file); err != nil {
		t.Fatalf("Delete: %v", err)
	}
	e, err := s.Load(file)
	if err != nil {
		t.Fatalf("Load after delete: %v", err)
	}
	if e.Offset != 0 {
		t.Errorf("expected zero offset after delete, got %d", e.Offset)
	}
}

func TestLoad_CorruptFile_ReturnsError(t *testing.T) {
	dir := t.TempDir()
	p := filepath.Join(dir, "checkpoints.json")
	_ = os.WriteFile(p, []byte("not json{"), 0o644)
	s := New(p)
	_, err := s.Load("/any")
	if err == nil {
		t.Fatal("expected error for corrupt checkpoint file")
	}
}

func TestSave_MultipleFiles_IndependentEntries(t *testing.T) {
	s := tempStore(t)
	_ = s.Save("/var/log/a.log", 111)
	_ = s.Save("/var/log/b.log", 222)

	a, _ := s.Load("/var/log/a.log")
	b, _ := s.Load("/var/log/b.log")

	if a.Offset != 111 {
		t.Errorf("a: got %d, want 111", a.Offset)
	}
	if b.Offset != 222 {
		t.Errorf("b: got %d, want 222", b.Offset)
	}
}
