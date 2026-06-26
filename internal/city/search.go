package city

import (
	"errors"
	"sort"
	"strings"
)

// maxResults caps the number of cities returned for a single query (FR-006).
const maxResults = 10

// ErrEmptyQuery is returned when the search query is empty or whitespace-only (FR-008).
var ErrEmptyQuery = errors.New("empty query")

// normalize trims surrounding whitespace and lowercases the query for
// case-insensitive prefix matching (FR-003).
func normalize(raw string) string {
	return strings.ToLower(strings.TrimSpace(raw))
}

// Search returns the cities whose name starts with the query (case-insensitive
// prefix match, FR-003). An empty or whitespace-only query is rejected with
// ErrEmptyQuery without scanning the dataset (FR-008).
func Search(cities []City, raw string) ([]City, error) {
	q := normalize(raw)
	if q == "" {
		return nil, ErrEmptyQuery
	}

	var matches []City
	for _, c := range cities {
		if strings.HasPrefix(strings.ToLower(c.Name), q) {
			matches = append(matches, c)
		}
	}

	// Rank by population (descending) with name as a stable tiebreaker, then cap
	// the result set at maxResults (FR-006).
	sort.SliceStable(matches, func(i, j int) bool {
		if matches[i].Population != matches[j].Population {
			return matches[i].Population > matches[j].Population
		}
		return matches[i].Name < matches[j].Name
	})

	if len(matches) > maxResults {
		matches = matches[:maxResults]
	}

	return matches, nil
}
