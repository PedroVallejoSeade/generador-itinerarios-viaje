package main

import (
	"bufio"
	"fmt"
	"io"
	"strings"

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

// closingMessage is printed once on any clean exit path — an exit keyword or
// end-of-input (FR-011).
const closingMessage = "Goodbye! Safe travels.\n"

// runInteractive drives the interactive session: it prints the welcome message
// once, then loops prompting for a city name, reading a full line, trimming it
// (FR-004), and rendering up to 10 population-ranked matches for a non-empty
// query (FR-006, FR-007). Entering "exit"/"quit" (case-insensitive, full line)
// or signaling end-of-input prints a closing message and returns 0 (FR-011).
func runInteractive(in io.Reader, out io.Writer, cities []city.City) int {
	fmt.Fprint(out, welcomeMessage)

	scanner := bufio.NewScanner(in)
	for {
		fmt.Fprint(out, prompt)
		if !scanner.Scan() {
			break
		}

		query := strings.TrimSpace(scanner.Text())

		if strings.EqualFold(query, "exit") || strings.EqualFold(query, "quit") {
			break
		}

		results, err := city.Search(cities, query)
		if err != nil {
			// Empty/whitespace-only query: friendly message, no search (FR-008).
			fmt.Fprintln(out, "Please enter a city name.")
			continue
		}
		if len(results) == 0 {
			// No matches: clear message, then re-prompt (FR-009).
			fmt.Fprintf(out, "No cities found matching %q.\n", query)
			continue
		}

		for i, c := range results {
			fmt.Fprintf(out, "%d. %s\n", i+1, format(c))
		}
	}

	fmt.Fprint(out, closingMessage)
	return 0
}
