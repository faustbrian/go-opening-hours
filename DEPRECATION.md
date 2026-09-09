# Deprecation Policy

Deprecations MUST identify the replacement, reason, migration steps, and
earliest removal version. Public Go identifiers use a valid `Deprecated:` doc
paragraph and corresponding changelog entry.

At `v1` and later, a supported replacement SHOULD exist for at least one minor
release before removal. Security or correctness defects MAY require faster
removal when continued support would be unsafe; the release notes must explain
the exception.

Silent behavior changes, undocumented aliases, and indefinite deprecated code
are prohibited. Deprecations are checked during compatibility and release
review.

## Opening-hours integration paths

`openinghourscalendar`, `openinghoursconfig`, `openinghourstemporal`,
`openinghoursvalidation`, and `openinghourswire` are supported compatibility
facades for the corresponding `adapters/<target>` packages introduced in
v1.1.0. New consumers should use the target-oriented paths.

The earliest permitted removal is the later of 180 days after v1.1.0 public
availability and two subsequently published stable minor releases containing
the successors. Removal also requires an authorized next major, migrated owned
consumers, a fresh public-usage audit, and clean external-consumer proof.
