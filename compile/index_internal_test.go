package compile

import (
	"testing"
	"time"

	openinghours "github.com/faustbrian/go-opening-hours/v3"
)

func TestParseIndexRejectsMalformedJSON(t *testing.T) {
	index, err := parseIndex([]byte(`{"version":`))
	if !openinghours.IsCode(err, openinghours.CodeInvalidEncoding) {
		t.Fatalf("malformed JSON error = %v, want invalid_encoding", err)
	}
	if index != (Index{}) {
		t.Fatal("malformed JSON published an index")
	}
}

func TestParseIndexPreservesCanonicalOpenQuery(t *testing.T) {
	schedule, err := openinghours.NewSchedule(openinghours.Config{
		Timezone: "UTC", Weekly: map[time.Weekday]openinghours.DayRule{
			time.Monday: openinghours.OpenAllDay(),
		},
	})
	if err != nil {
		t.Fatal(err)
	}
	encoded, err := schedule.CanonicalJSON()
	if err != nil {
		t.Fatal(err)
	}
	index, err := parseIndex(encoded)
	if err != nil {
		t.Fatal(err)
	}
	result, err := index.IsOpen(time.Date(2026, time.January, 5, 10, 0, 0, 0, time.UTC))
	if err != nil || !result.Open || !index.Schedule().Equal(schedule) {
		t.Fatalf("canonical index query = %#v, error=%v, equal=%t", result, err, index.Schedule().Equal(schedule))
	}
}
