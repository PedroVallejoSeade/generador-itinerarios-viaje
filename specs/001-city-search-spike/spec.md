# Feature Specification: City Search Spike

**Feature Branch**: `001-city-search-spike`

**Created**: 2026-06-26

**Status**: Draft

**Input**: User description: "City Search Spike — As a travel itinerary app user I want to search for a city by name and see matching results so that I can select the correct city when multiple cities share a similar name. When a user types a city name in the terminal, the application should query a public city database and return a list of matching cities with enough context (country, region) to distinguish between them. This spike will research and evaluate the best public API or dataset to power this city search functionality in a Go terminal application, focusing on simplicity, no authentication requirements, and fast response times."

## User Scenarios & Testing *(mandatory)*

### User Story 1 - Find a city by name (Priority: P1)

A traveler launches the application from the terminal, types a city name (for example, "Paris"), and receives a list of cities whose names match what they typed. This is the foundational capability the rest of the itinerary generator will build upon.

**Why this priority**: Without the ability to look up a city, no downstream itinerary functionality is possible. This is the minimum viable slice that delivers value on its own — a working, demonstrable city lookup.

**Independent Test**: Can be fully tested by running the application, entering a well-known city name, and confirming that the matching city appears in the returned list. Delivers value as a standalone city-lookup tool.

**Acceptance Scenarios**:

1. **Given** the application is running, **When** the user enters a city name that exists in the data source, **Then** the application displays a list containing that city.
2. **Given** the application is running, **When** the user enters a city name, **Then** each result is shown in human-readable text on standard output.

---

### User Story 2 - Disambiguate between similarly named cities (Priority: P1)

A traveler searches for a name shared by several places (for example, "Springfield" or "San José"). The application returns each match along with enough context — at minimum its country and region/state — so the traveler can confidently identify and select the one they intended.

**Why this priority**: Disambiguation is the core problem this feature exists to solve. A bare list of identical names provides no value; the contextual detail is what makes the search usable. It is bundled at P1 because the search result format must carry this context from the first deliverable.

**Independent Test**: Can be fully tested by searching for a name known to exist in multiple countries/regions and verifying that each result is distinguishable by its country and region/state.

**Acceptance Scenarios**:

1. **Given** multiple cities share the searched name, **When** results are returned, **Then** each result includes its country and region/state.
2. **Given** two results share both name and country, **When** results are returned, **Then** the region/state value differs so the user can tell them apart.

---

### User Story 3 - Evaluate and recommend a city data source (Priority: P1)

As the team building the itinerary generator, we need a documented evaluation of candidate public city data sources (APIs or datasets) so we can commit to one with confidence. The spike compares options against the project's constraints — no authentication, free to use, fast responses, and sufficient disambiguation data — and produces a recommendation.

**Why this priority**: This is a spike-first exploration (per the project constitution). The technical approach for sourcing city data is currently unproven, so an explicit, time-boxed evaluation must precede production implementation. The recommendation is the primary deliverable of this spike.

**Independent Test**: Can be fully tested by reviewing the spike findings document and confirming it compares at least three candidate sources against the stated criteria and names a recommended choice with rationale.

**Acceptance Scenarios**:

1. **Given** the spike is complete, **When** the findings document is reviewed, **Then** it lists at least three candidate data sources evaluated against the constraints (no authentication, free, response time, disambiguation coverage).
2. **Given** the spike is complete, **When** the findings document is reviewed, **Then** it names a single recommended data source with a clear rationale and any identified risks.

---

### Edge Cases

- **No matches**: When the entered name matches no city in the data source, the application displays a clear "no results found" message rather than failing or showing an empty, unexplained response.
- **Empty input**: When the user enters a blank or whitespace-only query, the application prompts for a valid city name instead of querying the data source.
- **Partial / misspelled names**: When the user enters a partial name (for example, "Franc"), the application returns reasonable matches where the entered text is contained in or prefixes a city name.
- **Data source unavailable**: When the underlying data source cannot be reached or returns an error, the application reports the problem on standard error and exits with a non-zero status, without crashing.
- **Large result sets**: When a common name matches many cities, the application limits the displayed results to a manageable number so the output stays readable.
- **Non-ASCII / accented names**: When a city name contains accents or non-Latin characters (for example, "São Paulo", "Zürich"), the search still returns the expected match.

## Requirements *(mandatory)*

### Functional Requirements

- **FR-001**: System MUST accept a city name as text input from the user via the terminal.
- **FR-002**: System MUST query a public city data source for entries matching the entered name.
- **FR-003**: System MUST return a list of matching cities, supporting case-insensitive and partial-name matches.
- **FR-004**: Each result MUST include disambiguating context, at minimum the city's country and region/state.
- **FR-005**: System MUST present results as human-readable text on standard output.
- **FR-006**: System MUST limit the number of displayed results to a manageable count when many cities match.
- **FR-007**: System MUST display a clear, explanatory message when no cities match the query.
- **FR-008**: System MUST reject empty or whitespace-only queries with a prompt for valid input, without querying the data source.
- **FR-009**: System MUST handle data-source or connectivity errors gracefully, reporting them on standard error and exiting with a non-zero status code.
- **FR-010**: The selected data source MUST be usable without authentication, API keys, or account registration.
- **FR-011**: The selected data source MUST be free to use for this application's purpose.
- **FR-012**: The spike MUST evaluate at least three candidate public city data sources against the project constraints (no authentication, free, response speed, disambiguation coverage).
- **FR-013**: The spike MUST produce a findings document recording the comparison, a single recommended data source, the rationale, and any identified risks or limitations.

### Key Entities *(include if feature involves data)*

- **City**: A populated place that a user can search for and select. Key attributes: name, country, region/state (province/administrative division). May also carry supporting context such as population or coordinates where available, used only to aid disambiguation.
- **Search Query**: The text the user enters to find a city. Key attributes: the raw input string and its normalized form (trimmed, case-folded) used for matching.
- **Search Result Set**: The ordered collection of cities returned for a query, bounded to a manageable display size.

## Success Criteria *(mandatory)*

### Measurable Outcomes

- **SC-001**: A user can distinguish between at least five cities sharing the same name within a single search result, using the country and region/state shown for each.
- **SC-002**: For a typical single-name query, results are displayed to the user in under 2 seconds.
- **SC-003**: For a representative set of well-known city searches, at least 95% return the intended city within the displayed result list.
- **SC-004**: A new user can run a city search successfully without creating any account, obtaining any API key, or performing any authentication setup.
- **SC-005**: The spike concludes with a findings document that compares at least three candidate data sources and names one recommended source with rationale.
- **SC-006**: A reviewer can determine the recommended data source and the reason for the choice in under 5 minutes by reading the findings document.

## Assumptions

- **Match strategy**: Matching is case-insensitive and includes prefix/substring matches on the city name; exact-only matching is not required for the spike. This was chosen as a reasonable default because the user's goal is to find a city even with partial input.
- **Result cap**: When many cities match, the application displays a bounded number of results (assumed up to 10, ordered by relevance such as population or name closeness) to keep terminal output readable.
- **Disambiguation fields**: Country and region/state are sufficient context to distinguish cities for this spike; finer detail (e.g., county, coordinates) is optional and used only if readily available from the chosen source.
- **Data source type**: Either a queryable public API or a bundled/downloadable open dataset is acceptable; the spike will determine which better satisfies the no-auth, free, and fast-response constraints.
- **Scope boundary**: Selecting a city, persisting selections, and any downstream itinerary generation are out of scope for this spike — its scope ends at returning and displaying disambiguated matches plus the data-source recommendation.
- **Interface**: Interaction is terminal/CLI-based with text input and output, consistent with the project's CLI Interface principle.
- **Spike disposability**: Code produced during the spike is exploratory and may be discarded; the durable deliverables are the working demonstration of city search and the data-source recommendation, after which the team returns to a test-first implementation workflow.
