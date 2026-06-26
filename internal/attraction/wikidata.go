package attraction

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strings"
)

// endpoint is the public Wikidata Query Service SPARQL endpoint (no auth, free).
const endpoint = "https://query.wikidata.org/sparql"

// userAgent identifies this client, as the WDQS usage policy requires a
// descriptive User-Agent for unauthenticated requests (research.md).
const userAgent = "generador-itinerarios-viaje/0.1 (city attractions spike; " +
	"https://github.com/PedroVallejoSeade/generador-itinerarios-viaje)"

// Fetcher retrieves raw SPARQL JSON results for a query string. Injecting it
// lets the durable mapping/ranking logic be tested with recorded fixtures and
// no live network (T006).
type Fetcher func(ctx context.Context, query string) ([]byte, error)

// sparqlValue is a single bound value in a SPARQL JSON result row.
type sparqlValue struct {
	Value string `json:"value"`
}

// sparqlResponse mirrors the subset of the SPARQL JSON results format this
// spike consumes: the rows under results.bindings, each a var→value map.
type sparqlResponse struct {
	Results struct {
		Bindings []map[string]sparqlValue `json:"bindings"`
	} `json:"results"`
}

// decodeSPARQL parses raw WDQS SPARQL JSON results (T014).
func decodeSPARQL(raw []byte) (*sparqlResponse, error) {
	var resp sparqlResponse
	if err := json.Unmarshal(raw, &resp); err != nil {
		return nil, fmt.Errorf("decode SPARQL JSON: %w", err)
	}
	return &resp, nil
}

// escapeSPARQLString escapes a value for safe inclusion inside a SPARQL string
// literal, neutralizing characters that would otherwise break out of the quotes
// (FR-002a keeps accented/Unicode names intact while staying injection-safe).
func escapeSPARQLString(s string) string {
	replacer := strings.NewReplacer(
		`\`, `\\`,
		`"`, `\"`,
		"\n", `\n`,
		"\r", `\r`,
		"\t", `\t`,
	)
	return replacer.Replace(s)
}

// attractionTypes is the curated set of Wikidata classes (instance-of, direct)
// treated as sightseeing attractions: museums, landmarks/towers, churches,
// castles/palaces, monuments, parks, squares, bridges, etc. Using direct P31
// membership (rather than a deep subclass walk) keeps the query fast on the
// public endpoint (research.md: keep the query tight).
const attractionTypes = "wd:Q33506 wd:Q570116 wd:Q12518 wd:Q16970 wd:Q23413 " +
	"wd:Q16560 wd:Q4989906 wd:Q22698 wd:Q174782 wd:Q12280 wd:Q207694 " +
	"wd:Q190603 wd:Q44613 wd:Q839954 wd:Q5773747 wd:Q2087181"

// buildQuery composes a SPARQL query for the most well-known attractions in the
// selected city. The city is resolved by name through the Wikidata entity-search
// MWAPI and constrained to a human settlement in the named country (FR-002a),
// then attractions located anywhere within it (transitive P131*) and typed as a
// sightseeing class are ranked by Wikipedia sitelink count and limited to the
// top results (FR-003, FR-004). The city name and country are preserved with
// full Unicode so accented/non-Latin names resolve correctly.
func buildQuery(sel CitySelection) string {
	name := escapeSPARQLString(sel.Name)
	country := escapeSPARQLString(sel.Country)
	return fmt.Sprintf(`SELECT DISTINCT ?item ?itemLabel ?itemDescription ?sitelinks WHERE {
  SERVICE wikibase:mwapi {
    bd:serviceParam wikibase:endpoint "www.wikidata.org" ;
                    wikibase:api "EntitySearch" ;
                    mwapi:search "%s" ;
                    mwapi:language "en" .
    ?city wikibase:apiOutputItem mwapi:item .
  }
  ?city wdt:P31/wdt:P279* wd:Q486972 ;
        wdt:P17 ?country .
  ?country rdfs:label|skos:altLabel "%s"@en .
  ?item wdt:P131* ?city ;
        wikibase:sitelinks ?sitelinks ;
        wdt:P31 ?type .
  VALUES ?type { %s }
  SERVICE wikibase:label { bd:serviceParam wikibase:language "en". }
}
ORDER BY DESC(?sitelinks)
LIMIT %d`, name, country, attractionTypes, MaxResults)
}

// DefaultFetcher issues the SPARQL query as an HTTPS GET to the WDQS endpoint
// with a descriptive User-Agent (WDQS policy) and requests JSON results. The
// per-lookup deadline is supplied by the caller's context (SC-002).
func DefaultFetcher(ctx context.Context, query string) ([]byte, error) {
	endpointURL, err := url.Parse(endpoint)
	if err != nil {
		return nil, fmt.Errorf("parse endpoint: %w", err)
	}
	q := endpointURL.Query()
	q.Set("query", query)
	q.Set("format", "json")
	endpointURL.RawQuery = q.Encode()

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, endpointURL.String(), nil)
	if err != nil {
		return nil, fmt.Errorf("build request: %w", err)
	}
	req.Header.Set("User-Agent", userAgent)
	req.Header.Set("Accept", "application/sparql-results+json")

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("query WDQS: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("WDQS returned status %s", resp.Status)
	}
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("read WDQS response: %w", err)
	}
	return body, nil
}
