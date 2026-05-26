// Package multipattern implements multi-pattern log line matching with
// configurable AND / OR combination logic and optional negation.
//
// Patterns are plain case-insensitive substrings by default. Wrapping a
// pattern in forward slashes (e.g. "/ERR[0-9]+/") treats it as a regular
// expression. Prefixing any pattern with '!' negates the match.
//
// Example usage:
//
//	mt, err := multipattern.New(
//		[]string{"error", "!timeout"},
//		multipattern.WithMode(multipattern.ModeAND),
//	)
//	if err != nil {
//		log.Fatal(err)
//	}
//	if mt.Match(line) {
//		// line contains "error" but not "timeout"
//	}
package multipattern
