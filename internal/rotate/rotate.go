// Package rotate provides support for detecting and reading rotated log files.
// It discovers companion files (e.g. app.log.1, app.log.2.gz) and presents
// them in chronological order alongside the primary file.
package rotate

import (
	"fmt"
	"os"
	"path/filepath"
	"sort"
	"strconv"
	"strings"
)

// Entry describes a single file in a rotation sequence.
type Entry struct {
	Path  string
	Index int // 0 = current, 1 = most recent rotated, etc.
}

// Discover returns all rotation entries for the given base log file in
// ascending chronological order (oldest first, current last).
func Discover(base string) ([]Entry, error) {
	if _, err := os.Stat(base); err != nil {
		return nil, fmt.Errorf("rotate: base file %q not found: %w", base, err)
	}

	dir := filepath.Dir(base)
	name := filepath.Base(base)

	entries, err := findRotated(dir, name)
	if err != nil {
		return nil, err
	}

	// Append the current (live) file last.
	entries = append(entries, Entry{Path: base, Index: 0})
	return entries, nil
}

// findRotated scans dir for files matching <name>.<N> or <name>.<N>.gz.
func findRotated(dir, name string) ([]Entry, error) {
	glob := filepath.Join(dir, name+".*")
	matches, err := filepath.Glob(glob)
	if err != nil {
		return nil, fmt.Errorf("rotate: glob error: %w", err)
	}

	var entries []Entry
	for _, m := range matches {
		base := filepath.Base(m)
		suffix := strings.TrimPrefix(base, name+".")
		suffix = strings.TrimSuffix(suffix, ".gz")
		idx, err := strconv.Atoi(suffix)
		if err != nil || idx < 1 {
			continue
		}
		entries = append(entries, Entry{Path: m, Index: idx})
	}

	// Sort descending by index so oldest rotated file comes first.
	sort.Slice(entries, func(i, j int) bool {
		return entries[i].Index > entries[j].Index
	})
	return entries, nil
}
