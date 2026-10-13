package openinghours_test

import (
	"crypto/sha256"
	"encoding/hex"
	"os"
	"testing"
	"time"

	openinghours "github.com/faustbrian/go-opening-hours/v5"
)

func TestR04FixedZoneAdmitsExactCollectedFragmentQuota(t *testing.T) {
	if r04RunWithPinnedZone(t) {
		return
	}
	_, _ = r04Fixture(t)
	raw, err := os.ReadFile("testdata/r04/zoneinfo/Etc/GMT-1")
	if err != nil {
		t.Fatal(err)
	}
	digest := sha256.Sum256(raw)
	if hex.EncodeToString(digest[:]) != "4bcd52f59d3e57ed01e54fb44b43e76f1f1fbf6887b701352eb95993e7242eda" {
		t.Fatal("fixed-zone input differs from the attributed Go 1.27.2 archive entry")
	}
	location, err := time.LoadLocation("Etc/GMT-1")
	if err != nil {
		t.Fatal(err)
	}
	start := time.Date(1, time.January, 1, 0, 30, 0, 0, location).UTC()
	end := start.Add(128 * 24 * time.Hour)
	for _, point := range []time.Time{start, {}, end} {
		_, offset := point.In(location).Zone()
		zoneStart, zoneEnd := point.In(location).ZoneBounds()
		if offset != 3600 || !zoneStart.IsZero() || !zoneEnd.IsZero() {
			t.Fatalf("fixed-zone control differs: offset=%d bounds=[%s,%s)", offset, zoneStart, zoneEnd)
		}
	}
	local := start.In(location)
	if local.Year() != 1 || local.Month() != time.January || local.Day() != 1 || local.Hour() != 0 || local.Minute() != 30 {
		t.Fatal("quota input must start at a valid year-one civil instant")
	}
	periods := []openinghours.Range{mustRange(t, 0, 30, 1, 30)}
	for index := range 63 {
		minute := 120 + 2*index
		periods = append(periods, mustRange(t, minute/60, minute%60, (minute+1)/60, (minute+1)%60))
	}
	rule, err := openinghours.OpenRanges(periods, openinghours.RejectOverlapAndAdjacent)
	if err != nil {
		t.Fatal(err)
	}
	weekly := make(map[time.Weekday]openinghours.DayRule)
	for weekday := time.Sunday; weekday <= time.Saturday; weekday++ {
		weekly[weekday] = rule
	}
	schedule, err := openinghours.NewSchedule(openinghours.Config{Timezone: "Etc/GMT-1", Weekly: weekly})
	if err != nil {
		t.Fatal(err)
	}
	ranges, err := schedule.EffectiveInstantRanges(start, end)
	if err != nil || len(ranges) != 8192 {
		t.Fatalf("exact collected-fragment quota: ranges=%d error=%v; want 8192", len(ranges), err)
	}
	for day := range 128 {
		dayStart := start.Add(time.Duration(day) * 24 * time.Hour)
		for period := range 64 {
			wantStart, wantEnd := dayStart, dayStart.Add(time.Hour)
			if period > 0 {
				wantStart = dayStart.Add(time.Duration(90+2*(period-1)) * time.Minute)
				wantEnd = wantStart.Add(time.Minute)
			}
			got := ranges[day*64+period]
			if !got.Start.Equal(wantStart) || !got.End.Equal(wantEnd) {
				t.Fatalf("fragment day=%d period=%d: %#v; want [%s,%s)", day, period, got, wantStart, wantEnd)
			}
		}
	}
	duration, err := schedule.OpenDuration(start, end)
	if err != nil || duration != 128*123*time.Minute {
		t.Errorf("exact-quota duration=%s error=%v; want 128 days of 123 open minutes", duration, err)
	}
	// One minute of the next day's first period adds a real 8193rd fragment.
	ranges, err = schedule.EffectiveInstantRanges(start, end.Add(time.Minute))
	if !openinghours.IsCode(err, openinghours.CodeLimitExceeded) || ranges != nil {
		t.Errorf("one extra fragment must reject without partial output: ranges=%d error=%v", len(ranges), err)
	}
	duration, err = schedule.OpenDuration(start, end.Add(time.Minute))
	if !openinghours.IsCode(err, openinghours.CodeLimitExceeded) || duration != 0 {
		t.Errorf("one extra fragment must discard partial duration: %s %v", duration, err)
	}
}
