# Implementation Plan: City Search Spike

**Branch**: `001-city-search-spike` | **Date**: 2026-06-26 | **Spec**: [spec.md](spec.md)

**Input**: Feature specification from `/specs/001-city-search-spike/spec.md`

**Note**: This template is filled in by the `/speckit.plan` command. See `.specify/templates/plan-template.md` for the execution workflow.

## Summary

Deliver a Go terminal application that takes a city name and returns matching cities with
disambiguating context (country + region/state), ranked by population and capped at 10
results. The primary unknown — which public city data source to use — is resolved by a
time-boxed (10-minute) spike that evaluates at least three candidate sources against the
project constraints (no authentication, free, fast, sufficient disambiguation data). The
research concludes that a **bundled open city dataset** (SimpleMaps World Cities, Basic /
CC BY 4.0) best satisfies all constraints because it is offline (sub-second, no network
dependency), needs no API key or account, and provides city, country, admin region, and
population directly. See [research.md](research.md) for the full comparison and rationale.

## Technical Context

**Language/Version**: Go 1.26 (latest stable, per constitution Technology Standards)

**Primary Dependencies**: Go standard library only — `flag`, `encoding/csv`, `strings`, `sort`, `embed`. No external modules (Simplicity First).

**Storage**: Bundled static dataset embedded into the binary via `//go:embed` (a single CSV of world cities). No database, no runtime file I/O for the dataset.

**Testing**: Go standard `testing` package, table-driven tests. Test fixture is a small CSV under `internal/city/testdata/`.

**Target Platform**: Cross-platform terminal/CLI (Linux primary; pure-Go, no cgo, so macOS/Windows build cleanly).

**Project Type**: Single-project CLI application.

**Performance Goals**: Results displayed in under 2 seconds for a typical single-name query (SC-002). With an in-memory embedded dataset and linear prefix scan, lookups complete in single-digit milliseconds.

**Constraints**: No authentication / API key / account (FR-010, SC-004); free to use (FR-011); offline-capable; bounded output (max 10 results, FR-006); case-insensitive prefix matching (FR-003).

**Scale/Scope**: Spike scope only — city lookup + disambiguated display + data-source recommendation. Dataset size on the order of 10^4–10^5 cities (cities with population ≥ 15,000 ≈ 25k rows), comfortably held in memory.

## Constitution Check

*GATE: Must pass before Phase 0 research. Re-check after Phase 1 design.*

| Principle | Status | Notes |
|-----------|--------|-------|
| I. Milestone-Driven Development | PASS | Single P1 deliverable slice (city search + disambiguation + recommendation); independently testable and demonstrable. |
| II. Test-First (TDD — NON-NEGOTIABLE) | PASS (documented spike exception) | Per Principle III, spike code is exploratory/disposable and exempt from strict TDD during the time-box. The durable matching logic is specified with table-driven tests first (see [quickstart.md](quickstart.md)); after the spike the team returns to the Red-Green-Refactor workflow. |
| III. Spike-First Exploration | PASS | This feature **is** the explicitly documented, time-boxed (10-minute) spike to evaluate unproven city-data sourcing. Output is [research.md](research.md) with a recommendation and risks. |
| IV. CLI Interface | PASS | Results to stdout, errors to stderr, standard exit codes (0 success, non-zero on error/no-source), `-h/--help` usage text. |
| V. Simplicity First (YAGNI) | PASS | Standard library only; single embedded CSV; linear prefix scan; no premature indexing, caching, or abstractions. |

**Result**: No violations. Complexity Tracking section intentionally left empty.

## Project Structure

### Documentation (this feature)

```text
specs/001-city-search-spike/
├── plan.md              # This file (/speckit.plan command output)
├── research.md          # Phase 0 output — spike findings & data-source recommendation
├── data-model.md        # Phase 1 output — City, SearchQuery, SearchResultSet
├── quickstart.md        # Phase 1 output — build/run/validate guide
├── contracts/
│   └── cli.md           # Phase 1 output — CLI input/output contract
├── checklists/
│   └── requirements.md  # Spec quality checklist (already present)
└── tasks.md             # Phase 2 output (/speckit.tasks — NOT created here)
```

### Source Code (repository root)

```text
cmd/
└── citysearch/
    └── main.go          # CLI entry point: arg parsing, wiring, stdout/stderr, exit codes

internal/
└── city/
    ├── city.go          # City type + dataset loading (embed) + parsing
    ├── search.go        # Normalize query + case-insensitive prefix match + population rank + cap
    ├── search_test.go   # Table-driven unit tests (matching, ranking, cap, edge cases)
    ├── cities.csv       # Embedded bundled dataset (//go:embed)
    └── testdata/
        └── cities_sample.csv  # Small deterministic fixture for tests
```

**Structure Decision**: Single Go module (`github.com/PedroVallejoSeade/generador-itinerarios-viaje`)
following the constitution's Technology Standards: `cmd/` for the CLI entry point and
`internal/` for application logic that should not be imported by external modules. No `pkg/`
is introduced (YAGNI). The dataset lives next to the package that owns it and is embedded at
build time so the binary is self-contained and offline-capable.

## Complexity Tracking

> No constitution violations to justify — section intentionally empty.
