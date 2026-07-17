# Formal normalization and boundary contract

All intervals use `[start,end)`. Local times have nanosecond precision and lie
in `[00:00,24:00)`. `00:00-00:00` is invalid; full-day opening is a state.

| Policy | Overlap | Adjacency |
| --- | --- | --- |
| `RejectOverlap` | error | preserve separately |
| `RejectOverlapAndAdjacent` | error | error |
| `MergeOverlap` | merge | preserve separately |
| `MergeAdjacent` | merge | merge |

Input is copied, sorted by start/end, and normalized in an owner-day linear
coordinate. Duplicate ranges count as overlap. A merge spanning exactly the
civil day from midnight becomes `DayOpenAllDay`; any ambiguous or longer wrap
returns `CodeDayBoundaryOverflow`.

Normalization is idempotent. Canonical JSON orders Sunday through Saturday,
then exceptions by date/precedence. Algebra results are immutable expression
trees capped by `MaxCompositionDepth`; query fragments are capped at 8,192.
`Schedule.Compare` orders complete schedule values by those canonical bytes, so
its ordering includes provenance and composition shape just like `Equal`.
