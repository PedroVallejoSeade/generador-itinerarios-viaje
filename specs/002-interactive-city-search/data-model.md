# Phase 1 Data Model — Interactive City Search CLI

Derived from the spec's **Key Entities** section. The matching entities (`City`,
`SearchQuery`, `SearchResultSet`) are **reused unchanged** from the
[001 data model](../001-city-search-spike/data-model.md); this feature adds the
interactive-session concepts only.

## Reused entities (from 001 — unchanged)

| Entity | Role in this feature |
|--------|----------------------|
| `City` | A populated place with `Name`, `Country`, `Region`, `Population`. Loaded once per session via `city.Load`. |
| `SearchQuery` | The trimmed, lower-cased query passed to `city.Search`. Normalization is reused, not redefined. |
| `SearchResultSet` | Up to 10 cities, population-descending, returned by `city.Search`. |

No fields, validation rules, or matching logic change for these entities.

## Entity: InteractiveSession

A single run of the application in interactive mode: welcome → one-or-more prompt/result
cycles → clean exit.

| Field | Type | Source | Notes |
|-------|------|--------|-------|
| `cities` | `[]city.City` | `city.Load()` once at session start | Held in memory for the whole session; no reload between searches (SC-005). |
| `in` | `io.Reader` | `os.Stdin` (real) / scripted reader (tests) | Wrapped in a `bufio.Scanner` for line reading. |
| `out` | `io.Writer` | `os.Stdout` (real) / buffer (tests) | Welcome, prompt, results, friendly messages, closing message. |
| `err` | `io.Writer` | `os.Stderr` (real) / buffer (tests) | Data-load failure messages only. |

**Lifecycle / states**

```text
Start ─► Welcome ─► Prompt ─► ReadLine ─┬─ EOF/exit-keyword ─► Closing ─► End(0)
                      ▲                  ├─ empty/whitespace ─► EmptyMsg ─┐
                      │                  └─ query ─► Results | NoMatchMsg ┤
                      └───────────────────────────────────────────────────┘
```

- The loop returns to **Prompt** after every empty, no-match, or results response (FR-010).
- A data-load failure before the loop short-circuits to **End(non-zero)** with an error on `err`
  (FR-012); the prompt is never shown in that case.

## Entity: InputCommand (classification of one input line)

The result of normalizing and classifying a single line read at the prompt. Not a persisted
record — a transient decision used to route the loop.

| Variant | Trigger (after `strings.TrimSpace`) | Effect |
|---------|-------------------------------------|--------|
| `Exit` | Trimmed line equals `exit` or `quit`, compared case-insensitively (full line); **or** end-of-input (Ctrl+D / scanner stops) | Print closing message; end session with status 0 (FR-011). |
| `Empty` | Trimmed line is `""` (covers whitespace-only input) | Print friendly "please enter a city name" message; re-prompt; no search performed (FR-008). |
| `Query` | Any other non-empty trimmed line | Pass to `city.Search`; render results or no-match message; re-prompt (FR-005, FR-006, FR-009). |

**Classification rules**

1. Trim surrounding whitespace first (FR-004) — makes whitespace-only collapse to `Empty` and
   makes keyword matching space-tolerant.
2. Empty check precedes keyword check precedes query.
3. Exit-keyword comparison is `strings.EqualFold(trimmed, "exit") || strings.EqualFold(trimmed, "quit")`
   — full-line equality, so a city whose name merely starts with "exit" is still searchable.

## Rendering: numbered result line

| Property | Value | Rule |
|----------|-------|------|
| Index | 1-based ordinal of the result within the displayed set | FR-007 / spec clarification |
| Line format (region present) | `"<n>. <Name>, <Region>, <Country>"` | FR-007 |
| Line format (region empty) | `"<n>. <Name>, <Country>"` | FR-007 (graceful region omission) |
| Order | As returned by `city.Search` (population descending) | FR-006 |
| Count | At most 10 lines; only actual matches (no padding) | FR-006, US1 scenario 3 |

> Note: the **one-shot** mode keeps its existing **un-numbered** output (FR-013). The index
> prefix applies to the interactive list only.

## Traceability

| Requirement | Model element |
|-------------|---------------|
| FR-001 / FR-002 / FR-002a | `InteractiveSession` Welcome state (welcome + prompt + exit hint) |
| FR-003 / FR-004 | `InputCommand` line read + trim; reused `SearchQuery` normalization |
| FR-005 / FR-006 | Reused `SearchResultSet` (prefix match, population rank, cap 10) |
| FR-007 | Numbered result line rendering |
| FR-008 | `InputCommand.Empty` variant |
| FR-009 | `InputCommand.Query` → no-match message |
| FR-010 | `InteractiveSession` loop returns to Prompt after each response |
| FR-011 | `InputCommand.Exit` variant (keywords + EOF) + Closing state |
| FR-012 | `InteractiveSession` data-load failure → End(non-zero); clean exit → End(0) |
| FR-013 | Mode selection: positional arg bypasses `InteractiveSession` (un-numbered one-shot) |
