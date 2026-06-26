# Implementation Plan: Interactive City Search CLI

**Branch**: `002-interactive-city-search` | **Date**: 2026-06-26 | **Spec**: [spec.md](spec.md)

**Input**: Feature specification from `/specs/002-interactive-city-search/spec.md`

**Note**: This template is filled in by the `/speckit.plan` command. See `.specify/templates/plan-template.md` for the execution workflow.

## Summary

Transform the existing one-shot `citysearch` command into a guided, interactive terminal
experience. When launched **without** a city argument, the app prints a welcome message and an
exit hint, then enters a read-evaluate-print loop: prompt → read a full line → trim → classify
the input → respond. Recognized inputs are a search query (display up to 10 ranked matches,
each prefixed with a 1-based index), an empty/whitespace line (friendly re-prompt), or an exit
action (`exit`/`quit` case-insensitively, or Ctrl+D / end-of-input) which prints a closing
message and exits 0. When launched **with** a city argument, the existing single-query behavior
is preserved unchanged (FR-013).

The matching engine (`city.Load` + `city.Search`) is reused verbatim from the
[001 spike](../001-city-search-spike/plan.md); this feature adds only the interactive
session loop and numbered rendering. No new dependencies are introduced — `bufio.Scanner`
from the standard library drives line reading and EOF detection. Per Constitution Principle II
this is production code (not a spike), so it follows the Red-Green-Refactor TDD workflow with
the session loop driven by injectable `io.Reader`/`io.Writer` for scripted, table-driven tests.

## Technical Context

**Language/Version**: Go 1.26.4 (latest stable, per constitution Technology Standards and `go.mod`).

**Primary Dependencies**: Go standard library only — `bufio` (line scanning + EOF), `flag`, `fmt`, `os`, `strings`, `io`. No external modules (Simplicity First).

**Storage**: Reuses the bundled world-cities CSV embedded at build time via `//go:embed` in `internal/city`. No new storage; no runtime file I/O for the dataset.

**Testing**: Go standard `testing` package, table-driven tests. The interactive session is tested by feeding scripted input through an `io.Reader` and asserting on captured `io.Writer` output (no real TTY needed).

**Target Platform**: Cross-platform terminal/CLI (Linux primary; pure-Go, no cgo → macOS/Windows build cleanly). Input over stdin, results to stdout, errors to stderr.

**Project Type**: Single-project CLI application.

**Performance Goals**: Each query is displayed effectively instantly (SC-005). The dataset is loaded once at startup; per-query lookups are an in-memory linear prefix scan completing in single-digit milliseconds.

**Constraints**: Reuse existing matching (case-insensitive prefix), ranking (population descending), and the 10-result cap (FR-005, FR-006). Exit on `exit`/`quit` (case-insensitive, full trimmed line) or Ctrl+D (FR-011). Clean exit returns 0; data-load failure returns non-zero (FR-012).

**Scale/Scope**: One interactive session = welcome + one-or-more prompt/result cycles + clean exit. Dataset ≈ 25k rows (cities ≥ 15,000 population), held in memory once per process.

## Constitution Check

*GATE: Must pass before Phase 0 research. Re-check after Phase 1 design.*

| Principle | Status | Notes |
|-----------|--------|-------|
| I. Milestone-Driven Development | PASS | Three prioritized, independently testable slices map directly to the spec's user stories: P1 guided greeting→prompt→results, P2 empty/no-match handling, P3 multi-search session + clean exit. Implemented one milestone at a time. |
| II. Test-First (TDD — NON-NEGOTIABLE) | PASS | Not a spike — full Red-Green-Refactor applies. The session loop is written behind `io.Reader`/`io.Writer` seams so behavior (welcome, prompt, numbered results, empty/no-match messages, exit keywords, EOF) is specified by table-driven tests written before implementation. |
| III. Spike-First Exploration | PASS (N/A) | No unproven approach: data sourcing and matching were resolved by the 001 spike; this feature reuses them. No new spike required. |
| IV. CLI Interface | PASS | stdin for input, stdout for results/prompts, stderr for errors; exit 0 on clean quit, non-zero on data-load failure; welcome text includes an exit hint (usage guidance). |
| V. Simplicity First (YAGNI) | PASS | Standard library only (`bufio.Scanner`); no readline/TUI library, no command framework, no history/autocomplete. Reuses existing search; adds only a loop, input classification, and index prefix. |

**Result**: No violations. Complexity Tracking section intentionally left empty.

## Project Structure

### Documentation (this feature)

```text
specs/002-interactive-city-search/
├── plan.md              # This file (/speckit.plan command output)
├── research.md          # Phase 0 output — interactive I/O & testability decisions
├── data-model.md        # Phase 1 output — InteractiveSession, InputCommand, numbered rendering
├── quickstart.md        # Phase 1 output — build/run/validate guide
├── contracts/
│   └── cli.md           # Phase 1 output — interactive CLI input/output contract
├── checklists/
│   └── requirements.md  # Spec quality checklist (already present)
└── tasks.md             # Phase 2 output (/speckit.tasks — NOT created here)
```

### Source Code (repository root)

```text
cmd/
└── citysearch/
    ├── main.go             # Entry point: mode selection (arg → one-shot; no arg → interactive),
    │                       #   dataset load, exit codes. Refactored to io.Reader/io.Writer seams.
    ├── interactive.go      # NEW — interactive session loop: welcome, prompt, read/classify/respond
    ├── interactive_test.go # NEW — table-driven tests driving the loop via scripted io.Reader
    └── main_test.go        # NEW (if absent) — one-shot mode + mode-selection tests

internal/
└── city/
    ├── city.go             # UNCHANGED — City type + embedded dataset load/parse
    ├── search.go           # UNCHANGED — normalize + prefix match + population rank + cap
    ├── search_test.go      # UNCHANGED — existing matching unit tests
    ├── cities.csv          # UNCHANGED — embedded dataset
    └── testdata/
        └── cities_sample.csv  # UNCHANGED — deterministic fixture (reused if package tests need it)
```

**Structure Decision**: Single Go module
(`github.com/PedroVallejoSeade/generador-itinerarios-viaje`), continuing the 001 layout:
`cmd/citysearch` owns the CLI entry point and the interactive session loop (application
wiring), while `internal/city` continues to own all matching/data logic and is reused
unchanged. The interactive loop lives in the `main` package rather than a new `internal`
package because it is presentation/wiring, not reusable domain logic (YAGNI). Testability is
achieved by depending on `io.Reader`/`io.Writer` interfaces (not `*os.File`), so the session
is fully exercised by `package main` tests with scripted input — no new package boundary is
needed. The one behavioral change to existing code is refactoring `run` and helpers from
`*os.File` to `io.Reader`/`io.Writer` to enable both modes to be tested.

## Complexity Tracking

> No constitution violations to justify — section intentionally empty.
