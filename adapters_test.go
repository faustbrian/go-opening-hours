package openinghours_test

import (
	"strings"
	"testing"
	"time"

	calendar "github.com/faustbrian/go-calendar/v2"
	clock "github.com/faustbrian/go-clock"
	openinghours "github.com/faustbrian/go-opening-hours/v4"
	openinghoursconfig "github.com/faustbrian/go-opening-hours/v4/adapters/config"
	openinghoursvalidation "github.com/faustbrian/go-opening-hours/v4/adapters/validation"
	openinghourswire "github.com/faustbrian/go-opening-hours/v4/adapters/wire"
	openinghoursencoding "github.com/faustbrian/go-opening-hours/v4/encoding"
)

type fixedClock struct{ now time.Time }

func (clock fixedClock) Now() time.Time { return clock.now }

type nilClock struct{}

func (*nilClock) Now() time.Time { panic("typed-nil clock invoked") }

type panicClock struct{}

func (panicClock) Now() time.Time { panic("clock invoked") }

var _ clock.Clock = fixedClock{}

func TestCivilDateIsOwnedByGoCalendar(t *testing.T) {
	owned := calendar.MustDate(2026, time.January, 5)
	date := acceptOpeningHoursDate(owned)
	roundTrip := acceptCalendarDate(date)
	if !roundTrip.Equal(owned) {
		t.Fatalf("date alias round trip = %v, want %v", roundTrip, owned)
	}
}

func acceptOpeningHoursDate(date openinghours.Date) openinghours.Date { return date }
func acceptCalendarDate(date calendar.Date) calendar.Date             { return date }

func TestScheduleUsesBoundedCalendarZoneLoader(t *testing.T) {
	_, err := openinghours.NewSchedule(openinghours.Config{
		Timezone: strings.Repeat("x", 256),
	})
	if !openinghours.IsCode(err, openinghours.CodeInvalidTimezone) {
		t.Fatalf("oversized timezone error = %v", err)
	}
}

func TestClockCapabilityIsExplicit(t *testing.T) {
	schedule := scheduleWithMonday(t, mustRange(t, 9, 0, 12, 0))
	result, err := schedule.IsOpenNow(fixedClock{
		now: time.Date(2026, time.January, 5, 10, 0, 0, 0, time.UTC),
	})
	if err != nil || !result.Open {
		t.Fatalf("IsOpenNow() = %#v, error=%v", result, err)
	}
	_, err = schedule.IsOpenNow(nil)
	if !openinghours.IsCode(err, openinghours.CodeInvalidClock) {
		t.Fatalf("IsOpenNow(nil) error = %v, want invalid clock", err)
	}
}

func TestIsOpenNowRejectsTypedNilClock(t *testing.T) {
	var clock *nilClock
	if _, err := (openinghours.Schedule{}).IsOpenNow(clock); !openinghours.IsCode(err, openinghours.CodeInvalidClock) {
		t.Fatalf("IsOpenNow(typed nil) error = %v", err)
	}
	if _, err := (openinghours.Schedule{}).IsOpenNow(panicClock{}); !openinghours.IsCode(err, openinghours.CodeInvalidClock) {
		t.Fatalf("IsOpenNow(panicking clock) error = %v", err)
	}
}

func TestEncodingWireConfigAndValidationAdaptersAgree(t *testing.T) {
	want := scheduleWithMonday(t, mustRange(t, 9, 0, 12, 0))
	encoded, err := openinghoursencoding.Marshal(want)
	if err != nil {
		t.Fatal(err)
	}
	fromEncoding, err := openinghoursencoding.Unmarshal(encoded)
	if err != nil {
		t.Fatal(err)
	}
	codec := openinghourswire.Codec{}
	wireBytes, err := codec.Encode(want)
	if err != nil {
		t.Fatal(err)
	}
	fromWire, err := codec.Decode(wireBytes)
	if err != nil {
		t.Fatal(err)
	}
	fromConfig, err := openinghoursconfig.Parse(string(encoded))
	if err != nil {
		t.Fatal(err)
	}
	if err := openinghoursvalidation.Validate(want); err != nil {
		t.Fatal(err)
	}
	for name, got := range map[string]openinghours.Schedule{
		"encoding": fromEncoding, "wire": fromWire, "config": fromConfig,
	} {
		if !want.Equal(got) {
			t.Errorf("%s adapter changed schedule", name)
		}
	}
}
