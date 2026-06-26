# Quickstart — City Search Spike

A runnable validation guide proving the city search works end-to-end. Implementation details
live in [data-model.md](data-model.md) and [contracts/cli.md](contracts/cli.md); this file is
the build/run/validate path.

## Prerequisites

- Go 1.26+ (`go version`)
- Repository checked out at branch `001-city-search-spike`
- The bundled dataset present at `internal/city/cities.csv` (sourced from SimpleMaps World
  Cities Basic, CC BY 4.0 — see [research.md](research.md)). A small fixture for tests lives at
  `internal/city/testdata/cities_sample.csv`.

## Build

```bash
go build ./...
go build -o bin/citysearch ./cmd/citysearch
```

## Run the test suite (TDD durable logic)

Tests are table-driven and written before/with the matching implementation
(`internal/city/search_test.go`):

```bash
go test ./...
```

Expected: all tests pass, covering case-insensitive prefix matching, population ranking, the
10-result cap, empty-query rejection, and the no-results path.

## Validate against the CLI contract

Run each scenario from [contracts/cli.md](contracts/cli.md) and confirm the outcome:

| Step | Command | Expected | Maps to |
|------|---------|----------|---------|
| 1 | `./bin/citysearch paris` | A line like `Paris, Île-de-France, France` on stdout, exit 0 | US1 / FR-003, FR-005 |
| 2 | `./bin/citysearch springfield` | Multiple lines, each distinguishable by region + country | US2 / FR-004 |
| 3 | `./bin/citysearch "São Paulo"` | Accented match returned | Edge case (non-ASCII) |
| 4 | `./bin/citysearch zzzzzz` | `No cities found matching "zzzzzz".`, exit 0 | FR-007 |
| 5 | `./bin/citysearch ""` | Prompt for valid input on stderr, exit non-zero | FR-008 |
| 6 | `./bin/citysearch london` | ≤ 10 lines, population-ordered (largest first) | FR-006 |

Check exit codes with `echo $?` after each run.

## Success Criteria Validation

| Criterion | How to verify |
|-----------|---------------|
| SC-001 (distinguish ≥ 5 same-name cities) | Run step 2 with a highly shared name; confirm ≥ 5 distinguishable rows. |
| SC-002 (< 2 s) | `time ./bin/citysearch paris` — well under 2 s (in-memory). |
| SC-003 (intended city present ≥ 95%) | Spot-check a set of well-known cities; intended city appears within the 10 results. |
| SC-004 (no account/key/auth) | Run with no network and no configuration — results still returned (offline embedded dataset). |
| SC-005 / SC-006 (findings document) | Review [research.md](research.md): ≥ 3 sources compared, one recommended with rationale and risks. |

## Spike Wrap-Up

Per Constitution Principle III and the spec's spike-disposability assumption: the durable
deliverables are this working demonstration and the data-source recommendation in
[research.md](research.md). After the 10-minute evaluation box, return to the Test-First
workflow for any production hardening.
