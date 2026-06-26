---
description: "Task list for City Attractions Lookup Spike"
---

# Tasks: City Attractions Lookup Spike

**Input**: Design documents from `/specs/003-city-attractions-spike/`

**Prerequisites**: plan.md (required), spec.md (required), research.md, data-model.md, contracts/cli.md, quickstart.md

**Tests**: Test tasks ARE included. The spec (FR-003..FR-009, FR-014) and plan ("durable logic ... specified with table-driven tests against golden JSON written first") explicitly call for test-first coverage of the ranking, cap, mapping, empty/short, and error paths. The time-boxed exploratory client code is exempt per Constitution Principle III, but the durable logic is tested first.

**Organization**: Tasks are grouped by user story (both P1) to enable independent implementation and testing.

## Format: `[ID] [P?] [Story] Description`

- **[P]**: Can run in parallel (different files, no dependencies)
- **[Story]**: Which user story this task belongs to (US1, US2)
- Include exact file paths in descriptions

## Path Conventions

Single Go module `github.com/PedroVallejoSeade/generador-itinerarios-viaje` (per plan.md). New
code lives under `internal/attraction/`; CLI wiring extends `cmd/citysearch/`. The existing
`internal/city` package is reused unchanged.

---

## Phase 1: Setup (Shared Infrastructure)

**Purpose**: Create the new package skeleton and test-fixture location used by both stories.

- [x] T001 Create the `internal/attraction/` package directory with a package doc comment in internal/attraction/attraction.go (package `attraction`)
- [x] T002 [P] Create the test-fixtures directory internal/attraction/testdata/ with a placeholder .gitkeep so recorded WDQS JSON responses have a home
- [x] T003 [P] Verify the toolchain builds the empty package: run `go build ./...` and `go vet ./...` from the repo root

---

## Phase 2: Foundational (Blocking Prerequisites)

**Purpose**: Core types and the fetcher seam that BOTH user stories depend on. No story work can begin until these exist.

**⚠️ CRITICAL**: No user story implementation can begin until this phase is complete.

- [x] T004 [P] Define the `Attraction` value type (`Name`, `Category`, `Description`, `Prominence int`) per data-model.md in internal/attraction/attraction.go
- [x] T005 [P] Define the `CitySelection` type (`Index int`, `Name`, `Country`, `Region`) and a constructor that maps a chosen `city.City` + 1-based index, per data-model.md, in internal/attraction/selection.go
- [x] T006 Define the fetcher seam — a `Fetcher` function/interface type `func(ctx context.Context, query string) ([]byte, error)` — and a default `net/http` HTTPS implementation (with User-Agent) in internal/attraction/wikidata.go

**Checkpoint**: Types and fetch seam exist — both user stories can now proceed.

---

## Phase 3: User Story 1 - See top attractions for a selected city (Priority: P1) 🎯 MVP

**Goal**: After the user selects a city by its number, fetch, rank, cap, and display that city's top-10 most-known attractions, handling the empty, short, invalid-selection, and source-error paths.

**Independent Test**: In interactive mode, search a well-known city (e.g. Paris), select it by number, and confirm a numbered list of ≤ 10 recognizable landmarks appears (Eiffel Tower, Louvre, …); verify the no-attractions, invalid-selection, and source-error paths behave per contracts/cli.md.

### Tests for User Story 1 (write first, ensure they FAIL) ⚠️

- [x] T007 [P] [US1] Record a real WDQS SPARQL JSON response for a well-known city (e.g. Paris) into internal/attraction/testdata/paris.json as the golden ranking/mapping fixture
- [x] T008 [P] [US1] Add a >10-result fixture internal/attraction/testdata/many.json and an empty-result fixture internal/attraction/testdata/empty.json for cap and no-attractions tests
- [x] T009 [P] [US1] Table-driven test for JSON→`Attraction` mapping (name, category, description, prominence) over golden fixtures in internal/attraction/attraction_test.go
- [x] T010 [P] [US1] Table-driven test for prominence ranking (descending sitelink count, ascending name tiebreaker) and the 10-result cap in internal/attraction/attraction_test.go
- [x] T011 [P] [US1] Test the <10 "show all" path and the empty "no attractions found" path in internal/attraction/attraction_test.go
- [x] T012 [P] [US1] Test the fetch-error path (fetcher returns an error → wrapped error, not an empty set) in internal/attraction/wikidata_test.go
- [x] T013 [P] [US1] Scripted CLI test for selection→attractions rendering, invalid selection (no lookup), no-attractions message, and fetch-error path in cmd/citysearch/interactive_test.go

### Implementation for User Story 1

- [x] T014 [US1] Implement the SPARQL query builder (city resolved by name + country, accent/Unicode-safe per FR-002a) and JSON-results parsing in internal/attraction/wikidata.go (depends on T006)
- [x] T015 [US1] Implement JSON→`Attraction` mapping (skip empty-name bindings, default `Prominence` to 0) in internal/attraction/attraction.go (depends on T004)
- [x] T016 [US1] Implement prominence ranking + 10-cap (`func Rank([]Attraction) []Attraction` / result-set builder) per data-model.md ordering rules in internal/attraction/attraction.go (depends on T015)
- [x] T017 [US1] Implement the top-level lookup `func Lookup(ctx, fetcher, CitySelection) ([]Attraction, error)` tying fetch → parse → map → rank → cap in internal/attraction/attraction.go (depends on T014, T016)
- [x] T018 [US1] Extend `runInteractive` to prompt for a city number after listing matches, validate the selection (reject out-of-range without a lookup, FR-008), and on a valid pick call `attraction.Lookup` in cmd/citysearch/interactive.go (depends on T017)
- [x] T019 [US1] Render the numbered attractions list to stdout (name + category/description when present, FR-005), the "No attractions found for …" message (FR-006), and report source/connectivity errors on stderr with non-zero exit (FR-009) in cmd/citysearch/interactive.go (depends on T018)

**Checkpoint**: User Story 1 is fully functional — selecting a city displays its ranked top-10 attractions, with all edge cases handled and tests green.

---

## Phase 4: User Story 2 - Evaluate and recommend an attractions data source (Priority: P1)

**Goal**: Deliver the documented data-source evaluation (≥ 3 candidates against the stated criteria), a single recommendation with rationale/risks, and a runnable proof of concept for a sample city.

**Independent Test**: Review research.md and confirm it compares ≥ 3 candidates (≥ 1 public API and ≥ 1 local dataset/library) against data quality, Go integration, auth, rate limits, and relevance; names one recommended approach with rationale and risks; and references a PoC that returns a plausible attractions list for a sample city.

> Note: the findings document already exists at [research.md](research.md). The tasks below verify it satisfies FR-012/FR-013 and bind the proof of concept (FR-014) to a reproducible artifact.

- [x] T020 [P] [US2] Verify research.md evaluates ≥ 3 candidates (WDQS, OpenTripMap, Overpass/OSM) against all five criteria — data quality, Go integration, auth, rate limits, relevance (FR-012) — and update if any criterion is missing
- [x] T021 [P] [US2] Verify research.md states a single recommended approach (WDQS SPARQL) with rationale and identified risks/limitations (FR-013); fill gaps if present
- [x] T022 [US2] Provide the proof of concept (FR-014): a reproducible sample-city lookup demonstrating the recommended approach returns a plausible attractions list — wired via the T007 golden fixture so the PoC is runnable offline through `go test` (depends on T007, T017)
- [x] T023 [US2] Cross-link the PoC and recommendation in quickstart.md "Proof of Concept" / "Spike Wrap-Up" sections so a reviewer can confirm feasibility in < 5 minutes (SC-006, SC-007)

**Checkpoint**: The spike's documented recommendation and proof of concept are complete and verifiable.

---

## Phase 5: Polish & Cross-Cutting Concerns

**Purpose**: Final validation and cleanup across both stories.

- [x] T024 [P] Run `gofmt -l` and `go vet ./...` across the repo; resolve any findings
- [x] T025 [P] Confirm full suite passes offline with `go test ./...` (no live network in tests)
- [x] T026 Walk through quickstart.md steps 1–7 against the built binary (`bin/citysearch`), confirming each maps to its contract scenario (A1–A7) and exit codes
- [x] T027 Remove disposable spike scaffolding (placeholder .gitkeep, dead exploratory code) per Constitution Principle III, keeping the durable client + fixtures

---

## Dependencies & Execution Order

### Phase Dependencies

- **Setup (Phase 1)**: No dependencies — start immediately.
- **Foundational (Phase 2)**: Depends on Setup — BLOCKS both user stories.
- **User Story 1 (Phase 3)**: Depends on Foundational. Independently testable and deliverable as the MVP.
- **User Story 2 (Phase 4)**: Depends on Foundational; its PoC task (T022) also depends on US1's fixture (T007) and lookup (T017). The documentation-verification tasks (T020, T021) are independent of US1.
- **Polish (Phase 5)**: Depends on both user stories being complete.

### User Story Dependencies

- **US1 (P1)**: Core MVP — no dependency on US2.
- **US2 (P1)**: Documentation/verification (T020, T021) is fully independent; the PoC (T022) reuses US1's lookup + fixture.

### Within Each User Story

- Tests (T007–T013) are written FIRST and must FAIL before implementation (T014–T019).
- Within implementation: parser/builder → mapping → ranking/cap → lookup → CLI wiring → rendering.

### Parallel Opportunities

- Setup: T002, T003 can run in parallel.
- Foundational: T004 and T005 in parallel; T006 follows.
- US1 tests T007–T013 are all [P] (distinct files/fixtures) and can be authored together before implementation.
- US2 doc-verification T020, T021 run in parallel and alongside all of US1.
- Polish: T024 and T025 in parallel.

---

## Parallel Example: User Story 1

```bash
# Author all US1 tests/fixtures together (before implementation):
Task: "Record Paris golden fixture in internal/attraction/testdata/paris.json"
Task: "Add many.json and empty.json fixtures in internal/attraction/testdata/"
Task: "Mapping test in internal/attraction/attraction_test.go"
Task: "Ranking + cap test in internal/attraction/attraction_test.go"
Task: "Short/empty path test in internal/attraction/attraction_test.go"
Task: "Fetch-error test in internal/attraction/wikidata_test.go"
Task: "CLI selection→attractions test in cmd/citysearch/interactive_test.go"
```

---

## Implementation Strategy

### MVP First (User Story 1)

1. Complete Phase 1: Setup.
2. Complete Phase 2: Foundational (types + fetcher seam — blocks everything).
3. Complete Phase 3: User Story 1 (write failing tests, then fetch→rank→cap→render).
4. **STOP and VALIDATE**: Run `go test ./...` and walk a Paris selection end-to-end.
5. This is the demoable MVP: pick a city, see its top-10 attractions.

### Incremental Delivery

1. Setup + Foundational → foundation ready.
2. Add User Story 1 → test independently → demo the attractions display (MVP).
3. Add User Story 2 → verify findings doc + bind reproducible PoC → spike deliverable complete.
4. Polish → format/vet, offline test run, quickstart walkthrough, remove disposable scaffolding.

### Parallel Team Strategy

With multiple developers, after Foundational completes:

- Developer A: User Story 1 (the `internal/attraction` logic + CLI wiring).
- Developer B: User Story 2 doc verification (T020, T021), then pairs on the PoC (T022) once US1's lookup lands.

---

## Notes

- [P] tasks = different files, no dependencies.
- [Story] label maps each task to US1 or US2 for traceability.
- Tests use recorded JSON fixtures only — no live network in `go test`.
- Verify tests fail before implementing (Test-First for durable logic).
- Commit after each task or logical group.
- Spike code under the time-box is exploratory (Constitution III); durable logic (ranking, cap, mapping, error paths) stays test-covered.
