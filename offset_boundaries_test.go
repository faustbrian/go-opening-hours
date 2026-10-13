package openinghours

import (
	"testing"
	"time"
)

func TestPinnedPOSIXYearBoundaryPreservesAvailability(t *testing.T) {
	// Reuse the exact provider seam already exercised by the public child probes.
	// Keep this native test independent of installed timezone editions.
	pinned := r04PinnedSchedule(t, 2, 30, 3, 30)
	from, err := NewLocalTime(2, 30, 0, 0)
	if err != nil {
		t.Fatal(err)
	}
	until, err := NewLocalTime(3, 30, 0, 0)
	if err != nil {
		t.Fatal(err)
	}
	interval, err := NewRange(from, until)
	if err != nil {
		t.Fatal(err)
	}
	rule, err := OpenRanges([]Range{interval}, RejectOverlap)
	if err != nil {
		t.Fatal(err)
	}
	exception, err := NewException(ExceptionConfig{
		Date: MustDate(2029, time.March, 11), Operation: ExceptionReplace, Rule: rule,
		Source: "posix-leap-regression", Revision: "1",
	})
	if err != nil {
		t.Fatal(err)
	}
	schedule, err := NewSchedule(Config{Timezone: "UTC", Exceptions: []Exception{exception}})
	if err != nil {
		t.Fatal(err)
	}
	schedule.data.location, schedule.data.timezone = pinned.data.location, "America/New_York"
	start := time.Date(2028, time.December, 31, 12, 0, 0, 0, time.UTC)
	annualEnd := time.Date(2029, time.January, 1, 0, 0, 0, 0, time.UTC)
	closed, err := schedule.EffectiveInstantRanges(start, annualEnd)
	if err != nil || len(closed) != 0 {
		t.Fatalf("annual bookkeeping is not availability: %#v %v", closed, err)
	}
	opening := time.Date(2029, time.March, 11, 7, 0, 0, 0, time.UTC)
	closing := opening.Add(30 * time.Minute)
	horizon := 75 * 24 * time.Hour
	ranges, err := schedule.EffectiveInstantRanges(start, start.Add(horizon))
	if err != nil || len(ranges) != 1 || !ranges[0].Start.Equal(opening) || !ranges[0].End.Equal(closing) {
		t.Fatalf("pinned POSIX expansion=%#v error=%v; want [%s,%s)", ranges, err, opening, closing)
	}
	duration, err := schedule.OpenDuration(start, start.Add(horizon))
	if err != nil || duration != 30*time.Minute {
		t.Fatalf("pinned POSIX duration=%s error=%v; want 30m", duration, err)
	}
	transition, err := schedule.NextTransition(start, horizon)
	if err != nil || !transition.Instant.Equal(opening) || transition.Kind != TransitionOpen {
		t.Fatalf("pinned POSIX first transition=%#v error=%v; want opening %s", transition, err, opening)
	}
}

func TestPinnedQueriesRejectOutOfCivilDomain(t *testing.T) {
	schedule := r04PinnedSchedule(t, 2, 30, 3, 30)
	for _, instant := range []time.Time{time.Unix(-(1 << 58), 0), time.Unix(1<<58, 0)} {
		if _, err := schedule.IsOpen(instant); !IsCode(err, CodeInvalidDate) {
			t.Fatalf("out-of-domain point error=%v", err)
		}
		if _, err := schedule.NextTransition(instant, time.Hour); !IsCode(err, CodeInvalidDate) {
			t.Fatalf("out-of-domain forward error=%v", err)
		}
		if _, err := schedule.PreviousTransition(instant, time.Hour); !IsCode(err, CodeInvalidDate) {
			t.Fatalf("out-of-domain reverse error=%v", err)
		}
		if ranges, err := schedule.EffectiveInstantRanges(instant, instant.Add(time.Hour)); !IsCode(err, CodeInvalidDate) || ranges != nil {
			t.Fatalf("out-of-domain ranges=%#v error=%v", ranges, err)
		}
	}
}
