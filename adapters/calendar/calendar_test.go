package openinghourscalendar_test

import (
	"errors"
	"testing"
	"time"

	calendar "github.com/faustbrian/go-calendar/v2"
	"github.com/faustbrian/go-calendar/v2/business"
	openinghours "github.com/faustbrian/go-opening-hours/v6"
	openinghourscalendar "github.com/faustbrian/go-opening-hours/v6/adapters/calendar"
)

func TestCalendarErrorsAndBounds(t *testing.T) {
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
	if _, err := openinghourscalendar.HolidayClosures(businessCalendar, date, date, 4097, 0, "source"); !errors.Is(err, openinghourscalendar.ErrExpansionLimit) {
		t.Fatalf("package maximum error = %v", err)
	}
	if got, err := openinghourscalendar.HolidayClosures(businessCalendar, date, date, openinghours.MaxExceptions, 0, "source"); err != nil || len(got) != 1 {
		t.Fatalf("exact package maximum = %#v, %v", got, err)
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
