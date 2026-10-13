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

## Point transitions and boundary expansion

Point queries test the instant's local wall-clock membership. Interval expansion
instead resolves an ambiguous opening earlier and closing later. These are not
universally equivalent across a fold: for New York's Sunday 01:30–02:30 rule on
2026-11-01, `IsOpen` is false at 06:00Z (the clock returns to 01:00 EST), while
`EffectiveInstantRanges(05:30Z,07:30Z)` includes that instant. `OpenDuration`
continues to sum this two-hour boundary-policy expansion, not the 90-minute
point-membership integral. The expansion and duration policies are unchanged.

Transition search follows actual point membership: it includes timezone-offset
jumps and both occurrences of a folded civil endpoint, then checks point values
immediately before/at each candidate. In this example the state changes are
05:30Z opening, 06:00Z closing, 06:30Z opening and 07:30Z closing.
`NextTransition`/typed searches return the first matching change strictly after
the cursor; `PreviousTransition` returns the nearest change strictly before it.
An offset jump that does not change availability is not a schedule transition.
Search horizons remain bounded to 366 elapsed days.

The regression corpus attributes exact New York TZif bytes to the official
Go 1.27.2 source archive and compares these selected cases with an independent
Python `zoneinfo` wall-membership oracle. That is a bounded regression proof,
not a claim of a frozen runtime IANA edition or universal timezone conformance.
Runtime schedules still use the deployment's installed timezone data.
