# CLI Contract — City Search Spike

The application's external interface is a command-line program. This contract defines its
inputs, outputs, and exit codes (Constitution Principle IV: CLI Interface).

## Invocation

```text
citysearch <city-name>
citysearch -h | --help
```

- The city name is supplied as a positional argument. A multi-word name may be quoted
  (e.g., `citysearch "san jose"`).

## Inputs

| Input | Source | Validation |
|-------|--------|-----------|
| City name query | First positional CLI argument | Trimmed of surrounding whitespace; must be non-empty after trimming (FR-001, FR-008). |
| `-h` / `--help` | Flag | Prints usage and exits 0. |

## Output Contract

### Success — one or more matches (stdout)

One result per line, ordered by population (descending), capped at 10 (FR-005, FR-006).
Each line carries disambiguating context (FR-004):

```text
<Name>, <Region>, <Country>
```

When a city's region is empty, it is omitted: `<Name>, <Country>`.

Example for `citysearch springfield`:

```text
Springfield, Missouri, United States
Springfield, Illinois, United States
Springfield, Massachusetts, United States
```

### Success — no matches (stdout)

A clear, explanatory message (FR-007); exit code remains 0 (a valid query that found nothing
is not an error):

```text
No cities found matching "zzzzzz".
```

### Empty / whitespace-only query (stderr)

Reject without querying the dataset (FR-008); exit non-zero:

```text
Please provide a city name to search for.
```

### Data-source / load failure (stderr)

If the embedded dataset cannot be parsed/loaded, report on stderr and exit non-zero (FR-009):

```text
error: unable to load city data: <reason>
```

## Exit Codes

| Code | Meaning |
|------|---------|
| 0 | Query processed successfully (including the valid "no results found" case). |
| 1 | Invalid usage — empty/whitespace query or missing argument. |
| 2 | Data-source/load error (FR-009). |

## Contract Test Scenarios

| ID | Input | Expected channel | Expected outcome | Exit |
|----|-------|------------------|------------------|------|
| C1 | `citysearch paris` | stdout | List includes a "Paris, …, France" line | 0 |
| C2 | `citysearch springfield` | stdout | ≥ 2 lines, each distinguishable by region/country | 0 |
| C3 | `citysearch "São Paulo"` | stdout | Includes the accented "São Paulo, …, Brazil" match | 0 |
| C4 | `citysearch zzzzzz` | stdout | "No cities found matching …" message | 0 |
| C5 | `citysearch ""` (or no arg) | stderr | Prompt for a valid city name | 1 |
| C6 | many-match common name | stdout | At most 10 lines, population-ordered | 0 |
