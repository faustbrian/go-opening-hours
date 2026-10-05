# Changelog

All notable changes follow Keep a Changelog. The project uses Semantic
Versioning after v1.0.0.

## [Unreleased]

## [3.0.0] - 2026-10-05

### Changed

- Adopt `github.com/faustbrian/go-opening-hours/v3` with Calendar v2,
  Validation v2 and public Temporal v2 dependencies. Migrate date,
  interval and validator imports together across canonical and retained
  adapters. Schedule semantics, resource limits and wire version 1 remain
  unchanged. Historical Opening v1 and v2 module graphs remain independently
  usable.

## [2.0.2] - 2026-10-05

### Changed

- Upgrade pgx from v5.10.0 to v5.11.0 while retaining native schedule
  JSONB mapping and nullable persistence contracts. Applications
  implementing pgx Rows must provide its new TypeMap method; upstream
  connection-string and date handling changes also apply to the selected
  pgx version.

- Upgrade go-clock from v1.0.0 to v1.2.0 while retaining the narrow
  injected current-time and elapsed-time capability contracts. The
  selected dependency adds manual-clock Close, deprecates Shutdown, and
  changes observed sleep outcomes, closed-clock sleep precedence, and
  ticker observations for applications using those features directly.
  Its Go 1.27.0 minimum matches this module's existing requirement.

- Upgrade go-config from v1.0.0 to v1.1.0 while retaining canonical
  and legacy configuration-value adapters and their unchanged decoder
  contract. Upstream adds the service adapter entry point and retains
  configservice as its compatibility facade.

- Adopt public `go-config/v2` v2.0.0 for canonical and legacy configuration
  integration tests, retaining schedule decoding and atomic validation
  assertions. Production adapters keep their structural value-unmarshal seam.

- Upgrade go-wire from v1.0.0 to v1.0.1 while retaining the public
  typed format identity and canonical and legacy schedule codecs.
  The selected module graph updates CBOR to v2.9.3 and MongoDB driver
  to v2.9.0; those optional codecs are not used by these adapters.
  The supplier's Go 1.27.0 minimum matches this module's requirement.

## [2.0.1] - 2026-10-01

### Released

- Publish v2.0.1 as the validated v2 release. The exposed v2.0.0 Git
  tag remains immutable but did not complete release validation and has no
  verified stable release. The v2 migration and runtime contracts below remain
  unchanged.
- The release source is `6ecdac5c7e268abac3ccc609c74711c9a316a2b0`.
  Source CI and release rehearsal passed; the clean public-proxy consumer
  compiled all 15 owned packages and passed focused public behavior tests.
  Published assets include the source module ZIP, module file, dependency
  SBOM, maintainer-local provenance, signed checksums, and signer key.
  The SBOM covers runtime module dependencies, excluding test dependencies,
  the standard library, and license detection; provenance is not a CI
  attestation or a claim of SLSA build assurance.

## [2.0.0] - 2026-10-01

Tagged snapshot only; release validation did not complete. No verified stable
GitHub release was published for this tag.

### Changed

- Move the root module and all owned package imports to
  `github.com/faustbrian/go-opening-hours/v2`. Exception source and revision
  identities must now be nonempty, valid UTF-8 within the existing 128-byte
  limits; malformed byte strings previously accepted by constructors are
  rejected with `CodeInvalidState`. See [v2 migration](docs/v2-migration.md).

### Added

- Expose the existing 64-range and 4,096-exception package limits so adapters
  and callers can reject oversized collections before conversion.

### Security

- Reject oversized SQL, configuration, exception, Spatie, Temporal, and
  holiday-expansion inputs before package-owned copies or conversions.
- Preserve constructor validation and first-invalid exception-set precedence
  while checking the aggregate exception budget before exception copies.
- Reject invalid zero-value exceptions in direct schedule input, matching
  named exception-set admission.
- Preserve exception identity through canonical and SQL persistence by
  rejecting malformed provenance before construction or collection admission.
- Keep canonical JSON round trips and compiled queries valid at the full
  16-level composition limit, including ranged weekly and exception rules;
  compilation now propagates decoding errors rather than returning a closed
  index after a failed decode.
- Treat typed-nil and panicking injected clocks as invalid or unmeasured so
  optional clock capabilities cannot suppress schedule query results.
- Document checksum-database enforcement and the remaining caller-owned input,
  SQL cancellation, and in-memory composition-shape boundaries.

## [1.1.0] - 2026-09-09

### Changed

- Add canonical `adapters/calendar`, `adapters/config`, `adapters/temporal`,
  `adapters/validation`, and `adapters/wire` entry points. The released
  top-level integration packages remain supported compatibility facades with
  unchanged signatures, behavior, sentinels, and named-type identities.
- Adopt Calendar v1.1.0, Temporal v1.1.0, and Validation v1.1.0 while
  retaining Config v1.0.0 and Wire v1.0.0.

### Deprecated

- Prefer the five target-oriented `adapters/<target>` packages for new code.
  The former integration paths remain supported for the longer of 180 days
  after successor public availability and two later stable minor releases.

- Adopt the `go-library-tools` v1.6.2 cohesion contract and local
  `make cohesion` gate without changing opening-hours API or runtime behavior.
- Pin reusable CI to the immutable v1.6.2 workflow, including proportional
  zero-mutant verification, and enforce cohesion metadata in the repository's
  required CI contract.
- Reconcile owned v1.0.0 dependency checksums with their transparency-log
  authenticated public module archives without changing dependency versions.

- Adopt released `go-library-tools` v1.0.5 for local and GitHub Actions
  verification while preserving package-owned evidence and fixtures.
- Replace copied repository tooling with the strict `.golib.yaml` contract and
  the centralized immutable CI workflow.

### Documentation

- Publish the module's family, capabilities, ownership, lifecycle, supported
  environments, package selection, and delivery status, and link package
  documentation to the immutable v1.6.2 ecosystem guidance.

- Document the standalone verification commands and shared gate behavior.

## [1.0.0] - 2026-08-26

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

[Unreleased]: https://github.com/faustbrian/go-opening-hours/compare/v1.1.0...HEAD
[1.1.0]: https://github.com/faustbrian/go-opening-hours/compare/v1.0.0...v1.1.0
[1.0.0]: https://github.com/faustbrian/go-opening-hours/releases/tag/v1.0.0
