package main

import (
	"bufio"
	"context"
	"fmt"
	"io"
	"strconv"
	"strings"
	"time"

	"github.com/PedroVallejoSeade/generador-itinerarios-viaje/internal/attraction"
	"github.com/PedroVallejoSeade/generador-itinerarios-viaje/internal/city"
)

// welcomeMessage frames the tool, instructs the user to enter a city name, and
// shows a visible exit hint. Printed once before the first prompt (FR-001,
// FR-002a).
const welcomeMessage = "Welcome to the Travel Itinerary Generator!\n" +
	"Enter a city name to find matching destinations.\n" +
	"(type 'exit' or press Ctrl+D to quit)\n"

// prompt is shown before each line of input (FR-002, FR-010).
const prompt = "city> "

// selectHint invites the user to pick a listed city by its number to see that
// city's attractions (contract step 2). It is printed after a result list.
const selectHint = "Enter a city's number to see its top attractions.\n"

// closingMessage is printed once on any clean exit path — an exit keyword or
// end-of-input (FR-011).
const closingMessage = "Goodbye! Safe travels.\n"

// lookupTimeout bounds a single attractions lookup so a slow/unreachable source
// fails fast rather than hanging the session (FR-009). It is generous relative
// to the SC-002 display target because the public WDQS endpoint's latency varies
// (a documented spike risk — see research.md).
const lookupTimeout = 30 * time.Second

// runInteractive drives the interactive session: it prints the welcome message
// once, then loops prompting for input. A non-numeric line is treated as a city
// search and renders up to 10 population-ranked, numbered matches (FR-006,
// FR-007). When matches are on screen, a numeric line selects one of them by its
// 1-based number and displays that city's top attractions (FR-001); an
// out-of-range number is rejected without a lookup (FR-008). Entering
// "exit"/"quit" or signaling end-of-input prints a closing message and returns 0
// (FR-011); a data-source/connectivity failure during a lookup is reported on
// errOut and returns a non-zero exit (FR-009).
func runInteractive(in io.Reader, out, errOut io.Writer, cities []city.City, fetch attraction.Fetcher) int {
	fmt.Fprint(out, welcomeMessage)

	scanner := bufio.NewScanner(in)
	var results []city.City // The most recently displayed matches, for selection.
	for {
		fmt.Fprint(out, prompt)
		if !scanner.Scan() {
			break
		}

		line := strings.TrimSpace(scanner.Text())

		if strings.EqualFold(line, "exit") || strings.EqualFold(line, "quit") {
			break
		}

		// With matches on screen, a numeric line is a city selection (FR-001);
		// otherwise the line is treated as a new search query.
		if n, err := strconv.Atoi(line); err == nil && len(results) > 0 {
			if n < 1 || n > len(results) {
				// Out-of-range selection: reject without a lookup (FR-008).
				fmt.Fprintf(out, "%q is not a valid selection. Enter the number of a listed city.\n", line)
				continue
			}
			sel := attraction.NewCitySelection(n, results[n-1])
			if code := showAttractions(out, errOut, fetch, sel); code != 0 {
				return code
			}
			continue
		}

		searchResults, err := city.Search(cities, line)
		if err != nil {
			// Empty/whitespace-only query: friendly message, no search (FR-008).
			fmt.Fprintln(out, "Please enter a city name.")
			continue
		}
		if len(searchResults) == 0 {
			// No matches: clear message, then re-prompt (FR-009).
			fmt.Fprintf(out, "No cities found matching %q.\n", line)
			continue
		}

		results = searchResults
		for i, c := range results {
			fmt.Fprintf(out, "%d. %s\n", i+1, format(c))
		}
		fmt.Fprint(out, selectHint)
	}

	fmt.Fprint(out, closingMessage)
	return 0
}

// showAttractions looks up and renders the selected city's top attractions. On
// success it prints a numbered, ranked list (FR-003, FR-004, FR-005) or a clear
// "no attractions found" message when the city has none (FR-006), returning 0 in
// both cases. A source/connectivity failure is reported on errOut and returns
// exit code 2 (FR-009).
func showAttractions(out, errOut io.Writer, fetch attraction.Fetcher, sel attraction.CitySelection) int {
	ctx, cancel := context.WithTimeout(context.Background(), lookupTimeout)
	defer cancel()

	items, err := attraction.Lookup(ctx, fetch, sel)
	if err != nil {
		fmt.Fprintf(errOut, "error: unable to fetch attractions: %v\n", err)
		return 2
	}

	if len(items) == 0 {
		fmt.Fprintf(out, "No attractions found for %s.\n", formatSelection(sel))
		return 0
	}

	fmt.Fprintf(out, "Top attractions in %s:\n", formatSelection(sel))
	for i, a := range items {
		fmt.Fprintf(out, "%d. %s\n", i+1, formatAttraction(a))
	}
	return 0
}

// formatSelection renders the chosen city as "<Name>, <Country>" for the
// attractions heading and messages (contract examples).
func formatSelection(sel attraction.CitySelection) string {
	return fmt.Sprintf("%s, %s", sel.Name, sel.Country)
}

// formatAttraction renders one attraction as its name, appending the category
// or short description when readily available (FR-005).
func formatAttraction(a attraction.Attraction) string {
	switch {
	case a.Category != "":
		return fmt.Sprintf("%s — %s", a.Name, a.Category)
	case a.Description != "":
		return fmt.Sprintf("%s — %s", a.Name, a.Description)
	default:
		return a.Name
	}
}
