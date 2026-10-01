# Security model

This model covers the current source and proposed v2.0.0 release, reviewed
as of 2026-10-01. Published
v1.1.0 remains the current public release until that release is delivered.

## Trust boundaries

Canonical parsing treats input as hostile: 1 MiB maximum, valid UTF-8,
duplicate-key rejection, unknown-field rejection, trailing-input rejection,
bounded nesting, strict values, and bounded provenance. Errors never quote the
document.

JSON validation permits at most 36 value-recursion edges: two per nested
composition plus six for the deepest ranged weekly or exception rule. This
matches all 16 permitted composition levels without changing the domain depth
or 1 MiB limits. Direct and named exception collections reject invalid
zero-value identities. Compiled indexes propagate canonical decoding failures
instead of silently substituting a closed schedule.

V2 exception source and revision identities are nonempty, valid UTF-8 with
the existing 128-byte ceilings. Constructor date, operation, priority, and
length errors retain their ordering before text admission. The same bounded
text predicate protects stored direct and named exception identities, avoiding
lossy JSON replacement and identity collisions during persistence. This is an
intentional validation change from v1; see [migration](v2-migration.md).

Fuzz targets exercise canonical JSON, text and SQL scanning, constructors,
timezone resolution, composition/search, structured Location/Spatie imports,
and pgx JSONB codecs.

Construction bounds daily ranges, exceptions, exception expansion, metadata,
composition depth, output fragments, and search horizon. Arithmetic stays in
validated nanosecond/day and bounded `time.Duration` ranges.

Exception cardinality and set validity are checked before exception ownership
copies, after the existing timezone, metadata, date, and weekly validation.
Invalid or empty sets stop admission immediately; each valid nonempty set
consumes the shared exception budget. Temporal and Spatie cardinality and
configuration/SQL input lengths are checked before their conversion or copy.
The public
`MaxRangesPerDay`, `MaxExceptions`, and `MaxJSONBytes` constants expose the
corresponding hard limits to adapters and callers. Holiday expansion is capped
at `MaxExceptions` dates even when a larger caller limit is supplied.

All parsers, scanners, constructors, composition methods, and queries return
typed errors for hostile or invalid runtime input. `MustDate` is the sole
intentional panic-on-error helper and is documented for static fixtures only;
runtime callers use `NewDate`. Observation callback panics are recovered after
the query result is complete and cannot alter that result. Typed-nil or
panicking injected clocks fail closed; failures in elapsed measurement produce
an explicitly unmeasured zero duration and cannot suppress a completed query.

## Threat controls

| Threat | Control | Executable evidence |
| --- | --- | --- |
| Oversized persistence/configuration values cause a second full-size copy | Length rejection precedes byte conversion or ownership copy | `TestInputLimitsAndAtomicSQLScan`, `TestParseRejectsOversizedStringBeforeConversion`; source ordering establishes copy avoidance |
| Oversized adapter collections allocate before reaching core limits | Spatie and Temporal adapters preflight the shared public cardinality bounds | `TestImportSpatieRejectsOversizedCollectionsBeforeConversion`, `TestRuleFromIntervalsRejectsOversizedCollectionBeforeConversion` |
| A hostile clock panics before or after a query | Typed-nil detection and panic containment isolate clock capabilities from query results | `TestIsOpenNowRejectsTypedNilClock`, `TestHostileElapsedClocksCannotChangeQueryResult` |
| A caller selects an effectively unbounded holiday expansion | Calendar expansion rejects limits above `MaxExceptions` | `TestCalendarErrorsAndBoundsMatchLegacy` |
| A replaced dependency archive is silently trusted | `go.sum` is checked against the public checksum database; authentication is never disabled | Fresh `GOMODCACHE` plus `GOSUMDB=sum.golang.org go mod download -json` and `go mod verify` |

## Concurrency

Schedules and indexes are immutable and safe for concurrent reads. Maps and
slices are copied on input; slices are copied on output. There is no unsafe,
cgo, `go:linkname`, reflection-based mutation, lock, global registry, cache,
goroutine, network call, or process clock read in core queries.

The optional compiled index owns only an immutable schedule value. It has no
close method, background worker, retained callback, or cache, so there is no
resource lifecycle to leak. Observation callbacks are invoked synchronously
without internal locks and are not retained after return.

## Denial of service

Callers must not retry typed limit/search errors without changing input. Avoid
attaching untrusted strings to application logs. Release verification follows
the repository's proportional assurance policy and selected release boundary.

The non-context query APIs perform only in-memory work and terminate at their
documented range, exception, composition-depth, output, and 366-day search
bounds. They cannot be canceled mid-call. SQL query cancellation remains the
database caller's responsibility before a driver invokes `Scan`; `Scan` only
receives an already-materialized value. The package prevents additional
oversized copies but cannot reclaim input memory already allocated by callers,
drivers, or configuration systems.

The composition-depth limit bounds expression height, not a separately
configurable total node count. A maximally broad in-memory tree can therefore
consume substantially more CPU than a typical schedule while remaining finite;
callers that assemble trees from untrusted plans should impose a smaller
application-level node budget.

## Accepted boundaries and review conditions

| Boundary | Owner and rationale | Mitigation | Review condition |
| --- | --- | --- | --- |
| JSON materialization | Package maintainer: parsing receives at most 1 MiB, but JSON decoding materializes wire collections before domain cardinality checks. Canonical output builds a wire tree and marshals before checking its byte length; the output ceiling is not an allocation ceiling. | Retain input/depth/domain limits; applications constrain serialized composition breadth. No universal pre-allocation guarantee is made. | Review when increasing limits, adding wire fields, or admitting caller-selected broader composition plans. |
| Composition breadth | Application owner, with package-maintained depth/query bounds: depth 16 bounds height, not total tree nodes. | Impose an application node budget for untrusted composition plans; do not repeatedly retry limit errors. | Review when composition shape becomes externally selectable or bounded-query cost no longer fits the application budget. |
| Materialized input and SQL cancellation | Caller/driver owner: scanner and config entry points receive values already allocated upstream and cannot cancel that acquisition. | Apply caller context and driver limits before Scan, and reject oversized values before the package makes a further copy. Failed decoding preserves the receiver. | Review when changing input providers, drivers, or context-aware acquisition contracts. |
| Synchronous capabilities | Caller owner: clocks, elapsed-measurement functions, and observers are trusted synchronous callbacks. Panic containment preserves error/query outcomes but cannot preempt a blocking callback. | Supply bounded nonblocking implementations; use an application-owned bounded handoff for slow observation export. The package retains no callbacks or background work. | Review when callbacks cross a network/process boundary or acquire blocking resources. |
| Dependencies | Package maintainer: timezone and adapter behavior relies on pinned public dependency modules, not a vendored complete implementation. | Keep version/checksum pins, public checksum-database authentication, and ordinary adapter/compatibility checks; evaluate dependency security findings before adoption. | Review when a pin changes, a relevant finding appears, or dependency behavior changes a documented contract. |
