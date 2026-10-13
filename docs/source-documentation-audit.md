# Source documentation audit

[Documentation map](index.md) · [Repository README](../README.md)
[Next: API reference](api.md)

## Scope and acceptance

Comments/documentation-only pass from fresh main
`b31229718005f68f508694699fb20bd76e457fb5`. This standalone repository owns
schedule values, rule algebra and canonical persistence; consumers include the
five canonical adapters, five intentionally distinct compatibility packages,
prepared compile handles, test builders and nullable PostgreSQL values.
No calendar behavior, executable source, fixture authority, dependency, SQL,
wire identity or compatibility type identity is changed.

Reproduce with `git ls-files '*.go' '*go.mod'` and
`golib source-docs inventory --json`: one `/v4` module, 74 tracked Go files
(31 production, 43 test), 15 source directories, no generated or vendored Go.
The sole build-tagged Go file is PostgreSQL `integration_test.go`; its live
fixture is distinct from untagged pgx codec tests; its background contexts
supply no per-operation cancellation, so service/process bounds remain external. Eight fuzz targets and seven
benchmarks are maintained. Native Example declarations are checked from tracked
source rather than assumed from README fragments. Existing deprecation blocks
retain compatibility-path retirement conditions; no generator is missing.

## Missing versus investigated documentation

Initial objective findings: 60 missing exported fields/methods across 11 files.
Package entries exist but need normal entry/ownership/limit boundaries. Core
algorithms and adapter implementations have been traced separately; no new
performance or official-standards-conformance claim is intended.

| Contract slice | Source-derived maintenance boundary |
| --- | --- |
| Civil values and immutable construction | Calendar v2 supplies Gregorian Date; midnight is zero LocalTime, equal Range endpoints are invalid, zero DayRule inherits and zero Schedule is closed without a zone. Constructors copy nested rule/exception slices and effective-date values; caller input must not change during construction. Pointer decode/scan methods replace caller storage only after success and require external synchronization. |
| Range normalization and exceptions | Start-date ownership is linearized through the next day before overlap/adjacency decisions. Exactly one midnight-owned full day may normalize to all-day; other spans >=24 civil hours are refused. Exception order is date/ascending priority/source/revision/operation, applied sequentially so later operations win. Per-date source/revision identity is unique; set names do not replace that identity. |
| Evaluation and algebra | Today unions prior-date owned spill before today's exceptions; one date's exceptions can close spill. Overlay distinguishes inherited from explicit closed rules. Composition shares immutable operands, bounds depth at 16 and requires exact nonzero timezone identity; outer metadata/effective configuration is not carried into the composition wrapper. Ranges returns weekly rules only, not dated/algebra availability. |
| Timezones and bounded search | Location loading and folds delegate to pinned Calendar v2 plus the installed Go/IANA database; canonical JSON stores a zone name, not frozen zone rules. Local folds select earlier/later; only ShiftForward handles gaps and it rejects folds. Interval boundary policy is earlier opening/later closing with gap shifting. Search skips candidates that do not change actual point-query state; duration uses elapsed rather than nominal civil hours. |
| Persistence and parser trust | Canonical wire edition 1 is application-owned, not an RFC conformance claim. Input size, UTF-8, duplicate keys, unknown fields, document nesting and evaluated rule cardinality are checked; zero-schedule envelopes ignore metadata/effective/outside values and composition uses operand policies rather than outer fields. Accepted documents need not have canonical byte ordering/whitespace. SHA-256/Equal identify representation, not authentication or algebraic equivalence. Strings/wire bytes include caller provenance and must not be logged as if automatically redacted. |
| Clock and observation | Clock and elapsed capabilities are borrowed; typed nil/panics are mapped to safe failures/zero measurements. Observer runs synchronously without locks and contains panics, but blocked callbacks/clock code have no cancellation bound. RangeCount is queried separately; SearchSteps is a horizon estimate, not an exact internal counter. |
| Optional and compatibility adapters | Calendar holiday expansion bounds traversed dates, not holiday count. Config values require strings and retain old state on failure. Temporal conversions reject lossy bounds/precision. Validation checks successful serialization/parsing, not equality or official compliance. Wire is stateless; PostgreSQL owns no connection lifecycle or server setup. Compatibility named types/sentinels retain their existing identities. |
| Tests and private helpers | Existing algebra, precedence, normalization, resource, successor and DST differential tests identify concrete invariants. Preserve installed-IANA assumptions and frozen corpora; no unchanged race, DST/conformance, fuzz, performance or live-database campaign is warranted by prose. Add only useful fuzz boundary and measured-workload descriptions. |

## Native proof and limits

`.golib.yaml` configures a full `go-doc` snapshot at `api/v4.txt` for 14
packages; `openinghourstest` is the 15th source package and is rendered separately.
Historical `api/baseline.txt`/`api/v3.txt` remain immutable. The snapshot must be
refreshed using the configured native package order and reproduced stably.
No source generation step or generated output requires editing.

Selected verification: Go 1.27.2 formatting, all-tracked runtime AST/token/literal/
directive equality, all 15 package renders (plus applicable tagged source audit),
configured snapshot stability, native examples only when present, objective
source/local-doc/specification checks, pinned changed-content secret scanning
and an ordinary complete-diff review. Broad Make routes can launch PostgreSQL
and runtime campaigns and are not required for this Tier A prose slice.
Coordinator owns shared enforcement, main delivery and exact-source hosted CI.

No new speculative rationale is supplied. Limits are structural, not latency
SLAs; synchronous callbacks may block, environment timezone data may vary,
metadata is bounded but not classified secret-free, and canonical persistence
is neither authentication nor a frozen civil-time ruleset.

## Scoped observations and retained runtime limitation

Go 1.27.2/darwin-arm64, America/New_York installed data: a Sunday 01:30–02:30
rule on 2026-11-01 returns closed at 06:00Z (EST -18000, previous minute EDT -14400),
while its 05:30Z–07:30Z interval expansion contains that instant. Forward search
from 05:59Z returns 07:30Z closing rather than the 06:00Z wall-membership change.
Current source comments and timezones/query guides state boundary-policy expansion
and endpoint-only search without promising exhaustive point equivalence. The
runtime limitation is preserved and reported to the coordinator, not repaired
or certified as a complete calendar/conformance implementation in this pass.

The same toolchain accepts a structured Spatie endpoint `9:00` through Go's
`15:04` parser. Input documentation therefore specifies actual parser/separator
admission, not an unsupported exact-two-digit-hour requirement. No fixtures or
accepted inputs were changed to make these observations pass.

## Local documentation result

Current objective source checking has no findings. Go 1.27.2 verified all 74
tracked files retain runtime AST/signatures, significant tokens, literals and
build/lint/other executable directives. One implicit struct-member terminator
is added solely to attach ImportError.Kind documentation. All 15 rendered source
packages and 370 parsed HTML comment groups were inspected; eight fuzz and seven
benchmark declarations are genuine AST entries and no native Examples exist.
No new executable example is introduced.

The configured fourteen-package native `go doc -all` concatenation was generated
in declared order, newline-normalized according to the shared native generator,
and reproduced byte-for-byte in `api/v4.txt`. The final changed root render was
refreshed twice after review; unchanged thirteen native renders were reused, not
reported as new executions. Historical snapshots and fixtures
are unchanged. AST preservation is not a runtime/cross-platform/conformance
qualification, and the scoped observations are not a new exhaustive DST suite.

The observed fold-expansion/endpoint-search limitation above is the concrete
remaining runtime caveat; addressing it needs separately authorized calendar
behavior and regression coverage. This source pass documents it truthfully.
Host IANA rules, trusted synchronous capability behavior and deployment/database
ownership remain external assumptions. Main publication and exact-source CI are
coordinator-owned.

## Objective Enforcement

The checksum-pinned public Tools v2.4.0 verifier now supplies
`make source-docs-check` as a prerequisite of the existing local and CI Make
routes. Ordinary shared CI selects strict source-documentation checking and
requires its result. Native prerequisites, operation budgets, security and
release controls are retained. Hosted exact-head evidence is recorded by the
coordination ledger; configuration alone does not establish green delivery.
