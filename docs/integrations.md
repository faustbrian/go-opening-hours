# Owned-module integrations

The package keeps each owned capability at an explicit adapter boundary:

- `Clock` and `ElapsedClock` alias narrow `clock` capabilities. Transition
  waiting remains caller-owned and composes `clock.Clock` with
  `clock.TimerFactory`; see the [cookbook](cookbook.md#wait-for-the-next-transition).
- `Date` aliases `calendar`'s immutable civil value. Schedule construction
  uses `calendar/timezone.LoadLocation`, including its 255-byte IANA identity
  bound, and delegates exact/fold occurrence resolution to the same module.
  Opening-hours retains only its explicit gap-shifting policy.
  `adapters/calendar.FromDate` and `ToDate` preserve exact civil values, while
  `HolidayClosures` expands a `business.Calendar` only inside the caller's
  inclusive date range and `maximumDates` bound.
- `adapters/temporal` converts ordinary and circular `timeofday.Interval`
  values only when their bounds are closed-open. Full-day and collapsed states
  map through `RuleFromIntervals`; mixed state collections fail with
  `ErrLossyMapping`.
- `adapters/config.Value` implements the `config` value-unmarshal seam,
  accepts canonical JSON text only, and leaves its prior schedule unchanged
  after a failed decode.
- `adapters/validation.Validator` returns a deterministic `validation`
  report with the stable `opening_hours.invalid_schedule` code.
- `adapters/wire.WireFormat` provides the typed `wire` registry identity;
  `Codec` still enforces the package's stricter canonical parser.

The five former top-level integration packages remain compatibility facades.
They delegate to these canonical implementations, share sentinel identities,
and preserve their released named types. The named types on the canonical and
legacy paths are intentionally distinct, so reflection and `errors.As` retain
the identity of the path selected by the consumer.

None of these adapters owns environment lookup, holiday datasets, provider
payload interpretation, process clocks, or global registries.
