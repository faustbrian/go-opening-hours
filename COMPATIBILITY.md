# Compatibility Policy

The repository has one releasable root Go module. Its packages, including
canonical adapters and compatibility facades, ship under root `v<version>`
tags and follow semantic versioning together.

Current source prepares v2.0.0 using the required `/v2` module/import suffix
on main, without a version-specific source directory. Published v1.1.0 remains
available under the original module path. V2 intentionally rejects malformed
UTF-8 exception source/revision identities previously accepted as structured
input; see [migration guidance](docs/v2-migration.md).

Before `v1`, minor releases MAY contain reviewed breaking changes, but every
break MUST be documented with migration guidance. Patch releases MUST remain
backward compatible. At and after `v1`, incompatible exported API or documented
behavior changes require a new major version.

Compatibility includes exported Go APIs, error classification, serialization,
protocol behavior, persistence schemas, environment variables, command output,
resource ownership, ordering, retry/idempotency semantics, and documented
defaults. A compile-compatible change can still be behaviorally breaking.

Specification-backed modules MUST NOT diverge from their declared standards.
Ambiguities require documented decisions and stable tests. Deprecated APIs
follow [`DEPRECATION.md`](DEPRECATION.md).

The `openinghourscalendar`, `openinghoursconfig`, `openinghourstemporal`,
`openinghoursvalidation`, and `openinghourswire` package paths remain supported
compatibility facades for their `adapters/<target>` successors. The facades
preserve released signatures, named-type identity, shared sentinels, encoding,
and error behavior for the interval defined in [`DEPRECATION.md`](DEPRECATION.md).
