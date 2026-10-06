# V3 preparation and migration

This is the historical v3 preparation record. Opening v3.0.0 is now publicly
available under `/v3`; the pending-release statements below describe that
earlier preparation state. Current main prepares v4; see
[the v4 migration guide](v4-migration.md).

Main prepares `github.com/faustbrian/go-opening-hours/v3`; it is not a
published or release-qualified version. Public Opening v2.0.2 remains
available under `/v2`, and historical v1 graphs remain independently usable.

The dependency baseline is public Calendar `/v2@v2.0.0`, Validation
`/v2@v2.0.0` and Temporal `/v2@v2.0.0`. Temporal v2 is published from source
`e79d4d85da711fdfa5b9de7618229842952cd044`; its public module archive and
go.mod are authenticated by the checksum database. Opening v3 still requires
its own release qualification. Earlier disposable
development composition used integrated Temporal source
`c69d4141997f4e43095420e9d414dec8bcef05c7`; that proof does not substitute for
public dependency resolution. Production modules contain no local replacements.

For v3 adoption, migrate all owned Opening
imports to `/v3`, Calendar date/business/timezone imports to `/v2`, Temporal
interval/time/bounds imports to `/v2`, and Validation context/report/validator
imports to `/v2` together. Opening Date remains an alias of the selected
Calendar Date. Different major versions have distinct named types and
sentinels; do not mix them within one schedule graph.

Canonical `adapters/<target>` and retained `openinghours<target>` packages
remain within the v3 root module. Existing conversion behavior, sentinel
parity and intentionally distinct wrapper/error identities remain supported.
Weekly and exception semantics, finite limits, diagnostic privacy, failed
receiver atomicity and canonical wire version 1 are unchanged.

Historical `api/baseline.txt` remains untouched. `api/v3.txt` was generated
with Go 1.27.1 `go doc -all` over the fourteen packages selected in
`.golib.yaml`, with the candidate Opening module and exact integrated
Temporal source above as the only local fixture replacements. Calendar2 and
Validation2 resolved through the public proxy and checksum database. This
development projection now matches documentation generated from the public
dependency graph without local replacements. Public module checksums and
canonical/retained adapter composition are verified, but this does not publish
Opening v3. Main CI, native security qualification, release gates and clean
public consumer verification remain required before publication claims.
