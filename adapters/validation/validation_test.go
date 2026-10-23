package openinghoursvalidation_test

import (
	"fmt"
	"strings"
	"testing"
	"time"

	openinghours "github.com/faustbrian/go-opening-hours/v6"
	openinghoursvalidation "github.com/faustbrian/go-opening-hours/v6/adapters/validation"
	validation "github.com/faustbrian/go-validation/v3"
)

func TestValidatorAndErrorContract(t *testing.T) {
	schedule, err := openinghours.NewSchedule(openinghours.Config{Timezone: "UTC"})
	if err != nil {
		t.Fatal(err)
	}
	if err := openinghoursvalidation.Validate(schedule); err != nil {
		t.Fatal(err)
	}
	ctx, err := validation.NewContext(validation.DefaultLimits())
	if err != nil {
		t.Fatal(err)
	}
	if report := openinghoursvalidation.Validator().Validate(ctx, schedule); !report.Empty() {
		t.Fatalf("valid schedule report = %s", report.String())
	}
	var nilError *openinghoursvalidation.ValidationError
	if nilError.Error() != "openinghoursvalidation: canonical round trip mismatch" {
		t.Fatalf("typed-nil Error() = %q", nilError.Error())
	}
}

func TestValidateRejectsScheduleBeyondCanonicalLimit(t *testing.T) {
	schedule := oversizedSchedule(t)
	if err := openinghoursvalidation.Validate(schedule); !openinghours.IsCode(err, openinghours.CodeLimitExceeded) {
		t.Fatalf("Validate() error = %v", err)
	}
	ctx, err := validation.NewContext(validation.DefaultLimits())
	if err != nil {
		t.Fatal(err)
	}
	report := openinghoursvalidation.Validator().Validate(ctx, schedule)
	if !report.HasErrors() || !report.HasCode(openinghoursvalidation.CodeInvalidSchedule) {
		t.Fatalf("invalid schedule report = %s", report.String())
	}
	violations := report.Violations()
	if len(violations) != 1 || violations[0].Severity() != validation.Error ||
		violations[0].Path().String() != ctx.Path().String() ||
		!openinghours.IsCode(violations[0].Cause(), openinghours.CodeLimitExceeded) {
		t.Fatalf("invalid schedule violation = %#v", violations)
	}
}

func oversizedSchedule(t *testing.T) openinghours.Schedule {
	t.Helper()
	exceptions := make([]openinghours.Exception, 4096)
	date := openinghours.MustDate(2026, time.January, 1)
	for index := range exceptions {
		revision := fmt.Sprintf("%04d%s", index, strings.Repeat("r", 124))
		var err error
		exceptions[index], err = openinghours.NewException(openinghours.ExceptionConfig{
			Date: date, Operation: openinghours.ExceptionClose,
			Source: strings.Repeat("s", 128), Revision: revision,
		})
		if err != nil {
			t.Fatal(err)
		}
	}
	schedule, err := openinghours.NewSchedule(openinghours.Config{
		Timezone: "UTC", Exceptions: exceptions,
		ConflictPolicy: openinghours.ResolveCanonical,
	})
	if err != nil {
		t.Fatal(err)
	}
	return schedule
}
