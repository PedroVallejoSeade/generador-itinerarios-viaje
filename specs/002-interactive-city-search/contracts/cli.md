# CLI Contract — Interactive City Search

This feature extends the [001 CLI contract](../001-city-search-spike/contracts/cli.md). The
program now has **two modes** selected by the presence of a positional argument. The one-shot
mode is unchanged; this document defines the **interactive** mode and the mode-selection rule.

## Invocation

```text
citysearch                # interactive mode (no positional argument)
citysearch <city-name>    # one-shot mode (unchanged — see 001 contract)
citysearch -h | --help    # usage; exit 0
```

**Mode selection**

| Condition | Mode | Behavior |
|-----------|------|----------|
| No positional argument | Interactive | Welcome → prompt loop → clean exit (this contract). |
| One positional argument | One-shot | Single lookup, un-numbered output, exit codes per 001 (FR-013). |

## Interactive session — I/O contract

### Startup (stdout)

On launch with no argument, before any input is read (FR-001, FR-002, FR-002a):

```text
Welcome to the Travel Itinerary Generator!
Enter a city name to find matching destinations.
(type 'exit' or press Ctrl+D to quit)
```

> Exact wording is illustrative; the contract requires: a welcome line that frames the tool, an
> instruction to enter a city name, and a visible hint on how to exit.

### Prompt (stdout)

Before each line of input, a prompt is shown (FR-002, FR-010), e.g.:

```text
city>
```

The prompt re-appears after every response, conveying that searching can repeat.

### Input

| Input | Source | Handling |
|-------|--------|----------|
| A line of text | stdin (one full line, Enter-terminated) | Trimmed of surrounding whitespace (FR-004); whole line is the query (multi-word allowed, no quoting). |
| Empty / whitespace-only line | stdin | Friendly re-prompt; no search (FR-008). |
| `exit` / `quit` (any case, full trimmed line) | stdin | Closing message; exit 0 (FR-011). |
| End-of-input (Ctrl+D) | stdin closed | Closing message; exit 0 (FR-011). |

### Response — query with matches (stdout)

Up to 10 lines, population-descending, each prefixed with a 1-based index (FR-006, FR-007):

```text
1. <Name>, <Region>, <Country>
2. <Name>, <Country>
```

Region is omitted when empty (line 2 form). Fewer than 10 matches print only the actual matches,
with no padding (US1 scenario 3).

### Response — query with no matches (stdout)

A clear message; the loop continues (FR-009, FR-010):

```text
No cities found matching "<query>".
```

### Response — empty / whitespace-only input (stdout)

A friendly message; the loop continues (FR-008):

```text
Please enter a city name.
```

### Closing (stdout)

On `exit`/`quit` or Ctrl+D, a brief closing message before terminating (FR-011):

```text
Goodbye! Safe travels.
```

### Data-load failure (stderr)

If the embedded dataset cannot be loaded at startup, report on stderr and exit non-zero **before**
showing the prompt (FR-012, edge case):

```text
error: unable to load city data: <reason>
```

## Exit Codes

| Code | Meaning |
|------|---------|
| 0 | Clean interactive exit (exit keyword or Ctrl+D), including sessions with no-match or empty inputs. |
| 2 | Data-source/load failure (consistent with the 001 contract). |

> One-shot mode retains code 1 for invalid usage (empty/whitespace argument); interactive mode
> never returns 1 because empty input is handled in-loop rather than terminating.

## Contract Test Scenarios

| ID | Scripted stdin | Expected (channel) | Exit |
|----|----------------|--------------------|------|
| I1 | *(none — startup only)* | Welcome + exit hint + prompt appear before any read (stdout) | — |
| I2 | `London`⏎ then `exit`⏎ | ≥1 numbered line incl. a "London, …, United Kingdom" match; prompt re-appears; closing message (stdout) | 0 |
| I3 | ` paris `⏎ then EOF | Trimmed query matches "Paris, …, France"; closing message on EOF (stdout) | 0 |
| I4 | `lOnDoN`⏎ then `quit`⏎ | Same matches as canonical case (stdout) | 0 |
| I5 | `   `⏎ (whitespace) then EOF | "Please enter a city name." then re-prompt; no search (stdout) | 0 |
| I6 | `zzzzzz`⏎ then EOF | `No cities found matching "zzzzzz".` then re-prompt (stdout) | 0 |
| I7 | `London`⏎ `Paris`⏎ `exit`⏎ | Two result blocks in one session, prompt between them, then closing (stdout) | 0 |
| I8 | `EXIT`⏎ | Closing message; full-line keyword match is case-insensitive (stdout) | 0 |
| I9 | `san jose`⏎ then EOF | Multi-word line treated as one query without quoting (stdout) | 0 |
| I10 | *(dataset load fails)* | `error: unable to load city data: …` on stderr; no prompt | 2 |
