package compile_test

import (
	"fmt"
	"strings"
	"testing"
	"time"

	openinghours "github.com/faustbrian/go-opening-hours/v2"
	"github.com/faustbrian/go-opening-hours/v2/compile"
)

func TestMaximumCompositionPreservesCanonicalAndCompiledQueries(t *testing.T) {
	start, err := openinghours.NewLocalTime(9, 0, 0, 0)
	if err != nil {
		t.Fatal(err)
	}
	end, err := openinghours.NewLocalTime(12, 0, 0, 0)
	if err != nil {
		t.Fatal(err)
	}
	rangeValue, err := openinghours.NewRange(start, end)
	if err != nil {
		t.Fatal(err)
	}
	rule, err := openinghours.OpenRanges([]openinghours.Range{rangeValue}, openinghours.RejectOverlap)
	if err != nil {
		t.Fatal(err)
	}
	exception, err := openinghours.NewException(openinghours.ExceptionConfig{
		Date:      openinghours.MustDate(2026, time.January, 6),
		Operation: openinghours.ExceptionReplace, Rule: rule,
		Source: "calendar", Revision: "2026",
	})
	if err != nil {
		t.Fatal(err)
	}
	leaf, err := openinghours.NewSchedule(openinghours.Config{
		Timezone: "UTC",
		Weekly: map[time.Weekday]openinghours.DayRule{
			time.Monday: openinghours.OpenAllDay(), time.Tuesday: rule,
		},
		Exceptions: []openinghours.Exception{exception},
	})
	if err != nil {
		t.Fatal(err)
	}
	plain, err := openinghours.NewSchedule(openinghours.Config{
		Timezone: "UTC", Weekly: map[time.Weekday]openinghours.DayRule{time.Monday: openinghours.OpenAllDay()},
	})
	if err != nil {
		t.Fatal(err)
	}
	schedule := leaf
	for depth := 1; depth < openinghours.MaxCompositionDepth; depth++ {
		schedule, err = schedule.Union(plain)
		if err != nil {
			t.Fatal(err)
		}
	}
	instant := time.Date(2026, time.January, 5, 10, 0, 0, 0, time.UTC)
	if got, err := schedule.IsOpen(instant); err != nil || !got.Open {
		t.Fatalf("original Monday query = %#v, error=%v", got, err)
	}
	encoded, err := schedule.CanonicalJSON()
	if err != nil {
		t.Fatal(err)
	}
	if len(encoded) > 12<<10 {
		t.Fatalf("fixture exceeds the intended 12 KiB scope: %d bytes", len(encoded))
	}
	t.Logf("maximum-depth fixture: %d bytes", len(encoded))
	t.Run("canonical round trip", func(t *testing.T) {
		decoded, err := openinghours.ParseJSON(encoded)
		if err != nil || !decoded.Equal(schedule) {
			t.Fatalf("maximum-depth round trip: error=%v, equal=%t", err, decoded.Equal(schedule))
		}
	})
	t.Run("compiled query", func(t *testing.T) {
		index, err := compile.New(schedule)
		if err != nil {
			t.Fatal(err)
		}
		got, err := index.IsOpen(instant)
		if err != nil || !got.Open || !index.Schedule().Equal(schedule) {
			t.Fatalf("compiled Monday query = %#v, error=%v, equal=%t", got, err, index.Schedule().Equal(schedule))
		}
	})
}

func TestIndexPreservesQueriesAndIsSafeForConcurrentReads(t *testing.T) {
	start, _ := openinghours.NewLocalTime(9, 0, 0, 0)
	end, _ := openinghours.NewLocalTime(12, 0, 0, 0)
	item, _ := openinghours.NewRange(start, end)
	rule, _ := openinghours.OpenRanges([]openinghours.Range{item}, openinghours.RejectOverlap)
	schedule, _ := openinghours.NewSchedule(openinghours.Config{
		Timezone: "UTC", Weekly: map[time.Weekday]openinghours.DayRule{time.Monday: rule},
	})
	index, err := compile.New(schedule)
	if err != nil {
		t.Fatal(err)
	}
	instant := time.Date(2026, time.January, 5, 10, 0, 0, 0, time.UTC)

	errors := make(chan error, 32)
	for range 32 {
		go func() {
			result, queryErr := index.IsOpen(instant)
			if queryErr != nil {
				errors <- queryErr
				return
			}
			if !result.Open {
				errors <- &closedError{}
				return
			}
			errors <- nil
		}()
	}
	for range 32 {
		if err := <-errors; err != nil {
			t.Fatal(err)
		}
	}
	local, _ := openinghours.NewLocalTime(10, 0, 0, 0)
	date := openinghours.MustDate(2026, time.January, 5)
	if result, queryErr := index.IsOpenLocal(date, local, openinghours.RejectDST); queryErr != nil || !result.Open {
		t.Fatalf("IsOpenLocal() = %#v error=%v", result, queryErr)
	}
	if transition, queryErr := index.NextTransition(
		time.Date(2026, time.January, 5, 8, 0, 0, 0, time.UTC), 2*time.Hour,
	); queryErr != nil || transition.Kind != openinghours.TransitionOpen {
		t.Fatalf("NextTransition() = %#v error=%v", transition, queryErr)
	}
	if !index.Schedule().Equal(schedule) {
		t.Fatal("Schedule() changed prepared schedule")
	}
}

func TestNewRejectsScheduleBeyondCanonicalLimit(t *testing.T) {
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
	if _, err := compile.New(schedule); !openinghours.IsCode(err, openinghours.CodeLimitExceeded) {
		t.Fatalf("compile.New error = %v", err)
	}
}

type closedError struct{}

func (*closedError) Error() string { return "compiled schedule unexpectedly closed" }
