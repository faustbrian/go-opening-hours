//nolint:staticcheck // This compatibility test intentionally exercises deprecated facades.
package openinghourscalendar_test

//lint:file-ignore SA1019 This compatibility test intentionally exercises deprecated facades.

import (
	"errors"
	"testing"
	"time"

	calendar "github.com/faustbrian/go-calendar"
	"github.com/faustbrian/go-calendar/business"
	openinghours "github.com/faustbrian/go-opening-hours"
	openinghourscalendar "github.com/faustbrian/go-opening-hours/adapters/calendar"
	legacy "github.com/faustbrian/go-opening-hours/openinghourscalendar"
)

func TestDateConversionAndHolidayClosuresMatchLegacy(t *testing.T) {
	date := calendar.MustDate(2026, time.December, 25)
	converted, err := openinghourscalendar.FromDate(date)
	if err != nil || converted != openinghours.MustDate(2026, time.December, 25) {
		t.Fatalf("FromDate() = %#v, %v", converted, err)
	}
	if got, err := openinghourscalendar.ToDate(converted); err != nil || !got.Equal(date) {
		t.Fatalf("ToDate() = %#v, %v", got, err)
	}

	holiday := business.MustHoliday(date, "Christmas", nil)
	businessCalendar, err := business.NewCalendar(business.Config{
		Revision: "2026", Holidays: []business.Holiday{holiday},
	})
	if err != nil {
		t.Fatal(err)
	}
	start := calendar.MustDate(2026, time.December, 24)
	end := calendar.MustDate(2026, time.December, 26)
	want, err := legacy.HolidayClosures(businessCalendar, start, end, 3, 100, "public-holidays")
	if err != nil {
		t.Fatal(err)
	}
	got, err := openinghourscalendar.HolidayClosures(businessCalendar, start, end, 3, 100, "public-holidays")
	if err != nil || len(got) != 1 || len(want) != 1 ||
		got[0].Date() != want[0].Date() || got[0].Operation() != want[0].Operation() ||
		got[0].Priority() != want[0].Priority() || got[0].Source() != want[0].Source() ||
		got[0].Revision() != want[0].Revision() || got[0].Set() != want[0].Set() {
		t.Fatalf("HolidayClosures() = %#v, %v; legacy = %#v", got, err, want)
	}
}

func TestCalendarErrorsAndBoundsMatchLegacy(t *testing.T) {
	//nolint:errorlint // Exact sentinel identity is the compatibility contract.
	if openinghourscalendar.ErrInvalidInput != legacy.ErrInvalidInput ||
		openinghourscalendar.ErrExpansionLimit != legacy.ErrExpansionLimit {
		t.Fatal("calendar sentinel identity differs between canonical and legacy paths")
	}
	if _, err := openinghourscalendar.FromDate(calendar.Date{}); !errors.Is(err, openinghourscalendar.ErrInvalidInput) {
		t.Fatalf("FromDate() error = %v", err)
	}
	if _, err := openinghourscalendar.ToDate(openinghours.Date{}); !errors.Is(err, openinghourscalendar.ErrInvalidInput) {
		t.Fatalf("ToDate() error = %v", err)
	}
	date := calendar.MustDate(2026, time.December, 25)
	holiday := business.MustHoliday(date, "Christmas", nil)
	businessCalendar, err := business.NewCalendar(business.Config{
		Revision: "2026", Holidays: []business.Holiday{holiday},
	})
	if err != nil {
		t.Fatal(err)
	}
	if _, err := openinghourscalendar.HolidayClosures(businessCalendar, date, date, 0, 0, "source"); !errors.Is(err, openinghourscalendar.ErrInvalidInput) {
		t.Fatalf("zero bound error = %v", err)
	}
	end := calendar.MustDate(2026, time.December, 26)
	if _, err := openinghourscalendar.HolidayClosures(businessCalendar, date, end, 1, 0, "source"); !errors.Is(err, openinghourscalendar.ErrExpansionLimit) {
		t.Fatalf("exhaustion error = %v", err)
	}
	if _, err := openinghourscalendar.HolidayClosures(businessCalendar, date, date, 1, 1_000_001, "source"); err == nil {
		t.Fatal("invalid exception provenance accepted")
	}
	maximum := calendar.MustDate(9999, time.December, 31)
	maximumCalendar, err := business.NewCalendar(business.Config{
		Revision: "max", Holidays: []business.Holiday{
			business.MustHoliday(maximum, "Maximum date", nil),
		},
	})
	if err != nil {
		t.Fatal(err)
	}
	if got, err := openinghourscalendar.HolidayClosures(maximumCalendar, maximum, maximum, 1, 0, "source"); err != nil || len(got) != 1 {
		t.Fatalf("maximum-date closure = %#v, %v", got, err)
	}
}
