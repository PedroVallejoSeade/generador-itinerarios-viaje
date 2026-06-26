---
description: "Task list for City Search Spike implementation"
---

# Tasks: City Search Spike

**Input**: Design documents from `/specs/001-city-search-spike/`

**Prerequisites**: [plan.md](plan.md) (required), [spec.md](spec.md) (required), [research.md](research.md), [data-model.md](data-model.md), [contracts/cli.md](contracts/cli.md), [quickstart.md](quickstart.md)

**Tests**: INCLUDED — the constitution mandates Test-First (TDD, non-negotiable) and the quickstart specifies table-driven tests written first. Test tasks precede their implementation within each story.

**Organization**: Tasks are grouped by user story (all three are P1) to enable independent implementation and testing.

## Format: `[ID] [P?] [Story] Description`

- **[P]**: Can run in parallel (different files, no dependencies on incomplete tasks)
- **[Story]**: Which user story the task belongs to (US1, US2, US3)
- All paths are repository-root relative (single Go project per plan.md)

## Path Conventions

- CLI entry: `cmd/citysearch/`
- Application logic: `internal/city/`
- Bundled dataset + test fixture: `internal/city/cities.csv`, `internal/city/testdata/`

---

## Phase 1: Setup (Shared Infrastructure)

**Purpose**: Project skeleton and tooling

- [X] T001 Create the Go project directory structure (`cmd/citysearch/`, `internal/city/`, `internal/city/testdata/`) per [plan.md](plan.md)
- [X] T002 Verify `go.mod` targets Go 1.26 with module `github.com/PedroVallejoSeade/generador-itinerarios-viaje`; run `go build ./...` to confirm an empty buildable module
- [X] T003 [P] Add a `Makefile` (or document commands) for `gofmt`, `go vet`, and `go test ./...` at repository root

**Checkpoint**: Empty project compiles and tooling commands are available.

---

## Phase 2: Foundational (Blocking Prerequisites)

**Purpose**: Dataset and the `City` type/loader that ALL user stories depend on

**⚠️ CRITICAL**: No user story work can begin until this phase is complete

- [X] T004 Acquire the SimpleMaps World Cities (Basic, CC BY 4.0) dataset and place the CSV at `internal/city/cities.csv` (per [research.md](research.md))
- [X] T005 [P] Create a small deterministic test fixture at `internal/city/testdata/cities_sample.csv` containing entries for Paris (FR), several Springfields (US, different regions), São Paulo (BR), and London (multiple) to exercise matching, disambiguation, ranking, and accents
- [X] T006 Define the `City` struct (`Name`, `Country`, `Region`, `Population`) in `internal/city/city.go` per [data-model.md](data-model.md)
- [X] T007 Implement dataset loading via `//go:embed` and CSV parsing in `internal/city/city.go`, skipping rows with empty names and defaulting blank/non-numeric population to 0; return a clear error on parse failure (FR-009)

**Checkpoint**: `City` records load in-memory from the embedded dataset; foundation ready for all stories.

---

## Phase 3: User Story 1 - Find a city by name (Priority: P1) 🎯 MVP

**Goal**: A user types a city name in the terminal and sees matching cities listed on stdout.

**Independent Test**: Run `citysearch paris` and confirm a Paris entry appears in the output.

### Tests for User Story 1 ⚠️ (write first, must FAIL before implementation)

- [X] T008 [P] [US1] Table-driven test for case-insensitive prefix matching (e.g., "par" → Paris; "PARIS" → Paris) against the fixture in `internal/city/search_test.go`
- [X] T009 [P] [US1] Table-driven test for empty/whitespace-only query rejection and for the no-results case in `internal/city/search_test.go`

### Implementation for User Story 1

- [X] T010 [US1] Implement `SearchQuery` normalization (trim + lowercase) and empty-query rejection in `internal/city/search.go` per [data-model.md](data-model.md)
- [X] T011 [US1] Implement case-insensitive prefix matching returning matched `City` records in `internal/city/search.go` (FR-003)
- [X] T012 [US1] Implement the CLI entry point in `cmd/citysearch/main.go`: parse the positional argument, load data, run the search, print matches to stdout, and `-h/--help` usage text (FR-001, FR-005; [contracts/cli.md](contracts/cli.md))
- [X] T013 [US1] Wire exit codes and messages in `cmd/citysearch/main.go`: prompt + exit 1 on empty query (FR-008), "No cities found…" + exit 0 on no matches (FR-007), data-load error + exit 2 on failure (FR-009)

**Checkpoint**: `citysearch <name>` returns matching cities — independently demonstrable MVP.

---

## Phase 4: User Story 2 - Disambiguate between similarly named cities (Priority: P1)

**Goal**: Each result carries country + region/state, ranked by population and capped at 10, so similarly named cities are distinguishable.

**Independent Test**: Run `citysearch springfield` and confirm multiple results, each distinguishable by region and country.

### Tests for User Story 2 ⚠️ (write first, must FAIL before implementation)

- [X] T014 [P] [US2] Table-driven test asserting each result exposes country and region, and that two same-name+same-country results differ by region, in `internal/city/search_test.go`
- [X] T015 [P] [US2] Table-driven test asserting population-descending ordering (name as tiebreaker) and the 10-result cap in `internal/city/search_test.go`

### Implementation for User Story 2

- [X] T016 [US2] Implement population-descending ranking with name tiebreaker and the max-10 truncation in `internal/city/search.go` (FR-006)
- [X] T017 [US2] Implement result formatting `"<Name>, <Region>, <Country>"` (omit region when empty) for stdout in `cmd/citysearch/main.go` (FR-004; [contracts/cli.md](contracts/cli.md))

**Checkpoint**: Same-name searches return ranked, capped, context-rich results — US1 and US2 both work.

---

## Phase 5: User Story 3 - Evaluate and recommend a city data source (Priority: P1)

**Goal**: A documented evaluation comparing ≥3 candidate sources against the constraints, naming one recommendation with rationale and risks.

**Independent Test**: Review [research.md](research.md) and confirm ≥3 sources compared and a single recommendation with rationale.

### Implementation for User Story 3

- [X] T018 [US3] Verify [research.md](research.md) satisfies FR-012/FR-013: ≥3 candidate sources evaluated against no-auth, free, response speed, and disambiguation coverage, with one recommendation, rationale, and identified risks; fill any gaps
- [X] T019 [US3] Add the required CC BY 4.0 attribution for the chosen dataset to `README.md` and to the `--help`/version output in `cmd/citysearch/main.go`

**Checkpoint**: Findings document is complete and the recommendation is reviewable in under 5 minutes (SC-006).

---

## Phase 6: Polish & Cross-Cutting Concerns

**Purpose**: Validation and finishing touches across stories

- [X] T020 [P] Update `README.md` with build/run usage examples and the dataset attribution
- [X] T021 Run the [quickstart.md](quickstart.md) validation scenarios (steps 1–6) and confirm outputs and exit codes
- [X] T022 [P] Run `gofmt -l .`, `go vet ./...`, and `go test ./...` and resolve any findings
- [X] T023 Verify SC-002 (< 2 s) with `time ./bin/citysearch paris` and SC-004 (works offline, no auth) by running with no network

---

## Dependencies & Execution Order

### Phase Dependencies

- **Setup (Phase 1)**: No dependencies — start immediately
- **Foundational (Phase 2)**: Depends on Setup — BLOCKS all user stories
- **User Stories (Phase 3–5)**: All depend on Foundational completion
  - US3 (research) is independent of US1/US2 code and may proceed in parallel once Setup is done
  - US1 → US2 share `search.go`/`main.go`, so run sequentially (US1 before US2)
- **Polish (Phase 6)**: Depends on the targeted user stories being complete

### User Story Dependencies

- **US1 (P1)**: After Foundational — no dependency on other stories
- **US2 (P1)**: After US1 (extends `search.go` ranking and `main.go` formatting); independently testable
- **US3 (P1)**: Independent — only depends on the data-source decision already captured in research.md

### Within Each User Story

- Tests written and failing before implementation (TDD)
- `City` model/loader (Phase 2) before search logic
- Search logic before CLI wiring
- Story complete before moving to the next

### Parallel Opportunities

- T003 (tooling) parallel with other Setup work
- T005 (fixture) parallel with T006 (model) — different files
- T008 and T009 (US1 tests) parallel — independent test cases
- T014 and T015 (US2 tests) parallel — independent test cases
- US3 (T018–T019) can run alongside US1/US2 if staffed separately
- Polish T020 and T022 parallel — different files

---

## Parallel Example: User Story 1

```bash
# Write both US1 test groups together (they fail first):
Task: "Prefix-matching table test in internal/city/search_test.go"
Task: "Empty-query + no-results table test in internal/city/search_test.go"
```

---

## Implementation Strategy

### MVP First (User Story 1 Only)

1. Complete Phase 1: Setup
2. Complete Phase 2: Foundational (CRITICAL — blocks all stories)
3. Complete Phase 3: User Story 1
4. **STOP and VALIDATE**: Run `citysearch paris` and confirm a result appears
5. Demo the standalone city-lookup tool

### Incremental Delivery

1. Setup + Foundational → foundation ready
2. US1 → test independently → demo (MVP city lookup)
3. US2 → test independently → demo (disambiguated, ranked results)
4. US3 → finalize findings document → recommendation delivered
5. Polish → quickstart validation + offline/perf checks

### Spike Note

Per Constitution Principle III and the spec's spike-disposability assumption, code here is
exploratory within a 10-minute evaluation box; the durable deliverables are the working city
search demonstration and the data-source recommendation in [research.md](research.md). After
the box, return to the Test-First workflow for any production hardening.

---

## Notes

- [P] tasks = different files, no dependencies
- [Story] label maps each task to its user story for traceability
- Verify tests fail before implementing (TDD, non-negotiable)
- Commit after each task or logical group
