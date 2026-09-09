// Package openinghourscalendar is the compatibility path for
// [github.com/faustbrian/go-opening-hours/adapters/calendar].
//
// Deprecated: use github.com/faustbrian/go-opening-hours/adapters/calendar.
// This package remains supported for the longer of 180 days after successor
// public availability and two subsequently published stable minor releases.
package openinghourscalendar

import (
	calendar "github.com/faustbrian/go-calendar"
	"github.com/faustbrian/go-calendar/business"
	openinghours "github.com/faustbrian/go-opening-hours"
	canonical "github.com/faustbrian/go-opening-hours/adapters/calendar"
)

var (
	// ErrInvalidInput reports an invalid date, calendar, source, or range.
	ErrInvalidInput = canonical.ErrInvalidInput
	// ErrExpansionLimit reports a non-positive or exhausted date bound.
	ErrExpansionLimit = canonical.ErrExpansionLimit
)

// FromDate converts a valid calendar civil date without timezone inference.
func FromDate(date calendar.Date) (openinghours.Date, error) {
	return canonical.FromDate(date)
}

// ToDate converts a valid opening-hours civil date without timezone inference.
func ToDate(date openinghours.Date) (calendar.Date, error) {
	return canonical.ToDate(date)
}

// HolidayClosures resolves holidays in an inclusive, explicitly bounded civil
// date range to deterministic exact-date closure exceptions.
func HolidayClosures(businessCalendar business.Calendar, start, end calendar.Date,
	maximumDates, priority int, source string,
) ([]openinghours.Exception, error) {
	return canonical.HolidayClosures(
		businessCalendar, start, end, maximumDates, priority, source,
	)
}
