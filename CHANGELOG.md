# Changelog

All notable changes follow Keep a Changelog. The project uses Semantic
Versioning after v1.0.0.

## [Unreleased]

### Added

- Immutable weekly rules, overnight ranges, full-day and closed states.
- Dated replace, add, subtract, and closure exceptions with named sets.
- DST-explicit local resolution and bounded instant/transition queries.
- Union, intersection, subtraction, and authoritative overlay algebra.
- Canonical comparison and bounded, provenance-safe human summaries.
- Strict canonical JSON, SQL/pgx JSONB persistence, adapters, and test helpers.
- `go-calendar` civil-date ownership, bounded zone loading, and fold resolution.
- Explicit DST policy on local queries and injected elapsed observation clocks.
- Structured Location, Track, Postal, and Spatie migration fixtures.
- Transition-waiting guidance using injected `go-clock` timer capabilities.
- Fuzz, race, mutation, coverage, benchmark, documentation, API, security, and
  PostgreSQL automation.

### Fixed

- Reject duplicate exception source revisions even when another priority sorts
  between them.

[Unreleased]: https://github.com/faustbrian/go-opening-hours/commits/main
