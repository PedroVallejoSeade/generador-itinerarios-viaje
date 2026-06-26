# Phase 0 Research — City Search Spike Findings

**Spike time-box**: 10 minutes (per spec Clarifications + Constitution Principle III)
**Objective**: Evaluate at least three public city data sources against the project
constraints and recommend one (FR-012, FR-013, US3).

## Evaluation Criteria

Derived from the spec constraints:

1. **No authentication** — no API key, account, or registration (FR-010, SC-004).
2. **Free** — free for this application's use (FR-011).
3. **Fast** — results in under 2 seconds (SC-002).
4. **Disambiguation coverage** — provides at least city name, country, and region/state, plus a signal for population-based ranking (FR-004, FR-006).
5. **Simplicity** — least code/integration effort (Constitution Principle V).

## Candidates Evaluated

| # | Source | Type | No-Auth | Free | Speed | Disambiguation fields | Notes |
|---|--------|------|---------|------|-------|-----------------------|-------|
| 1 | **SimpleMaps World Cities (Basic)** | Downloadable CSV (bundled) | ✅ download needs no account | ✅ CC BY 4.0 | ✅ offline, in-memory | city, country, `admin_name` (region), population, lat/lng | `admin_name` is human-readable directly — no code-to-name join needed. |
| 2 | **GeoNames `cities15000` dump** | Downloadable TSV (bundled) | ✅ dump needs no account | ✅ CC BY 4.0 | ✅ offline, in-memory | name, country code, `admin1 code`, population, coords | Region is an **admin1 code**; needs a join against `admin1CodesASCII.txt` to show a readable region. Country is an ISO code → needs `countryInfo.txt` join for a readable name. More moving parts. |
| 3 | **OpenStreetMap Nominatim API** | Live HTTP API | ✅ no key | ✅ free | ⚠️ network round-trip; usage policy caps ~1 req/s, requires User-Agent | display name, address (country, state), importance | Online dependency violates offline goal; rate limits and usage policy make it unsuitable as the primary lookup for an interactive CLI. |
| 4 (ref) | **GeoNames live search API** | Live HTTP API | ❌ requires free `username` registration | ✅ free tier | ⚠️ network | rich | **Disqualified on no-auth**: every request needs a registered `username` parameter. |
| 5 (ref) | **OpenWeather / Mapbox / Google geocoding** | Live HTTP API | ❌ API key required | ⚠️ quota/billing | ⚠️ network | rich | **Disqualified on no-auth + free** constraints. |

## Decision

- **Decision**: Bundle the **SimpleMaps World Cities (Basic, CC BY 4.0)** dataset as an
  embedded CSV and perform all lookups in-memory.
- **Rationale**:
  - Satisfies **no-auth** and **free** absolutely — the dataset is downloaded once at build
    time and embedded; the running app makes no network calls and needs no key or account
    (directly supports SC-004).
  - **Fastest** option — an in-memory scan over ~25k rows resolves in milliseconds, trivially
    meeting SC-002's 2-second bound and removing any network-failure surface.
  - **Best disambiguation ergonomics** — `admin_name` and `country` are already
    human-readable, so results show "City, Region, Country" with no secondary lookup tables;
    `population` directly powers the population ranking (FR-006).
  - **Simplest** integration (Principle V) — one CSV, one parser, no join tables, no HTTP
    client, no rate-limit handling.
- **Alternatives considered**:
  - *GeoNames `cities15000` dump* — equally open and offline and a strong fallback, but
    requires joining admin1/country code tables to render readable region/country, adding code
    for no spike-level benefit. Kept as the **documented backup** if SimpleMaps licensing or
    coverage proves insufficient.
  - *Nominatim API* — rejected as primary because the online dependency, latency, and
    1 req/s usage policy conflict with the offline/fast goals and interactive use.
  - *GeoNames live API / commercial geocoders* — rejected outright for violating the no-auth
    (and, for commercial APIs, free) constraints.

## Identified Risks & Limitations

- **Attribution obligation**: CC BY 4.0 requires crediting SimpleMaps. Mitigation — include
  attribution in the README and `--version`/`--help` output.
- **Dataset staleness**: Bundled data is a point-in-time snapshot; new/renamed cities require
  re-embedding. Acceptable for a spike; for production, document a refresh procedure.
- **Coverage floor**: The Basic tier omits very small towns (population threshold). Acceptable
  for disambiguating well-known cities (SC-003); note as a known limitation.
- **Binary size**: Embedding the CSV grows the binary by a few MB. Negligible for a CLI.
- **Accented names (FR-013 edge case)**: Matching must normalize case but should compare on
  the original Unicode names; the dataset stores accented names (e.g., "São Paulo", "Zürich"),
  so accent-preserving prefix matching returns expected results.

## Resolved Unknowns

All Technical Context items are resolved; **no `NEEDS CLARIFICATION` markers remain**. The
single open decision (data source) is now committed to the bundled SimpleMaps dataset, with
GeoNames as a documented fallback.

## Implementation Note (post-spike)

During implementation the SimpleMaps Basic download endpoint required interactive access
(returned HTTP 403 to scripted download), so the **documented GeoNames `cities15000` fallback
(CC BY 4.0)** was bundled instead. The dump was converted into the same
`city,country,admin_name,population` shape by joining `admin1CodesASCII.txt` (readable region)
and `countryInfo.txt` (readable country), giving ~34k cities embedded via `//go:embed`. This
validates the research conclusion: the fallback satisfies every constraint (no-auth, free,
offline/fast, full disambiguation fields) with only a one-time build-step join. Attribution is
credited to GeoNames in the README and `--help` output per the CC BY 4.0 obligation.
