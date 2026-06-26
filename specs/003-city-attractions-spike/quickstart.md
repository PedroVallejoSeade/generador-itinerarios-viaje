# Quickstart — City Attractions Lookup Spike

A runnable validation guide proving the attractions lookup works end-to-end. Implementation
details live in [data-model.md](data-model.md) and [contracts/cli.md](contracts/cli.md); this
file is the build/run/validate path. The data-source rationale is in [research.md](research.md)
(recommended: Wikidata Query Service SPARQL).

## Prerequisites

- Go 1.26+ (`go version`)
- Repository checked out at branch `003-city-attractions-spike`
- Network access for live attraction lookups (Wikidata Query Service — no API key/account
  needed; see [research.md](research.md)). Tests run **offline** against recorded JSON fixtures
  under `internal/attraction/testdata/`.

## Build

```bash
go build ./...
go build -o bin/citysearch ./cmd/citysearch
```

## Run the test suite (durable logic)

Tests are table-driven over recorded WDQS responses (no live network) and cover the durable
ranking/mapping logic:

```bash
go test ./...
```

Expected: all tests pass, covering prominence ranking (sitelink count, descending), the
10-result cap, the < 10 "show all" path, the empty "no attractions found" path, JSON→Attraction
mapping (name, category, description), accented city resolution, and the fetch-error path.

## Validate against the CLI contract

Start interactive mode, search for a city, then select it by number and confirm attractions
appear. Each step maps to [contracts/cli.md](contracts/cli.md):

| Step | Action | Expected | Maps to |
|------|--------|----------|---------|
| 1 | `./bin/citysearch`, search `paris`, select Paris by its number | Numbered list with recognizable landmarks (Eiffel Tower, Louvre, …), ≤ 10 rows | A1 / US1, FR-001..FR-005 |
| 2 | Select a city with many attractions | Exactly 10 rows, most-known first | A2 / FR-003, FR-004 |
| 3 | Select a small city with few attractions | All available rows, no error | A3 / FR-007 |
| 4 | Select a city with no attraction data | "No attractions found for …" | A4 / FR-006 |
| 5 | Enter a number not in the list | Invalid-selection message, no lookup, session continues | A5 / FR-008 |
| 6 | Search `São Paulo`, select it | Resolves correctly, returns attractions | A6 / FR-002a |
| 7 | Run with network disabled, then select a city | `error: unable to fetch attractions: …` on stderr, exit 2 | A7 / FR-009 |

Check exit codes with `echo $?` where relevant.

## Proof of Concept (FR-014)

The recommended approach is demonstrated by selecting a well-known sample city (e.g. Paris) and
confirming a plausible top-10 attractions list is returned from Wikidata — establishing
feasibility in Go (SC-007). A recorded response of this query is the test fixture that keeps the
PoC reproducible offline.

Run the PoC deterministically (no network) in under a minute:

```bash
go test ./internal/attraction/ -run TestLookup_FromGoldenFixture -v
```

This drives the full `fetch → decode → map → rank → cap` path (`attraction.Lookup`) against the
recorded `internal/attraction/testdata/paris.json` response, asserting a ranked list led by the
most-known landmark. The same logic powers the live CLI path (step 1 above) when network access
is available, confirming the recommendation in [research.md](research.md) is feasible.

## Success Criteria Validation

| Criterion | How to verify |
|-----------|---------------|
| SC-001 (≥ 5 recognizable attractions for well-known cities) | Run step 1–2 across a few major cities; confirm ≥ 5 recognizable rows each. |
| SC-002 (< 3 s display) | `time` a selection for a typical city; well under 3 s. |
| SC-003 (identify ≥ 3 must-sees) | Inspect the list for a well-known city; ≥ 3 obvious must-see attractions present. |
| SC-004 (no account/key/auth) | Run without any API key or account configured — lookups still succeed (WDQS is no-auth). |
| SC-005 / SC-006 (findings document) | Review [research.md](research.md): ≥ 3 candidates compared, one recommended with rationale and risks, readable in < 5 min. |
| SC-007 (PoC feasible in Go) | Step 1 returns a plausible attractions list for the sample city. |

## Spike Wrap-Up

Per Constitution Principle III and the spec's spike-disposability assumption: the durable
deliverables are this working demonstration and the data-source recommendation in
[research.md](research.md). The exploratory client code is disposable; after the time-box,
return to the Test-First workflow for any production hardening (e.g. caching, OpenTripMap
upgrade if richer ranking is needed).
