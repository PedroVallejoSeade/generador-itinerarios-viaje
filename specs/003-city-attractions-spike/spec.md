# Feature Specification: City Attractions Lookup Spike

**Feature Branch**: `003-city-attractions-spike`

**Created**: 2026-06-26

**Status**: Draft

**Input**: User description: "City Attractions Lookup Spike — As a traveler who has selected a city from the search results I want to see the top 10 most known attractions for that city so that I can discover the must-see sights, museums, and landmarks for my chosen destination. After the user selects a city by entering its number from the search results, the application should display the top 10 most famous attractions for that city (landmarks, museums, sightseeing spots). This is a spike to research and evaluate the best tool to power this feature: a free no-auth public API (preferred), a Go library, or an open dataset queried locally. Evaluate each option on data quality, ease of integration in Go, authentication requirements, rate limits, and result relevance, and document the recommended approach with a brief proof of concept."

## Clarifications

### Session 2026-06-26

- Q: Must the attractions lookup work fully offline, or is network connectivity acceptable? → A: Network access is acceptable; the recommended source MAY require internet connectivity (e.g., a public API).
- Q: Which identity should the attractions lookup use to resolve the selected city? → A: Resolve by name + country (plus region/state when available) to disambiguate same-named cities.
- Q: What detail should be displayed per attraction? → A: Numbered list of attraction names, plus category/short description when readily available from the source.

## User Scenarios & Testing *(mandatory)*

### User Story 1 - See top attractions for a selected city (Priority: P1)

A traveler has searched for a city and selected it by entering its number from the result list. The application then displays the top 10 most well-known attractions for that city — landmarks, museums, and sightseeing spots — so the traveler can start discovering what to visit at their destination.

**Why this priority**: This is the core value of the feature. Without showing attractions for the chosen city, the selection step has no payoff. It is the minimum viable slice that turns "I picked a city" into "here is what to see there."

**Independent Test**: Can be fully tested by selecting a well-known city (for example, New York City) and confirming the application returns a list of recognizable attractions (such as the Statue of Liberty, Central Park, Times Square).

**Acceptance Scenarios**:

1. **Given** the user has selected a city by its number, **When** the application looks up attractions, **Then** it displays a list of attractions for that city in human-readable text on standard output.
2. **Given** a well-known city is selected, **When** attractions are returned, **Then** the list contains recognizable landmarks, museums, or sightseeing spots associated with that city.
3. **Given** more than ten attractions exist for the selected city, **When** results are returned, **Then** the application displays at most the top 10 most well-known ones.

---

### User Story 2 - Evaluate and recommend an attractions data source (Priority: P1)

As the team building the itinerary generator, we need a documented evaluation of candidate ways to obtain city attraction data — a free no-authentication public API, a Go library, or an open dataset queried locally — so we can commit to one with confidence. The spike compares the options against the project's constraints and produces a recommendation backed by a brief proof of concept.

**Why this priority**: This is a spike-first exploration (per the project constitution). The technical approach for sourcing attraction data is currently unproven, so an explicit, time-boxed evaluation must precede production implementation. The recommendation plus proof of concept is the primary deliverable of this spike.

**Independent Test**: Can be fully tested by reviewing the spike findings document and confirming it compares at least three candidate options against the stated criteria, names a recommended choice with rationale, and references a working proof of concept.

**Acceptance Scenarios**:

1. **Given** the spike is complete, **When** the findings document is reviewed, **Then** it evaluates at least three candidate options (covering at least one public API and one local dataset or library) against data quality, ease of Go integration, authentication requirements, rate limits, and result relevance.
2. **Given** the spike is complete, **When** the findings document is reviewed, **Then** it names a single recommended approach with a clear rationale and any identified risks or limitations.
3. **Given** the recommended approach is identified, **When** the proof of concept is run for a sample city, **Then** it returns a plausible list of attractions, demonstrating feasibility.

---

### Edge Cases

- **No attractions found**: When the selected city has no attractions available from the chosen source, the application displays a clear "no attractions found" message rather than failing or showing an empty, unexplained response.
- **Fewer than ten attractions**: When a city has fewer than ten known attractions, the application displays all available ones without padding or error.
- **Invalid selection**: When the user enters a number that does not correspond to any city in the result list, the application reports the invalid selection and does not attempt a lookup.
- **Data source unavailable**: When the underlying data source cannot be reached or returns an error, the application reports the problem on standard error and exits with a non-zero status, without crashing.
- **Ambiguous or sparse city data**: When the selected city is small or lesser-known, the application returns whatever relevant attractions are available and clearly indicates if the list is short.
- **Non-ASCII / accented city names**: When the selected city name contains accents or non-Latin characters (for example, "São Paulo", "Zürich"), the attractions lookup still resolves to the correct city.

## Requirements *(mandatory)*

### Functional Requirements

- **FR-001**: System MUST accept the user's selection of a city by its number from the previously displayed search results.
- **FR-002**: System MUST look up attractions for the selected city using the chosen attraction data source.
- **FR-002a**: System MUST resolve the selected city for the attractions lookup using its name plus country (and region/state when available) to disambiguate same-named cities, and MUST correctly handle accented or non-Latin city names.
- **FR-003**: System MUST return a list of attractions for the selected city, ranked so the most well-known attractions appear first.
- **FR-004**: System MUST limit the displayed attractions to a maximum of 10 (the top 10 most well-known) when more are available.
- **FR-005**: System MUST present the attractions as human-readable text on standard output, formatted as a numbered list of attraction names, including a category and/or short description for an attraction when that supporting context is readily available from the chosen source.
- **FR-006**: System MUST display a clear, explanatory message when no attractions are found for the selected city.
- **FR-007**: System MUST display all available attractions when a city has fewer than 10, without error.
- **FR-008**: System MUST reject an invalid city selection (a number not present in the result list) with a clear message and without performing a lookup.
- **FR-009**: System MUST handle data-source or connectivity errors gracefully, reporting them on standard error and exiting with a non-zero status code.
- **FR-010**: The selected data source SHOULD be usable without authentication, API keys, or account registration (preferred, per the Simplicity First principle); if the recommended option requires authentication, the spike MUST justify the trade-off.
- **FR-011**: The selected data source MUST be free to use for this application's purpose.
- **FR-012**: The spike MUST evaluate at least three candidate options for sourcing attraction data, covering at least one public API and at least one local dataset or Go library, against these criteria: data quality, ease of integration in Go, authentication requirements, rate limits, and result relevance.
- **FR-013**: The spike MUST produce a findings document recording the comparison, a single recommended approach, the rationale, and any identified risks or limitations.
- **FR-014**: The spike MUST include a brief proof of concept that retrieves attractions for at least one sample city using the recommended approach.

### Key Entities *(include if feature involves data)*

- **City Selection**: The city the user chose from the search results, identified by its position/number in that list. Carries enough identity (name, country, and region/state when available) to resolve the correct city for the attractions lookup; resolution is keyed on name + country (plus region/state when available) to disambiguate same-named cities.
- **Attraction**: A notable place to visit in a city — a landmark, museum, or sightseeing spot. Key attributes: name and a measure of prominence/popularity used for ranking. May carry optional supporting context (category, short description, location) where readily available.
- **Attraction Result Set**: The ordered collection of attractions returned for a selected city, ranked by prominence and bounded to a maximum display size of 10.

## Success Criteria *(mandatory)*

### Measurable Outcomes

- **SC-001**: For a representative set of well-known cities, at least 90% of lookups return at least 5 recognizable, relevant attractions for the city.
- **SC-002**: After selecting a city, the attractions list is displayed to the user in under 3 seconds for a typical lookup.
- **SC-003**: A user can identify at least 3 "must-see" attractions for a well-known selected city from the displayed list without needing any external reference.
- **SC-004**: A new user can view a city's attractions without creating any account, obtaining any API key, or performing any authentication setup (for the recommended no-auth approach).
- **SC-005**: The spike concludes with a findings document that compares at least three candidate options and names one recommended approach with rationale.
- **SC-006**: A reviewer can determine the recommended approach and the reason for the choice in under 5 minutes by reading the findings document.
- **SC-007**: The proof of concept successfully returns attractions for at least one sample city, demonstrating the recommended approach is feasible in Go.

## Assumptions

- **Trigger**: The attractions lookup begins after the city-selection step (entering a number from the search results) introduced by the prior interactive city-search feature; this spike assumes that selection mechanism already exists and produces a resolvable city.
- **Ranking basis**: "Top 10 most known" is interpreted as ranking attractions by a prominence/popularity signal available from the chosen source (for example, number of reviews, ratings, or a built-in popularity rank); the exact signal depends on the source selected during the spike.
- **Attraction types**: Landmarks, museums, and sightseeing spots are all in scope; restaurants, hotels, and other non-sightseeing categories are out of scope for this spike.
- **Data source type**: A free no-authentication public API is preferred per Simplicity First, but a Go library or a bundled/local open dataset is acceptable if it better satisfies data quality and relevance; the spike determines the best fit. Network/internet connectivity is acceptable: the recommended source MAY require online access (full offline operation is not a requirement).
- **Result cap**: The application displays at most 10 attractions to keep terminal output readable.
- **Scope boundary**: Building a full itinerary, persisting selected attractions, mapping/routing between attractions, and any booking functionality are out of scope for this spike — its scope ends at returning and displaying the top attractions plus the data-source recommendation and proof of concept.
- **Interface**: Interaction is terminal/CLI-based with text input and output, consistent with the project's CLI Interface principle.
- **Spike time-box**: The data-source evaluation (US2 / FR-012, FR-013, FR-014) is time-boxed per the constitution's Spike-First Exploration principle; the findings document captures whatever comparison, recommendation, and proof of concept is reached within that box.
- **Spike disposability**: Code produced during the spike is exploratory and may be discarded; the durable deliverables are the working proof of concept and the documented recommendation, after which the team returns to a test-first implementation workflow.
