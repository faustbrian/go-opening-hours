# Observability

[Documentation map](index.md) · [Repository README](../README.md)
[Next: API reference](api.md)

`ObserveIsOpen` and `ObserveNextTransition` execute normal queries and call an
optional observer afterward. They accept a separate `clock` `ElapsedClock`
for monotonic duration measurement. Passing a nil elapsed clock explicitly
disables measurement and reports a zero duration; the package never reads the
process clock implicitly. Observation fields are limited to operation, outcome,
range count, bounded search-step estimate, and duration. `SearchSteps` is derived
from the horizon (capped at 367), not counted candidate visits; `RangeCount` comes
from a separate effective-range query. Duration excludes that extra query and
observer delivery, and trusts the injected elapsed function's value.

No label, source, revision, date, timezone, customer identifier, or schedule
content is present. Callbacks run without locks. A callback panic is recovered
and cannot change the returned query value or error. There is no exporter,
registry, background refresh, or hidden goroutine.

Applications may translate observations into their own metrics, but should keep
dimensions low-cardinality and must not attach schedule or customer data.
