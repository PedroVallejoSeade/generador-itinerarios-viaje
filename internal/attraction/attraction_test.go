package attraction

import (
	"os"
	"path/filepath"
	"testing"
)

// loadFixture reads a recorded WDQS SPARQL JSON response from testdata.
func loadFixture(t *testing.T, name string) []byte {
	t.Helper()
	b, err := os.ReadFile(filepath.Join("testdata", name))
	if err != nil {
		t.Fatalf("read fixture %s: %v", name, err)
	}
	return b
}

// TestMapBindings_Golden covers JSON→Attraction mapping over the Paris golden
// fixture: name, category, description, and prominence are populated (T009,
// FR-005). It is table-driven over a representative set of bindings.
func TestMapBindings_Golden(t *testing.T) {
	resp, err := decodeSPARQL(loadFixture(t, "paris.json"))
	if err != nil {
		t.Fatalf("decode paris.json: %v", err)
	}
	got := mapBindings(resp)
	if len(got) != 7 {
		t.Fatalf("mapped %d attractions, want 7", len(got))
	}

	byName := make(map[string]Attraction, len(got))
	for _, a := range got {
		byName[a.Name] = a
	}

	cases := []struct {
		name       string
		category   string
		prominence int
		wantDesc   bool
	}{
		{"Eiffel Tower", "tower", 196, true},
		{"Louvre", "art museum", 168, true},
		{"Notre-Dame de Paris", "cathedral", 152, true},
		{"Sacré-Cœur, Paris", "minor basilica", 96, true},
	}
	for _, tc := range cases {
		a, ok := byName[tc.name]
		if !ok {
			t.Errorf("attraction %q missing from mapping", tc.name)
			continue
		}
		if a.Category != tc.category {
			t.Errorf("%s category = %q, want %q", tc.name, a.Category, tc.category)
		}
		if a.Prominence != tc.prominence {
			t.Errorf("%s prominence = %d, want %d", tc.name, a.Prominence, tc.prominence)
		}
		if tc.wantDesc && a.Description == "" {
			t.Errorf("%s description should be mapped", tc.name)
		}
	}
}

// TestMapBindings_SkipsEmptyNameAndDefaultsProminence covers the validation
// rules: a binding without a usable label is skipped, and a missing sitelink
// count defaults Prominence to 0 (T009, data-model.md).
func TestMapBindings_SkipsEmptyNameAndDefaultsProminence(t *testing.T) {
	resp := &sparqlResponse{}
	resp.Results.Bindings = []map[string]sparqlValue{
		{"itemLabel": {Value: "  "}, "sitelinks": {Value: "5"}},        // skipped: blank name
		{"itemLabel": {Value: "Keep"}},                                 // no sitelinks → 0
		{"itemLabel": {Value: "Also Keep"}, "sitelinks": {Value: "x"}}, // non-numeric → 0
	}
	got := mapBindings(resp)
	if len(got) != 2 {
		t.Fatalf("mapped %d attractions, want 2 (blank skipped)", len(got))
	}
	for _, a := range got {
		if a.Prominence != 0 {
			t.Errorf("%s prominence = %d, want 0 default", a.Name, a.Prominence)
		}
	}
}

// TestRank_OrdersByProminenceThenName covers the ranking rule: descending by
// prominence, ascending by name as a stable tiebreaker (T010, FR-003).
func TestRank_OrdersByProminenceThenName(t *testing.T) {
	in := []Attraction{
		{Name: "Beta", Prominence: 10},
		{Name: "Gamma", Prominence: 50},
		{Name: "Alpha", Prominence: 10},
	}
	got := Rank(in)
	want := []string{"Gamma", "Alpha", "Beta"}
	if len(got) != len(want) {
		t.Fatalf("ranked %d, want %d", len(got), len(want))
	}
	for i, name := range want {
		if got[i].Name != name {
			t.Errorf("rank[%d] = %q, want %q", i, got[i].Name, name)
		}
	}
}

// TestRank_CapsAtTen covers the 10-result cap from a >10 source set, and that
// the kept rows are the highest-prominence ones in descending order (T010,
// FR-004).
func TestRank_CapsAtTen(t *testing.T) {
	resp, err := decodeSPARQL(loadFixture(t, "many.json"))
	if err != nil {
		t.Fatalf("decode many.json: %v", err)
	}
	got := Rank(mapBindings(resp))
	if len(got) != MaxResults {
		t.Fatalf("ranked %d, want cap of %d", len(got), MaxResults)
	}
	for i := 1; i < len(got); i++ {
		if got[i-1].Prominence < got[i].Prominence {
			t.Errorf("rank not descending at %d: %d before %d", i, got[i-1].Prominence, got[i].Prominence)
		}
	}
	// The two lowest-prominence entries (10, 20) must have been dropped.
	for _, a := range got {
		if a.Prominence < 30 {
			t.Errorf("entry %q (prominence %d) should have been capped out", a.Name, a.Prominence)
		}
	}
}

// TestRank_ShortReturnsAll covers FR-007: fewer than the cap returns all rows
// without padding (T011).
func TestRank_ShortReturnsAll(t *testing.T) {
	resp, err := decodeSPARQL(loadFixture(t, "paris.json"))
	if err != nil {
		t.Fatalf("decode paris.json: %v", err)
	}
	got := Rank(mapBindings(resp))
	if len(got) != 7 {
		t.Fatalf("ranked %d, want all 7 (short set, no padding)", len(got))
	}
	if got[0].Name != "Eiffel Tower" {
		t.Errorf("top attraction = %q, want Eiffel Tower", got[0].Name)
	}
}

// TestRank_EmptyReturnsEmpty covers FR-006: an empty source set ranks to an
// empty result (the CLI turns this into a "no attractions found" message) (T011).
func TestRank_EmptyReturnsEmpty(t *testing.T) {
	resp, err := decodeSPARQL(loadFixture(t, "empty.json"))
	if err != nil {
		t.Fatalf("decode empty.json: %v", err)
	}
	if got := Rank(mapBindings(resp)); len(got) != 0 {
		t.Fatalf("ranked %d, want 0 for empty source", len(got))
	}
}
