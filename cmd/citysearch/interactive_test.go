package main

import (
	"bytes"
	"context"
	"io"
	"strings"
	"testing"

	"github.com/PedroVallejoSeade/generador-itinerarios-viaje/internal/city"
)

// noLookupFetcher fails the test if it is ever called: the existing search-only
// scenarios never select a city, so no attractions lookup must occur.
func noLookupFetcher(t *testing.T) func(context.Context, string) ([]byte, error) {
	return func(context.Context, string) ([]byte, error) {
		t.Helper()
		t.Fatal("attraction fetcher called unexpectedly (no selection was made)")
		return nil, nil
	}
}

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
	code := runInteractive(strings.NewReader(""), &out, io.Discard, fixtureCities(), noLookupFetcher(t))
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
	code := runInteractive(strings.NewReader("London\n"), &out, io.Discard, fixtureCities(), noLookupFetcher(t))
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
	code := runInteractive(strings.NewReader("   \n"), &out, io.Discard, fixtureCities(), noLookupFetcher(t))
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
	code := runInteractive(strings.NewReader("zzzzzz\n"), &out, io.Discard, fixtureCities(), noLookupFetcher(t))
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
	code := runInteractive(strings.NewReader("London\nParis\nexit\n"), &out, io.Discard, fixtureCities(), noLookupFetcher(t))
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
			code := runInteractive(strings.NewReader(input), &out, io.Discard, fixtureCities(), noLookupFetcher(t))
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
	code := runInteractive(strings.NewReader("London\n"), &out, io.Discard, fixtureCities(), noLookupFetcher(t))
	if code != 0 {
		t.Fatalf("runInteractive exit = %d, want 0", code)
	}
	if got := out.String(); !strings.Contains(got, "Goodbye") {
		t.Errorf("output = %q, want a closing message on EOF", got)
	}
}

// fixedFetcher returns the same recorded SPARQL JSON for any query.
func fixedFetcher(body string) func(context.Context, string) ([]byte, error) {
	return func(context.Context, string) ([]byte, error) {
		return []byte(body), nil
	}
}

const parisAttractionsJSON = `{"results":{"bindings":[
	{"itemLabel":{"value":"Eiffel Tower"},"typeLabel":{"value":"tower"},"sitelinks":{"value":"196"}},
	{"itemLabel":{"value":"Louvre"},"typeLabel":{"value":"art museum"},"sitelinks":{"value":"168"}}
]}}`

const emptyAttractionsJSON = `{"results":{"bindings":[]}}`

// TestRunInteractive_SelectionShowsAttractions covers contract scenario A1
// (FR-001, FR-003, FR-005): after a search, selecting a city by its number
// renders that city's ranked, numbered attractions.
func TestRunInteractive_SelectionShowsAttractions(t *testing.T) {
	var out, errOut bytes.Buffer
	code := runInteractive(strings.NewReader("Paris\n1\nexit\n"), &out, &errOut, fixtureCities(), fixedFetcher(parisAttractionsJSON))
	if code != 0 {
		t.Fatalf("runInteractive exit = %d, want 0", code)
	}
	got := out.String()
	if !strings.Contains(got, "Top attractions in Paris, France:") {
		t.Errorf("output = %q, want the attractions heading", got)
	}
	for _, want := range []string{"Eiffel Tower — tower", "Louvre — art museum"} {
		if !strings.Contains(got, want) {
			t.Errorf("output = %q, want attraction line %q", got, want)
		}
	}
	// The most-known attraction (highest sitelinks) ranks first.
	if ei, lo := strings.Index(got, "Eiffel Tower"), strings.Index(got, "Louvre"); ei < 0 || lo < 0 || ei > lo {
		t.Errorf("Eiffel Tower should rank before Louvre; got %q", got)
	}
	if errOut.Len() != 0 {
		t.Errorf("stderr = %q, want empty on success", errOut.String())
	}
}

// TestRunInteractive_InvalidSelection covers contract scenario A5 (FR-008): a
// number outside the result list is rejected with a clear message and no lookup
// is attempted.
func TestRunInteractive_InvalidSelection(t *testing.T) {
	var out, errOut bytes.Buffer
	// fixtureCities has 2 Paris matches; 9 is out of range.
	code := runInteractive(strings.NewReader("Paris\n9\nexit\n"), &out, &errOut, fixtureCities(), noLookupFetcher(t))
	if code != 0 {
		t.Fatalf("runInteractive exit = %d, want 0", code)
	}
	got := out.String()
	if !strings.Contains(got, "is not a valid selection") {
		t.Errorf("output = %q, want an invalid-selection message", got)
	}
	if strings.Contains(got, "Top attractions") {
		t.Errorf("output = %q, want no attractions rendered for an invalid selection", got)
	}
}

// TestRunInteractive_NoAttractions covers contract scenario A4 (FR-006): a valid
// selection that finds nothing yields a clear message and the session continues
// (exit 0).
func TestRunInteractive_NoAttractions(t *testing.T) {
	var out, errOut bytes.Buffer
	code := runInteractive(strings.NewReader("Paris\n1\nexit\n"), &out, &errOut, fixtureCities(), fixedFetcher(emptyAttractionsJSON))
	if code != 0 {
		t.Fatalf("runInteractive exit = %d, want 0", code)
	}
	got := out.String()
	if !strings.Contains(got, "No attractions found for Paris, France.") {
		t.Errorf("output = %q, want the no-attractions message", got)
	}
}

// TestRunInteractive_FetchError covers contract scenario A7 (FR-009): a
// source/connectivity failure is reported on stderr and exits non-zero.
func TestRunInteractive_FetchError(t *testing.T) {
	var out, errOut bytes.Buffer
	failing := func(context.Context, string) ([]byte, error) {
		return nil, errFetch
	}
	code := runInteractive(strings.NewReader("Paris\n1\n"), &out, &errOut, fixtureCities(), failing)
	if code != 2 {
		t.Fatalf("runInteractive exit = %d, want 2 on fetch error", code)
	}
	if !strings.Contains(errOut.String(), "unable to fetch attractions") {
		t.Errorf("stderr = %q, want the fetch-error message", errOut.String())
	}
	if strings.Contains(out.String(), "Top attractions") {
		t.Errorf("stdout = %q, want no attractions on a fetch error", out.String())
	}
}

// errFetch is a sentinel error for the fetch-error scenario.
var errFetch = errorString("connection refused")

type errorString string

func (e errorString) Error() string { return string(e) }
