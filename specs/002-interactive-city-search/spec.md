# Feature Specification: Interactive City Search CLI

**Feature Branch**: `002-interactive-city-search`

**Created**: 2026-06-26

**Status**: Draft

**Input**: User description: "Interactive City Search CLI — When the user launches the application, they should be greeted with a welcome message that sets the tone for the travel itinerary generator. The terminal then prompts the user to enter a city name. Once the user types a city name and presses Enter, the application queries the city database and displays the top 10 matching cities with their country and region, formatted in a clean, readable way. This transforms the current one-shot CLI command into an interactive terminal experience."

## User Scenarios & Testing *(mandatory)*

### User Story 1 - Guided interactive city search (Priority: P1)

A traveler launches the application without supplying any arguments. They are greeted by a welcome message that frames the tool as a travel itinerary helper, then prompted to type a city name. After typing a name and pressing Enter, they see the top matching cities, each shown with its region and country, so they can identify and choose their destination.

**Why this priority**: This is the core of the feature and the minimum viable product. Without the guided greeting → prompt → results flow, there is no interactive experience to deliver. It transforms the existing one-shot command into a welcoming, self-explanatory tool.

**Independent Test**: Launch the application with no arguments, observe the welcome message and prompt, type a known city name (e.g. "London"), press Enter, and verify a clean, readable list of up to 10 matching cities — each showing name, region (when known), and country — is displayed.

**Acceptance Scenarios**:

1. **Given** the application is launched with no city argument, **When** it starts, **Then** a welcome message and a prompt to enter a city name are displayed before any input is read.
2. **Given** the prompt is shown, **When** the user types a valid city name and presses Enter, **Then** up to 10 matching cities are displayed, ordered with the most relevant (largest) first, each formatted as name, region, and country on its own line.
3. **Given** the user enters a city name that has fewer than 10 matches, **When** results are shown, **Then** only the actual matches are displayed without padding or placeholder rows.
4. **Given** a matching city has no known region, **When** it is displayed, **Then** the region is omitted gracefully and the line remains readable.

---

### User Story 2 - Helpful handling of empty input and no matches (Priority: P2)

When a traveler presses Enter without typing anything, or types a name that matches no cities, the application responds with a clear, friendly message instead of failing silently or crashing, and re-prompts them to try again.

**Why this priority**: Friendly error handling is essential for a welcoming interactive experience, but the tool still delivers core value (P1) without it. It prevents user confusion and dead-ends.

**Independent Test**: At the prompt, press Enter with no text and confirm a friendly "please enter a city name" message and a re-prompt; then type a nonsense string (e.g. "zzzzzz") and confirm a clear "no cities found" message and a re-prompt.

**Acceptance Scenarios**:

1. **Given** the prompt is shown, **When** the user presses Enter without typing a name, **Then** a friendly message asks them to enter a city name and the prompt is shown again.
2. **Given** the prompt is shown, **When** the user types a name with no matches, **Then** a clear message states no cities were found for that query and the prompt is shown again.

---

### User Story 3 - Search multiple destinations in one session (Priority: P3)

After viewing results for one city, a traveler can immediately search for another destination without relaunching the application, and can exit the session cleanly when finished.

**Why this priority**: Repeated searching adds convenience for trip planning across multiple destinations, but a single guided search (P1) already delivers value. This is an enhancement layered on top.

**Independent Test**: Perform one search, confirm the prompt returns for another search, perform a second search, then issue the exit action and confirm the application ends cleanly with a closing message.

**Acceptance Scenarios**:

1. **Given** results for a search have been displayed, **When** the display completes, **Then** the user is prompted to search for another city.
2. **Given** the prompt is shown, **When** the user issues the exit action (an exit keyword or end-of-input signal), **Then** the application displays a brief closing message and terminates with a success status.

---

### Edge Cases

- **Leading/trailing whitespace**: A name typed with surrounding spaces (e.g. " paris ") is treated the same as the trimmed name.
- **Mixed case**: A name typed in any letter case (e.g. "lOnDoN") matches the same cities as the canonical case.
- **Whitespace-only input**: Input consisting only of spaces is treated as empty input (friendly re-prompt), not as a search.
- **Multi-word city names**: Names with spaces (e.g. "san jose") are accepted as a single query without requiring quotes, since the entire input line is the query.
- **End-of-input (Ctrl+D) at the prompt**: The session ends cleanly with a closing message rather than erroring.
- **Unreadable or missing city data**: If the underlying city dataset cannot be loaded, the application reports a clear error and exits with a non-zero status rather than hanging at the prompt.
- **Very large match sets**: When far more than 10 cities match, exactly the top 10 are shown so the list stays scannable.

## Requirements *(mandatory)*

### Functional Requirements

- **FR-001**: On launch with no city provided as an argument, the application MUST display a welcome message that introduces the travel itinerary tool before reading any input.
- **FR-002**: After the welcome message, the application MUST display a prompt inviting the user to enter a city name.
- **FR-003**: The application MUST read a full line of user input as the city query after the prompt.
- **FR-004**: The application MUST trim surrounding whitespace from the entered query before searching.
- **FR-005**: The application MUST match cities by name using case-insensitive matching, consistent with the existing search behavior.
- **FR-006**: The application MUST display at most 10 matching cities for a query, ordered most-relevant first (largest population first, consistent with existing search ranking).
- **FR-007**: Each displayed city MUST show its name, region, and country in a clean, readable single-line format, omitting the region when it is unknown.
- **FR-008**: When the entered query is empty or whitespace-only, the application MUST display a friendly message asking for a city name and MUST NOT perform a search for that input.
- **FR-009**: When a query produces no matches, the application MUST display a clear message indicating no cities were found for that query.
- **FR-010**: The application MUST allow the user to perform another search after results (or a no-match/empty message) are shown, without relaunching.
- **FR-011**: The application MUST provide a clear way to end the session (an exit keyword and/or end-of-input signal) and MUST display a brief closing message on exit.
- **FR-012**: On a clean exit, the application MUST terminate with a success status code; on a data-load failure it MUST terminate with a non-zero status code and a clear error message.
- **FR-013**: When the application is invoked with a city name argument (existing one-shot usage), it MUST continue to behave as a single-query lookup and MUST NOT enter interactive mode.

### Key Entities *(include if feature involves data)*

- **City**: A populated place a user can search for, identified by its name and disambiguated by its region (administrative area) and country. Ranked by population for ordering search results.
- **Search query**: The free-text city name entered by the user at the prompt; normalized (trimmed, case-insensitive) before matching.
- **Interactive session**: A single run of the application in interactive mode, consisting of a welcome, one or more prompt-and-result cycles, and a clean exit.

## Success Criteria *(mandatory)*

### Measurable Outcomes

- **SC-001**: A first-time user, given no instructions, can launch the app and successfully see matching cities for a destination within their first interaction, guided solely by the on-screen welcome and prompt.
- **SC-002**: Every search for a name with matches returns results, capped at 10 cities, each line clearly showing name, region (when present), and country.
- **SC-003**: Empty input and no-match scenarios always produce a friendly, understandable message and never crash or leave the user without guidance.
- **SC-004**: A user can complete searches for at least two different destinations in a single session and then exit cleanly without relaunching the application.
- **SC-005**: Search results for any single query are displayed effectively instantly from the user's perspective (no perceptible wait).

## Assumptions

- The interactive experience is triggered when the application is launched without a city-name argument; supplying a city name preserves the existing one-shot lookup behavior (FR-013).
- The city dataset, matching rules (case-insensitive prefix match), ranking (by population, largest first), and 10-result cap are reused from the existing search implementation rather than redefined.
- The entire input line typed at the prompt is treated as the city query, so multi-word city names do not require quoting.
- "Clean, readable format" reuses the existing single-line "Name, Region, Country" rendering (region omitted when unknown) for visual consistency with the current tool.
- Session continuation (multiple searches per launch) and an exit mechanism are in scope as a P3 enhancement; recognized exit actions include an exit keyword (e.g. "exit"/"quit") and an end-of-input (Ctrl+D) signal.
- Input and output occur over a standard terminal (stdin/stdout), with errors reported to the standard error stream, consistent with the project's CLI principles.
