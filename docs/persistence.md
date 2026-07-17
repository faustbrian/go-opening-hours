# Persistence

`Schedule.Value` emits canonical JSON bytes and `Schedule.Scan` accepts JSONB
bytes, strings, or SQL `NULL`. `NULL` becomes the fail-closed zero schedule.
Invalid values leave a typed error and do not expose input.

The `postgres.JSONB` wrapper distinguishes nullable database state from a valid
zero schedule. The root type implements the interfaces selected by pgx JSONB's
native codec, so no global connection registration is necessary.

Recommended schema:

```sql
create table resource_availability (
    resource_id bigint primary key,
    schedule jsonb not null,
    revision text not null,
    check (jsonb_typeof(schedule) = 'object')
);
```

The application owns migrations from old columns. Decode legacy values, build a
schedule, canonicalize it, verify representative instants, then write JSONB in
a reversible migration. See [legacy migration](legacy-migration.md).
