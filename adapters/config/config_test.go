//nolint:staticcheck // This compatibility test intentionally exercises deprecated facades.
package openinghoursconfig_test

//lint:file-ignore SA1019 This compatibility test intentionally exercises deprecated facades.

import (
	"reflect"
	"strings"
	"testing"

	configdecode "github.com/faustbrian/go-config/decode"
	openinghours "github.com/faustbrian/go-opening-hours"
	openinghoursconfig "github.com/faustbrian/go-opening-hours/adapters/config"
	legacy "github.com/faustbrian/go-opening-hours/openinghoursconfig"
)

func TestValueRoundTripAndAtomicFailure(t *testing.T) {
	schedule, err := openinghours.NewSchedule(openinghours.Config{Timezone: "UTC"})
	if err != nil {
		t.Fatal(err)
	}
	encoded, err := schedule.CanonicalJSON()
	if err != nil {
		t.Fatal(err)
	}
	var value openinghoursconfig.Value
	if err := configdecode.Value(string(encoded), &value); err != nil {
		t.Fatal(err)
	}
	if !value.Schedule().Equal(schedule) {
		t.Fatal("decoded value changed schedule")
	}
	if err := value.UnmarshalConfigValue(`{}`); err == nil {
		t.Fatal("invalid value accepted")
	}
	if !value.Schedule().Equal(schedule) {
		t.Fatal("failed decode changed the prior schedule")
	}
	if text, err := openinghoursconfig.NewValue(schedule).MarshalText(); err != nil || string(text) != string(encoded) {
		t.Fatalf("MarshalText() = %s, %v", text, err)
	}
}

func TestValueErrorsAndNamedIdentity(t *testing.T) {
	//nolint:errorlint // Exact sentinel identity is the compatibility contract.
	if openinghoursconfig.ErrInvalidValue != legacy.ErrInvalidValue {
		t.Fatal("config sentinel identity differs between canonical and legacy paths")
	}
	var nilValue *openinghoursconfig.Value
	//nolint:errorlint // Exact sentinel return is the compatibility contract.
	if err := nilValue.UnmarshalConfigValue("{}"); err != openinghoursconfig.ErrInvalidValue {
		t.Fatalf("nil receiver error = %v", err)
	}
	var value openinghoursconfig.Value
	//nolint:errorlint // Exact sentinel return is the compatibility contract.
	if err := value.UnmarshalConfigValue(42); err != openinghoursconfig.ErrInvalidValue {
		t.Fatalf("non-string error = %v", err)
	}
	if reflect.TypeOf(value) == reflect.TypeOf(legacy.Value{}) {
		t.Fatal("canonical and legacy Value reflection identities match")
	}
	if _, err := openinghoursconfig.Parse(`{"secret":"do-not-repeat"}`); err == nil || strings.Contains(err.Error(), "do-not-repeat") {
		t.Fatalf("invalid config error = %v", err)
	}
	zeroText, err := (openinghoursconfig.Value{}).MarshalText()
	if err != nil {
		t.Fatalf("zero Value MarshalText() error = %v", err)
	}
	zero, err := openinghoursconfig.Parse(string(zeroText))
	if err != nil || !zero.Equal((openinghoursconfig.Value{}).Schedule()) {
		t.Fatalf("zero Value round trip = %#v, %v", zero, err)
	}
}
