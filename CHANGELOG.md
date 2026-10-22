# Changelog

All notable changes follow Keep a Changelog. The project uses Semantic
Versioning after v1.0.0.

## 5.0.1 - 2026-10-22

### Changed

- Publish the current source and documentation with a fresh immutable tag.

## [Unreleased]

## [5.0.0] - 2026-10-13

Recorded as a prepared release candidate on this date; publication is a separate
step. Public adoption requires the v5.0.0 tag and release.

### Fixed

- Preserve earlier supported transition candidates when an admitted forward
  search horizon extends beyond civil year 9999, matching reverse search's
  existing year-one boundary handling.

- Keep bounded transition searches, instant ranges and elapsed duration aligned
  with point membership across POSIX timezone-extension year boundaries,
  including the last day of a leap year. Same-offset annual bookkeeping does
  not consume collected-fragment quota. Reject out-of-civil-domain query starts
  before offset traversal, including huge instants on 32-bit hosts; canonical
  storage and explicit local resolution policies are unchanged.
  ([78702b0b0a](https://github.com/faustbrian/go-opening-hours/commit/78702b0b0ad9ab68fb57cda09fe144a9937531bc))

- Restore first/nearest actual availability transitions across timezone-offset
  jumps and repeated fold endpoints, including typed searches and prepared
  indexes. ([9268ec2577](https://github.com/faustbrian/go-opening-hours/commit/9268ec25777270b83a8bdcef8da5c621cbfaca3c))
- Make instant expansion and elapsed duration agree with civil point membership
  across folds and gaps. New York's 2026-11-01 Sunday 01:30–02:30 rule expands
  to two open intervals totaling 90 minutes, not the former contiguous two-hour
  envelope. Consumers relying on that envelope must adopt the point-membership
  policy; explicit local resolution policies and stored canonical bytes remain
  unchanged. ([e1979a85e9](https://github.com/faustbrian/go-opening-hours/commit/e1979a85e99b159f57accfafbcc2b887125b76db))

### Changed

- Prepare `github.com/faustbrian/go-opening-hours/v5` for the breaking
  point-membership expansion policy. Migrate all owned root and adapter imports
  from `/v4` to `/v5` together after publication; public v4 stays unchanged.
  Go 1.27, dependency versions, explicit local resolution and canonical JSON
  version 1 are retained. See the [v5 migration guide](docs/v5-migration.md).
  ([ba3ea5c821](https://github.com/faustbrian/go-opening-hours/commit/ba3ea5c82173085f08cc332ba5d06ce744ba213e))
- Clarify copied schedule ownership, point-membership DST expansion, adapter
  admission, representation equality and synchronous observation limits.
  Stored formats and calendar behavior are unchanged by this documentation-only
  item. ([f9539194e6](https://github.com/faustbrian/go-opening-hours/commit/f9539194e65b56c4eed3d051364d3140658e22fb))

### Documentation

- Improve package adoption navigation and reconcile the independently published
  v4 cohort; this is documentation history, not another v5 runtime correction.
  ([43f3227a1f](https://github.com/faustbrian/go-opening-hours/commit/43f3227a1f116f51ad5caec791f4fbe107a65a41), [72a64ab20d](https://github.com/faustbrian/go-opening-hours/commit/72a64ab20d55d1ff70838f916f919f75503851da))

### Maintenance

- Standardize the risk-based repository development policy without changing
  schedule semantics. ([0474b88589](https://github.com/faustbrian/go-opening-hours/commit/0474b88589c1bb079b5d3ea0dc009e3a339e994b))
- Update verification through patched documentation tooling, bounded public
  Tools v2.1.1, supply-chain Tools v2.3.0 and source-documentation enforcement
  at Tools v2.4.0. Development/CI selects Go 1.27.2 while the module minimum
  stays Go 1.27.0. Ordinary bounded checks and explicit full/release gates remain
  separate; these changes are verification maintenance, not supplier repins or
  runtime fixes. ([f3beee84e8](https://github.com/faustbrian/go-opening-hours/commit/f3beee84e889c09a0c838ccfa8aad9fda1eb39ee), [3d7170435c](https://github.com/faustbrian/go-opening-hours/commit/3d7170435cce222c3770693731741f6e655c4f50), [b312297180](https://github.com/faustbrian/go-opening-hours/commit/b31229718005f68f508694699fb20bd76e457fb5), [46d206ac47](https://github.com/faustbrian/go-opening-hours/commit/46d206ac471a6985105eb0578bec036583890528))

Prepared source accountability covers the 12 non-merge commits from public
`v4.0.0` (`a78aeaf2f39171fe079d14e962f7a4ecf8ffcc12`) through
`78702b0b0ad9ab68fb57cda09fe144a9937531bc`, each linked once above.
The current forward-horizon repair and release alignment are outside that
frozen history range.

## [4.0.0] - 2026-10-06

### Changed

- Prepare `github.com/faustbrian/go-opening-hours/v4` with public Wire v3.0.0.
  Migrate all Opening imports to `/v4` and Wire imports to `/v3` together:
  canonical and retained `WireFormat` constants now use Wire v3's named
  `Format` type. Canonical JSON version 1, strict decoding, schedule behavior,
  and distinct compatibility codecs remain unchanged; see the
  [v4 migration guide](docs/v4-migration.md).

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

[Unreleased]: https://github.com/faustbrian/go-opening-hours/compare/e1979a85e99b159f57accfafbcc2b887125b76db...HEAD
[5.0.0]: https://github.com/faustbrian/go-opening-hours/compare/v4.0.0...e1979a85e99b159f57accfafbcc2b887125b76db
[1.1.0]: https://github.com/faustbrian/go-opening-hours/compare/v1.0.0...v1.1.0
[1.0.0]: https://github.com/faustbrian/go-opening-hours/releases/tag/v1.0.0
