# V5 point-membership migration

[Documentation map](index.md) · [Repository README](../README.md)

Current main prepares v5.0.0 of `github.com/faustbrian/go-opening-hours/v5`.
Public adoption requires the v5.0.0 tag and release; source-tree presence alone
does not establish availability. Published v4.0.0 remains available unchanged at
`/v4`.

Once the v5.0.0 tag and release are published, select it explicitly:

```sh
go get github.com/faustbrian/go-opening-hours/v5@v5.0.0
```

Migrate all owned Opening imports from `/v4` to `/v5` together, including
canonical `adapters/<target>` packages and retained `openinghours<target>`
facades. The new module's named schedule and adapter types are not interchangeable
with v4 types. Calendar v2, Config v2, Temporal v2, Validation v2, Wire v3,
pgx v5 and Go 1.27 remain selected without supplier repins.

## Changed timezone expansion policy

Instant expansion and elapsed duration now agree with `IsOpen` and actual
transition search. New York's 2026-11-01 Sunday 01:30–02:30 rule expands to
`[05:30Z,06:00Z)` and `[06:30Z,07:30Z)`, totaling 90 minutes. The repeated
01:00–01:30 period is closed. V4's contiguous earliest-open/latest-close
envelope instead counted two hours. Gap expansion includes only existing
matching wall times, not shifted missing endpoints.

Review duration, expanded-interval and business expectations that depended on
the old envelope. All-day fold availability still spans 25 elapsed hours.
Caller clipping, start-inclusive/end-exclusive membership, the 366-day horizon,
8192 collected-fragment limit and no partial results on error remain bounded.
Explicit `ResolveLocal`/`IsOpenLocal` policies remain unchanged.

## Stored representations

The canonical JSON version remains `1`; the schedule fields, parser and
serialization are unchanged by the query-policy repair. Existing canonical
values remain readable, but querying them under v5 uses the new expansion and
duration semantics. The Wire registry name remains `opening-hours+json;v=1`;
codecs and SQL values do not encode the Go module major or a frozen timezone
database. Migration must therefore review semantics as well as nominal types.

The maintained symbol projection is `api/v5.txt`. Historical v4 tags, source
and the published `api/v4.txt` remain intact. See [timezones](timezones.md) for
the attributable bounded regression corpus and runtime timezone-data limits.
