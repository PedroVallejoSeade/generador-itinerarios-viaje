# Quickstart — Interactive City Search CLI

A build/run/validate guide for the interactive mode. For the data model and exact I/O contract,
see [data-model.md](data-model.md) and [contracts/cli.md](contracts/cli.md). This guide does not
contain implementation code; the loop is built during the implementation phase under TDD.

## Prerequisites

- Go 1.26.4+ on `PATH` (`go version`).
- Repository checked out; run commands from the repo root.

## Build

```bash
go build -o bin/citysearch ./cmd/citysearch
```

## Run — interactive mode (no argument)

```bash
./bin/citysearch
```

Expected: a welcome message, an exit hint, and a `city>`-style prompt appear before any input is
read (FR-001, FR-002, FR-002a). Type a city name and press Enter to see up to 10 numbered matches.

Example session:

```text
Welcome to the Travel Itinerary Generator!
Enter a city name to find matching destinations.
(type 'exit' or press Ctrl+D to quit)
city> london
1. London, England, United Kingdom
2. London, Ontario, Canada
...
city> exit
Goodbye! Safe travels.
```

## Run — one-shot mode (unchanged, FR-013)

```bash
./bin/citysearch "san jose"
```

Expected: the existing single-query, **un-numbered** output; the app does **not** enter
interactive mode.

## Validation scenarios

Run each and confirm the expected outcome (maps to [contracts/cli.md](contracts/cli.md) IDs).

| # | Steps | Expected |
|---|-------|----------|
| 1 (I1) | Launch with no args | Welcome + exit hint + prompt shown before any input (SC-001) |
| 2 (I2) | Type `London`, Enter | Numbered list; top match is the largest London; prompt returns (SC-002) |
| 3 (I3) | Type ` paris ` (spaces), Enter | Same result as `paris` — whitespace trimmed |
| 4 (I4) | Type `lOnDoN`, Enter | Same matches as `london` — case-insensitive |
| 5 (I5) | Press Enter on empty line | "Please enter a city name." + re-prompt; no search (SC-003) |
| 6 (I6) | Type `zzzzzz`, Enter | "No cities found matching …" + re-prompt (SC-003) |
| 7 (I7) | Search two cities, then `exit` | Two result blocks, then closing message; exit 0 (SC-004) |
| 8 (I9) | Type `san jose`, Enter | Multi-word query works without quotes |
| 9 | Press Ctrl+D at the prompt | Closing message; exit 0 |

Quick exit-code check after a clean quit:

```bash
printf 'london\nexit\n' | ./bin/citysearch; echo "exit=$?"   # expect exit=0
```

## Run the tests

```bash
go test ./...
```

Expected (after implementation): the interactive-session tests in
`cmd/citysearch/interactive_test.go` and the existing `internal/city` tests pass. Per the
constitution's TDD principle, these tests are written **before** the loop implementation and must
fail first (Red), then pass (Green).

## Scripted (non-interactive) end-to-end check

Because the loop reads stdin line-by-line and exits on EOF, you can drive it without a TTY:

```bash
printf 'london\n  \nzzzzzz\nquit\n' | ./bin/citysearch
```

Expected, in order: numbered London results → "Please enter a city name." → "No cities found
matching \"zzzzzz\"." → closing message; process exits 0.
