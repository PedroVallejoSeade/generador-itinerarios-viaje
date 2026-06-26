package attraction

import (
	"context"
	"errors"
	"strings"
	"testing"
)

// TestLookup_WrapsFetchError covers the error path (T012, FR-009): when the
// fetcher fails, Lookup returns a wrapped error — never an empty result set.
func TestLookup_WrapsFetchError(t *testing.T) {
	sentinel := errors.New("network down")
	fetch := func(ctx context.Context, query string) ([]byte, error) {
		return nil, sentinel
	}

	got, err := Lookup(context.Background(), fetch, CitySelection{Name: "Paris", Country: "France"})
	if err == nil {
		t.Fatal("Lookup error = nil, want a wrapped fetch error")
	}
	if !errors.Is(err, sentinel) {
		t.Errorf("Lookup error = %v, want it to wrap the sentinel", err)
	}
	if got != nil {
		t.Errorf("Lookup result = %v, want nil on error", got)
	}
}

// TestLookup_NilFetcher covers the misconfiguration guard: a nil fetcher yields
// an error rather than a panic (T012).
func TestLookup_NilFetcher(t *testing.T) {
	if _, err := Lookup(context.Background(), nil, CitySelection{Name: "Paris"}); err == nil {
		t.Fatal("Lookup with nil fetcher = nil error, want an error")
	}
}

// TestLookup_FromGoldenFixture ties fetch→decode→map→rank end to end against the
// recorded Paris response, proving the recommended approach returns a plausible,
// ranked top-N list with no live network (T012, FR-014 PoC).
func TestLookup_FromGoldenFixture(t *testing.T) {
	fetch := func(ctx context.Context, query string) ([]byte, error) {
		return loadFixture(t, "paris.json"), nil
	}

	got, err := Lookup(context.Background(), fetch, CitySelection{Name: "Paris", Country: "France"})
	if err != nil {
		t.Fatalf("Lookup: %v", err)
	}
	if len(got) == 0 {
		t.Fatal("Lookup returned no attractions, want a ranked list")
	}
	if got[0].Name != "Eiffel Tower" {
		t.Errorf("top attraction = %q, want Eiffel Tower", got[0].Name)
	}
}

// TestBuildQuery_PreservesNameAndCountry covers FR-002a: accented/Unicode city
// names and the country are carried into the SPARQL query intact (T012/T014).
func TestBuildQuery_PreservesNameAndCountry(t *testing.T) {
	q := buildQuery(CitySelection{Name: "São Paulo", Country: "Brazil"})
	if !strings.Contains(q, "São Paulo") {
		t.Errorf("query missing accented city name; got:\n%s", q)
	}
	if !strings.Contains(q, "Brazil") {
		t.Errorf("query missing country; got:\n%s", q)
	}
}

// TestEscapeSPARQLString_NeutralizesQuotes guards against a quote in the city
// name breaking out of the SPARQL string literal.
func TestEscapeSPARQLString_NeutralizesQuotes(t *testing.T) {
	got := escapeSPARQLString(`a"b\c`)
	if got != `a\"b\\c` {
		t.Errorf("escapeSPARQLString = %q, want %q", got, `a\"b\\c`)
	}
}
