# Phase 0 Research — Interactive City Search CLI

This feature reuses the data source and matching engine already proven by the
[001 spike](../001-city-search-spike/research.md), so there are **no `NEEDS CLARIFICATION`
items** about data sourcing. The open questions are limited to *how* to deliver an interactive
terminal loop in Go that is simple, standard-library-only, and testable. Each decision below
follows the Red-Green-Refactor and YAGNI principles.

## Decision 1: Line reading & end-of-input (Ctrl+D) detection

**Decision**: Read input with `bufio.NewScanner(stdin)` and loop while `scanner.Scan()` returns
`true`. A `false` return with no scanner error signals end-of-input (Ctrl+D), which ends the
session cleanly (FR-011).

**Rationale**:

- `bufio.Scanner` with the default `ScanLines` split function reads one full line per `Scan()`
  call and strips the trailing newline — exactly the "read a full line as the query" behavior
  required by FR-003, including multi-word names without quoting.
- `Scan()` returning `false` cleanly distinguishes EOF from a populated line, giving a natural,
  branch-free way to detect Ctrl+D and exit with a closing message and status 0.
- It is standard library, zero dependencies (Principle V).

**Alternatives considered**:

- `bufio.Reader.ReadString('\n')` — viable but forces manual handling of the `io.EOF` sentinel
  intermixed with partial final lines; more error-prone than `Scanner`'s boolean loop.
- `fmt.Fscan` / `fmt.Fscanln` — splits on whitespace, breaking multi-word city names; rejected.
- A readline library (e.g. `chzyer/readline`) — adds history/editing/autocomplete the spec does
  not require; violates YAGNI and the standard-library-first standard; rejected.

**Caveat noted for implementation**: very long lines exceeding `bufio.Scanner`'s default 64 KB
token limit return a scan error. City queries are far shorter than this; the loop treats any
scanner error as end-of-session (exit cleanly) rather than crashing.

## Decision 2: Input classification (query vs. empty vs. exit)

**Decision**: After trimming surrounding whitespace, classify each line in this order:
(1) empty string → empty-input branch (FR-008); (2) case-insensitive equality with `exit` or
`quit` → exit branch (FR-011); (3) otherwise → search query (FR-005, FR-006).

**Rationale**:

- Trim-first makes whitespace-only input collapse to empty (edge case: whitespace-only → friendly
  re-prompt) and makes exit-keyword matching robust to surrounding spaces.
- Matching the **full trimmed line** case-insensitively against `exit`/`quit` (per the spec
  clarification) avoids accidentally treating a city literally named with those words as an exit
  unless the entire input is exactly that keyword.
- Reuses the existing `city.Search` normalization for the query path — no duplicate normalization
  logic.

**Alternatives considered**:

- Substring/prefix exit detection (e.g. line starts with `exit`) — would prevent searching for a
  city whose name begins with "exit"; rejected in favor of full-line equality.
- A command registry/parser abstraction — over-engineered for two keywords; YAGNI; rejected.

## Decision 3: Testability via `io.Reader` / `io.Writer` seams

**Decision**: Implement the session loop as a function taking `stdin io.Reader`,
`stdout io.Writer`, `stderr io.Writer`, and the pre-loaded `[]city.City`. Refactor the existing
`run`/helpers from `*os.File` to these interfaces. `main` wires in `os.Stdin`/`os.Stdout`/
`os.Stderr`; tests wire in `strings.Reader` (scripted input) and `bytes.Buffer` (captured output).

**Rationale**:

- TDD is non-negotiable here (not a spike). Driving the loop with scripted input and asserting on
  captured output exercises every branch — welcome text, prompt, numbered results, empty message,
  no-match message, exit keywords, and EOF — without a real TTY.
- `io.Reader`/`io.Writer` are the idiomatic Go seams; no mocking framework needed (Principle V).
- Keeping the loop in `package main` (tested by `package main` tests) avoids inventing a new
  `internal` package for what is presentation/wiring, not reusable domain logic.

**Alternatives considered**:

- Spawning the built binary and piping stdin in tests (integration-only) — slower, harder to
  assert on intermediate prompts, and would not satisfy unit-level Red-Green-Refactor; kept only
  as an optional end-to-end check in `quickstart.md`.
- Extracting an `internal/session` package — unnecessary boundary for non-reusable wiring; YAGNI.

## Decision 4: Numbered result rendering

**Decision**: Render each result line as `"<n>. <Name>, <Region>, <Country>"` using a 1-based
index `n`, reusing the existing region-omitting format (`"<n>. <Name>, <Country>"` when region
is empty) from the 001 one-shot renderer.

**Rationale**:

- Matches the spec clarification and FR-007 exactly while preserving visual consistency with the
  existing single-line format.
- The one-shot mode keeps its un-numbered output (FR-013 — existing behavior unchanged); only the
  interactive list is numbered, so the index prefix is applied in the interactive renderer.

**Alternatives considered**:

- Numbering one-shot output too — would change existing behavior the spec says to preserve;
  rejected.
- Column-aligned/tabular output (`text/tabwriter`) — nicer alignment but beyond the "clean,
  readable single line" the spec asks for; YAGNI; can be revisited if requested.

## Decision 5: Mode selection (interactive vs. one-shot)

**Decision**: In `main`/`run`, after flag parsing: if a positional city argument is present, run
the existing one-shot path (FR-013); if no argument is present, load the dataset once and enter
the interactive session loop.

**Rationale**:

- A single, explicit branch on `flag.NArg()` cleanly separates the two modes with no new flags.
- The dataset is loaded once before the loop so repeated searches in a session incur no reload
  (SC-005), and a load failure is reported before the prompt appears (FR-012 / edge case).

**Alternatives considered**:

- A dedicated `-i`/`--interactive` flag — redundant; "no argument" already unambiguously signals
  intent and keeps the surface minimal (YAGNI); rejected.

## Summary of resolved unknowns

| Unknown | Resolution |
|---------|-----------|
| How to read full lines & detect Ctrl+D | `bufio.Scanner`; `Scan()==false` ⇒ end-of-input, clean exit |
| How to classify input | Trim → empty? → `exit`/`quit` (case-insensitive full line)? → else query |
| How to make the loop testable under TDD | `io.Reader`/`io.Writer` seams; scripted input + captured output in `package main` tests |
| How to render results | 1-based index prefix on the existing `Name, Region, Country` line |
| How to choose mode | Positional arg ⇒ one-shot (unchanged); no arg ⇒ interactive |

No `NEEDS CLARIFICATION` items remain. Proceed to Phase 1.
