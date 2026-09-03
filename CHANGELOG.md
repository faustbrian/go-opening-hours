# Changelog

All notable changes follow Keep a Changelog. The project uses Semantic
Versioning after v1.0.0.

## [Unreleased]

### Changed

- Adopt the `go-library-tools` v1.3.0 schema-v2 cohesion contract and local
  `make cohesion` gate without changing opening-hours API or runtime behavior.
- Pin reusable CI to the immutable v1.3.0 workflow and enforce cohesion
  metadata in the repository's required CI contract.

- Adopt released `go-library-tools` v1.0.5 for local and GitHub Actions
  verification while preserving package-owned evidence and fixtures.
- Replace copied repository tooling with the strict `.golib.yaml` contract and
  the centralized immutable CI workflow.

### Documentation

- Publish the module's family, capabilities, ownership, lifecycle, supported
  environments, package selection, and delivery status, and link package
  documentation to the immutable v1.3.0 ecosystem guidance.

- Document the standalone verification commands and shared gate behavior.

## [1.0.0] - 2026-08-25

### Fixed

- Bind the reviewed zero-mutant `openinghourswire` delegation package to its
  exact standalone source identity.

### Changed

- Exclude intentional nested modules from root local-proxy archives so local,
  bootstrap, CI, and public module checksums describe the same source
  boundary.

- Track the pinned documentation-tool lockfile so clean CI checkouts install
  the exact validated cspell dependency.

- Reconcile standalone dependency checksums against deterministic current
  module archives so CI, local verification, and release consumers resolve
  identical content.

- Harden standalone documentation validation with deterministic spelling and
  link checks, package-specific documentation gates, and repository-local
  contributor guidance.

### Documentation

- Link the package README to package-owned documentation.

### Changed

- Publish the module from its standalone `github.com/faustbrian/go-opening-hours` identity while preserving its documented API and behavior.
- Refresh local `v0.0.0` owned-module checksums after dependency manifests and
  release notes were normalized; runtime behavior and public APIs are
  unchanged.
- Require owned sibling modules at local `v0.0.0`; clean external consumers
  pin each module to an exact main pseudo-version.

- Refresh owned-module checksums against the final consolidated archives.
- Refreshed the generated API baseline with the current Go documentation
  formatter without changing exported declarations.
- Normalized standalone module metadata against the canonical owned dependency
  graph, including complete checksums for clean consumer resolution.

### Added

- Immutable weekly rules, overnight ranges, full-day and closed states.
- Dated replace, add, subtract, and closure exceptions with named sets.
- DST-explicit local resolution and bounded instant/transition queries.
- Union, intersection, subtraction, and authoritative overlay algebra.
- Canonical comparison and bounded, provenance-safe human summaries.
- Strict canonical JSON, SQL/pgx JSONB persistence, adapters, and test helpers.
- `calendar` civil-date ownership, bounded zone loading, and fold resolution.
- Explicit DST policy on local queries and injected elapsed observation clocks.
- Structured Location, Track, Postal, and Spatie migration fixtures.
- Transition-waiting guidance using injected `clock` timer capabilities.
- Fuzz, race, mutation, coverage, benchmark, documentation, API, security, and
  PostgreSQL automation.
- Pairwise algebra/conservation properties, exception permutation proof, and a
  broad differential against Go timezone rules.
- A disposable mutation runner with machine-readable evidence, zero-error
  enforcement, and a blocking minimum score.

### Fixed

- Report overnight-spill provenance only when the queried point is within the
  preceding day's spill, rather than whenever any spill exists on that date.
- Reject duplicate exception source revisions even when another priority sorts
  between them.
- Replace unreachable owned-module pseudo-versions with published revisions so
  clean checkouts can reproduce every gate without local replacements.

[Unreleased]: https://github.com/faustbrian/go-opening-hours/compare/v1.0.0...HEAD
[1.0.0]: https://github.com/faustbrian/go-opening-hours/releases/tag/v1.0.0
