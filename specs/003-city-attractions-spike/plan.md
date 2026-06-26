# Implementation Plan: City Attractions Lookup Spike

**Branch**: `003-city-attractions-spike` | **Date**: 2026-06-26 | **Spec**: [spec.md](spec.md)

**Input**: Feature specification from `/specs/003-city-attractions-spike/spec.md`

**Note**: This template is filled in by the `/speckit.plan` command. See `.specify/templates/plan-template.md` for the execution workflow.

## Summary

After the traveler selects a city by its number from the existing search results, the
application looks up and displays the **top 10 most well-known attractions** for that city —
landmarks, museums, and sightseeing spots — as a numbered list (name, plus category/short
description when readily available). The primary unknown — *which attraction data source to
use* — is resolved by a time-boxed spike that evaluates at least three candidates (a free
no-auth public API, a Go library, and a locally queried open dataset) against data quality,
ease of Go integration, authentication requirements, rate limits, and result relevance
(FR-012, US2).

The research concludes that the **Wikidata Query Service (WDQS) SPARQL endpoint** best
satisfies the constraints: it is **no-auth and free**, returns attraction **name + type
(category)** for a city resolved by name + country, and — critically — provides a built-in
**prominence signal** (the count of Wikipedia language sitelinks) that directly powers the
"top 10 most known" ranking (FR-003, FR-004). Network access is acceptable per the spec
clarification. **OpenTripMap** (richer popularity `rate`, but requires a free API key → auth
trade-off) and the no-auth **Overpass / OpenStreetMap API** (good POI coverage, weak ranking
signal) are the documented alternatives. See [research.md](research.md) for the full
comparison, rationale, and risks.

This is a spike (Constitution Principle III): the data-source evaluation is time-boxed and the
exploratory client code is disposable. The durable deliverables are [research.md](research.md)
(the recommendation) and a brief proof of concept that fetches attractions for at least one
sample city (FR-014). The existing `internal/city` selection data is reused unchanged; the
spike adds a new `internal/attraction` package (data source client + prominence ranking + cap)
and a selection→attractions step in the interactive loop, both written behind an injectable
HTTP-fetch seam so behaviour is exercised with stubbed responses and recorded JSON fixtures.

## Technical Context

**Language/Version**: Go 1.26 (latest stable, per constitution Technology Standards and `go.mod`).

**Primary Dependencies**: Go standard library only — `net/http` (WDQS request), `net/url` (query/User-Agent), `encoding/json` (SPARQL JSON results), `context` (per-lookup timeout), plus existing `fmt`/`io`/`strings`/`sort`. No external modules (Simplicity First); the SPARQL query is sent as a plain HTTPS GET.

**Storage**: No new persistent storage. The bundled world-cities dataset (`internal/city`, embedded via `//go:embed`) is reused to resolve the selected city. Attractions are fetched live over the network per lookup; recorded JSON responses are saved under `internal/attraction/testdata/` as deterministic test fixtures (no live network in tests).

**Testing**: Go standard `testing` package, table-driven tests. The data source is accessed through a small fetcher interface (e.g. `func(ctx, query) ([]byte, error)`), so ranking, capping, mapping, the empty/short-list paths, and error handling are tested against stubbed/golden JSON with no real HTTP calls.

**Target Platform**: Cross-platform terminal/CLI (Linux primary; pure-Go, no cgo → macOS/Windows build cleanly). Selection over stdin, attraction list to stdout, errors to stderr.

**Project Type**: Single-project CLI application.

**Performance Goals**: The attractions list is displayed in under 3 seconds for a typical lookup (SC-002). A single WDQS SPARQL query per selection, bounded by a `context` deadline; results are rendered immediately on return.

**Constraints**: Recommended source SHOULD need no auth/API key/account (FR-010, SC-004) and MUST be free (FR-011); network connectivity is acceptable (spec clarification). Resolve the city by name + country (plus region/state when available), correctly handling accented/non-Latin names (FR-002a). Rank by a prominence signal and cap the display at 10 (FR-003, FR-004); show all when fewer than 10 (FR-007); clear "no attractions found" message (FR-006); data-source/connectivity errors reported on stderr with a non-zero exit (FR-009).

**Scale/Scope**: Spike scope only — resolve the selected city, fetch + rank + cap its attractions, display them, and document the data-source recommendation with a proof of concept. One network round-trip per selection; result set bounded to 10 rows.

## Constitution Check

*GATE: Must pass before Phase 0 research. Re-check after Phase 1 design.*

| Principle | Status | Notes |
|-----------|--------|-------|
| I. Milestone-Driven Development | PASS | Two prioritized, independently testable P1 slices map to the spec's user stories: US1 = display top-10 attractions for a selected city; US2 = the data-source evaluation + recommendation + proof of concept. Each delivers standalone value and is tested on its own. |
| II. Test-First (TDD — NON-NEGOTIABLE) | PASS (documented spike exception) | Per Principle III, the time-boxed source-evaluation/PoC code is exploratory and exempt from strict TDD during the box. The durable logic (prominence ranking, 10-cap, JSON→Attraction mapping, empty/short/error paths) is specified with table-driven tests against golden JSON written first; after the spike the team returns to Red-Green-Refactor. |
| III. Spike-First Exploration | PASS | This feature **is** the explicitly documented, time-boxed spike to evaluate the unproven attraction-data sourcing. Output is [research.md](research.md) (≥ 3 candidates, one recommendation, risks) plus a working proof of concept (FR-012..FR-014). |
| IV. CLI Interface | PASS | City selection read from stdin; numbered attraction list to stdout; data-source/connectivity errors to stderr with a non-zero exit; clean lookups exit 0. Consistent with the existing interactive CLI. |
| V. Simplicity First (YAGNI) | PASS | Standard-library `net/http` + `encoding/json` only — no SDK, no GraphQL/SPARQL client library, no caching/persistence layer. A single SPARQL query per selection; reuses existing city resolution; adds only a fetch-rank-cap-render path. No-auth source preferred precisely to avoid key-management complexity. |

**Result**: No violations. Complexity Tracking section intentionally left empty.

## Project Structure

### Documentation (this feature)

```text
specs/003-city-attractions-spike/
├── plan.md              # This file (/speckit.plan command output)
├── research.md          # Phase 0 output — attraction data-source evaluation & recommendation
├── data-model.md        # Phase 1 output — CitySelection, Attraction, AttractionResultSet
├── quickstart.md        # Phase 1 output — build/run/validate guide
├── contracts/
│   └── cli.md           # Phase 1 output — selection→attractions CLI input/output contract
├── checklists/
│   └── requirements.md  # Spec quality checklist (already present)
└── tasks.md             # Phase 2 output (/speckit.tasks — NOT created here)
```

### Source Code (repository root)

```text
cmd/
└── citysearch/
    ├── main.go             # UNCHANGED behavior — mode selection, dataset load, exit codes.
    ├── interactive.go      # EXTENDED — after listing numbered matches, prompt for a number;
    │                       #   on a valid selection, fetch + render the city's attractions.
    ├── interactive_test.go # EXTENDED — scripted selection → attractions rendering, invalid
    │                       #   selection, no-attractions, and fetch-error paths.
    └── main_test.go        # UNCHANGED (one-shot + mode-selection tests).

internal/
├── city/                   # UNCHANGED — City type + embedded dataset + Search (reused for selection).
└── attraction/             # NEW (spike) — attraction lookup
    ├── attraction.go       # Attraction type + JSON→Attraction mapping + prominence ranking + 10-cap.
    ├── wikidata.go         # WDQS SPARQL request builder + HTTPS fetch (behind a fetcher seam) + parse.
    ├── attraction_test.go  # Table-driven tests over golden JSON: ranking, cap, short/empty, errors.
    └── testdata/
        └── *.json          # Recorded WDQS SPARQL responses (deterministic fixtures, no live network).
```

**Structure Decision**: Single Go module
(`github.com/PedroVallejoSeade/generador-itinerarios-viaje`), continuing the 001/002 layout.
`cmd/citysearch` owns the CLI wiring and the interactive selection step (presentation), while a
new `internal/attraction` package owns the data-source client, JSON mapping, prominence
ranking, and the 10-result cap — mirroring how `internal/city` owns search. The attraction
package depends on a small fetcher seam (an `http.Client` call wrapped behind a function/
interface) rather than calling the network directly, so all behaviour is exercised by
package tests using recorded JSON fixtures — no live network and no new external dependency
(YAGNI). City resolution (name + country, FR-002a) reuses the existing `internal/city` records;
no change to that package is required for the spike.

## Complexity Tracking

> No constitution violations to justify — section intentionally empty.
