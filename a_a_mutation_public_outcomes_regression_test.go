package openinghours

import (
	"context"
	"crypto/sha256"
	"encoding/hex"
	"math"
	"os"
	"os/exec"
	"path/filepath"
	"regexp"
	"strconv"
	"strings"
	"testing"
	"time"
)

// Each potentially nonadvancing query stays in the owning process. Its one-second
// outcome assertion must precede the child and mutation-engine watchdogs. A parent
// kill is cleanup, not a successful mutation assertion.
func r04MutationChild(t *testing.T) bool {
	t.Helper()
	if os.Getenv("OPENING_HOURS_R04_MUTATION_PROBE") == t.Name() {
		return false
	}
	directory, err := filepath.Abs("testdata/r04/zoneinfo")
	if err != nil {
		t.Fatal(err)
	}
	executable, err := os.Executable()
	if err != nil {
		t.Fatal(err)
	}
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()
	// #nosec G204 -- this running owning test binary and exact test name are controlled; no shell executes.
	command := exec.CommandContext(ctx, executable, "-test.run=^"+regexp.QuoteMeta(t.Name())+"$", "-test.count=1", "-test.timeout=2s")
	for _, entry := range os.Environ() {
		if strings.HasPrefix(entry, "GOLIB_GREMLINS_COVERAGE_PROFILE=") || strings.HasPrefix(entry, "ZONEINFO=") || strings.HasPrefix(entry, "OPENING_HOURS_R04_MUTATION_PROBE=") {
			continue
		}
		command.Env = append(command.Env, entry)
	}
	command.Env = append(command.Env, "OPENING_HOURS_R04_MUTATION_PROBE="+t.Name(), "ZONEINFO="+directory)
	output, err := command.CombinedOutput()
	if err != nil {
		if len(output) > 8192 {
			output = output[:8192]
		}
		t.Fatalf("joined owning mutation probe: %v: %s", err, output)
	}
	return true
}

func r04BoundedTransition(t *testing.T, query func() (Transition, error)) Transition {
	t.Helper()
	type result struct {
		transition Transition
		err        error
	}
	results := make(chan result, 1)
	done := make(chan struct{})
	go func() { defer close(done); value, err := query(); results <- result{value, err} }()
	timer := time.NewTimer(time.Second)
	defer timer.Stop()
	select {
	case value := <-results:
		<-done
		if value.err != nil {
			t.Fatalf("R04_OUTCOME: finite public query failed: %v", value.err)
		}
		return value.transition
	case <-timer.C:
		t.Fatal("R04_OUTCOME: finite public query did not return before owning operation deadline")
	}
	return Transition{}
}

func r04DailyMutationSchedule(t *testing.T, zone string, startHour, endHour int) Schedule {
	t.Helper()
	rule := OpenAllDay()
	if startHour != 0 || endHour != 24 {
		if endHour == 24 {
			endHour = 0
		}
		start, err := NewLocalTime(startHour, 0, 0, 0)
		if err != nil {
			t.Fatal(err)
		}
		end, err := NewLocalTime(endHour, 0, 0, 0)
		if err != nil {
			t.Fatal(err)
		}
		period, err := NewRange(start, end)
		if err != nil {
			t.Fatal(err)
		}
		rule, err = OpenRanges([]Range{period}, RejectOverlap)
		if err != nil {
			t.Fatal(err)
		}
	}
	weekly := make(map[time.Weekday]DayRule)
	for weekday := time.Sunday; weekday <= time.Saturday; weekday++ {
		weekly[weekday] = rule
	}
	schedule, err := NewSchedule(Config{Timezone: zone, Weekly: weekly})
	if err != nil {
		t.Fatal(err)
	}
	return schedule
}

func TestR04MutationMidnightEndpointReturnsPublicClose(t *testing.T) {
	if r04MutationChild(t) {
		return
	}
	schedule := r04DailyMutationSchedule(t, "UTC", 23, 24)
	from := time.Date(2026, time.January, 1, 23, 30, 0, 0, time.UTC)
	want := time.Date(2026, time.January, 2, 0, 0, 0, 0, time.UTC)
	transition := r04BoundedTransition(t, func() (Transition, error) { return schedule.NextTransition(from, 2*time.Hour) })
	if !transition.Instant.Equal(want) || transition.Kind != TransitionClose {
		t.Fatalf("R04_OUTCOME: midnight close=%#v want%s", transition, want)
	}
	for _, point := range []struct {
		instant time.Time
		open    bool
	}{{want.Add(-time.Nanosecond), true}, {want, false}} {
		state, err := schedule.IsOpen(point.instant)
		if err != nil || state.Open != point.open {
			t.Fatalf("midnight point control=%#v %v", state, err)
		}
	}
}

func TestR04MutationNegativeOffsetRetainsLastCivilDate(t *testing.T) {
	if r04MutationChild(t) {
		return
	}
	raw, err := os.ReadFile("testdata/r04/zoneinfo/Etc/GMT+1")
	if err != nil {
		t.Fatal(err)
	}
	digest := sha256.Sum256(raw)
	if hex.EncodeToString(digest[:]) != "e4bf68f1311482d075d69a086a0f39bd176ad3c2cc0d9999e833e7ed4a8f2ff8" {
		t.Fatal("negative fixed zone input differs from authenticated archive")
	}
	location, err := time.LoadLocation("Etc/GMT+1")
	if err != nil {
		t.Fatal(err)
	}
	instant := time.Date(9999, time.December, 31, 23, 30, 0, 0, location)
	_, offset := instant.Zone()
	if offset != -3600 {
		t.Fatal("negative-offset fixture control differs")
	}
	schedule := r04DailyMutationSchedule(t, "Etc/GMT+1", 23, 24)
	state, err := schedule.IsOpen(instant)
	if err != nil || !state.Open {
		t.Fatalf("R04_OUTCOME: valid last civil instant rejected: %#v %v", state, err)
	}
}

func TestR04MutationCivilGuardRejectsMachineYearAlias(t *testing.T) {
	if strconv.IntSize != 32 {
		t.Skip("specific observable year-narrowing witness requires actual 32-bit runtime")
	}
	schedule := r04DailyMutationSchedule(t, "UTC", 0, 24)
	// Gregorian UTC 4294969322-01-01, whose mathematical year is 2^32+2026.
	// Keep the absolute instant literal representable on both host widths.
	instant := time.Unix(135536078568643200, 0)
	state, err := schedule.IsOpen(instant)
	if !IsCode(err, CodeInvalidDate) || state.Open {
		t.Fatalf("R04_OUTCOME: huge absolute instant admitted as machine-year alias: %#v %v", state, err)
	}
}

func TestR04MutationCivilMaximumRetainsPublicError(t *testing.T) {
	schedule := r04DailyMutationSchedule(t, "UTC", 0, 24)
	maximum := time.Date(10000, time.January, 1, 0, 0, 0, 0, time.UTC)
	var message string
	for _, instant := range []time.Time{maximum, maximum.Add(time.Second)} {
		state, err := schedule.IsOpen(instant)
		if !IsCode(err, CodeInvalidDate) || state.Open {
			t.Fatalf("R04_OUTCOME: excluded civil maximum admitted: %#v %v", state, err)
		}
		if message != "" && message != err.Error() {
			t.Fatalf("R04_OUTCOME: civil maximum public error projection differs: %q vs %q", message, err.Error())
		}
		message = err.Error()
	}
	last := maximum.Add(-time.Nanosecond)
	if state, err := schedule.IsOpen(last); err != nil || !state.Open {
		t.Fatalf("civil maximum predecessor control: %#v %v", state, err)
	}
}

func TestR04MutationCivilGuardRejectsWidePublicLocalOffset(t *testing.T) {
	if strconv.IntSize != 64 {
		t.Skip("specific wide public offset requires a 64-bit int")
	}
	if r04MutationChild(t) {
		return
	}
	previous := time.Local
	wideOffset := int64(math.MinInt64)
	time.Local = time.FixedZone("wide", int(wideOffset))
	defer func() { time.Local = previous }()
	schedule := r04DailyMutationSchedule(t, "Local", 0, 24)
	instant := time.Unix(-9223372035087550208, 0).UTC()
	if localized := instant.In(time.Local); localized.Year() != 2026 || localized.Month() != time.January || localized.Day() != 1 {
		t.Fatalf("wide public offset control differs: %s", localized)
	}
	state, err := schedule.IsOpen(instant)
	if !IsCode(err, CodeInvalidDate) || state.Open {
		t.Fatalf("R04_OUTCOME: wide public offset bypassed absolute civil admission: %#v %v", state, err)
	}
}

func TestR04MutationLateCivilSearchRetainsNearBoundary(t *testing.T) {
	if r04MutationChild(t) {
		return
	}
	r04MutationPinnedNY(t)
	location, err := time.LoadLocation("America/New_York")
	if err != nil {
		t.Fatal(err)
	}
	schedule := r04DailyMutationSchedule(t, "America/New_York", 23, 24)
	from := time.Date(9999, time.December, 30, 23, 30, 0, 0, location)
	want := time.Date(9999, time.December, 31, 0, 0, 0, 0, location)
	// A longer admitted horizon reaches POSIX offset boundaries in year10000;
	// it must retain the same earlier in-domain closing as the short control.
	for _, horizon := range []time.Duration{2 * time.Hour, 150 * 24 * time.Hour} {
		transition := r04BoundedTransition(t, func() (Transition, error) { return schedule.NextTransition(from, horizon) })
		if !transition.Instant.Equal(want) || transition.Kind != TransitionClose {
			t.Fatalf("R04_OUTCOME: near civil-boundary close=%#v want%s", transition, want)
		}
	}
}

func TestR04MutationStaleYearSearchAdvancesBeforeWatchdog(t *testing.T) {
	if r04MutationChild(t) {
		return
	}
	r04MutationPinnedNY(t)
	schedule := r04DailyMutationSchedule(t, "America/New_York", 8, 9)
	from := time.Date(2028, time.December, 31, 12, 0, 0, 0, time.UTC)
	transition := r04BoundedTransition(t, func() (Transition, error) { return schedule.NextTransition(from, 3*time.Hour) })
	want := time.Date(2028, time.December, 31, 13, 0, 0, 0, time.UTC)
	if !transition.Instant.Equal(want) || transition.Kind != TransitionOpen {
		t.Fatalf("R04_OUTCOME: stale-year next=%#v want%s", transition, want)
	}
}

func TestR04MutationYearEndHorizonRetainsNextValidCivilOpening(t *testing.T) {
	if r04MutationChild(t) {
		return
	}
	r04MutationPinnedNY(t)
	schedule := r04DailyMutationSchedule(t, "America/New_York", 12, 13)
	from := time.Date(9999, time.January, 1, 20, 0, 0, 0, time.UTC)
	want := time.Date(9999, time.January, 2, 17, 0, 0, 0, time.UTC)
	// Both horizons contain this next supported civil day. Extending the outer
	// deadline beyond the last supported year must not lose an earlier opening.
	for _, horizon := range []time.Duration{24 * time.Hour, 365 * 24 * time.Hour} {
		transition := r04BoundedTransition(t, func() (Transition, error) { return schedule.NextTransition(from, horizon) })
		if !transition.Instant.Equal(want) || transition.Kind != TransitionOpen {
			t.Fatalf("R04_OUTCOME: admitted long horizon lost next valid civil opening: %#v want%s", transition, want)
		}
	}
}

func TestR04MutationLateCivilHorizonRetainsLastDayOpening(t *testing.T) {
	if r04MutationChild(t) {
		return
	}
	r04MutationPinnedNY(t)
	schedule := r04DailyMutationSchedule(t, "America/New_York", 12, 13)
	from := time.Date(9999, time.December, 30, 20, 0, 0, 0, time.UTC)
	want := time.Date(9999, time.December, 31, 17, 0, 0, 0, time.UTC)
	for _, horizon := range []time.Duration{24 * time.Hour, 150 * 24 * time.Hour} {
		transition := r04BoundedTransition(t, func() (Transition, error) { return schedule.NextTransition(from, horizon) })
		if !transition.Instant.Equal(want) || transition.Kind != TransitionOpen {
			t.Fatalf("R04_OUTCOME: later unsupported offset boundary lost last supported-day opening: %#v want%s", transition, want)
		}
	}
}

func r04MutationPinnedNY(t *testing.T) {
	t.Helper()
	directory, err := filepath.Abs("testdata/r04/zoneinfo")
	if err != nil || os.Getenv("ZONEINFO") != directory {
		t.Fatal("probe must start with pinned public zoneinfo")
	}
	raw, err := os.ReadFile("testdata/r04/zoneinfo/America/New_York")
	if err != nil {
		t.Fatal(err)
	}
	hash := sha256.Sum256(raw)
	if hex.EncodeToString(hash[:]) != "d7f2206b3a45989fc9ad63d558922532fa7352280d5f87176bf1db79cb1d1fa9" {
		t.Fatal("NY fixture identity changed")
	}
}

func TestR04MutationForwardCivilClampIsObservedInOwningProfile(t *testing.T) {
	// Public UTC construction and query, directly instrumented in the owning
	// process. Earlier bounded owner tests kill nonadvancing mutants before this
	// control in the native failfast phase; this is not a new private-state route.
	schedule := r04DailyMutationSchedule(t, "UTC", 12, 13)
	from := time.Date(9999, time.January, 1, 20, 0, 0, 0, time.UTC)
	want := time.Date(9999, time.January, 2, 12, 0, 0, 0, time.UTC)
	transition, err := schedule.NextTransition(from, 365*24*time.Hour)
	if err != nil || !transition.Instant.Equal(want) || transition.Kind != TransitionOpen {
		t.Fatalf("R04_OUTCOME: direct public clamp=%#v %v want%s", transition, err, want)
	}
}
