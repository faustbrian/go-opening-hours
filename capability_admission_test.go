package openinghours_test

import (
	"testing"
	"time"

	openinghours "github.com/faustbrian/go-opening-hours/v4"
)

// Methods deliberately do not dereference the receiver: rejecting a nil
// capability must not depend on a collaborator panicking when invoked.
type nilSafeFunctionClock func()

func (nilSafeFunctionClock) Now() time.Time {
	return time.Date(2026, time.January, 5, 10, 0, 0, 0, time.UTC)
}

func (nilSafeFunctionClock) Since(time.Time) time.Duration { return time.Millisecond }

func (nilSafeFunctionClock) Measure() func() time.Duration {
	return func() time.Duration { return time.Millisecond }
}

func TestNilFunctionClockIsRejectedBeforeInvocation(t *testing.T) {
	schedule := scheduleWithMonday(t, mustRange(t, 9, 0, 12, 0))
	var clock nilSafeFunctionClock
	got, err := schedule.IsOpenNow(clock)
	if !openinghours.IsCode(err, openinghours.CodeInvalidClock) || got != (openinghours.Availability{}) {
		t.Fatalf("nil function clock = %#v, %v; want zero availability and invalid_clock", got, err)
	}
}

func TestNilFunctionElapsedClockIsUnmeasured(t *testing.T) {
	schedule := scheduleWithMonday(t, mustRange(t, 9, 0, 12, 0))
	instant := time.Date(2026, time.January, 5, 10, 0, 0, 0, time.UTC)
	var clock nilSafeFunctionClock
	var observed openinghours.Observation
	observer := func(value openinghours.Observation) { observed = value }
	availability, err := schedule.ObserveIsOpen(instant, clock, observer)
	if err != nil || !availability.Open || observed.Duration != 0 || observed.Outcome != openinghours.OutcomeOpen {
		t.Fatalf("nil elapsed clock availability = %#v, %#v, %v", availability, observed, err)
	}
	transition, err := schedule.ObserveNextTransition(instant, 3*time.Hour, clock, observer)
	want := time.Date(2026, time.January, 5, 12, 0, 0, 0, time.UTC)
	if err != nil || !transition.Instant.Equal(want) || observed.Duration != 0 || observed.Outcome != openinghours.OutcomeFound {
		t.Fatalf("nil elapsed clock transition = %#v, %#v, %v", transition, observed, err)
	}
}
