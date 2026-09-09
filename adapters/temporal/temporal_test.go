//nolint:staticcheck // This compatibility test intentionally exercises deprecated facades.
package openinghourstemporal_test

//lint:file-ignore SA1019 This compatibility test intentionally exercises deprecated facades.

import (
	"errors"
	"testing"

	openinghours "github.com/faustbrian/go-opening-hours"
	openinghourstemporal "github.com/faustbrian/go-opening-hours/adapters/temporal"
	legacy "github.com/faustbrian/go-opening-hours/openinghourstemporal"
	temporal "github.com/faustbrian/go-temporal"
	"github.com/faustbrian/go-temporal/timeofday"
)

func TestRangeAndRuleConversionsMatchLegacy(t *testing.T) {
	start, _ := timeofday.New(22, 0, 0, 0, 0)
	end, _ := timeofday.New(2, 0, 0, 0, 0)
	interval, _ := timeofday.Between(start, end, temporal.ClosedOpen)
	got, err := openinghourstemporal.RangeFromInterval(interval)
	if err != nil || !got.Overnight() {
		t.Fatalf("RangeFromInterval() = %#v, %v", got, err)
	}
	want, err := legacy.RangeFromInterval(interval)
	if err != nil || got != want {
		t.Fatalf("canonical = %#v; legacy = %#v, %v", got, want, err)
	}
	roundTrip, err := openinghourstemporal.IntervalFromRange(got, 0)
	if err != nil || !roundTrip.Equal(interval) {
		t.Fatalf("IntervalFromRange() = %#v, %v", roundTrip, err)
	}
	rule, err := openinghourstemporal.RuleFromIntervals([]timeofday.Interval{interval}, openinghours.RejectOverlap)
	if err != nil || rule.State() != openinghours.DayOpenRanges {
		t.Fatalf("RuleFromIntervals() = %#v, %v", rule, err)
	}
}

func TestLossyMappingsAndSpecialRules(t *testing.T) {
	//nolint:errorlint // Exact sentinel identity is the compatibility contract.
	if openinghourstemporal.ErrLossyMapping != legacy.ErrLossyMapping {
		t.Fatal("temporal sentinel identity differs between canonical and legacy paths")
	}
	closed, _ := timeofday.Between(timeofday.Midnight(), timeofday.Noon(), temporal.Closed)
	if _, err := openinghourstemporal.RangeFromInterval(closed); !errors.Is(err, openinghourstemporal.ErrLossyMapping) {
		t.Fatalf("closed-bounds error = %v", err)
	}
	endBoundaryStart, _ := timeofday.Between(timeofday.EndOfDay(), timeofday.Noon(), temporal.ClosedOpen)
	if _, err := openinghourstemporal.RangeFromInterval(endBoundaryStart); !errors.Is(err, openinghourstemporal.ErrLossyMapping) {
		t.Fatalf("end-boundary-start error = %v", err)
	}
	if rule, err := openinghourstemporal.RuleFromIntervals(nil, openinghours.RejectOverlap); err != nil || rule.State() != openinghours.DayClosed {
		t.Fatalf("empty rule = %#v, %v", rule, err)
	}
	if rule, err := openinghourstemporal.RuleFromIntervals([]timeofday.Interval{timeofday.FullDay()}, openinghours.RejectOverlap); err != nil || rule.State() != openinghours.DayOpenAllDay {
		t.Fatalf("full-day rule = %#v, %v", rule, err)
	}
	if rule, err := openinghourstemporal.RuleFromIntervals(
		[]timeofday.Interval{timeofday.Collapsed(timeofday.Noon())}, openinghours.RejectOverlap,
	); err != nil || rule.State() != openinghours.DayClosed {
		t.Fatalf("collapsed rule = %#v, %v", rule, err)
	}
	if _, err := openinghourstemporal.RuleFromIntervals(
		[]timeofday.Interval{timeofday.FullDay(), closed}, openinghours.RejectOverlap,
	); !errors.Is(err, openinghourstemporal.ErrLossyMapping) {
		t.Fatalf("mixed full-day error = %v", err)
	}
	ordinary, _ := timeofday.Between(timeofday.Midnight(), timeofday.Noon(), temporal.ClosedOpen)
	if _, err := openinghourstemporal.RuleFromIntervals(
		[]timeofday.Interval{ordinary, closed}, openinghours.RejectOverlap,
	); !errors.Is(err, openinghourstemporal.ErrLossyMapping) {
		t.Fatalf("mixed bounds error = %v", err)
	}
	toEnd, _ := timeofday.Between(timeofday.Noon(), timeofday.EndOfDay(), temporal.ClosedOpen)
	if value, err := openinghourstemporal.RangeFromInterval(toEnd); err != nil || !value.Overnight() {
		t.Fatalf("end-boundary range = %#v, %v", value, err)
	}
	if _, err := openinghourstemporal.IntervalFromRange(openinghours.Range{}, 10); err == nil {
		t.Fatal("invalid temporal precision accepted")
	}
	preciseStart, _ := openinghours.NewLocalTime(1, 0, 0, 0)
	preciseEnd, _ := openinghours.NewLocalTime(2, 0, 0, 1)
	preciseRange, _ := openinghours.NewRange(preciseStart, preciseEnd)
	if _, err := openinghourstemporal.IntervalFromRange(preciseRange, 0); err == nil {
		t.Fatal("lossy end precision accepted")
	}
}
