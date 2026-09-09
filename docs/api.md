# API reference

The authoritative symbol reference is the package documentation generated from
the source and checked by `golib api check` against `api/baseline.txt`. The
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

Every constructor returns a typed package `Error` category through `IsCode`.
Errors contain an operation and stable code only; they do not embed source
documents or schedules. Limits include 64 ranges/day, 4,096 exceptions, 1 MiB
JSON, 16 composition levels, 366 elapsed search days, 8,192 output fragments,
128-byte exception provenance, and 256-byte metadata fields.

The zero `Date` is invalid, zero `LocalTime` is midnight, zero `DayRule` is
inherited, and zero `Schedule` is closed without a timezone.

Canonical and compatibility adapter paths share the same package-owned error
sentinels and observable behavior. Their exported named `Value`,
`ValidationError`, and `Codec` types retain distinct reflection identities;
use the type from the import path your code selected. Failed configuration
decoding is atomic and leaves the prior schedule unchanged.
