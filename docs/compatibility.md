# Compatibility

[Documentation map](index.md) · [Repository README](../README.md)
[Next: API reference](api.md)

The minimum supported compiler is Go 1.27.0. Current main prepares v6 with
`api/v6.txt`; the published v5 and v4 projections are retained at `api/v5.txt`
and `api/v4.txt`. Historical `api/baseline.txt` and `api/v3.txt` remain
unchanged. Breaking
changes require a new major release. Public v2.0.2 is available at `github.com/faustbrian/go-opening-hours/v2`; published v1.1.0 remains
available under its original module path. The intentional provenance validation
change and import migration are described in [v2 migration](v2-migration.md).
The exposed v2.0.0 Git tag remains immutable; it is not a verified stable
release. V2.0.1 is the validated replacement release, not a rewritten tag.

Prepared v6 uses Wire3, Calendar2, Config3, Temporal3 and Validation3 nominal
types and removes the facades; published v5 and v4.0.0 used Temporal2 and
Validation2. V5 changed interval expansion/duration to exact civil point
membership; see the [v5 migration guide](v5-migration.md).
Opening v3.0.0 remains available under `/v3`; see the
[v4 migration guide](v4-migration.md) before changing the imported cohort.

Canonical wire version is `1`. Decoders reject other versions. Canonical bytes
are stable for semantically identical construction order within the same wire
version. Timezone results depend on the installed IANA database, so historical
rule updates may change instant expansion without changing stored bytes.

PostgreSQL JSONB is tested on maintained server releases 14 through 18. pgx v5
uses the root scanner/valuer directly. Location, Track, Postal, and Spatie
fixture compatibility is documented in [legacy migration](legacy-migration.md).
Track and Postal use the shared Location representation; no unverified provider
schema is claimed.

## Adapter paths

The integration entry points are `adapters/calendar`, `adapters/config`,
`adapters/temporal`, `adapters/validation`, and `adapters/wire`. Version 6
removed their top-level predecessors before the announced interval ended; see
[`DEPRECATION.md`](../DEPRECATION.md).

| Fixture | Imported contract | Proven behavior |
| --- | --- | --- |
| Location | weekday `{from,to}` slots | split ranges and explicit closure |
| Track | shared Location slots | start-inclusive/end-exclusive boundary |
| Postal | shared Location slots | owner-day overnight spill and closing |
| Spatie | single-hyphen time-only arrays parsed with Go `15:04` endpoints | weekday opening and closed empty day |
