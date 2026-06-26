---
description: "Task list for Interactive City Search CLI implementation"
---

# Tasks: Interactive City Search CLI

**Input**: Design documents from `/specs/002-interactive-city-search/`

**Prerequisites**: plan.md (required), spec.md (required for user stories), research.md, data-model.md, contracts/cli.md

**Tests**: TDD is NON-NEGOTIABLE for this feature (Constitution Principle II; the loop is production code, not a spike). Test tasks are therefore INCLUDED and MUST be written first and fail before implementation (Red-Green-Refactor).

**Organization**: Tasks are grouped by user story to enable independent implementation and testing of each story.

## Format: `[ID] [P?] [Story] Description`

- **[P]**: Can run in parallel (different files, no dependencies)
- **[Story]**: Which user story this task belongs to (US1, US2, US3)
- Include exact file paths in descriptions

## Path Conventions

- **Single Go project** at repository root: entry point + interactive loop in `cmd/citysearch/`; reused matching/data logic in `internal/city/` (unchanged this feature).
- Module: `github.com/PedroVallejoSeade/generador-itinerarios-viaje`.

---

## Phase 1: Setup (Shared Infrastructure)

**Purpose**: Confirm the existing project builds and tests green before changing anything.

- [ ] T001 Verify baseline builds and tests pass by running `go build ./...` and `go test ./...` from the repository root; record any pre-existing failures before starting.
- [ ] T002 [P] Verify formatting/vet baseline by running `gofmt -l .` and `go vet ./...` from the repository root; the working tree must be clean before edits.

---

## Phase 2: Foundational (Blocking Prerequisites)

**Purpose**: Refactor the CLI entry point to `io.Reader`/`io.Writer` seams so BOTH one-shot and interactive modes are unit-testable with scripted input and captured output (research Decision 3). This blocks all user stories because the interactive loop and its tests depend on these seams.

**⚠️ CRITICAL**: No user story work can begin until this phase is complete.

- [ ] T003 Add `cmd/citysearch/main_test.go` (package `main`) with table-driven tests for the existing one-shot mode driven through `io.Writer` seams: assert stdout for a matching query (un-numbered `Name, Region, Country` lines, region omitted when empty), the no-match message, the empty/missing-argument usage error on stderr, and the data-load/exit-code contract (0 success, 1 invalid usage, 2 data-load failure per FR-013 / 001 contract). These tests MUST fail to compile/run first (Red).
- [ ] T004 Refactor `run` in `cmd/citysearch/main.go` to accept `stdout io.Writer, stderr io.Writer` (replacing `*os.File`), update the `format` helper and all `fmt.Fprint*` call sites accordingly, and update `main()` to pass `os.Stdout`/`os.Stderr`; keep one-shot behavior identical so T003 passes (Green). Add the `io` import.
- [ ] T005 Extend `run` in `cmd/citysearch/main.go` to also accept `stdin io.Reader` (wired from `os.Stdin` in `main()`), threading it through unused for now so the signature supports the interactive loop added in Phase 3; ensure T003 still passes.

**Checkpoint**: Entry point exposes `io.Reader`/`io.Writer` seams; one-shot mode unchanged and fully tested. User stories can now begin.

---

## Phase 3: User Story 1 - Guided interactive city search (Priority: P1) 🎯 MVP

**Goal**: Launching with no argument shows a welcome message, an exit hint, and a prompt before reading input; typing a city name and pressing Enter displays up to 10 population-ranked matches, each prefixed with a 1-based index and formatted as `Name, Region, Country` (region omitted when unknown).

**Independent Test**: Run with no args, feed `London\n` then EOF via a scripted `io.Reader`; assert the welcome + exit hint + prompt appear before any read, and that a numbered list (`1. London, ..., United Kingdom`, ≤10 lines, no padding) is captured on stdout.

### Tests for User Story 1 (write FIRST — must fail before implementation) ⚠️

- [ ] T006 [P] [US1] Add table-driven tests in `cmd/citysearch/interactive_test.go` (package `main`) driving the interactive session through a scripted `strings.Reader` and a captured `bytes.Buffer`: cover I1 (welcome + exit hint + prompt emitted before any input is read — contract scenario I1 / FR-001, FR-002, FR-002a) and the numbered-results rendering for a known query, asserting 1-based index prefixes, ≤10 lines, population-descending order, and graceful region omission (FR-006, FR-007). Use `internal/city/testdata/cities_sample.csv` or a small in-memory `[]city.City` fixture for determinism. Tests MUST fail first (Red).

### Implementation for User Story 1

- [ ] T007 [US1] Create `cmd/citysearch/interactive.go` (package `main`) with a `runInteractive(in io.Reader, out io.Writer, cities []city.City) int` session function that prints the welcome message + exit hint (FR-001, FR-002a), prints the prompt (FR-002), reads one full line via `bufio.NewScanner` (FR-003), trims it (FR-004), and for a non-empty query calls `city.Search` and renders results; loop terminates on scanner EOF returning 0. Add an `interactive` rendering helper producing `"<n>. <Name>, <Region>, <Country>"` (region-omitted form when empty) per FR-007.
- [ ] T008 [US1] Wire mode selection in `run` (`cmd/citysearch/main.go`): when `fs.NArg() == 0`, load the dataset once via `city.Load()` (report load failure on stderr and return the data-load exit code before any prompt — FR-012), then call `runInteractive(stdin, stdout, cities)`; when an argument is present, keep the existing one-shot path (FR-013). Ensure T006 passes (Green).

**Checkpoint**: User Story 1 fully functional — guided welcome → prompt → numbered results works end-to-end and is independently testable.

---

## Phase 4: User Story 2 - Helpful handling of empty input and no matches (Priority: P2)

**Goal**: Pressing Enter with no text (or whitespace-only) shows a friendly "please enter a city name" message and re-prompts without searching; a query with no matches shows a clear "no cities found" message and re-prompts.

**Independent Test**: Feed `   \n` (whitespace) then `zzzzzz\n` then EOF; assert "Please enter a city name." (no search performed) followed by `No cities found matching "zzzzzz".`, each followed by a re-prompt (contract I5, I6).

### Tests for User Story 2 (write FIRST — must fail before implementation) ⚠️

- [ ] T009 [P] [US2] Add table-driven cases to `cmd/citysearch/interactive_test.go` covering empty/whitespace-only input → friendly message + re-prompt with no search (FR-008, contract I5) and a no-match query → `No cities found matching "<query>".` + re-prompt (FR-009, contract I6). Tests MUST fail first (Red).

### Implementation for User Story 2

- [ ] T010 [US2] Extend the input classification in `runInteractive` (`cmd/citysearch/interactive.go`) so that, after trimming, an empty string prints the friendly "Please enter a city name." message and re-prompts without calling `city.Search` (FR-008), and a query yielding zero results prints `No cities found matching "<query>".` then re-prompts (FR-009). Ensure T009 passes (Green).

**Checkpoint**: User Stories 1 AND 2 both work independently — empty and no-match inputs are handled gracefully and the loop continues.

---

## Phase 5: User Story 3 - Search multiple destinations in one session (Priority: P3)

**Goal**: After any response the prompt returns so the user can search again without relaunching; entering `exit`/`quit` (case-insensitive, full trimmed line) or signaling EOF (Ctrl+D) prints a brief closing message and exits with success.

**Independent Test**: Feed `London\nParis\nexit\n`; assert two distinct numbered result blocks with a prompt between them, then a closing message, then exit 0 (contract I7); separately assert `EXIT\n` (I8) and EOF both close cleanly.

### Tests for User Story 3 (write FIRST — must fail before implementation) ⚠️

- [ ] T011 [P] [US3] Add table-driven cases to `cmd/citysearch/interactive_test.go` covering: two sequential searches in one session with the prompt re-appearing between them (FR-010, contract I7); exit keyword recognition for `exit`/`quit` case-insensitively against the full trimmed line including `EXIT` (FR-011, contract I8); EOF/Ctrl+D clean close (contract I3); and that a closing message is printed and the function returns 0 on every exit path (FR-011, FR-012). Tests MUST fail first (Red).

### Implementation for User Story 3

- [ ] T012 [US3] Complete the `runInteractive` loop (`cmd/citysearch/interactive.go`): classify a trimmed line as Exit when it equals `exit` or `quit` via `strings.EqualFold` (full-line match — FR-011) or when the scanner stops/EOF occurs, printing the brief closing message and returning 0; ensure the prompt is re-emitted after every empty, no-match, and results response so repeated searching works (FR-010). Ensure T011 passes (Green).

**Checkpoint**: All three user stories independently functional — multi-search sessions and all clean-exit paths (keywords + EOF) work.

---

## Phase 6: Polish & Cross-Cutting Concerns

**Purpose**: Final consistency, documentation, and full-suite validation across all stories.

- [ ] T013 [P] Run `gofmt -w cmd/citysearch/` and `go vet ./...`; resolve any formatting or vet findings introduced by the new/edited files.
- [ ] T014 Update `cmd/citysearch/main.go` package/usage doc comments to describe the two modes (interactive vs. one-shot) so `-h/--help` and source docs reflect FR-013 without changing one-shot output.
- [ ] T015 Run the full suite `go test ./...` and confirm all `cmd/citysearch` and `internal/city` tests pass (Refactor step complete, no regressions).
- [ ] T016 Execute the `quickstart.md` validation scenarios — including `printf 'london\n  \nzzzzzz\nquit\n' | ./bin/citysearch` and `printf 'london\nexit\n' | ./bin/citysearch; echo "exit=$?"` (expect `exit=0`) — and confirm outputs match contracts/cli.md scenarios I1–I10.

---

## Dependencies & Execution Order

### Phase Dependencies

- **Setup (Phase 1)**: No dependencies — start immediately.
- **Foundational (Phase 2)**: Depends on Setup. The `io.Reader`/`io.Writer` refactor (T003–T005) BLOCKS all user stories.
- **User Stories (Phase 3–5)**: All depend on Foundational completion. They are layered on the same `runInteractive` function, so they are most naturally implemented in priority order (P1 → P2 → P3); each remains independently testable.
- **Polish (Phase 6)**: Depends on all targeted user stories being complete.

### User Story Dependencies

- **User Story 1 (P1)**: Starts after Foundational. Establishes the session skeleton (welcome, prompt, query→numbered results, EOF-terminated loop). No dependency on US2/US3.
- **User Story 2 (P2)**: Starts after Foundational. Extends classification in the same loop; independently testable. (Code-wise it edits the US1 loop, so sequence after US1 to avoid merge churn.)
- **User Story 3 (P3)**: Starts after Foundational. Adds exit keywords + closing message + multi-search validation to the same loop; independently testable. (Sequence after US1.)

### Within Each User Story

- The test task (T006 / T009 / T011) MUST be written and FAIL before its implementation task (T007–T008 / T010 / T012) — Red-Green-Refactor.
- Implementation: session function before mode wiring (US1: T007 before T008).

### Parallel Opportunities

- Setup: T002 can run in parallel with T001.
- The three user-story TEST tasks (T006, T009, T011) target the same file `cmd/citysearch/interactive_test.go`; the `[P]` marks them as independently authored test blocks, but if edited in the same working copy, sequence them to avoid file conflicts.
- Across stories, implementation tasks (T008, T010, T012) edit the shared `runInteractive`/`run` functions and therefore are NOT parallel with each other.
- Polish: T013 can run in parallel with documentation task T014.

---

## Parallel Example: User Story 1

```bash
# Write the failing tests first (Red), then implement (Green):
Task: "T006 [US1] Add interactive session tests in cmd/citysearch/interactive_test.go"
# After T006 fails as expected:
Task: "T007 [US1] Create runInteractive loop in cmd/citysearch/interactive.go"
Task: "T008 [US1] Wire no-arg mode selection in cmd/citysearch/main.go"
```

---

## Implementation Strategy

### MVP First (User Story 1 Only)

1. Complete Phase 1: Setup.
2. Complete Phase 2: Foundational (`io` seams — CRITICAL, blocks all stories).
3. Complete Phase 3: User Story 1.
4. **STOP and VALIDATE**: Run no-arg session, confirm welcome → prompt → numbered results.
5. This is a shippable MVP: the one-shot mode is preserved (FR-013) and interactive guided search works.

### Incremental Delivery

1. Setup + Foundational → seams ready, one-shot still green.
2. Add User Story 1 → guided numbered search → MVP.
3. Add User Story 2 → empty/no-match friendliness.
4. Add User Story 3 → multi-search session + clean exit.
5. Polish → format, docs, full quickstart validation.

### Notes

- `[P]` tasks = different files / independent edits, no dependencies.
- `[Story]` label maps each task to its user story for traceability.
- `internal/city/` is reused UNCHANGED this feature (no tasks touch it).
- Verify each story's tests fail before implementing (TDD — non-negotiable here).
- Commit after each task or logical group.
