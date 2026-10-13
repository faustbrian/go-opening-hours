// Package openinghoursvalidation is the compatibility path for
// [github.com/faustbrian/go-opening-hours/v5/adapters/validation].
// Calls delegate to the canonical adapter with the same ownership/resource
// limits; path-specific named types and shared sentinel identities are retained.
// No external source or background lifecycle is added here.
//
// Deprecated: use github.com/faustbrian/go-opening-hours/v5/adapters/validation.
// This package remains supported for the longer of 180 days after successor
// public availability and two subsequently published stable minor releases.
package openinghoursvalidation

import (
	openinghours "github.com/faustbrian/go-opening-hours/v5"
	canonical "github.com/faustbrian/go-opening-hours/v5/adapters/validation"
	validation "github.com/faustbrian/go-validation/v2"
)

// CodeInvalidSchedule is the stable validation violation code.
const CodeInvalidSchedule = canonical.CodeInvalidSchedule

// Validate requires successful canonical serialization followed by accepted
// parsing. It does not compare the reparsed value or assert external conformance.
func Validate(schedule openinghours.Schedule) error {
	return canonical.Validate(schedule)
}

// Validator returns a deterministic validation adapter.
func Validator() validation.Validator[openinghours.Schedule] {
	return canonical.Validator()
}

// ValidationError retains a distinct public diagnostic type for this import
// path. Validate currently returns encoding/parsing errors, not this type.
type ValidationError struct{}

// Error returns the stable diagnostic without rendering a schedule.
func (*ValidationError) Error() string {
	return "openinghoursvalidation: canonical round trip mismatch"
}
