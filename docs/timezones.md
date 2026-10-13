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

Range expansion represents exactly the instants whose local wall fields match
the effective civil ranges. It visits both fold occurrences independently and
excludes missing gap times. A skipped civil date contributes no elapsed interval.
This does not change explicit `ResolveLocal`/`IsOpenLocal` resolution policies.

The package delegates bounded zone loading and exact/fold occurrence resolution
to `calendar/timezone`, which in turn uses Go's installed IANA database.
Opening-hours adds only the explicit `ShiftForward` schedule policy for gaps.
Canonical encoding stores the zone identity, not a frozen rule snapshot.

## Agreement of point queries, expansion and duration

Point queries test the instant's local wall-clock membership. Interval expansion
uses that same membership. For New York's Sunday 01:30–02:30 rule on 2026-11-01,
`IsOpen` is false at 06:00Z when the clock returns to 01:00 EST.
`EffectiveInstantRanges(05:30Z,07:30Z)` returns `[05:30Z,06:00Z)` and
`[06:30Z,07:30Z)`, and `OpenDuration` is 90 minutes. Clipping to the closed
`[06:00Z,06:30Z)` period returns no intervals and zero duration. An all-day rule
still spans the entire 25-hour fold day, because the offset jump does not close it.

The prepared v5.0.0 point-membership policy supersedes the former contiguous
earlier-opening/later-closing envelope and gap-boundary shifting used by interval
expansion. Applications depending on that envelope must change their duration
and expansion expectations when adopting v5 after publication. Canonical bytes
and explicit local resolution policies are unchanged.

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
