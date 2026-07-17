# Owned-module integrations

The package keeps each owned capability at an explicit adapter boundary:

- `Clock` and `ElapsedClock` alias narrow `go-clock` capabilities. Transition
  waiting remains caller-owned and composes `clock.Clock` with
  `clock.TimerFactory`; see the [cookbook](cookbook.md#wait-for-the-next-transition).
- `Date` aliases `go-calendar`'s immutable civil value. Schedule construction
  uses `go-calendar/timezone.LoadLocation`, including its 255-byte IANA identity
  bound, and delegates exact/fold occurrence resolution to the same module.
  Opening-hours retains only its explicit gap-shifting policy.
  `openinghourscalendar.FromDate` and `ToDate` preserve compatibility, while
  `HolidayClosures` expands a `business.Calendar` only inside the caller's
  inclusive date range and `maximumDates` bound.
- `openinghourstemporal` converts ordinary and circular `timeofday.Interval`
  values only when their bounds are closed-open. Full-day and collapsed states
  map through `RuleFromIntervals`; mixed state collections fail with
  `ErrLossyMapping`.
- `openinghoursconfig.Value` implements the `go-config` value-unmarshal seam and
  accepts canonical JSON text only.
- `openinghoursvalidation.Validator` returns a deterministic `go-validation`
  report with the stable `opening_hours.invalid_schedule` code.
- `openinghourswire.WireFormat` provides the typed `go-wire` registry identity;
  `Codec` still enforces the package's stricter canonical parser.

None of these adapters owns environment lookup, holiday datasets, provider
payload interpretation, process clocks, or global registries.
