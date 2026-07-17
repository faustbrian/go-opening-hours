# Performance and complexity

| Operation | Complexity | Bound |
| --- | --- | --- |
| construction | `O(R log R + E log E)` | 64/day, 4,096 exceptions |
| daily query | `O(log E + R + F)` | bounded fragments |
| point query | daily query plus timezone conversion | one date |
| composition | `O(F_left + F_right)` per date | depth 16 |
| transition search | dates in horizon times daily query | 366 days |
| canonical encoding | schedule cardinality | 1 MiB output |

The prepared `compile.Index` owns an immutable canonical copy and relies on the
root's sorted exception index. It has no cache, lock, or cleanup lifecycle.

Run `make benchmark` for allocations. The seven categories cover construction,
normalization, daily lookup, transition search, large exception sets,
composition, and canonical encoding. Results are machine-specific and are not
contractual latency guarantees.
