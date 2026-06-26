# Phase 1 Data Model — City Attractions Lookup Spike

Derived from the spec's **Key Entities** section and the resolved data source
([research.md](research.md): Wikidata Query Service SPARQL, ranked by Wikipedia sitelink count).

## Entity: CitySelection

The city the user chose from the search results, carrying enough identity to resolve the
correct city for the attractions lookup (FR-001, FR-002a).

| Field | Type | Source | Required | Notes |
|-------|------|--------|----------|-------|
| `Index` | int | 1-based position the user typed | Yes | Maps the typed number back to a displayed `city.City`. Out-of-range → invalid selection (FR-008). |
| `Name` | string | `city.City.Name` | Yes | May contain accents/non-ASCII (e.g. "São Paulo"); used to resolve the city entity (FR-002a). |
| `Country` | string | `city.City.Country` | Yes | Disambiguates same-named cities (FR-002a). |
| `Region` | string | `city.City.Region` | No (may be empty) | Further disambiguation when available (FR-002a). |

**Validation rules**

- `Index` MUST correspond to a city in the most recently displayed result set; otherwise the
  selection is rejected with a clear message and **no lookup is attempted** (FR-008).
- Resolution is keyed on `Name` + `Country` (plus `Region` when present) to pick the correct
  city entity for the data source (FR-002a).

**Relationships**: References exactly one `city.City` from the prior search result set.

**Lifecycle**: Transient — exists only for the duration of one selection→attractions exchange.

## Entity: Attraction

A notable place to visit in the selected city — a landmark, museum, or sightseeing spot.

| Field | Type | Source (WDQS) | Required | Notes |
|-------|------|---------------|----------|-------|
| `Name` | string | `?itemLabel` | Yes | Display name; full Unicode. |
| `Category` | string | `?typeLabel` (instance-of label, e.g. "museum", "landmark") | No | Optional supporting context shown "when readily available" (FR-005). |
| `Description` | string | `?itemDescription` (Wikidata short description) | No | Optional short description shown when present (FR-005). |
| `Prominence` | int | sitelink count (`?sitelinks`) | Yes (ranking key) | Number of Wikipedia language editions; higher = more well-known (FR-003). |

**Validation rules**

- An entity with an empty `Name` (no usable label) is skipped.
- `Category` / `Description` are optional; rendering omits them gracefully when empty (FR-005).
- `Prominence` defaults to 0 when absent and sorts last.

**Lifecycle**: Immutable value object, materialized per lookup from the SPARQL JSON response.

## Entity: AttractionResultSet

The ordered, bounded collection of attractions returned for a selected city.

| Property | Value | Rule |
|----------|-------|------|
| Ordering | Descending by `Prominence` (sitelink count), then ascending by `Name` as a stable tiebreaker | FR-003 (rank most-known first) |
| Cap | At most **10** entries | FR-004 |
| Short set | Fewer than 10 → display **all** available, no padding, no error | FR-007 |
| Emptiness | An empty set triggers the clear "no attractions found" message | FR-006 |
| Error | Source/connectivity failure is **not** an empty set — it is reported on stderr, exit non-zero | FR-009 |

**Mapping rule (SPARQL JSON → result set)**

Each SPARQL result binding maps to one `Attraction` (`itemLabel`→`Name`, `typeLabel`→
`Category`, `itemDescription`→`Description`, `sitelinks`→`Prominence`). The set is ordered by
the rule above and truncated to the cap of 10.

## Traceability

| Requirement | Model element |
|-------------|---------------|
| FR-001 (accept numbered selection) | `CitySelection.Index` + validation |
| FR-002 / FR-002a (lookup by name + country, accent-safe) | `CitySelection.Name/Country/Region` resolution |
| FR-003 (rank most-known first) | `AttractionResultSet` ordering on `Attraction.Prominence` |
| FR-004 (cap 10) | `AttractionResultSet` cap |
| FR-005 (name + category/description when available) | `Attraction.Name/Category/Description` |
| FR-006 (no-attractions message) | `AttractionResultSet` emptiness rule |
| FR-007 (show all when < 10) | `AttractionResultSet` short-set rule |
| FR-008 (reject invalid selection) | `CitySelection.Index` validation |
| FR-009 (source/connectivity errors) | `AttractionResultSet` error rule |
