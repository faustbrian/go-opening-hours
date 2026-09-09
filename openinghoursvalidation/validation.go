// Package openinghoursvalidation is the compatibility path for
// [github.com/faustbrian/go-opening-hours/adapters/validation].
//
// Deprecated: use github.com/faustbrian/go-opening-hours/adapters/validation.
// This package remains supported for the longer of 180 days after successor
// public availability and two subsequently published stable minor releases.
package openinghoursvalidation

import (
	openinghours "github.com/faustbrian/go-opening-hours"
	canonical "github.com/faustbrian/go-opening-hours/adapters/validation"
	validation "github.com/faustbrian/go-validation"
)

// CodeInvalidSchedule is the stable validation violation code.
const CodeInvalidSchedule = canonical.CodeInvalidSchedule

// Validate proves that a schedule has a lossless strict canonical round trip.
func Validate(schedule openinghours.Schedule) error {
	return canonical.Validate(schedule)
}

// Validator returns a deterministic validation adapter.
func Validator() validation.Validator[openinghours.Schedule] {
	return canonical.Validator()
}

// ValidationError reports a canonical round-trip mismatch without data.
type ValidationError struct{}

func (*ValidationError) Error() string {
	return "openinghoursvalidation: canonical round trip mismatch"
}
