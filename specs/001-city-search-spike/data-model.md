# Phase 1 Data Model — City Search Spike

Derived from the spec's **Key Entities** section and the resolved data source
([research.md](research.md): bundled SimpleMaps World Cities CSV).

## Entity: City

A populated place a user can search for and disambiguate.

| Field | Type | Source column | Required | Notes |
|-------|------|---------------|----------|-------|
| `Name` | string | `city` | Yes | Display name, may contain accents/non-ASCII (e.g., "São Paulo"). Used for matching and display. |
| `Country` | string | `country` | Yes | Human-readable country name. Disambiguation field (FR-004). |
| `Region` | string | `admin_name` | Yes (may be empty for a few entries) | State/province/administrative division. Disambiguation field (FR-004). Empty string allowed where the source has none. |
| `Population` | int64 | `population` | No (defaults to 0) | Ranking key (FR-006). Missing/blank parses to 0 and sorts last. |

**Validation rules**

- A row with an empty `Name` is skipped at load time (cannot be matched or displayed).
- `Population` that is blank or non-numeric defaults to `0` (does not fail the load).
- `Region` may be empty; display logic omits an empty region gracefully (renders "Name, Country").

**Relationships**: None. Cities are independent records in a flat in-memory slice.

**Lifecycle / state**: Immutable, read-only. Loaded once from the embedded CSV at process
start; no create/update/delete.

## Entity: SearchQuery

The user's input text and its normalized form used for matching.

| Field | Type | Derivation | Notes |
|-------|------|-----------|-------|
| `Raw` | string | The argument the user typed | Preserved for messages/echo. |
| `Normalized` | string | `strings.TrimSpace` then `strings.ToLower(Raw)` | Used for case-insensitive prefix matching (FR-003). |

**Validation rules**

- If `Normalized` is empty (blank or whitespace-only input), the query is **rejected** without
  touching the dataset; the app prompts for a valid city name (FR-008).

## Entity: SearchResultSet

The ordered, bounded collection of cities returned for a query.

| Property | Value | Rule |
|----------|-------|------|
| Ordering | Descending by `Population`, then ascending by `Name` as a stable tiebreaker | FR-006 (population ranking) |
| Cap | At most **10** entries | FR-006 / spec Clarifications |
| Emptiness | An empty set triggers the "no results found" message (FR-007) | — |

**Matching rule (query → result set)**

A `City` matches when `strings.HasPrefix(strings.ToLower(city.Name), query.Normalized)` is
true — case-insensitive **prefix** match (FR-003, spec Clarifications). The full matching set
is sorted by the ordering rule above and truncated to the cap.

## Traceability

| Requirement | Model element |
|-------------|---------------|
| FR-003 (case-insensitive prefix match) | `SearchQuery.Normalized` + prefix matching rule |
| FR-004 (country + region context) | `City.Country`, `City.Region` |
| FR-006 (cap 10, population rank) | `SearchResultSet` ordering + cap |
| FR-007 (no-results message) | `SearchResultSet` emptiness rule |
| FR-008 (reject empty query) | `SearchQuery` validation rule |
