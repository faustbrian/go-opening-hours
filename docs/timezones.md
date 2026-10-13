# Timezones and DST in five minutes

[Documentation map](index.md) · [Repository README](../README.md)
[Next: API reference](api.md)

Every non-zero schedule loads an explicit IANA location. `IsOpen(time.Time)`
starts from an instant and is unambiguous. `IsOpenLocal` and `ResolveLocal`
require a policy and resolve civil fields to an instant before evaluation.

| Local condition | `RejectDST` | `PreferEarlier` | `PreferLater` | `ShiftForward` |
| --- | --- | --- | --- | --- |
| exact | exact instant | exact | exact | exact |
| spring gap | error | error | error | shift by gap |
| autumn fold | ambiguous error | first instant | second instant | ambiguous error |

Range expansion uses earlier resolution for an ambiguous opening and later
resolution for an ambiguous closing, so `01:30-02:30` across a fall fold covers
both repeated occurrences. Gap boundaries shift forward. A skipped civil date
therefore contributes no elapsed interval when both boundaries collapse.

The package delegates bounded zone loading and exact/fold occurrence resolution
to `calendar/timezone`, which in turn uses Go's installed IANA database.
Opening-hours adds only the explicit `ShiftForward` schedule policy for gaps.
Canonical encoding stores the zone identity, not a frozen rule snapshot.

## Boundary expansion and search limit

Point queries test the instant's local wall-clock membership. Interval expansion
instead resolves an ambiguous opening earlier and closing later. These are not
universally equivalent across a fold: for New York's Sunday 01:30–02:30 rule on
2026-11-01, `IsOpen` is false at 06:00Z (the clock returns to 01:00 EST), while
`EffectiveInstantRanges(05:30Z,07:30Z)` includes that instant. `OpenDuration`
sums this boundary-policy expansion, not an exhaustive point-membership integral.

Transition search enumerates those resolved range endpoints and checks point
values immediately before/after them. It does not generate all offset jumps or
repeat every folded endpoint. In the same example, `NextTransition` from 05:59Z
returns the 07:30Z closing endpoint rather than the 06:00Z wall-membership change.
Do not treat this API as an exhaustive DST transition stream. These observations
use Go 1.27.2 on darwin/arm64 with the installed New York timezone data; no frozen
IANA edition or general conformance is asserted. Correcting that runtime boundary
would require a separate calendar/search change and focused behavioral coverage.
