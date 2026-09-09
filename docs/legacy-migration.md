# Legacy Laravel and Spatie migration

## Go adapter import migration

New consumers should select the target-oriented paths:

| Compatibility path | Canonical path |
| --- | --- |
| `openinghourscalendar` | `adapters/calendar` |
| `openinghoursconfig` | `adapters/config` |
| `openinghourstemporal` | `adapters/temporal` |
| `openinghoursvalidation` | `adapters/validation` |
| `openinghourswire` | `adapters/wire` |

Change the import and any explicit named `Value`, `ValidationError`, or `Codec`
references in one consumer change. Functions retain their signatures and
behavior, while sentinels are shared across both paths. The named types are
deliberately distinct so reflection and `errors.As` keep reporting the path
the consumer selected. Failed config decoding remains atomic on both paths.

The compatibility paths remain supported for the longer of 180 days after
successor public availability and two subsequent stable minor releases. Their
eventual removal additionally requires a new major release and a fresh
consumer audit.

## Location structured data

The verified Location shape is weekday-keyed arrays of `{from,to}` strings.
`encoding/testdata/location.json` proves split ranges and closed days. Import it
with `encoding.ImportLocation` and an explicit `ImportLimits` value.

Track and Postal consume that same application-owned Location boundary rather
than defining provider parsers here. Their separate compatibility fixtures are
`encoding/testdata/track.json` and `encoding/testdata/postal.json`; the latter
also proves owner-day overnight spill into the next date.

## Spatie opening-hours

Lossless weekly arrays such as `"09:00-17:00"` are supported by
`encoding.ImportSpatie`; `encoding/testdata/spatie.json` is the fixture. Empty
arrays become closed. Application code must resolve Spatie dated exceptions to
exact dates and construct package exceptions with source/revision.

Spatie's `24:00` endpoint cannot be passed as `LocalTime`. Map a range ending at
`24:00` to an overnight endpoint of `00:00` only when its start is non-midnight;
map `00:00-24:00` to `OpenAllDay`. Record and reject any ambiguous mapping.

## Carrier prose

Location's legacy `OpeningHours::createFromRawCarrierInfo` handles many
provider/language formats. That parser remains in the application. Parse there,
emit structured weekday slots, then import. Raw strings, unexpected-format
fallbacks, and translation behavior are intentionally not portable here.

## Rollout

1. Read legacy and new values side by side.
2. Convert with explicit timezone and normalization policy.
3. Compare representative instants, overnight spill, and known holidays.
4. Persist canonical JSONB plus revision.
5. Dual-read until mismatches are resolved; retain a reversible rollback.
