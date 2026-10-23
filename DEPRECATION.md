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

Version 6 removed the `openinghourscalendar`, `openinghoursconfig`,
`openinghourstemporal`, `openinghoursvalidation`, and `openinghourswire`
compatibility facades. Import the corresponding `adapters/<target>` package.
The removal precedes the announced support window; the v5 line keeps the
facades.
