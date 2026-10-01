package openinghours_test

import (
	"encoding/json"
	"errors"
	"strings"
	"testing"
	"time"

	openinghours "github.com/faustbrian/go-opening-hours/v2"
)

func TestExceptionProvenanceRequiresValidUTF8(t *testing.T) {
	for _, field := range []string{"source", "revision"} {
		t.Run(field, func(t *testing.T) {
			config := openinghours.ExceptionConfig{
				Date: openinghours.MustDate(2026, time.January, 1), Operation: openinghours.ExceptionClose,
				Source: "calendar", Revision: "2026",
			}
			if field == "source" {
				config.Source = string([]byte{0xff})
			} else {
				config.Revision = string([]byte{0xff})
			}
			got, err := openinghours.NewException(config)
			if !openinghours.IsCode(err, openinghours.CodeInvalidState) || got.Date().IsValid() || got.Source() != "" || got.Revision() != "" {
				t.Fatalf("invalid %s accepted: error=%v; want zero exception and invalid_state", field, err)
			}
		})
	}
}

func TestExceptionProvenanceValidationPreservesErrorPriority(t *testing.T) {
	invalidText := string([]byte{0xff})
	for _, test := range []struct {
		name   string
		config openinghours.ExceptionConfig
		code   openinghours.Code
	}{
		{"date", openinghours.ExceptionConfig{Operation: openinghours.ExceptionClose, Source: invalidText, Revision: "r"}, openinghours.CodeInvalidDate},
		{"operation", openinghours.ExceptionConfig{Date: openinghours.MustDate(2026, time.January, 1), Operation: openinghours.ExceptionClose + 1, Source: invalidText, Revision: "r"}, openinghours.CodeInvalidState},
		{"priority", openinghours.ExceptionConfig{Date: openinghours.MustDate(2026, time.January, 1), Operation: openinghours.ExceptionClose, Priority: 1_000_001, Source: invalidText, Revision: "r"}, openinghours.CodeInvalidState},
		{"length", openinghours.ExceptionConfig{Date: openinghours.MustDate(2026, time.January, 1), Operation: openinghours.ExceptionClose, Source: invalidText + strings.Repeat("x", 128), Revision: "r"}, openinghours.CodeLimitExceeded},
	} {
		t.Run(test.name, func(t *testing.T) {
			if _, err := openinghours.NewException(test.config); !openinghours.IsCode(err, test.code) {
				t.Fatalf("provenance validation priority = %v; want %s", err, test.code)
			}
		})
	}
}

func TestUnicodeExceptionIdentitiesPersistLosslessly(t *testing.T) {
	sources := []string{"café/北", "café/南"}
	revisions := []string{"révision-一", "révision-二"}
	input := make([]openinghours.Exception, len(sources))
	for index := range input {
		var err error
		input[index], err = openinghours.NewException(openinghours.ExceptionConfig{
			Date: openinghours.MustDate(2026, time.January, 1), Operation: openinghours.ExceptionClose,
			Priority: index + 1, Source: sources[index], Revision: revisions[index],
		})
		if err != nil || input[index].Source() != sources[index] || input[index].Revision() != revisions[index] {
			t.Fatalf("Unicode identity %d changed: %v", index, err)
		}
	}
	schedule, err := openinghours.NewSchedule(openinghours.Config{Timezone: "UTC", Exceptions: input})
	if err != nil {
		t.Fatal(err)
	}
	encoded, err := schedule.CanonicalJSON()
	if err != nil {
		t.Fatal(err)
	}
	var document struct {
		Exceptions []struct{ Source, Revision string }
	}
	if err := json.Unmarshal(encoded, &document); err != nil {
		t.Fatal(err)
	}
	if len(document.Exceptions) != len(input) {
		t.Fatal("canonical encoding lost an exception")
	}
	for index, identity := range document.Exceptions {
		if identity.Source != sources[index] || identity.Revision != revisions[index] {
			t.Fatalf("canonical encoding changed identity %d", index)
		}
	}
	decoded, err := openinghours.ParseJSON(encoded)
	if err != nil || !decoded.Equal(schedule) {
		t.Fatalf("Unicode canonical round trip changed schedule: %v", err)
	}
	value, err := schedule.Value()
	if err != nil {
		t.Fatal(err)
	}
	var scanned openinghours.Schedule
	if err := scanned.Scan(value); err != nil || !scanned.Equal(schedule) {
		t.Fatalf("Unicode SQL round trip changed schedule: %v", err)
	}
}

func TestScheduleRejectsInvalidDirectException(t *testing.T) {
	got, err := openinghours.NewSchedule(openinghours.Config{
		Timezone: "UTC", Exceptions: []openinghours.Exception{{}},
	})
	if !openinghours.IsCode(err, openinghours.CodeInvalidState) || !got.Equal(openinghours.Schedule{}) {
		t.Fatalf("invalid direct exception = %v, error=%v; want zero schedule and invalid_state", got, err)
	}
}

func TestExceptionAdmissionPreservesValidationOrder(t *testing.T) {
	oversized := make([]openinghours.Exception, openinghours.MaxExceptions+1)
	invalidDate := openinghours.Date{}
	for _, test := range []struct {
		name   string
		config openinghours.Config
		code   openinghours.Code
		op     string
	}{
		{"timezone", openinghours.Config{Timezone: "invalid", Exceptions: oversized}, openinghours.CodeInvalidTimezone, "new schedule"},
		{"metadata", openinghours.Config{Timezone: "UTC", Metadata: openinghours.Metadata{Label: strings.Repeat("x", 257)}, Exceptions: oversized}, openinghours.CodeLimitExceeded, "validate metadata"},
		{"date", openinghours.Config{Timezone: "UTC", EffectiveStart: &invalidDate, Exceptions: oversized}, openinghours.CodeInvalidDate, "new schedule"},
		{"weekday", openinghours.Config{Timezone: "UTC", Weekly: map[time.Weekday]openinghours.DayRule{8: openinghours.Closed()}, Exceptions: oversized}, openinghours.CodeInvalidWeekday, "new schedule"},
		{"invalid set", openinghours.Config{Timezone: "UTC", Exceptions: oversized, ExceptionSets: []openinghours.ExceptionSet{{}}}, openinghours.CodeInvalidState, "new schedule"},
		{"direct limit", openinghours.Config{Timezone: "UTC", Exceptions: oversized}, openinghours.CodeLimitExceeded, "new schedule"},
	} {
		t.Run(test.name, func(t *testing.T) {
			got, err := openinghours.NewSchedule(test.config)
			var typed *openinghours.Error
			if !errors.As(err, &typed) || typed.Code != test.code || typed.Op != test.op || !got.Equal(openinghours.Schedule{}) {
				t.Fatalf("admission = %v, error=%v; want zero schedule and %s/%s", got, err, test.op, test.code)
			}
		})
	}
}

func TestInvalidFirstExceptionSetPrecedesLaterLimit(t *testing.T) {
	closure, err := openinghours.NewException(openinghours.ExceptionConfig{
		Date:      openinghours.MustDate(2026, time.January, 1),
		Operation: openinghours.ExceptionClose, Source: "calendar", Revision: "2026",
	})
	if err != nil {
		t.Fatal(err)
	}
	input := make([]openinghours.Exception, openinghours.MaxExceptions)
	for index := range input {
		input[index] = closure
	}
	later, err := openinghours.NewExceptionSet("later", input)
	if err != nil {
		t.Fatal(err)
	}
	got, err := openinghours.NewSchedule(openinghours.Config{
		Timezone: "UTC", Exceptions: []openinghours.Exception{closure},
		ExceptionSets: []openinghours.ExceptionSet{{}, later},
	})
	if !openinghours.IsCode(err, openinghours.CodeInvalidState) || !got.Equal(openinghours.Schedule{}) {
		t.Fatalf("invalid first set = %v, error=%v; want zero schedule and invalid_state", got, err)
	}
}

func TestNamedExceptionSetIsCopiedAndPreserved(t *testing.T) {
	date := openinghours.MustDate(2026, time.December, 24)
	closure, _ := openinghours.NewException(openinghours.ExceptionConfig{
		Date: date, Operation: openinghours.ExceptionClose,
		Priority: 100, Source: "holidays", Revision: "2026",
	})
	input := []openinghours.Exception{closure}
	set, err := openinghours.NewExceptionSet("public-holidays", input)
	if err != nil {
		t.Fatal(err)
	}
	input[0] = openinghours.Exception{}
	schedule, err := openinghours.NewSchedule(openinghours.Config{
		Timezone: "UTC", ExceptionSets: []openinghours.ExceptionSet{set},
	})
	if err != nil {
		t.Fatal(err)
	}
	encoded, err := schedule.CanonicalJSON()
	if err != nil {
		t.Fatal(err)
	}
	decoded, err := openinghours.ParseJSON(encoded)
	if err != nil || !schedule.Equal(decoded) {
		t.Fatalf("set round trip error=%v equal=%t", err, schedule.Equal(decoded))
	}
	if set.Name() != "public-holidays" || len(set.Exceptions()) != 1 ||
		set.Exceptions()[0].Set() != "public-holidays" {
		t.Fatalf("exception set = %#v", set)
	}
}

func TestExceptionRangeExpansionIsInclusiveAndBounded(t *testing.T) {
	start := openinghours.MustDate(2026, time.December, 24)
	end := openinghours.MustDate(2026, time.December, 26)
	set, err := openinghours.ExpandExceptionRange(openinghours.ExceptionRangeConfig{
		Name: "holiday", Start: start, End: end, MaximumDates: 3,
		Operation: openinghours.ExceptionClose, Priority: 10,
		Source: "calendar", Revision: "2026",
	})
	if err != nil {
		t.Fatal(err)
	}
	if len(set.Exceptions()) != 3 {
		t.Fatalf("expanded dates = %d, want 3", len(set.Exceptions()))
	}

	_, err = openinghours.ExpandExceptionRange(openinghours.ExceptionRangeConfig{
		Name: "holiday", Start: start, End: end, MaximumDates: 2,
		Operation: openinghours.ExceptionClose, Priority: 10,
		Source: "calendar", Revision: "2026",
	})
	if !openinghours.IsCode(err, openinghours.CodeLimitExceeded) {
		t.Fatalf("ExpandExceptionRange() error = %v, want limit exceeded", err)
	}
}
