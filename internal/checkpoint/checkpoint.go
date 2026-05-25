// Package checkpoint persists and restores the last successfully processed
// byte offset for a log file, enabling resumable log processing across runs.
package checkpoint

import (
	"encoding/json"
	"errors"
	"os"
	"path/filepath"
	"time"
)

// Entry holds the persisted state for a single log file.
type Entry struct {
	File      string    `json:"file"`
	Offset    int64     `json:"offset"`
	UpdatedAt time.Time `json:"updated_at"`
}

// Store reads and writes checkpoint entries to a JSON file on disk.
type Store struct {
	path string
	now  func() time.Time
}

// New returns a Store that persists checkpoints to the given file path.
func New(path string) *Store {
	return &Store{path: path, now: time.Now}
}

// Load returns the Entry for the given log file, or a zero-value Entry if none
// exists yet. A missing checkpoint file is not treated as an error.
func (s *Store) Load(file string) (Entry, error) {
	all, err := s.loadAll()
	if err != nil {
		return Entry{}, err
	}
	return all[filepath.Clean(file)], nil
}

// Save persists the byte offset reached for the given log file.
func (s *Store) Save(file string, offset int64) error {
	all, err := s.loadAll()
	if err != nil {
		return err
	}
	all[filepath.Clean(file)] = Entry{
		File:      file,
		Offset:    offset,
		UpdatedAt: s.now(),
	}
	return s.persist(all)
}

// Delete removes the checkpoint entry for the given log file.
func (s *Store) Delete(file string) error {
	all, err := s.loadAll()
	if err != nil {
		return err
	}
	delete(all, filepath.Clean(file))
	return s.persist(all)
}

func (s *Store) loadAll() (map[string]Entry, error) {
	data, err := os.ReadFile(s.path)
	if errors.Is(err, os.ErrNotExist) {
		return make(map[string]Entry), nil
	}
	if err != nil {
		return nil, err
	}
	var m map[string]Entry
	if err := json.Unmarshal(data, &m); err != nil {
		return nil, err
	}
	return m, nil
}

func (s *Store) persist(m map[string]Entry) error {
	data, err := json.MarshalIndent(m, "", "  ")
	if err != nil {
		return err
	}
	return os.WriteFile(s.path, data, 0o644)
}
