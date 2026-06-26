# Phase 0 Research — City Attractions Lookup Spike Findings

**Spike time-box**: ~30 minutes (per Constitution Principle III)
**Objective**: Evaluate at least three candidate ways to source city attraction data — a free
no-auth public API, a Go library, and a locally queried open dataset — against the project
constraints, then recommend one with a rationale and a proof of concept (FR-012, FR-013,
FR-014, US2).

## Evaluation Criteria

Derived directly from the spec (FR-010, FR-011, FR-012, SC-001..SC-004) and the constitution:

1. **Data quality / relevance** — returns recognizable landmarks, museums, and sightseeing
   spots for a city (SC-001, SC-003).
2. **Prominence / ranking signal** — exposes a popularity/prominence measure so the "top 10
   most known" can be ranked (FR-003, FR-004). This is the decisive criterion: without a
   ranking signal we cannot meaningfully pick the *top* 10.
3. **Ease of Go integration** — minimal client code using the standard library (Principle V).
4. **Authentication** — no API key/account/registration preferred (FR-010, SC-004).
5. **Free** — free for this application's use (FR-011).
6. **Rate limits / availability** — usable for an interactive CLI; network access is acceptable
   per the spec clarification (no offline requirement).

## Candidates Evaluated

| # | Source | Type | Data quality / relevance | Prominence signal | No-Auth | Free | Go integration | Rate limits |
|---|--------|------|--------------------------|-------------------|---------|------|----------------|-------------|
| 1 | **Wikidata Query Service (WDQS) SPARQL** | Public HTTP API (SPARQL) | ✅ Strong — entities typed as tourist attraction / museum / landmark / historic site, resolved within the city's administrative entity | ✅ **Yes** — count of Wikipedia sitelinks (language editions) is a robust global fame signal; can also order by it in SPARQL | ✅ No key/account | ✅ Free | ✅ `net/http` GET + `encoding/json`; SPARQL JSON results are simple to parse | ⚠️ ~60 s query timeout, throttles aggressive use, **requires a descriptive `User-Agent`** |
| 2 | **OpenTripMap API** | Public HTTP API (REST) | ✅ Strong — purpose-built POI/attraction data with kinds, names, and Wikidata/Wikipedia cross-links | ✅ Yes — built-in `rate` popularity score (0–7), ideal for ranking | ❌ **Requires a free API key** | ✅ Free tier | ✅ Simple REST + JSON | ⚠️ Free tier daily/throughput caps; key must be managed |
| 3 | **Overpass API (OpenStreetMap)** | Public HTTP API (Overpass QL) | ✅ Good coverage — `tourism=attraction/museum/artwork`, `historic=*` nodes/ways within a city area | ⚠️ **Weak** — OSM has no native popularity rank; would need a proxy (e.g. presence of `wikidata`/`wikipedia` tags) and that proxy is coarse | ✅ No key/account | ✅ Free (ODbL) | ⚠️ Overpass QL is a custom query language; area resolution adds steps | ⚠️ Shared public instances throttle; heavier queries can be slow |
| 4 | **Bundled OSM/Wikidata extract (local dataset)** | Downloadable dataset (bundled) | ✅ Good — same underlying data as 1/3 | ⚠️ Needs a precomputed prominence column; data must be curated/joined offline | ✅ No auth at runtime | ✅ Free (ODbL / CC0) | ⚠️ Large extract + offline ETL pipeline to build a per-city attractions table | ✅ None at runtime (offline) |
| 5 (ref) | **Foursquare / Google Places** | Public HTTP API | ✅ Excellent | ✅ Excellent (ratings, popularity) | ❌ API key + billing/account | ⚠️ Quota/billing | ✅ SDKs | ❌ **Disqualified on no-auth + free** |
| 6 (ref) | **Dedicated Go attractions library** | Go module | n/a | n/a | n/a | n/a | ❌ No maintained, no-auth Go library wraps a free attractions source; community libraries are thin REST wrappers around keyed APIs (e.g. OpenTripMap) | n/a |

## Decision

- **Decision**: Use the **Wikidata Query Service (WDQS) SPARQL endpoint**
  (`https://query.wikidata.org/sparql`, JSON results) as the attraction data source. For a
  city resolved by name + country (FR-002a), query entities located within (or administratively
  inside) the city that are instances/subclasses of tourist-attraction-like types, **ordered by
  Wikipedia sitelink count descending**, and take the top 10 (FR-003, FR-004).
- **Rationale**:
  - **No-auth + free** (FR-010, FR-011, SC-004) — WDQS needs no API key or account; the only
    obligation is sending a descriptive `User-Agent` header (good-citizen policy), which is
    trivial with `net/http`.
  - **Built-in prominence signal** — the count of Wikipedia language sitelinks is a strong,
    globally consistent proxy for "most known," directly satisfying the "top 10 most known"
    ranking requirement (FR-003) without inventing a heuristic. This is the criterion that
    separates WDQS from Overpass/OSM.
  - **Category for free** — each result's `instance of` label supplies the optional
    category/short context the spec asks to show "when readily available" (FR-005).
  - **Accent/Unicode safe** (FR-002a) — Wikidata labels are full Unicode; resolving the city to
    a QID by name + country handles "São Paulo", "Zürich", etc. correctly.
  - **Simplest no-auth integration with ranking** (Principle V) — a single HTTPS GET returning
    SPARQL JSON, parsed with the standard library; no SDK, no key management, no offline ETL.
  - **Meets latency** — a single bounded query comfortably fits the 3-second target (SC-002)
    for typical cities.
- **Alternatives considered**:
  - *OpenTripMap* — arguably the **best data quality + ranking** (built-in `rate` score) and a
    strong production candidate, but it **requires a free API key**, conflicting with the no-auth
    preference (FR-010). Documented as the recommended upgrade **if** the team accepts the
    auth/key-management trade-off for richer popularity ranking and per-attraction detail.
  - *Overpass / OSM API* — fully no-auth and free with excellent POI coverage, but **lacks a
    native popularity rank**; "top 10 most known" would rely on a coarse proxy. Kept as the
    **no-auth fallback** if WDQS availability/timeouts prove problematic.
  - *Bundled local extract* — removes the runtime network dependency, but offline operation is
    explicitly **not required** (spec clarification) and it adds a heavy offline ETL pipeline to
    precompute prominence for no spike-level benefit. Rejected for the spike.
  - *Foursquare / Google Places* — rejected outright for violating no-auth + free.
  - *Dedicated Go library* — no maintained no-auth library exists; available wrappers target
    keyed APIs, so adopting one would reintroduce the auth constraint and an external dependency.

## Identified Risks & Limitations

- **WDQS rate limits / timeouts**: The public endpoint enforces a ~60 s query timeout and
  throttles heavy use. Mitigation — keep the SPARQL query tight (filter by type, limit to top
  results, order in-query), set a short client `context` deadline (SC-002), and send a
  descriptive `User-Agent`.
- **Sitelink count ≠ tourist popularity**: Sitelinks measure encyclopedic notability, which
  correlates with fame but is not a visitor-popularity metric. Acceptable proxy for a spike;
  note OpenTripMap's `rate` as the richer alternative if relevance is judged insufficient.
- **City→QID resolution ambiguity**: Same-named cities require disambiguating by country (and
  region/state when available) to pick the correct QID (FR-002a). Mitigation — constrain the
  city lookup by country; document the residual ambiguity risk.
- **Coverage for small/lesser-known cities**: Sparse Wikidata coverage may yield fewer than 10
  (or zero) attractions. This is expected and handled by FR-006/FR-007 (show all / clear
  "no attractions found" message), not an error.
- **Network dependency**: Unlike the 001 city dataset, attractions require connectivity; a
  source/connectivity failure is reported on stderr with a non-zero exit (FR-009).
- **Attribution**: Wikidata content is CC0 (no attribution required); OSM/Overpass would
  require ODbL attribution if adopted as the fallback.

## Proof of Concept (FR-014)

A brief PoC issues one WDQS SPARQL query for a sample well-known city (e.g. resolve "Paris,
France" → QID, then fetch attractions located within it ordered by sitelink count, limited to
10) and prints the names + category. A recorded response is saved under
`internal/attraction/testdata/` so the durable ranking/mapping logic is validated
deterministically without live network. This demonstrates the recommended approach returns a
plausible top-10 list (SC-007).

## Resolved Unknowns

All Technical Context items are resolved; **no `NEEDS CLARIFICATION` markers remain**. The
single open decision (attraction data source) is committed to **WDQS SPARQL**, with
**OpenTripMap** (if the auth trade-off is accepted) and **Overpass/OSM** (no-auth fallback) as
documented alternatives.
