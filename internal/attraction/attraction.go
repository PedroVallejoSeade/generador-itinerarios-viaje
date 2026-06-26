// Package attraction looks up a selected city's most well-known attractions
// from the Wikidata Query Service (WDQS) SPARQL endpoint, ranks them by a
// prominence signal (the count of Wikipedia language sitelinks), and caps the
// result at the top 10. The HTTP fetch is injected through a Fetcher seam so the
// durable mapping/ranking logic is exercised with recorded JSON fixtures and no
// live network (research.md, data-model.md).
package attraction

import (
	"context"
	"fmt"
	"sort"
	"strconv"
	"strings"
)

// MaxResults is the cap on how many attractions are displayed (FR-004).
const MaxResults = 10

// Attraction is a notable place to visit in a city — a landmark, museum, or
// sightseeing spot — materialized from one SPARQL result binding.
type Attraction struct {
	Name        string // Display name (?itemLabel); full Unicode, required.
	Category    string // Optional instance-of label (?typeLabel), e.g. "museum".
	Description string // Optional Wikidata short description (?itemDescription).
	Prominence  int    // Wikipedia sitelink count (?sitelinks); ranking key.
}

// mapBindings turns decoded SPARQL bindings into Attraction values (FR-005). A
// binding without a usable label is skipped, and a missing/blank sitelink count
// defaults Prominence to 0 so it sorts last (data-model.md).
func mapBindings(resp *sparqlResponse) []Attraction {
	if resp == nil {
		return nil
	}
	var out []Attraction
	for _, b := range resp.Results.Bindings {
		name := strings.TrimSpace(b["itemLabel"].Value)
		if name == "" {
			continue
		}
		prominence, _ := strconv.Atoi(strings.TrimSpace(b["sitelinks"].Value))
		out = append(out, Attraction{
			Name:        name,
			Category:    strings.TrimSpace(b["typeLabel"].Value),
			Description: strings.TrimSpace(b["itemDescription"].Value),
			Prominence:  prominence,
		})
	}
	return out
}

// Rank orders attractions most-known first — descending by Prominence, then
// ascending by Name as a stable tiebreaker — and caps the result at MaxResults
// (FR-003, FR-004). Fewer than MaxResults are returned in full (FR-007); a
// nil/empty input yields an empty slice (FR-006).
func Rank(items []Attraction) []Attraction {
	ranked := make([]Attraction, len(items))
	copy(ranked, items)
	sort.SliceStable(ranked, func(i, j int) bool {
		if ranked[i].Prominence != ranked[j].Prominence {
			return ranked[i].Prominence > ranked[j].Prominence
		}
		return ranked[i].Name < ranked[j].Name
	})
	if len(ranked) > MaxResults {
		ranked = ranked[:MaxResults]
	}
	return ranked
}

// Lookup resolves the selected city's attractions: it builds the SPARQL query,
// fetches the JSON results through the injected fetcher, decodes and maps the
// bindings to Attraction values, and ranks/caps them. A fetcher or decode
// failure is wrapped and returned — never reported as an empty set — so the CLI
// can surface a source/connectivity error (FR-009).
func Lookup(ctx context.Context, fetch Fetcher, sel CitySelection) ([]Attraction, error) {
	if fetch == nil {
		return nil, fmt.Errorf("attraction lookup for %q: no fetcher configured", sel.Name)
	}
	raw, err := fetch(ctx, buildQuery(sel))
	if err != nil {
		return nil, fmt.Errorf("fetch attractions for %q: %w", sel.Name, err)
	}
	resp, err := decodeSPARQL(raw)
	if err != nil {
		return nil, fmt.Errorf("parse attractions for %q: %w", sel.Name, err)
	}
	return Rank(mapBindings(resp)), nil
}
