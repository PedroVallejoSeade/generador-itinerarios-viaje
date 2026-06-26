// Command citysearch is a terminal tool that looks up cities by name and lists
// matching cities with disambiguating country/region context.
package main

import (
	"flag"
	"fmt"
	"os"
	"strings"

	"github.com/PedroVallejoSeade/generador-itinerarios-viaje/internal/city"
)

const attribution = "City data © GeoNames (cities15000), licensed under CC BY 4.0 — https://www.geonames.org/."

const usage = `citysearch — look up cities by name.

Usage:
  citysearch <city-name>
  citysearch -h | --help

Arguments:
  <city-name>   Name (or prefix) of the city to search for. Quote multi-word
                names, e.g. citysearch "san jose".

Output:
  One match per line as "<Name>, <Region>, <Country>" (region omitted when
  unknown), ordered by population (largest first) and capped at 10 results.

` + attribution + "\n"

func main() {
	os.Exit(run(os.Args[1:], os.Stdout, os.Stderr))
}

// run executes the CLI and returns the process exit code:
//
//	0 success (including a valid query with no matches)
//	1 invalid usage (empty/whitespace query or missing argument)
//	2 data-source/load error
func run(args []string, stdout, stderr *os.File) int {
	fs := flag.NewFlagSet("citysearch", flag.ContinueOnError)
	fs.SetOutput(stderr)
	fs.Usage = func() { fmt.Fprint(stderr, usage) }
	if err := fs.Parse(args); err != nil {
		if err == flag.ErrHelp {
			// -h/--help: usage already printed; success.
			return 0
		}
		// flag already printed the error and usage.
		return 1
	}

	if fs.NArg() < 1 {
		fmt.Fprintln(stderr, "Please provide a city name to search for.")
		return 1
	}

	query := fs.Arg(0)

	cities, err := city.Load()
	if err != nil {
		fmt.Fprintf(stderr, "error: unable to load city data: %v\n", err)
		return 2
	}

	results, err := city.Search(cities, query)
	if err != nil {
		// Only ErrEmptyQuery is returned here; treat as invalid usage.
		fmt.Fprintln(stderr, "Please provide a city name to search for.")
		return 1
	}

	if len(results) == 0 {
		fmt.Fprintf(stdout, "No cities found matching %q.\n", strings.TrimSpace(query))
		return 0
	}

	for _, c := range results {
		fmt.Fprintln(stdout, format(c))
	}
	return 0
}

// format renders a city as "<Name>, <Region>, <Country>", omitting the region
// when it is empty (FR-004).
func format(c city.City) string {
	if c.Region == "" {
		return fmt.Sprintf("%s, %s", c.Name, c.Country)
	}
	return fmt.Sprintf("%s, %s, %s", c.Name, c.Region, c.Country)
}
