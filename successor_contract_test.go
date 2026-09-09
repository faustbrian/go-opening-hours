//nolint:staticcheck // This compatibility test intentionally exercises deprecated facades.
package openinghours_test

//lint:file-ignore SA1019 This compatibility test intentionally exercises deprecated facades.

import (
	"errors"
	"reflect"
	"testing"
	"time"

	calendar "github.com/faustbrian/go-calendar"
	"github.com/faustbrian/go-calendar/business"
	openinghours "github.com/faustbrian/go-opening-hours"
	calendaradapter "github.com/faustbrian/go-opening-hours/adapters/calendar"
	configadapter "github.com/faustbrian/go-opening-hours/adapters/config"
	temporaladapter "github.com/faustbrian/go-opening-hours/adapters/temporal"
	validationadapter "github.com/faustbrian/go-opening-hours/adapters/validation"
	wireadapter "github.com/faustbrian/go-opening-hours/adapters/wire"
	calendarlegacy "github.com/faustbrian/go-opening-hours/openinghourscalendar"
	configlegacy "github.com/faustbrian/go-opening-hours/openinghoursconfig"
	temporallegacy "github.com/faustbrian/go-opening-hours/openinghourstemporal"
	validationlegacy "github.com/faustbrian/go-opening-hours/openinghoursvalidation"
	wirelegacy "github.com/faustbrian/go-opening-hours/openinghourswire"
	temporal "github.com/faustbrian/go-temporal"
	"github.com/faustbrian/go-temporal/timeofday"
)

func TestSuccessorAndLegacyPathsPreserveContracts(t *testing.T) {
	date := calendar.MustDate(2026, time.December, 25)
	holiday := business.MustHoliday(date, "Christmas", nil)
	businessCalendar, err := business.NewCalendar(business.Config{
		Revision: "2026", Holidays: []business.Holiday{holiday},
	})
	if err != nil {
		t.Fatal(err)
	}
	canonicalClosures, canonicalErr := calendaradapter.HolidayClosures(
		businessCalendar, date, date, 1, 0, "holidays",
	)
	legacyClosures, legacyErr := calendarlegacy.HolidayClosures(
		businessCalendar, date, date, 1, 0, "holidays",
	)
	if canonicalErr != nil || legacyErr != nil || len(canonicalClosures) != 1 ||
		len(legacyClosures) != 1 || canonicalClosures[0].Date() != legacyClosures[0].Date() {
		t.Fatalf("calendar parity = %#v/%v, %#v/%v", canonicalClosures, canonicalErr, legacyClosures, legacyErr)
	}

	schedule, err := openinghours.NewSchedule(openinghours.Config{Timezone: "UTC"})
	if err != nil {
		t.Fatal(err)
	}
	encoded, _ := schedule.CanonicalJSON()
	canonicalValue := configadapter.NewValue(schedule)
	legacyValue := configlegacy.NewValue(schedule)
	canonicalText, canonicalErr := canonicalValue.MarshalText()
	legacyText, legacyErr := legacyValue.MarshalText()
	if canonicalErr != nil || legacyErr != nil || string(canonicalText) != string(encoded) ||
		string(legacyText) != string(encoded) || reflect.TypeOf(canonicalValue) == reflect.TypeOf(legacyValue) {
		t.Fatal("config value behavior or named identity differs")
	}

	start, _ := timeofday.New(9, 0, 0, 0, 0)
	end, _ := timeofday.New(17, 0, 0, 0, 0)
	interval, _ := timeofday.Between(start, end, temporal.ClosedOpen)
	canonicalRange, canonicalErr := temporaladapter.RangeFromInterval(interval)
	legacyRange, legacyErr := temporallegacy.RangeFromInterval(interval)
	if canonicalErr != nil || legacyErr != nil || canonicalRange != legacyRange {
		t.Fatalf("temporal parity = %#v/%v, %#v/%v", canonicalRange, canonicalErr, legacyRange, legacyErr)
	}

	if err := validationadapter.Validate(schedule); err != nil {
		t.Fatal(err)
	}
	if err := validationlegacy.Validate(schedule); err != nil {
		t.Fatal(err)
	}
	canonicalValidationError := &validationadapter.ValidationError{}
	legacyValidationError := &validationlegacy.ValidationError{}
	if canonicalValidationError.Error() != legacyValidationError.Error() ||
		reflect.TypeOf(canonicalValidationError) == reflect.TypeOf(legacyValidationError) {
		t.Fatal("validation error behavior or named identity differs")
	}
	var canonicalTarget *validationadapter.ValidationError
	if errors.As(legacyValidationError, &canonicalTarget) {
		t.Fatal("legacy validation error matched successor named type")
	}

	canonicalCodec := wireadapter.Codec{}
	legacyCodec := wirelegacy.Codec{}
	canonicalBytes, canonicalErr := canonicalCodec.Encode(schedule)
	legacyBytes, legacyErr := legacyCodec.Encode(schedule)
	if canonicalErr != nil || legacyErr != nil || string(canonicalBytes) != string(legacyBytes) ||
		reflect.TypeOf(canonicalCodec) == reflect.TypeOf(legacyCodec) {
		t.Fatal("wire codec behavior or named identity differs")
	}
	//nolint:errorlint // Exact sentinel identity is the compatibility contract.
	if calendaradapter.ErrInvalidInput != calendarlegacy.ErrInvalidInput ||
		configadapter.ErrInvalidValue != configlegacy.ErrInvalidValue ||
		temporaladapter.ErrLossyMapping != temporallegacy.ErrLossyMapping ||
		wireadapter.Format != wirelegacy.Format || wireadapter.WireFormat != wirelegacy.WireFormat {
		t.Fatal("shared sentinels or constants differ")
	}
}
