# API reference

[Documentation map](index.md) · [Repository README](../README.md)
[Next: Security model](security.md)

This reference describes the prepared v6 API at
`github.com/faustbrian/go-opening-hours/v6`.
V6 is not published yet. Public v5, v4.0.0, v3.0.0, v2.0.2 and historical v1.1.0 remain independently available.
The exposed v2.0.0 Git tag is immutable but did not complete release validation.

The authoritative symbol reference is the package documentation generated from
the source and selected for `golib api check` against `api/v6.txt`. Historical
`api/baseline.txt`, `api/v3.txt` and the published `api/v4.txt` are retained. The
baseline and package selection are declared in `.golib.yaml`. Key contracts
are:

| Area | Exported API |
| --- | --- |
| civil values | `calendar` `Date`, `LocalTime`, `Range`, constructors/accessors |
| day rules | `DayState`, `DayRule`, `OverlapPolicy`, rule constructors |
| schedules | `Config`, `Schedule`, metadata, equality, comparison, hash, revision |
| exceptions | `Exception`, `ExceptionSet`, operations, conflict policy |
| queries | `Availability`, `DailyRange`, `InstantRange`, `Transition` |
| DST | `LocalResolutionPolicy`, `ResolvedLocal`, `LocalKind` |
| algebra | `Union`, `Intersection`, `Subtract`, `Overlay` |
| formatting | `HumanSummary` display text, canonical JSON/text wire encoding |
| encoding | strict parse, SQL scanner/valuer |
| capabilities | `Clock`, `ElapsedClock`, `Observer`, `Observation` |
| canonical adapters | `adapters/calendar`, `adapters/config`, `adapters/temporal`, `adapters/validation`, `adapters/wire` |
| compatibility paths | `openinghourscalendar`, `openinghoursconfig`, `openinghourstemporal`, `openinghoursvalidation`, `openinghourswire` |

Root constructors with error results return typed package `Error` categories
through `IsCode`; `MustDate` panics on invalid static input and adapter errors
retain their own documented identities.
Errors contain an operation and stable code only; they do not embed source
documents or schedules. Limits include `MaxRangesPerDay` (64 ranges/day),
`MaxExceptions` (4,096 exceptions), `MaxJSONBytes` (1 MiB JSON), 16 composition
levels, 366 elapsed search days, 8,192 output fragments, 128-byte exception
provenance, and 256-byte metadata fields.

The zero `Date` is invalid, zero `LocalTime` is midnight, zero `DayRule` is
inherited, and zero `Schedule` is closed without a timezone.

Canonical and compatibility adapter paths share the same package-owned error
sentinels and observable behavior. Their exported named `Value`,
`ValidationError`, and `Codec` types retain distinct reflection identities;
use the type from the import path your code selected. Failed configuration
decoding is atomic and leaves the prior schedule unchanged.
