package openinghours_test

import (
	"testing"
	"time"

	openinghours "github.com/faustbrian/go-opening-hours/v5"
	prepared "github.com/faustbrian/go-opening-hours/v5/compile"
)

func TestR04POSIXLeapYearExtensionKeepsQueryAgreement(t *testing.T) {
	// Reuse the authenticated NY corpus in a joined, bounded child process.
	// A nonadvancing candidate mutation must not strand the parent test process.
	if r04RunWithPinnedZone(t) {
		return
	}
	_, _ = r04Fixture(t)
	rule, err := openinghours.OpenRanges([]openinghours.Range{
		mustRange(t, 2, 30, 3, 30),
	}, openinghours.RejectOverlap)
	if err != nil {
		t.Fatal(err)
	}
	exception, err := openinghours.NewException(openinghours.ExceptionConfig{
		Date:      openinghours.MustDate(2029, time.March, 11),
		Operation: openinghours.ExceptionReplace, Rule: rule,
		Source: "posix-leap-regression", Revision: "1",
	})
	if err != nil {
		t.Fatal(err)
	}
	schedule, err := openinghours.NewSchedule(openinghours.Config{
		Timezone: "America/New_York", Exceptions: []openinghours.Exception{exception},
	})
	if err != nil {
		t.Fatal(err)
	}
	encoded, err := schedule.CanonicalJSON()
	if err != nil {
		t.Fatal(err)
	}
	decoded, err := openinghours.ParseJSON(encoded)
	if err != nil {
		t.Fatal(err)
	}

	// Python zoneinfo over the existing exact TZif independently places the
	// jump at 07:00Z (01:59:59 EST -> 03:00 EDT), closing at 07:30Z.
	opening := time.Date(2029, time.March, 11, 7, 0, 0, 0, time.UTC)
	closing := opening.Add(30 * time.Minute)
	start := time.Date(2028, time.December, 31, 12, 0, 0, 0, time.UTC)
	horizon := 75 * 24 * time.Hour
	for _, candidate := range []struct {
		name     string
		schedule openinghours.Schedule
	}{{"constructed", schedule}, {"canonical", decoded}} {
		t.Run(candidate.name, func(t *testing.T) {
			for _, point := range []struct {
				instant time.Time
				open    bool
			}{{opening.Add(-time.Nanosecond), false}, {opening, true}, {closing.Add(-time.Nanosecond), true}, {closing, false}} {
				state, err := candidate.schedule.IsOpen(point.instant)
				if err != nil || state.Open != point.open {
					t.Errorf("point %s open=%t error=%v; want %t", point.instant, state.Open, err, point.open)
				}
			}
			index, err := prepared.New(candidate.schedule)
			if err != nil {
				t.Fatal(err)
			}
			indexed, err := index.NextTransition(start, horizon)
			if err != nil || !indexed.Instant.Equal(opening) || indexed.Kind != openinghours.TransitionOpen {
				t.Errorf("prepared forward=%#v error=%v; want opening %s", indexed, err, opening)
			}
			// A same-offset annual boundary at the horizon is not availability.
			annualEnd := time.Date(2029, time.January, 1, 0, 0, 0, 0, time.UTC)
			closed, err := candidate.schedule.EffectiveInstantRanges(start, annualEnd)
			if err != nil || len(closed) != 0 {
				t.Errorf("closed annual horizon=%#v error=%v; want no intervals", closed, err)
			}
			ranges, err := candidate.schedule.EffectiveInstantRanges(start, start.Add(horizon))
			if err != nil || len(ranges) != 1 || !ranges[0].Start.Equal(opening) || !ranges[0].End.Equal(closing) {
				t.Errorf("POSIX expansion=%#v error=%v; want [%s,%s)", ranges, err, opening, closing)
			}
			duration, err := candidate.schedule.OpenDuration(start, start.Add(horizon))
			if err != nil || duration != 30*time.Minute {
				t.Errorf("POSIX duration=%s error=%v; want 30m", duration, err)
			}
			for _, query := range []struct {
				name string
				call func(time.Time, time.Duration) (openinghours.Transition, error)
				from time.Time
				want time.Time
				kind openinghours.TransitionKind
			}{
				{"first-forward", candidate.schedule.NextTransition, start, opening, openinghours.TransitionOpen},
				{"typed-opening", candidate.schedule.NextOpening, start, opening, openinghours.TransitionOpen},
				{"nearest-reverse-opening", candidate.schedule.PreviousTransition, opening.Add(time.Minute), opening, openinghours.TransitionOpen},
				{"typed-closing-control", candidate.schedule.NextClosing, start, closing, openinghours.TransitionClose},
				{"nearest-reverse-closing-control", candidate.schedule.PreviousTransition, closing.Add(time.Minute), closing, openinghours.TransitionClose},
			} {
				transition, err := query.call(query.from, horizon)
				if err != nil || !transition.Instant.Equal(query.want) || transition.Kind != query.kind {
					t.Errorf("%s=%#v error=%v; want %s kind=%v", query.name, transition, err, query.want, query.kind)
				}
			}
		})
	}
}

func TestR04SearchReturnsBoundaryAtStalePOSIXYearEnd(t *testing.T) {
	if r04RunWithPinnedZone(t) {
		return
	}
	_, _ = r04Fixture(t)
	rule, err := openinghours.OpenRanges([]openinghours.Range{
		mustRange(t, 8, 30, 9, 30),
	}, openinghours.RejectOverlap)
	if err != nil {
		t.Fatal(err)
	}
	schedule, err := openinghours.NewSchedule(openinghours.Config{
		Timezone: "America/New_York", Weekly: map[time.Weekday]openinghours.DayRule{time.Sunday: rule},
	})
	if err != nil {
		t.Fatal(err)
	}
	from := time.Date(2028, time.December, 31, 12, 0, 0, 0, time.UTC)
	want := from.Add(90 * time.Minute)
	type queryResult struct {
		transition openinghours.Transition
		err        error
	}
	result := make(chan queryResult, 1)
	done := make(chan struct{})
	// This operation runs only in the owned pinned child. On failure the child
	// exits and the parent joins it; no noncooperative worker escapes that process.
	go func() {
		defer close(done)
		transition, err := schedule.NextTransition(from, 3*time.Hour)
		result <- queryResult{transition: transition, err: err}
	}()
	deadline := time.NewTimer(time.Second)
	defer deadline.Stop()
	var observed queryResult
	select {
	case observed = <-result:
		<-done
	case <-deadline.C:
		t.Fatal("finite stale-year search did not return its public boundary within the owned operation bound")
	}
	transition, err := observed.transition, observed.err
	if err != nil || !transition.Instant.Equal(want) || transition.Kind != openinghours.TransitionOpen {
		t.Fatalf("finite stale-year search=%#v error=%v; want opening %s", transition, err, want)
	}
}

func TestR04QueriesRejectOutOfCivilDomainBeforeOffsetTraversal(t *testing.T) {
	if r04RunWithPinnedZone(t) {
		return
	}
	schedule, _ := r04Fixture(t)
	for _, instant := range []time.Time{
		time.Date(0, time.January, 1, 0, 0, 0, 0, time.UTC),
		time.Date(10000, time.January, 2, 0, 0, 0, 0, time.UTC),
		time.Unix(1<<58, 0),
		time.Unix(-(1 << 58), 0),
	} {
		if state, err := schedule.IsOpen(instant); !openinghours.IsCode(err, openinghours.CodeInvalidDate) || state.Open {
			t.Errorf("out-of-domain point=%#v error=%v; want invalid date", state, err)
		}
		for _, query := range []struct {
			name string
			call func(time.Time, time.Duration) (openinghours.Transition, error)
		}{{"forward", schedule.NextTransition}, {"reverse", schedule.PreviousTransition}, {"opening", schedule.NextOpening}, {"closing", schedule.NextClosing}} {
			transition, err := query.call(instant, time.Hour)
			if !openinghours.IsCode(err, openinghours.CodeInvalidDate) || !transition.Instant.IsZero() {
				t.Errorf("out-of-domain %s=%#v error=%v; want invalid date without partial transition", query.name, transition, err)
			}
		}
		ranges, err := schedule.EffectiveInstantRanges(instant, instant.Add(time.Hour))
		if !openinghours.IsCode(err, openinghours.CodeInvalidDate) || ranges != nil {
			t.Errorf("out-of-domain ranges=%#v error=%v; want invalid date without partial ranges", ranges, err)
		}
		duration, err := schedule.OpenDuration(instant, instant.Add(time.Hour))
		if !openinghours.IsCode(err, openinghours.CodeInvalidDate) || duration != 0 {
			t.Errorf("out-of-domain duration=%s error=%v; want invalid date without partial duration", duration, err)
		}
	}
	// A zero schedule remains an unconditional closed value, not a zone query.
	var zero openinghours.Schedule
	huge := time.Unix(1<<58, 0)
	if state, err := zero.IsOpen(huge); err != nil || state.Open {
		t.Fatalf("zero point control=%#v %v", state, err)
	}
	if _, err := zero.NextTransition(huge, time.Hour); !openinghours.IsCode(err, openinghours.CodeSearchExhausted) {
		t.Fatalf("zero search control=%v", err)
	}
	if ranges, err := zero.EffectiveInstantRanges(huge, huge.Add(time.Hour)); err != nil || len(ranges) != 0 {
		t.Fatalf("zero expansion control=%#v %v", ranges, err)
	}
}
