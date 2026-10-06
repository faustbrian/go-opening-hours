package openinghours_test

import (
	"encoding/json"
	"strconv"
	"testing"
	"time"

	openinghours "github.com/faustbrian/go-opening-hours/v4"
)

func budgetException(t *testing.T, source string) openinghours.Exception {
	t.Helper()
	value, err := openinghours.NewException(openinghours.ExceptionConfig{
		Date:      openinghours.MustDate(2026, time.January, 1),
		Operation: openinghours.ExceptionClose, Source: source, Revision: "2026",
	})
	if err != nil {
		t.Fatal(err)
	}
	return value
}

func TestConsecutiveExceptionSetsPreserveAllIdentities(t *testing.T) {
	first, err := openinghours.NewExceptionSet("first", []openinghours.Exception{budgetException(t, "one")})
	if err != nil {
		t.Fatal(err)
	}
	second, err := openinghours.NewExceptionSet("second", []openinghours.Exception{
		budgetException(t, "two"), budgetException(t, "three"),
	})
	if err != nil {
		t.Fatal(err)
	}
	schedule, err := openinghours.NewSchedule(openinghours.Config{
		Timezone: "UTC", ExceptionSets: []openinghours.ExceptionSet{first, second},
		ConflictPolicy: openinghours.ResolveCanonical,
	})
	if err != nil {
		t.Fatal(err)
	}
	encoded, err := schedule.CanonicalJSON()
	if err != nil {
		t.Fatal(err)
	}
	var document struct {
		Exceptions []struct{ Source, Revision, Set string }
	}
	if err := json.Unmarshal(encoded, &document); err != nil {
		t.Fatal(err)
	}
	want := map[string]string{"one": "first", "two": "second", "three": "second"}
	if len(document.Exceptions) != len(want) {
		t.Fatalf("exception count = %d, want %d", len(document.Exceptions), len(want))
	}
	for _, value := range document.Exceptions {
		set, ok := want[value.Source]
		if !ok || value.Set != set || value.Revision != "2026" {
			t.Fatalf("unexpected exception identity: %#v", value)
		}
		delete(want, value.Source)
	}
	if len(want) != 0 {
		t.Fatalf("lost exception identities: %v", want)
	}
}

func TestAggregateExceptionLimitPrecedesLaterMalformedSet(t *testing.T) {
	input := make([]openinghours.Exception, openinghours.MaxExceptions)
	for index := range input {
		input[index] = budgetException(t, strconv.Itoa(index))
	}
	full, err := openinghours.NewExceptionSet("full", input)
	if err != nil {
		t.Fatal(err)
	}
	later, err := openinghours.NewExceptionSet("later", []openinghours.Exception{budgetException(t, "later")})
	if err != nil {
		t.Fatal(err)
	}
	got, err := openinghours.NewSchedule(openinghours.Config{
		Timezone: "UTC", ExceptionSets: []openinghours.ExceptionSet{full, later, {}},
	})
	if !openinghours.IsCode(err, openinghours.CodeLimitExceeded) || !got.Equal(openinghours.Schedule{}) {
		t.Fatalf("aggregate admission = %v, %v; want zero schedule and limit_exceeded", got, err)
	}
}
