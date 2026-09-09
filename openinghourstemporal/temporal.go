// Package openinghourstemporal is the compatibility path for
// [github.com/faustbrian/go-opening-hours/adapters/temporal].
//
// Deprecated: use github.com/faustbrian/go-opening-hours/adapters/temporal.
// This package remains supported for the longer of 180 days after successor
// public availability and two subsequently published stable minor releases.
package openinghourstemporal

import (
	openinghours "github.com/faustbrian/go-opening-hours"
	canonical "github.com/faustbrian/go-opening-hours/adapters/temporal"
	"github.com/faustbrian/go-temporal/timeofday"
)

// ErrLossyMapping reports an interval whose state or bounds cannot be
// represented without changing semantics.
var ErrLossyMapping = canonical.ErrLossyMapping

// RangeFromInterval converts an ordinary or circular start-inclusive,
// end-exclusive interval.
func RangeFromInterval(interval timeofday.Interval) (openinghours.Range, error) {
	return canonical.RangeFromInterval(interval)
}

// IntervalFromRange converts a range while making fractional precision
// explicit. Digits must exactly represent both endpoints.
func IntervalFromRange(value openinghours.Range, digits int) (timeofday.Interval, error) {
	return canonical.IntervalFromRange(value, digits)
}

// RuleFromIntervals converts an explicitly bounded interval collection. An
// empty or single collapsed interval maps to closed; full day must stand alone.
func RuleFromIntervals(intervals []timeofday.Interval,
	policy openinghours.OverlapPolicy,
) (openinghours.DayRule, error) {
	return canonical.RuleFromIntervals(intervals, policy)
}
