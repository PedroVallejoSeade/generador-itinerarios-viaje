package city

import (
	"errors"
	"os"
	"strconv"
	"testing"
)

// loadFixture parses the deterministic test dataset.
func loadFixture(t *testing.T) []City {
	t.Helper()
	f, err := os.Open("testdata/cities_sample.csv")
	if err != nil {
		t.Fatalf("open fixture: %v", err)
	}
	defer f.Close()
	cities, err := parse(f)
	if err != nil {
		t.Fatalf("parse fixture: %v", err)
	}
	return cities
}

// hasName reports whether any result has the given name.
func hasName(results []City, name string) bool {
	for _, c := range results {
		if c.Name == name {
			return true
		}
	}
	return false
}

func TestSearchPrefixCaseInsensitive(t *testing.T) {
	cities := loadFixture(t)

	tests := []struct {
		name      string
		query     string
		wantName  string
		wantFound bool
	}{
		{"lowercase prefix", "par", "Paris", true},
		{"uppercase exact", "PARIS", "Paris", true},
		{"mixed case", "PaRiS", "Paris", true},
		{"prefix matches multiple names", "par", "Parma", true},
		{"accented name lower prefix", "são", "São Paulo", true},
		{"surrounding whitespace trimmed", "  london  ", "London", true},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			results, err := Search(cities, tc.query)
			if err != nil {
				t.Fatalf("Search(%q) returned error: %v", tc.query, err)
			}
			if got := hasName(results, tc.wantName); got != tc.wantFound {
				t.Errorf("Search(%q): found %q = %v, want %v", tc.query, tc.wantName, got, tc.wantFound)
			}
		})
	}
}

func TestSearchEmptyQueryRejected(t *testing.T) {
	cities := loadFixture(t)

	for _, q := range []string{"", "   ", "\t", "\n  "} {
		t.Run("query="+q, func(t *testing.T) {
			results, err := Search(cities, q)
			if !errors.Is(err, ErrEmptyQuery) {
				t.Errorf("Search(%q): err = %v, want ErrEmptyQuery", q, err)
			}
			if results != nil {
				t.Errorf("Search(%q): results = %v, want nil", q, results)
			}
		})
	}
}

func TestSearchNoResults(t *testing.T) {
	cities := loadFixture(t)

	results, err := Search(cities, "zzzzzz")
	if err != nil {
		t.Fatalf("Search returned error: %v", err)
	}
	if len(results) != 0 {
		t.Errorf("Search(\"zzzzzz\"): got %d results, want 0", len(results))
	}
}

func TestSearchResultsCarryContext(t *testing.T) {
	cities := loadFixture(t)

	results, err := Search(cities, "springfield")
	if err != nil {
		t.Fatalf("Search returned error: %v", err)
	}
	if len(results) < 2 {
		t.Fatalf("Search(\"springfield\"): got %d results, want >= 2", len(results))
	}

	// Every result must expose country and region for disambiguation (FR-004).
	for _, c := range results {
		if c.Country == "" {
			t.Errorf("result %+v missing country", c)
		}
		if c.Region == "" {
			t.Errorf("result %+v missing region", c)
		}
	}

	// Two same-name, same-country results must differ by region.
	seen := map[string]bool{}
	for _, c := range results {
		if c.Name == "Springfield" && c.Country == "United States" {
			if seen[c.Region] {
				t.Errorf("duplicate region %q for Springfield, United States", c.Region)
			}
			seen[c.Region] = true
		}
	}
	if len(seen) < 2 {
		t.Errorf("expected >= 2 distinct US Springfield regions, got %d", len(seen))
	}
}

func TestSearchRankingAndCap(t *testing.T) {
	// Build a fixture with more than maxResults matches to exercise the cap.
	var cities []City
	for i := 0; i < maxResults+5; i++ {
		cities = append(cities, City{
			Name:       "Testville",
			Country:    "Testland",
			Region:     strconv.Itoa(i),
			Population: int64((maxResults + 5 - i) * 1000),
		})
	}
	// Add a zero-population entry that must sort last.
	cities = append(cities, City{Name: "Testville", Country: "Testland", Region: "zero", Population: 0})

	results, err := Search(cities, "test")
	if err != nil {
		t.Fatalf("Search returned error: %v", err)
	}

	if len(results) != maxResults {
		t.Fatalf("got %d results, want cap of %d", len(results), maxResults)
	}

	for i := 1; i < len(results); i++ {
		if results[i-1].Population < results[i].Population {
			t.Errorf("results not population-descending at %d: %d before %d",
				i, results[i-1].Population, results[i].Population)
		}
	}
}

func TestSearchNameTiebreaker(t *testing.T) {
	// Same population — name ascending breaks the tie.
	cities := []City{
		{Name: "Abravo", Country: "X", Region: "r", Population: 100},
		{Name: "Aalpha", Country: "X", Region: "r", Population: 100},
	}
	results, err := Search(cities, "a")
	if err != nil {
		t.Fatalf("Search returned error: %v", err)
	}
	if len(results) != 2 {
		t.Fatalf("got %d results, want 2", len(results))
	}
	if results[0].Name != "Aalpha" || results[1].Name != "Abravo" {
		t.Errorf("name tiebreaker failed: got %q then %q", results[0].Name, results[1].Name)
	}
}
