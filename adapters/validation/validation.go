// Package openinghoursvalidation adapts schedule validation to the shared
// validation contract.
// Use Validate or Validator for serialization/parse admission. The adapter
// owns neither source IO nor an application validation lifecycle and does not
// establish official conformance or compare arbitrary algebraic equivalence.
package openinghoursvalidation

import (
	openinghours "github.com/faustbrian/go-opening-hours/v4"
	validation "github.com/faustbrian/go-validation/v2"
)

// CodeInvalidSchedule is the stable validation violation code.
const CodeInvalidSchedule = "opening_hours.invalid_schedule"

// Validate requires successful canonical serialization followed by accepted
// parsing. It does not compare the reparsed value or assert external conformance.
func Validate(schedule openinghours.Schedule) error {
	encoded, err := schedule.CanonicalJSON()
	if err != nil {
		return err
	}
	_, err = openinghours.ParseJSON(encoded)
	return err
}

// Validator returns a deterministic validation adapter.
func Validator() validation.Validator[openinghours.Schedule] {
	return validation.ValidatorFunc[openinghours.Schedule](func(
		ctx validation.Context, schedule openinghours.Schedule,
	) validation.Report {
		report := validation.NewReport(ctx.Limits())
		if err := Validate(schedule); err != nil {
			return report.Add(validation.NewViolation(
				ctx.Path(), CodeInvalidSchedule, validation.Error, nil, err,
			))
		}
		return report
	})
}

// ValidationError retains a distinct public diagnostic type for this import
// path. Validate currently returns encoding/parsing errors, not this type.
type ValidationError struct{}

// Error returns the stable, data-free diagnostic.
func (*ValidationError) Error() string {
	return "openinghoursvalidation: canonical round trip mismatch"
}
