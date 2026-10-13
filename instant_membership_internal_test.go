package openinghours

import (
	"os"
	"testing"
	"time"
)

func TestInstantMembershipGapAndFullDayControls(t *testing.T) {
	schedule := r04PinnedSchedule(t, 2, 30, 3, 30)
	start := time.Date(2026, time.March, 8, 6, 0, 0, 0, time.UTC)
	end := start.Add(2 * time.Hour)
	ranges, err := schedule.EffectiveInstantRanges(start, end)
	wantStart, wantEnd := start.Add(time.Hour), start.Add(90*time.Minute)
	if err != nil || len(ranges) != 1 || !ranges[0].Start.Equal(wantStart) || !ranges[0].End.Equal(wantEnd) {
		t.Fatalf("gap excludes nonexistent wall time only: %#v %v", ranges, err)
	}
	schedule.data.weekly[time.Sunday] = OpenAllDay()
	start = time.Date(2026, time.November, 1, 4, 0, 0, 0, time.UTC)
	end = start.Add(25 * time.Hour)
	ranges, err = schedule.EffectiveInstantRanges(start, end)
	if err != nil || len(ranges) != 1 || !ranges[0].Start.Equal(start) || !ranges[0].End.Equal(end) {
		t.Fatalf("all-day fold must normalize to 25h without a false closure: %#v %v", ranges, err)
	}
	duration, err := schedule.OpenDuration(start, end)
	if err != nil || duration != 25*time.Hour {
		t.Fatalf("all-day fold duration=%s %v", duration, err)
	}
}

func TestInstantMembershipRetainsRevisitedCaseyDates(t *testing.T) {
	// Exact official provider bytes; public construction is exercised separately.
	raw, err := os.ReadFile("testdata/r04/zoneinfo/Antarctica/Casey")
	if err != nil {
		t.Fatal(err)
	}
	location, err := time.LoadLocationFromTZData("Antarctica/Casey", raw)
	if err != nil {
		t.Fatal(err)
	}
	base := time.Date(2010, time.March, 4, 0, 0, 0, 0, time.UTC)
	for _, tc := range []struct {
		name                string
		weekday             time.Weekday
		hour                int
		start, end, opening time.Time
	}{
		{"reversed-outer-dates", time.Thursday, 23, base.Add(14*time.Hour + 59*time.Minute), base.Add(15*time.Hour + 46*time.Minute), base.Add(15*time.Hour + 30*time.Minute)},
		{"intermediate-earlier-date", time.Thursday, 23, base.Add(14*time.Hour + 59*time.Minute), base.Add(16 * time.Hour), base.Add(15*time.Hour + 30*time.Minute)},
		{"intermediate-later-date", time.Friday, 1, base.Add(10*time.Hour + 59*time.Minute), base.Add(15*time.Hour + 46*time.Minute), base.Add(14*time.Hour + 30*time.Minute)},
	} {
		t.Run(tc.name, func(t *testing.T) {
			from, _ := NewLocalTime(tc.hour, 30, 0, 0)
			until, _ := NewLocalTime(tc.hour, 45, 0, 0)
			interval, _ := NewRange(from, until)
			rule, _ := OpenRanges([]Range{interval}, RejectOverlap)
			schedule, err := NewSchedule(Config{Timezone: "UTC", Weekly: map[time.Weekday]DayRule{tc.weekday: rule}})
			if err != nil {
				t.Fatal(err)
			}
			schedule.data.location, schedule.data.timezone = location, "Antarctica/Casey"
			ranges, err := schedule.EffectiveInstantRanges(tc.start, tc.end)
			if err != nil || len(ranges) != 1 || !ranges[0].Start.Equal(tc.opening) || !ranges[0].End.Equal(tc.opening.Add(15*time.Minute)) {
				t.Fatalf("revisited-date expansion=%#v %v; want %s+15min", ranges, err, tc.opening)
			}
			duration, err := schedule.OpenDuration(tc.start, tc.end)
			if err != nil || duration != 15*time.Minute {
				t.Fatalf("revisited-date duration=%s %v", duration, err)
			}
		})
	}
}

func TestInstantMembershipDiscardsCollectedRangesOnLaterDateError(t *testing.T) {
	date := MustDate(2026, time.January, 1)
	schedule := internalSchedule(t, Config{
		Weekly:       map[time.Weekday]DayRule{time.Thursday: OpenAllDay()},
		EffectiveEnd: &date, OutsideEffective: OutsideError,
	})
	start := time.Date(2026, time.January, 1, 0, 0, 0, 0, time.UTC)
	if state, err := schedule.IsOpen(start); err != nil || !state.Open {
		t.Fatal("positive first-day control must collect an open fragment")
	}
	ranges, err := schedule.EffectiveInstantRanges(start, start.Add(48*time.Hour))
	if !IsCode(err, CodeOutsideEffectiveRange) || ranges != nil {
		t.Fatalf("late date error must discard the earlier open fragment: %#v %v", ranges, err)
	}
	duration, err := schedule.OpenDuration(start, start.Add(48*time.Hour))
	if !IsCode(err, CodeOutsideEffectiveRange) || duration != 0 {
		t.Fatalf("late date error must not return partial duration: %s %v", duration, err)
	}
}
