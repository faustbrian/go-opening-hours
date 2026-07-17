# Security model

## Trust boundaries

Canonical parsing treats input as hostile: 1 MiB maximum, valid UTF-8,
duplicate-key rejection, unknown-field rejection, trailing-input rejection,
bounded nesting, strict values, and bounded provenance. Errors never quote the
document.

Fuzz targets exercise canonical JSON, text and SQL scanning, constructors,
timezone resolution, composition/search, structured Location/Spatie imports,
and pgx JSONB codecs.

Construction bounds daily ranges, exceptions, exception expansion, metadata,
composition depth, output fragments, and search horizon. Arithmetic stays in
validated nanosecond/day and bounded `time.Duration` ranges.

## Concurrency

Schedules and indexes are immutable and safe for concurrent reads. Maps and
slices are copied on input; slices are copied on output. There is no unsafe,
cgo, `go:linkname`, reflection-based mutation, lock, global registry, cache,
goroutine, network call, or process clock read in core queries.

## Denial of service

Callers must not retry typed limit/search errors without changing input. Avoid
attaching untrusted strings to application logs. Run fuzz, race, mutation,
vulnerability, and PostgreSQL gates before release.
