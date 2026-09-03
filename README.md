# opening-hours

[![CI](https://github.com/faustbrian/go-opening-hours/actions/workflows/ci.yml/badge.svg?branch=main)](https://github.com/faustbrian/go-opening-hours/actions/workflows/ci.yml)
[![CodeQL](https://img.shields.io/badge/CodeQL-required-blue)](https://github.com/faustbrian/go-opening-hours/actions/workflows/ci.yml)
[![Coverage](https://img.shields.io/badge/coverage-100%25_required-blue)](CONTRIBUTING.md#verification)
[![Mutation](https://img.shields.io/badge/mutation-100%25_required-blue)](CONTRIBUTING.md#verification)
[![Documentation](https://img.shields.io/badge/docs-checked_in_CI-blue)](docs/)
[![Go Reference](https://pkg.go.dev/badge/github.com/faustbrian/go-opening-hours.svg)](https://pkg.go.dev/github.com/faustbrian/go-opening-hours)
[![Release](https://img.shields.io/github/v/release/faustbrian/go-opening-hours?sort=semver)](https://github.com/faustbrian/go-opening-hours/releases)
[![Go](https://img.shields.io/badge/go-1.26.6-00ADD8?logo=go)](https://go.dev/)
[![License](https://img.shields.io/badge/license-MIT-blue.svg)](LICENSE)

Immutable, deterministic, timezone-safe recurring opening hours and dated
exceptions for Go 1.26.6 and later.

The package models generic availability for service points, storefronts,
offices, pickup locations, and support desks. It does not parse carrier prose,
book appointments, plan workforces, or decide whether an order is eligible.

## Five-minute start

```go
package main

import (
	"fmt"
	"time"

	openinghours "github.com/faustbrian/go-opening-hours"
)

func main() {
	start, _ := openinghours.NewLocalTime(9, 0, 0, 0)
	end, _ := openinghours.NewLocalTime(17, 0, 0, 0)
	dayRange, _ := openinghours.NewRange(start, end)
	monday, _ := openinghours.OpenRanges(
		[]openinghours.Range{dayRange},
		openinghours.RejectOverlapAndAdjacent,
	)

	schedule, _ := openinghours.NewSchedule(openinghours.Config{
		Timezone: "Europe/Helsinki",
		Weekly: map[time.Weekday]openinghours.DayRule{
			time.Monday: monday,
		},
	})

	result, _ := schedule.IsOpen(
		time.Date(2026, time.January, 5, 10, 0, 0, 0, time.UTC),
	)
	fmt.Println(result.Open, result.Explanation.Timezone)
}
```

All boundaries are start-inclusive and end-exclusive. A range such as
`22:00-02:00` belongs to its start date. The zero `Schedule` is closed and has
no timezone; it never means always open.

## What is explicit

- IANA timezone identity and DST gap/fold policy
- inherited, ranged, all-day, and closed day states
- overlap, adjacency, and normalization policy
- exact-date replace, add, subtract, and close operations
- exception priority, source, revision, and optional named set
- inclusive effective dates and outside-range behavior
- bounded transition horizons, output counts, parsing, and composition depth
- canonical JSON, stable comparison/hash, separate human display summaries
- SQL/JSONB persistence and native pgx behavior
- injected clocks and privacy-safe observation callbacks

## Packages

| Package | Purpose |
| --- | --- |
| root | Values, rules, exceptions, algebra, queries, encoding, SQL |
| `compile` | Immutable prepared query handle |
| `encoding` | Canonical Location/Spatie imports; Track/Postal fixtures |
| `postgres` | Nullable JSONB wrapper and pgx compatibility |
| `openinghourswire` | Byte-codec adapter |
| `openinghoursvalidation` | Canonical validation adapter |
| `openinghoursconfig` | Strict configuration adapter |
| `openinghourscalendar` | `calendar` dates and holiday closures |
| `openinghourstemporal` | Lossless `temporal/timeofday` conversion |
| `openinghourstest` | Panic-on-error test builders |

## Documentation

Start at the [documentation index](docs/README.md). The five-minute guides cover
[weekly schedules](docs/weekly-schedules.md), [exceptions](docs/exceptions.md),
[overnight ranges](docs/overnight.md), [timezones](docs/timezones.md), and
[queries](docs/queries.md). The formal contracts are in
[precedence](docs/precedence.md) and [normalization](docs/normalization.md).
Owned-module integration is covered in [integrations](docs/integrations.md).

For ecosystem-wide selection and ownership guidance, see the versioned
[Golib ecosystem index](https://github.com/faustbrian/go-library-tools/blob/v1.4.0/docs/ecosystem/README.md)
and its [Domain utilities family](https://github.com/faustbrian/go-library-tools/blob/v1.4.0/docs/ecosystem/design-language.md#package-families-and-selection).

## Local verification

```sh
make check
```

`make check` runs the complete shared-library contract, including formatting,
tests, race detection, exact coverage, mutation verification, fuzzing,
benchmarks, API compatibility, documentation, security, and the typed timezone
regression operation. Package-specific services are task-owned by the shared
tool; PostgreSQL integration is enabled by the module manifest.

## Support and policy

See [security policy](SECURITY.md), [contribution guide](CONTRIBUTING.md),
[compatibility policy](docs/compatibility.md), and [changelog](CHANGELOG.md).
The project is available under the [MIT License](LICENSE).
