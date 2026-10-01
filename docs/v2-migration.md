# V2 migration

V2.0.1 is the published, verified stable release. Published
v1.1.0 remains available at its original module path.

The v2.0.0 Git tag was exposed before release validation completed and remains
immutable. It has no verified stable release; use validated v2.0.1 rather
than adopting that tagged snapshot.

Require
`github.com/faustbrian/go-opening-hours/v2@v2.0.1` and add `/v2` after
`go-opening-hours` in every owned package import. This includes the root,
canonical adapters, compatibility facades, compile, encoding, PostgreSQL,
and test helpers. All packages remain in one root module on main; package
directories and the Go 1.27.0 minimum do not change.

Exception source and revision identities must now each be nonempty, valid
UTF-8 within the existing 128-byte bounds. V1 accepted nonempty malformed
byte strings, but JSON encoding could replace distinct identities with the
same replacement character, preventing canonical or SQL round trips. V2
rejects malformed text with `CodeInvalidState`; existing date, operation,
priority, and length errors retain their precedence.

Validate application provenance before construction. If it represents binary
data, choose an application-owned lossless textual identity, such as an
appropriate hexadecimal or base64 encoding, and keep the resulting value
within 128 bytes. Do not replace invalid bytes or guess an original identity
from already-lossy stored JSON. Any affected stored data needs an explicit
application migration based on the original identities.

Canonical wire version 1, valid Unicode identities, timezone/DST semantics,
exception precedence, and SQL receiver atomicity remain unchanged. V2 also
rejects zero exceptions in direct schedule input and preserves canonical
round trips and compiled availability throughout the existing 16-level
composition limit. Compiled indexes return decoding errors instead of silently
substituting a closed schedule.

Published v1 consumers and historical compatibility fixtures are not
automatically migrated. Update application imports and named types together;
do not mix v1 and v2 values within one schedule graph. Maintained external
consumer adoption follows actual v2 publication.
