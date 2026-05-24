package dedupe

import (
	"strings"
	"testing"
)

func TestFilterReader_UniqueLines(t *testing.T) {
	input := "alpha\nbeta\ngamma\n"
	fr := NewFilterReader(strings.NewReader(input), New())
	var got []string
	for fr.Next() {
		got = append(got, fr.Line())
	}
	if fr.Err() != nil {
		t.Fatalf("unexpected error: %v", fr.Err())
	}
	if len(got) != 3 {
		t.Fatalf("expected 3 lines, got %d", len(got))
	}
}

func TestFilterReader_RemovesDuplicates(t *testing.T) {
	input := "dup\ndup\nunique\ndup\n"
	fr := NewFilterReader(strings.NewReader(input), New())
	var got []string
	for fr.Next() {
		got = append(got, fr.Line())
	}
	// "dup" appears first, then "unique"; second "dup" is within window.
	if len(got) != 2 {
		t.Fatalf("expected 2 unique lines, got %d: %v", len(got), got)
	}
	if got[0] != "dup" || got[1] != "unique" {
		t.Fatalf("unexpected lines: %v", got)
	}
}

func TestFilterReader_EmptyInput(t *testing.T) {
	fr := NewFilterReader(strings.NewReader(""), New())
	for fr.Next() {
		t.Fatal("expected no lines from empty input")
	}
}

func TestLines_Helper(t *testing.T) {
	input := "x\ny\nx\nz\ny\n"
	got, err := Lines(strings.NewReader(input), New())
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	want := []string{"x", "y", "z"}
	if len(got) != len(want) {
		t.Fatalf("expected %v, got %v", want, got)
	}
	for i, w := range want {
		if got[i] != w {
			t.Errorf("index %d: want %q, got %q", i, w, got[i])
		}
	}
}

func TestLines_AllDuplicates(t *testing.T) {
	input := "same\nsame\nsame\n"
	got, err := Lines(strings.NewReader(input), New())
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(got) != 1 {
		t.Fatalf("expected 1 unique line, got %d", len(got))
	}
}
