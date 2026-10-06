# V4 preparation and Wire v3 migration

Main prepares `github.com/faustbrian/go-opening-hours/v4`. The candidate
requires its own release qualification; public Opening v3.0.0 remains
available under `/v3`.

Migrate every owned Opening import from `/v3` to `/v4`, including canonical
`adapters/<target>` packages and retained `openinghours<target>` facades.
Migrate Wire imports to `github.com/faustbrian/go-wire/v3` and select public
`v3.0.0` in the same consumer change. Both exported `WireFormat` constants
now have the nominal Wire v3 `Format` identity. Wire v1 and v3 named types
are not interchangeable; do not mix the old and new adapter graphs.

The registry name remains `opening-hours+json;v=1`. Opening's `Codec` still
uses the stricter package-owned canonical parser, not a general Wire JSON
codec. Canonical bytes, invalid-input errors, and schedule round trips are
unchanged. Canonical and retained `Codec` types remain deliberately distinct,
and the compatibility retention interval is unchanged. Calendar v2,
Config v2, Temporal v2, Validation v2, pgx v5, and Go 1.27 remain selected.

Historical `api/baseline.txt` and `api/v3.txt` remain unchanged; the prepared
v4 API projection is `api/v4.txt`. This migration does not rewrite historical
tags or module graphs and does not by itself publish Opening v4.
