<!--
SYNC IMPACT REPORT
Version Change: Initial Constitution (1.0.0)
Principles Defined:
  - I. Milestone-Driven Development
  - II. Test-First (TDD - NON-NEGOTIABLE)
  - III. Spike-First Exploration
  - IV. CLI Interface
  - V. Simplicity First (YAGNI)
Sections Added:
  - Core Principles (5 principles)
  - Technology Standards
  - Development Workflow
  - Governance
Templates Status:
  ✅ spec-template.md - Verified alignment with TDD and milestone principles
  ✅ tasks-template.md - Verified alignment with milestone-driven and TDD principles
  ✅ plan-template.md - Constitution Check section ready for principle validation
Follow-up TODOs: None
-->

# Travel Itinerary Generator Constitution

## Core Principles

### I. Milestone-Driven Development

Every feature MUST be decomposed into independent, deliverable milestones.

- Each milestone delivers standalone value and can be tested independently
- Milestones are prioritized (P1, P2, P3...) with P1 being the minimum viable implementation
- Implementation proceeds one milestone at a time - no parallel milestone work
- Each milestone completion requires passing tests before moving to the next
- Progress is measured by completed milestones, not individual tasks

**Rationale**: Milestone-driven development ensures continuous delivery of working functionality, reduces integration risk, and provides clear progress indicators. Users receive value incrementally rather than waiting for complete feature delivery.

### II. Test-First (TDD - NON-NEGOTIABLE)

Test-Driven Development is mandatory for every milestone.

- Tests MUST be written before implementation code
- Tests MUST fail initially (proving they test something)
- Implementation proceeds only after tests are written and approved
- Red-Green-Refactor cycle strictly enforced:
  1. **Red**: Write failing test
  2. **Green**: Write minimal code to pass
  3. **Refactor**: Improve code while keeping tests green
- No milestone is considered complete without passing tests

**Rationale**: TDD ensures code correctness from the start, provides living documentation, enables confident refactoring, and reduces debugging time. The discipline of writing tests first forces clear thinking about requirements and interface design.

### III. Spike-First Exploration

When technical approach is unclear, MUST conduct time-boxed spike before implementation.

- Spike is a time-boxed investigation (typically 1-4 hours)
- Spike goals: evaluate feasibility, compare approaches, identify risks
- Spike output: brief findings document with recommended approach
- Spike code is exploratory and disposable - not production code
- Decision to spike must be explicit and documented in plan.md
- After spike, return to TDD workflow for actual implementation

**Rationale**: Spikes prevent premature commitment to unproven approaches, reduce waste from false starts, and provide evidence-based decision making. Time-boxing prevents analysis paralysis while ensuring adequate exploration.

### IV. CLI Interface

Application MUST provide a clean command-line interface with text-based I/O.

- Primary interface is terminal-based with clear text output
- Input via command-line arguments and flags
- Output to stdout (results) and stderr (errors/logs)
- Human-readable output as default format
- Exit codes follow standard conventions (0=success, non-zero=error)
- Help text and usage examples required for all commands

**Rationale**: CLI interfaces are simple, scriptable, testable, and composable with other tools. Text-based I/O ensures debuggability, automation capability, and integration with standard Unix tooling.

### V. Simplicity First (YAGNI)

Start simple and add complexity only when proven necessary.

- Implement only what is needed for current milestone - no speculative features
- Choose the simplest solution that solves the immediate problem
- Prefer standard library over external dependencies when reasonable
- Avoid premature abstractions - wait for patterns to emerge
- Question complexity: every abstraction must justify its existence
- Refactor toward simplicity during the Green-Refactor cycle

**Rationale**: YAGNI (You Aren't Gonna Need It) reduces code volume, maintenance burden, and cognitive load. Complexity is expensive - every abstraction, dependency, and pattern has a cost. Start simple and evolve based on actual need, not anticipated need.

## Technology Standards

**Language**: Go (latest stable version)

**Dependencies**: Favor standard library; external dependencies require justification

**Project Structure**: Single Go module with clear package organization
- `cmd/` for CLI entry points
- `internal/` for application-specific code
- `pkg/` for potentially reusable libraries (use sparingly)
- `testdata/` for test fixtures

**Testing**: Go standard testing package (`testing`), table-driven tests preferred

**Documentation**: Package documentation via Go doc comments, README for user guide

## Development Workflow

### Pre-Implementation Gate

Before starting any milestone implementation:
1. Milestone scope clearly defined in spec.md with acceptance criteria
2. If technical approach unclear → conduct spike first
3. Tests written and approved
4. Tests run and fail (proving they test the right thing)

### Implementation Process

1. Write tests for milestone (TDD Red phase)
2. Verify tests fail with expected failure message
3. Implement minimal code to make tests pass (TDD Green phase)
4. Refactor for clarity and simplicity (TDD Refactor phase)
5. Verify all tests still pass
6. Update documentation if interfaces changed

### Milestone Completion Gate

Milestone is complete when:
- All tests pass
- Code reviewed for adherence to constitution principles
- Documentation updated
- Committed to version control with clear commit message

## Governance

This constitution supersedes all other development practices and guides all technical decisions.

### Amendment Process

Constitution amendments require:
1. Clear rationale for the change (what problem does it solve?)
2. Impact analysis on existing features and workflows
3. Update to constitution version following semantic versioning:
   - **MAJOR**: Backward-incompatible principle removal or fundamental redefinition
   - **MINOR**: New principle added or existing principle materially expanded
   - **PATCH**: Clarifications, wording improvements, non-semantic refinements
4. Propagation of changes to all dependent templates and documentation

### Compliance

- All feature specifications must reference constitution principles
- All implementation plans must include Constitution Check section
- Code reviews verify adherence to constitution principles
- Complexity must be explicitly justified when deviating from Simplicity First principle

### Development Guidance

For day-to-day development guidance and best practices, refer to Go documentation and this constitution. When in doubt, refer back to the five core principles.

**Version**: 1.0.0 | **Ratified**: 2026-06-26 | **Last Amended**: 2026-06-26
