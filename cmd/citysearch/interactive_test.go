package main

import (
	"bytes"
	"strings"
	"testing"

	"github.com/PedroVallejoSeade/generador-itinerarios-viaje/internal/city"
)

// fixtureCities is a small, deterministic dataset for driving the interactive
// session in tests. It includes three "London" entries (one with an empty
// region) to exercise numbering, population-descending ordering, and graceful
// region omission, plus a multi-word "San Jose" entry.
func fixtureCities() []city.City {
	return []city.City{
		{Name: "London", Region: "England", Country: "United Kingdom", Population: 8961989},
		{Name: "London", Region: "Ontario", Country: "Canada", Population: 383822},
		{Name: "London", Region: "", Country: "Kiribati", Population: 1829},
		{Name: "Paris", Region: "Île-de-France", Country: "France", Population: 2138551},
		{Name: "Paris", Region: "Texas", Country: "United States", Population: 25171},
		{Name: "San Jose", Region: "California", Country: "United States", Population: 1026908},
	}
}

// countNumberedLines counts lines of the form "<n>. ..." (a 1-based result
// line). Because the prompt is written without a trailing newline, the first
// result can share a line with the prompt in scripted (non-echoing) input, so
// any leading prompt text is stripped before classification.
func countNumberedLines(s string) int {
	n := 0
	for _, line := range strings.Split(s, "\n") {
		line = strings.TrimSpace(strings.ReplaceAll(line, prompt, ""))
		if len(line) < 3 {
			continue
		}
		dot := strings.IndexByte(line, '.')
		if dot <= 0 {
			continue
		}
		allDigits := true
		for _, r := range line[:dot] {
			if r < '0' || r > '9' {
				allDigits = false
				break
			}
		}
		if allDigits {
			n++
		}
	}
	return n
}

// TestRunInteractive_WelcomeAndPromptBeforeInput covers contract scenario I1
// (FR-001, FR-002, FR-002a): the welcome message, an exit hint, and a prompt
// are emitted before any input is read. Driving the session with an
// immediately-EOF reader proves the greeting does not depend on input.
func TestRunInteractive_WelcomeAndPromptBeforeInput(t *testing.T) {
	var out bytes.Buffer
	code := runInteractive(strings.NewReader(""), &out, fixtureCities())
	if code != 0 {
		t.Fatalf("runInteractive exit = %d, want 0", code)
	}
	got := out.String()

	for _, want := range []string{"Welcome", "city name", "exit", "Ctrl+D"} {
		if !strings.Contains(got, want) {
			t.Errorf("welcome output = %q, want substring %q", got, want)
		}
	}
	if !strings.Contains(got, "city>") {
		t.Errorf("output = %q, want a prompt containing %q", got, "city>")
	}
	// Welcome must precede the prompt (emitted before any read).
	if wi, pi := strings.Index(got, "Welcome"), strings.Index(got, "city>"); wi < 0 || pi < 0 || wi > pi {
		t.Errorf("welcome must appear before the prompt; got %q", got)
	}
}

// TestRunInteractive_NumberedResults covers US1 numbered rendering (FR-006,
// FR-007): a known query yields 1-based, population-descending lines, ≤10 of
// them, with the region gracefully omitted when empty.
func TestRunInteractive_NumberedResults(t *testing.T) {
	var out bytes.Buffer
	code := runInteractive(strings.NewReader("London\n"), &out, fixtureCities())
	if code != 0 {
		t.Fatalf("runInteractive exit = %d, want 0", code)
	}
	got := out.String()

	wantOrdered := []string{
		"1. London, England, United Kingdom",
		"2. London, Ontario, Canada",
		"3. London, Kiribati", // region omitted when empty
	}
	prev := -1
	for _, want := range wantOrdered {
		idx := strings.Index(got, want)
		if idx < 0 {
			t.Errorf("output = %q, want line %q", got, want)
			continue
		}
		if idx < prev {
			t.Errorf("line %q out of order in %q", want, got)
		}
		prev = idx
	}

	if n := countNumberedLines(got); n != 3 {
		t.Errorf("got %d numbered lines, want 3 (no padding); output=%q", n, got)
	}
}

// TestRunInteractive_EmptyInput covers US2 / contract I5 (FR-008): an
// empty/whitespace-only line shows a friendly message and re-prompts without
// performing a search.
func TestRunInteractive_EmptyInput(t *testing.T) {
	var out bytes.Buffer
	code := runInteractive(strings.NewReader("   \n"), &out, fixtureCities())
	if code != 0 {
		t.Fatalf("runInteractive exit = %d, want 0", code)
	}
	got := out.String()

	if !strings.Contains(got, "Please enter a city name.") {
		t.Errorf("output = %q, want the empty-input message", got)
	}
	if n := countNumberedLines(got); n != 0 {
		t.Errorf("got %d numbered lines, want 0 (no search performed); output=%q", n, got)
	}
	// The prompt is shown again after the friendly message (re-prompt).
	if strings.Count(got, prompt) < 2 {
		t.Errorf("output = %q, want the prompt to re-appear after the message", got)
	}
}

// TestRunInteractive_NoMatch covers US2 / contract I6 (FR-009): a query with no
// matches shows a clear message and re-prompts.
func TestRunInteractive_NoMatch(t *testing.T) {
	var out bytes.Buffer
	code := runInteractive(strings.NewReader("zzzzzz\n"), &out, fixtureCities())
	if code != 0 {
		t.Fatalf("runInteractive exit = %d, want 0", code)
	}
	got := out.String()

	if !strings.Contains(got, `No cities found matching "zzzzzz".`) {
		t.Errorf("output = %q, want the no-match message", got)
	}
	if strings.Count(got, prompt) < 2 {
		t.Errorf("output = %q, want the prompt to re-appear after the message", got)
	}
}

// TestRunInteractive_MultipleSearches covers US3 / contract I7 (FR-010): two
// sequential searches in one session, with the prompt re-appearing between
// them, then a closing message and a 0 exit on the exit keyword.
func TestRunInteractive_MultipleSearches(t *testing.T) {
	var out bytes.Buffer
	code := runInteractive(strings.NewReader("London\nParis\nexit\n"), &out, fixtureCities())
	if code != 0 {
		t.Fatalf("runInteractive exit = %d, want 0", code)
	}
	got := out.String()

	li := strings.Index(got, "1. London, England, United Kingdom")
	pi := strings.Index(got, "1. Paris, Île-de-France, France")
	if li < 0 {
		t.Errorf("output = %q, want a London result block", got)
	}
	if pi < 0 {
		t.Errorf("output = %q, want a Paris result block", got)
	}
	if li >= 0 && pi >= 0 && li > pi {
		t.Errorf("London block should precede Paris block; got %q", got)
	}
	if !strings.Contains(got, "Goodbye") {
		t.Errorf("output = %q, want a closing message", got)
	}
	// Three prompts: before London, before Paris, before the exit keyword.
	if c := strings.Count(got, prompt); c != 3 {
		t.Errorf("got %d prompts, want 3; output=%q", c, got)
	}
}

// TestRunInteractive_ExitKeywords covers US3 / contract I8 (FR-011): the exit
// keywords are matched against the full trimmed line, case-insensitively, and
// each prints a closing message and returns 0 without searching.
func TestRunInteractive_ExitKeywords(t *testing.T) {
	for _, input := range []string{"exit\n", "quit\n", "EXIT\n", "Quit\n", "  exit  \n"} {
		t.Run(strings.TrimSpace(input), func(t *testing.T) {
			var out bytes.Buffer
			code := runInteractive(strings.NewReader(input), &out, fixtureCities())
			if code != 0 {
				t.Fatalf("runInteractive(%q) exit = %d, want 0", input, code)
			}
			got := out.String()
			if !strings.Contains(got, "Goodbye") {
				t.Errorf("input %q: output = %q, want a closing message", input, got)
			}
			if n := countNumberedLines(got); n != 0 {
				t.Errorf("input %q: got %d numbered lines, want 0 (no search)", input, n)
			}
		})
	}
}

// TestRunInteractive_EOFClosesCleanly covers US3 / contract I3 (FR-011): an
// end-of-input (Ctrl+D) prints the closing message and returns 0.
func TestRunInteractive_EOFClosesCleanly(t *testing.T) {
	var out bytes.Buffer
	code := runInteractive(strings.NewReader("London\n"), &out, fixtureCities())
	if code != 0 {
		t.Fatalf("runInteractive exit = %d, want 0", code)
	}
	if got := out.String(); !strings.Contains(got, "Goodbye") {
		t.Errorf("output = %q, want a closing message on EOF", got)
	}
}
