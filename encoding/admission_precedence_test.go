package encoding_test

import (
	"testing"

	openinghoursencoding "github.com/faustbrian/go-opening-hours/v2/encoding"
)

func TestImportLimitsRejectBeforeMalformedInput(t *testing.T) {
	for _, test := range []struct {
		name   string
		limits openinghoursencoding.ImportLimits
	}{
		{"negative days", openinghoursencoding.ImportLimits{MaximumDays: -1, MaximumRangesPerDay: 64}},
		{"zero days", openinghoursencoding.ImportLimits{MaximumDays: 0, MaximumRangesPerDay: 64}},
		{"excess days", openinghoursencoding.ImportLimits{MaximumDays: 8, MaximumRangesPerDay: 64}},
		{"negative ranges", openinghoursencoding.ImportLimits{MaximumDays: 7, MaximumRangesPerDay: -1}},
		{"zero ranges", openinghoursencoding.ImportLimits{MaximumDays: 7, MaximumRangesPerDay: 0}},
		{"excess ranges", openinghoursencoding.ImportLimits{MaximumDays: 7, MaximumRangesPerDay: 65}},
	} {
		t.Run(test.name, func(t *testing.T) {
			t.Run("location", func(t *testing.T) {
				_, err := openinghoursencoding.ImportLocation(
					"invalid-zone",
					map[string][]openinghoursencoding.Slot{"monday": {{From: "invalid", To: "invalid"}}},
					test.limits,
				)
				assertImportErrorKind(t, err, "limit")
			})
			t.Run("spatie", func(t *testing.T) {
				_, err := openinghoursencoding.ImportSpatie(
					"invalid-zone", map[string][]string{"monday": {"invalid"}}, test.limits,
				)
				assertImportErrorKind(t, err, "limit")
			})
		})
	}
}
