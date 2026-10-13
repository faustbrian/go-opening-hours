# Selected fold regression input

`zoneinfo/America/New_York` is copied byte-for-byte from official Go1.27.2
`lib/time/zoneinfo.zip` at commit022c8636110ebac86a9b88724cda168ed713c21f.
`source.json` records the authenticated archive and entry SHA256. `README.upstream`
preserves the upstream IANA/public-domain attribution. No separate IANA release
edition is inferred. Input identity is checked by the Go regression and oracle.

Regenerate only after separately authorizing a timezone input update. Acquire
the immutable official Go archive, verify its recorded SHA256, extract the named
entry unchanged, then run `python3 oracle.py zoneinfo/America/New_York > oracle.json`.
The independent Python standard-library zoneinfo oracle visits exactly10801
second ticks in a fixed three-hour window, with before/at/after microsecond
controls. It evaluates the selected Sunday01:30–02:30 wall predicate without
calling opening-hours, Calendar or the Go timezone implementation. This selected
oracle does not assert global timezone conformance.

Tests launch their own bounded process with this ZONEINFO directory before
startup so cached environment selection cannot silently use installed zone data.
The corpus covers actual point transitions, disjoint instant expansion and
the 90-minute point-membership integral. The original contiguous two-hour
envelope was reproduced before the selected expansion-policy repair; its
historical source remains in Git rather than as a contradictory current assertion.

`zoneinfo/Antarctica/Casey` comes from the same authenticated official archive;
its entry hash is recorded in `source.json`. The selected 2010-03-04T15:00Z
three-hour fold moves local time from March5 02:00+11 to March4 23:00+08.
Python zoneinfo independently confirms local March4 23:30/23:45 correspond to
15:30Z/15:45Z. In-process private-provider tests use those exact bytes to prove
forward/backward search and instant expansion over reversed and intermediate
civil date bounds; they do not claim a new public timezone-injection API.
