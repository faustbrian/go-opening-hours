package openinghoursvalidation_test

import (
	"testing"

	openinghours "github.com/faustbrian/go-opening-hours/v6"
	openinghoursvalidation "github.com/faustbrian/go-opening-hours/v6/adapters/validation"
	validation "github.com/faustbrian/go-validation/v3"
)

func TestValidateAndError(t *testing.T) {
	schedule, _ := openinghours.NewSchedule(openinghours.Config{Timezone: "UTC"})
	if err := openinghoursvalidation.Validate(schedule); err != nil {
		t.Fatal(err)
	}
	if (&openinghoursvalidation.ValidationError{}).Error() != "openinghoursvalidation: canonical round trip mismatch" {
		t.Fatal("unexpected validation error text")
	}
}

func TestValidatorIntegratesWithGoValidation(t *testing.T) {
	ctx, err := validation.NewContext(validation.DefaultLimits())
	if err != nil {
		t.Fatal(err)
	}
	schedule, _ := openinghours.NewSchedule(openinghours.Config{Timezone: "UTC"})
	if report := openinghoursvalidation.Validator().Validate(ctx, schedule); !report.Empty() {
		t.Fatalf("valid schedule report = %s", report.String())
	}
	oversized := oversizedSchedule(t)
	report := openinghoursvalidation.Validator().Validate(ctx, oversized)
	if !report.HasErrors() || !report.HasCode(openinghoursvalidation.CodeInvalidSchedule) {
		t.Fatalf("invalid schedule report = %s", report.String())
	}
}
