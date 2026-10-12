# Contributing

## Before Editing

1. Read [`AGENTS.md`](AGENTS.md) and the affected module's docs.
2. Run standalone `go test ./...` or a focused affected-package test.
3. Identify owned dependencies and reverse dependants in `modules.json`.
4. Preserve unrelated work and generated/corpus provenance.

## Changes

Keep commits focused and conventional. Update every affected changelog with
the behavior and migration impact. Public API changes require compatibility
evidence and documentation. Specification behavior requires a decision record,
fixture coverage, and interoperability evidence.

New direct dependencies and dependency updates must follow the
[dependency governance policy](AGENTS.md#dependencies-and-supply-chain). Package-local
update bots are forbidden; the root policy owns every module and action update.

Specification-backed changes must follow the
[specification governance contract](AGENTS.md#design), update
the affected stable decision entries, and complete the Specification Decisions
section of the pull request template. An unresolved interpretation or stale
source pin is release-blocking; peer behavior cannot silently select policy.

Full scheduled and release mutation gates retain the configured mutation
contract. Local mutation is optional, not a universal contribution prerequisite.

Do not add package-local workflows, permanent replacements, machine-specific
paths, bypass flags, broad mutation exclusions, or aggregate quality metrics
that hide a failing package.

## Verification

Use Go 1.27.2 for development and CI; the module minimum remains Go 1.27.0.
Standalone `go test ./...` runs deterministic tests, fuzz seeds and saved corpus
regressions without sibling checkouts or ecosystem ledgers. It does not actively
explore new fuzz inputs. Untagged Go tests exclude live PostgreSQL integration.

With the checksum-pinned public verifier selected by `.golib.yaml`, use:

```bash
make inventory
make check
make ci
```

Both Make routes use bounded local checks. The manifest still requires
PostgreSQL and the `integration` test tag, so the service-aware Make route may
require Docker and is not equivalent to untagged Go tests. Base-free local
checks run deterministic seeds/corpora; changed-source CI selects bounded active
smoke for causally affected targets. Ordinary CI consumes the public verifier.
One existing `compile` race assertion checks 32 concurrent readers against exact
open/local/transition results and retained canonical schedule identity. It runs
once with a two-minute timeout; no broad stress campaign is added.

Short active exploration selects one target; optional local mutation selects
one package:

```bash
GOMAXPROCS=2 go test -run '^$' -fuzz '^FuzzLocalResolution$' -fuzztime=2s -parallel=1 -timeout=2m .
golib mutation --module . --package ./compile
```

Timezone, DST/fold, overnight, composition and search changes retain affected
boundary and differential assertions. PostgreSQL codec or persistence changes
also require a disposable database and tagged live round-trip checks:

```bash
POSTGRES_URL='<disposable test database URL>' go test -tags=integration -count=1 -timeout=2m ./postgres
```

`make check-full` retains full assurance, the DST/fold/timezone operation, all
eight native 10000x fuzz budgets and benchmarks. Canonical and legacy adapters
remain packages in this module; preserve their API identity and compatibility
checks along with affected direct consumers. Relevant release qualification is
`golib release dry-run --module .`. The public verifier owns disposable caches,
service fixtures and gate thresholds. Report unavailable, skipped and failed
checks honestly; local tests are not live-database, hosted-CI or public-release
qualification.

## Adding A Module

Follow [repository structure policy](AGENTS.md#repository-structure). New modules
require an explicit purpose, ownership boundary, dependency review, package
catalog entry, full quality gates, documentation, changelog, license, security
policy, compatibility plan, and release dry-run.

## Source Documentation

Install the checksum-pinned public verifier declared in `.golib.yaml`.
`make source-docs-check` validates objective Go comment, marker, link and
generated-producer structure; `make check`, `make check-full` and `make ci`
require it before their existing checks. Shared CI runs the same strict
contract and includes its result in Required. This is structural enforcement,
not a substitute for source-truth review or behavioral verification.
