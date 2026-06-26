# CLI Contract — City Attractions Lookup Spike

This contract extends the interactive `citysearch` CLI (Constitution Principle IV). It defines
the selection→attractions exchange: after the app lists numbered city matches, the user picks a
city by its number and the app displays that city's top attractions.

## Invocation

```text
citysearch                 Interactive mode (no argument): search → select → attractions.
citysearch -h | --help
```

The attractions step is reached **inside** interactive mode (no new flags). One-shot mode
(`citysearch <city-name>`) is unchanged by this feature.

## Interaction Flow

1. User enters a city name; the app prints up to 10 numbered matches (existing behavior).
2. App prompts the user to **select a city by its number**.
3. User enters a number:
   - **Valid** (matches a listed result) → app looks up and displays that city's top
     attractions (FR-001, FR-002).
   - **Invalid** (out of range / not a number) → app reports the invalid selection and does
     **not** perform a lookup (FR-008).

## Inputs

| Input | Source | Validation |
|-------|--------|-----------|
| City selection number | A line of stdin after results are shown | Must be an integer matching a displayed result's 1-based index; otherwise rejected (FR-008). |
| `-h` / `--help` | Flag | Prints usage and exits 0. |

## Output Contract

### Success — attractions found (stdout)

A numbered list, ranked most-known first, capped at 10 (FR-003, FR-004, FR-005). Category
and/or short description are appended when readily available:

```text
Top attractions in Paris, France:
1. Eiffel Tower — tower
2. Louvre — art museum
3. Notre-Dame de Paris — cathedral
4. Arc de Triomphe — triumphal arch
...
```

When fewer than 10 attractions exist, **all** available are shown without error (FR-007).

### Success — no attractions found (stdout)

A clear, explanatory message; the session continues (a valid selection that found nothing is
not a failure) (FR-006):

```text
No attractions found for Foober, Nowhereland.
```

### Invalid selection (stdout/stderr)

Reject without performing a lookup; the session continues for another selection (FR-008):

```text
"42" is not a valid selection. Enter the number of a listed city.
```

### Data-source / connectivity failure (stderr)

If the attraction source cannot be reached or returns an error, report on stderr and exit
non-zero (FR-009):

```text
error: unable to fetch attractions: <reason>
```

## Exit Codes

| Code | Meaning |
|------|---------|
| 0 | Selection processed successfully (including the valid "no attractions found" case). |
| 2 | Data-source/connectivity error while fetching attractions (FR-009). |

> Note: an invalid selection (FR-008) is an in-session, recoverable message, not a process
> exit — the loop continues, consistent with the existing interactive design.

## Contract Test Scenarios

| ID | Setup / Input | Expected channel | Expected outcome | Exit |
|----|----------------|------------------|------------------|------|
| A1 | Select a well-known city (e.g. Paris) | stdout | Numbered list with recognizable landmarks (Eiffel Tower, Louvre, …), ≤ 10 rows, most-known first | 0 |
| A2 | Selected city has > 10 attractions | stdout | Exactly 10 rows, ranked by prominence | 0 |
| A3 | Selected city has < 10 attractions | stdout | All available rows, no padding, no error | 0 |
| A4 | Selected city has none | stdout | "No attractions found for …" message | 0 |
| A5 | Enter a number not in the result list | stdout/stderr | Invalid-selection message, no lookup attempted | (session continues) |
| A6 | Accented city (e.g. "São Paulo") selected | stdout | Resolves correctly and returns its attractions | 0 |
| A7 | Source unreachable / error response | stderr | "unable to fetch attractions: …" | 2 |
