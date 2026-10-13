package openinghours

import (
	"os"
	"testing"
	"time"
)

// Inject only the exact provider value at the private immutable owning boundary.
// Public construction/environment routing is proved separately by isolated probes.
func r04PinnedSchedule(t *testing.T, startHour, startMinute, endHour, endMinute int) Schedule {
	t.Helper()
	raw, err := os.ReadFile("testdata/r04/zoneinfo/America/New_York")
	if err != nil {
		t.Fatal(err)
	}
	location, err := time.LoadLocationFromTZData("America/New_York", raw)
	if err != nil {
		t.Fatal(err)
	}
	start, err := NewLocalTime(startHour, startMinute, 0, 0)
	if err != nil {
		t.Fatal(err)
	}
	end, err := NewLocalTime(endHour, endMinute, 0, 0)
	if err != nil {
		t.Fatal(err)
	}
	interval, err := NewRange(start, end)
	if err != nil {
		t.Fatal(err)
	}
	day, err := OpenRanges([]Range{interval}, RejectOverlap)
	if err != nil {
		t.Fatal(err)
	}
	schedule, err := NewSchedule(Config{Timezone: "UTC", Weekly: map[time.Weekday]DayRule{time.Sunday: day}})
	if err != nil {
		t.Fatal(err)
	}
	schedule.data.location = location
	schedule.data.timezone = "America/New_York"
	return schedule
}

func TestSearchCandidatesRetainBothFoldEndpointKinds(t *testing.T) {
	for _, hours := range [][4]int{{1, 30, 2, 30}, {1, 0, 1, 30}} {
		schedule := r04PinnedSchedule(t, hours[0], hours[1], hours[2], hours[3])
		date := MustDate(2026, time.November, 1)
		endpoint := LocalTime{nanosecond: int64(90 * time.Minute)}
		candidates := schedule.searchBoundaryCandidates(date, endpoint.nanosecond)
		if len(candidates) != 2 || !candidates[0].Instant.Equal(time.Date(2026, time.November, 1, 5, 30, 0, 0, time.UTC)) || !candidates[1].Instant.Equal(time.Date(2026, time.November, 1, 6, 30, 0, 0, time.UTC)) {
			t.Fatalf("fold endpoint occurrences=%#v", candidates)
		}
		cursor := time.Date(2026, time.November, 1, 4, 59, 0, 0, time.UTC)
		for range 4 {
			next, err := schedule.NextTransition(cursor, 3*time.Hour)
			if err != nil {
				t.Fatal(err)
			}
			before, _ := schedule.IsOpen(next.Instant.Add(-time.Nanosecond))
			after, _ := schedule.IsOpen(next.Instant)
			if !next.Instant.After(cursor) || before.Open == after.Open {
				t.Fatal("candidate must advance actual point state")
			}
			previous, err := schedule.PreviousTransition(next.Instant.Add(time.Nanosecond), 3*time.Hour)
			if err != nil || !previous.Instant.Equal(next.Instant) || previous.Kind != next.Kind {
				t.Fatal("nearest reverse transition lost forward boundary")
			}
			cursor = next.Instant
		}
	}
}

func TestSearchOffsetGapAndStrictHorizonControls(t *testing.T) {
	schedule := r04PinnedSchedule(t, 2, 30, 3, 30)
	before := time.Date(2026, time.March, 8, 6, 59, 0, 0, time.UTC)
	jump := before.Add(time.Minute)
	opening, err := schedule.NextOpening(before, time.Minute)
	if err != nil || !opening.Instant.Equal(jump) || opening.Kind != TransitionOpen {
		t.Fatalf("gap jump opening=%#v %v", opening, err)
	}
	if _, err := schedule.NextTransition(before, time.Minute-time.Nanosecond); !IsCode(err, CodeSearchExhausted) {
		t.Fatal("candidate beyond horizon admitted")
	}
	closing, err := schedule.NextTransition(jump, time.Hour)
	if err != nil || !closing.Instant.Equal(jump.Add(30*time.Minute)) || closing.Kind != TransitionClose {
		t.Fatal("strict cursor did not select later actual close")
	}
	gapDate := MustDate(2026, time.March, 8)
	if candidates := schedule.searchBoundaryCandidates(gapDate, int64(150*time.Minute)); len(candidates) != 0 {
		t.Fatal("nonexistent civil endpoint gained a realization")
	}
	if candidates := schedule.searchBoundaryCandidates(MustDate(9999, time.December, 31), nanosecondsPerDay); len(candidates) != 0 {
		t.Fatal("out-of-domain next-date endpoint admitted")
	}
	for weekday := time.Sunday; weekday <= time.Saturday; weekday++ {
		schedule.data.weekly[weekday] = OpenAllDay()
	}
	if _, err := schedule.NextTransition(before, time.Hour); !IsCode(err, CodeSearchExhausted) {
		t.Fatal("offset-only change created phantom all-day transition")
	}
}

// Casey's three-hour fold revisits the preceding civil date. Chronological
// instant bounds need not have chronologically ordered localized dates.
func TestSearchFoldAcrossCivilMidnight(t *testing.T) {
	raw, err := os.ReadFile("testdata/r04/zoneinfo/Antarctica/Casey")
	if err != nil {
		t.Fatal(err)
	}
	location, err := time.LoadLocationFromTZData("Antarctica/Casey", raw)
	if err != nil {
		t.Fatal(err)
	}
	start, _ := NewLocalTime(23, 30, 0, 0)
	end, _ := NewLocalTime(23, 45, 0, 0)
	interval, _ := NewRange(start, end)
	day, _ := OpenRanges([]Range{interval}, RejectOverlap)
	schedule, err := NewSchedule(Config{Timezone: "UTC", Weekly: map[time.Weekday]DayRule{time.Thursday: day}})
	if err != nil {
		t.Fatal(err)
	}
	schedule.data.location = location
	schedule.data.timezone = "Antarctica/Casey"
	before := time.Date(2010, time.March, 4, 14, 59, 0, 0, time.UTC)
	opening := time.Date(2010, time.March, 4, 15, 30, 0, 0, time.UTC)
	closing := opening.Add(15 * time.Minute)
	next, err := schedule.NextTransition(before, time.Hour)
	if err != nil || !next.Instant.Equal(opening) || next.Kind != TransitionOpen {
		t.Errorf("cross-date next=%#v %v; want %s", next, err, opening)
	}
	// Both localized outer bounds can be Friday while the fold revisits Thursday.
	longer, err := schedule.NextTransition(before, 2*time.Hour)
	if err != nil || !longer.Instant.Equal(opening) {
		t.Errorf("intermediate earlier date opening=%#v %v", longer, err)
	}
	previous, err := schedule.PreviousTransition(closing.Add(time.Minute), 2*time.Hour)
	if err != nil || !previous.Instant.Equal(closing) || previous.Kind != TransitionClose {
		t.Errorf("cross-date previous=%#v %v; want %s", previous, err, closing)
	}
	for _, point := range []struct {
		instant time.Time
		open    bool
	}{{opening.Add(-time.Nanosecond), false}, {opening, true}, {closing.Add(-time.Nanosecond), true}, {closing, false}} {
		state, err := schedule.IsOpen(point.instant)
		if err != nil || state.Open != point.open {
			t.Error("cross-date fold point control disagrees")
		}
	}
	fridayStart, _ := NewLocalTime(1, 30, 0, 0)
	fridayEnd, _ := NewLocalTime(1, 45, 0, 0)
	fridayRange, _ := NewRange(fridayStart, fridayEnd)
	fridayRule, _ := OpenRanges([]Range{fridayRange}, RejectOverlap)
	schedule.data.weekly[time.Thursday] = Closed()
	schedule.data.weekly[time.Friday] = fridayRule
	// Both outer bounds are Thursday, but the intermediate instant span reaches Friday.
	intermediate, err := schedule.PreviousTransition(closing.Add(time.Minute), 4*time.Hour)
	want := time.Date(2010, time.March, 4, 14, 45, 0, 0, time.UTC)
	if err != nil || !intermediate.Instant.Equal(want) || intermediate.Kind != TransitionClose {
		t.Errorf("intermediate later date closing=%#v %v; want %s", intermediate, err, want)
	}

}
