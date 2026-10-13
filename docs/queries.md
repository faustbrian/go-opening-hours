# Availability queries in five minutes

[Documentation map](index.md) · [Repository README](../README.md)
[Next: API reference](api.md)

- `IsOpen` evaluates an instant in the schedule timezone.
- `IsOpenLocal` resolves civil fields using its required DST policy and then
  evaluates the selected instant.
- `EffectiveRanges` returns fragments for one civil date.
- `EffectiveInstantRanges` clips absolute intervals to a caller range.
- `NextTransition`, `NextOpening`, `NextClosing`, and `PreviousTransition`
  search actual point-availability changes strictly beyond/before the supplied
  instant, including offset jumps and both occurrences of folded endpoints.
  They use the same civil point membership as interval expansion and duration.
- `OpenDuration` sums elapsed, not wall-clock, duration.

All instant intervals must be positive and no longer than 366 elapsed days.
Search exhaustion returns `CodeSearchExhausted`; invalid bounds return
`CodeInvalidHorizon` or `CodeInvalidInterval`. No query reads the process clock.
Use `IsOpenNow` only with an injected `Clock`. Use `RejectDST` when a caller
wants gaps and folds to fail instead of selecting or shifting an occurrence.

`Availability.Explanation` identifies weekly, spill, exception, composition,
or outside-effective-range provenance and includes the timezone. Exception
source/revision are bounded; labels and complete schedule data are never added.
