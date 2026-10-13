package openinghours_test

import (
	"testing"
	"time"

	openinghours "github.com/faustbrian/go-opening-hours/v5"
)

func TestR04PreviousTransitionIncludesExactOldestOffsetBoundary(t *testing.T) {
	if r04RunWithPinnedZone(t) {
		return
	}
	schedule, _ := r04Fixture(t)
	closing := time.Date(2026, time.November, 1, 6, 0, 0, 0, time.UTC)
	from := closing.Add(15 * time.Minute)
	for _, horizon := range []time.Duration{15 * time.Minute, 15*time.Minute + 2*time.Nanosecond} {
		transition, err := schedule.PreviousTransition(from, horizon)
		if err != nil || !transition.Instant.Equal(closing) || transition.Kind != openinghours.TransitionClose {
			t.Errorf("reverse horizon %s=%#v error=%v; want exact-oldest close %s", horizon, transition, err, closing)
		}
	}
	if _, err := schedule.PreviousTransition(from, 15*time.Minute-time.Nanosecond); !openinghours.IsCode(err, openinghours.CodeSearchExhausted) {
		t.Errorf("offset boundary outside horizon admitted: %v", err)
	}
	before, beforeErr := schedule.IsOpen(closing.Add(-time.Nanosecond))
	after, afterErr := schedule.IsOpen(closing)
	if beforeErr != nil || afterErr != nil || !before.Open || after.Open {
		t.Fatalf("exact-oldest boundary is not an actual closure: before=%#v %v after=%#v %v", before, beforeErr, after, afterErr)
	}
}
