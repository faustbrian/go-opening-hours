# Compatibility

The minimum supported compiler is Go 1.26.6, the latest stable release at
implementation time. Public API changes are checked against `api/baseline.txt`.
Breaking changes require a new major release; v1 follows
Semantic Versioning.

Canonical wire version is `1`. Decoders reject other versions. Canonical bytes
are stable for semantically identical construction order within the same wire
version. Timezone results depend on the installed IANA database, so historical
rule updates may change instant expansion without changing stored bytes.

PostgreSQL JSONB is tested on maintained server releases 14 through 18. pgx v5
uses the root scanner/valuer directly. Location, Track, Postal, and Spatie
fixture compatibility is documented in [legacy migration](legacy-migration.md).
Track and Postal use the shared Location representation; no unverified provider
schema is claimed.

## Adapter migration interval

The canonical integration entry points are `adapters/calendar`,
`adapters/config`, `adapters/temporal`, `adapters/validation`, and
`adapters/wire`. Their released top-level predecessors remain supported for
the longer of 180 days after v1.1.0 becomes publicly available and two later
stable minor releases containing the successors. Both paths preserve behavior
and shared sentinels during that interval. Named types keep path-specific
reflection identity, so consumers should migrate imports and type references
together rather than mixing the two named types.

| Fixture | Imported contract | Proven behavior |
| --- | --- | --- |
| Location | weekday `{from,to}` slots | split ranges and explicit closure |
| Track | shared Location slots | start-inclusive/end-exclusive boundary |
| Postal | shared Location slots | owner-day overnight spill and closing |
| Spatie | strict `HH:MM-HH:MM` arrays | weekday opening and closed empty day |
