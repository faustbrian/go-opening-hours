package openinghours_test

import (
	"testing"
	"time"

	openinghours "github.com/faustbrian/go-opening-hours/v5"
	prepared "github.com/faustbrian/go-opening-hours/v5/compile"
)

// The selected point-membership policy supersedes the original contiguous
// fold envelope retained as a historical countercase in the audit regression.
// The authenticated independent oracle supplies both disjoint open intervals.
func TestR04InstantExpansionMatchesPointMembership(t *testing.T) {
	if r04RunWithPinnedZone(t) {
		return
	}
	schedule, oracle := r04Fixture(t)
	points := oracle.Transitions
	cases := []struct {
		name       string
		start, end time.Time
		want       []openinghours.InstantRange
	}{
		{"entire-fold", points[0].Instant, points[3].Instant, []openinghours.InstantRange{
			{Start: points[0].Instant, End: points[1].Instant},
			{Start: points[2].Instant, End: points[3].Instant},
		}},
		{"closed-repeated-period", points[1].Instant, points[2].Instant, nil},
		{"clip-both-open-components", points[0].Instant.Add(15 * time.Minute), points[2].Instant.Add(15 * time.Minute), []openinghours.InstantRange{
			{Start: points[0].Instant.Add(15 * time.Minute), End: points[1].Instant},
			{Start: points[2].Instant, End: points[2].Instant.Add(15 * time.Minute)},
		}},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			ranges, err := schedule.EffectiveInstantRanges(tc.start, tc.end)
			if err != nil || len(ranges) != len(tc.want) {
				t.Fatalf("point-membership expansion=%#v error=%v; want %#v", ranges, err, tc.want)
			}
			for i, want := range tc.want {
				if !ranges[i].Start.Equal(want.Start) || !ranges[i].End.Equal(want.End) {
					t.Fatalf("range %d=%#v; want %#v", i, ranges[i], want)
				}
			}
			for _, point := range oracle.States {
				if point.Instant.Before(tc.start) || !point.Instant.Before(tc.end) {
					continue
				}
				contained := false
				for _, interval := range ranges {
					contained = contained || (!point.Instant.Before(interval.Start) && point.Instant.Before(interval.End))
				}
				if contained != point.Open {
					t.Fatalf("expanded membership at %s=%t; independent point oracle=%t", point.Instant, contained, point.Open)
				}
			}
		})
	}
}

func TestR04PointExpansionSurvivesCanonicalAndCompositionSeams(t *testing.T) {
	if r04RunWithPinnedZone(t) {
		return
	}
	schedule, oracle := r04Fixture(t)
	encoded, err := schedule.CanonicalJSON()
	if err != nil {
		t.Fatal(err)
	}
	decoded, err := openinghours.ParseJSON(encoded)
	if err != nil {
		t.Fatal(err)
	}
	index, err := prepared.New(schedule)
	if err != nil {
		t.Fatal(err)
	}
	union, err := schedule.Union(decoded)
	if err != nil {
		t.Fatal(err)
	}
	intersection, err := schedule.Intersection(decoded)
	if err != nil {
		t.Fatal(err)
	}
	for name, candidate := range map[string]openinghours.Schedule{
		"canonical": decoded, "compiled": index.Schedule(), "union": union, "intersection": intersection,
	} {
		t.Run(name, func(t *testing.T) {
			ranges, err := candidate.EffectiveInstantRanges(oracle.Transitions[0].Instant, oracle.Transitions[3].Instant)
			if err != nil || len(ranges) != 2 || !ranges[0].End.Equal(oracle.Transitions[1].Instant) || !ranges[1].Start.Equal(oracle.Transitions[2].Instant) {
				t.Fatalf("owned seam lost disjoint fold membership: %#v %v", ranges, err)
			}
			duration, err := candidate.OpenDuration(oracle.Transitions[0].Instant, oracle.Transitions[3].Instant)
			if err != nil || duration != 90*time.Minute {
				t.Fatalf("owned seam duration=%s %v", duration, err)
			}
		})
	}
	closed, err := schedule.Subtract(decoded)
	if err != nil {
		t.Fatal(err)
	}
	ranges, err := closed.EffectiveInstantRanges(oracle.Transitions[0].Instant, oracle.Transitions[3].Instant)
	if err != nil || len(ranges) != 0 {
		t.Fatalf("self-subtraction expansion=%#v %v", ranges, err)
	}
}

func TestR04OpenDurationIntegratesActualOpenComponents(t *testing.T) {
	if r04RunWithPinnedZone(t) {
		return
	}
	schedule, oracle := r04Fixture(t)
	points := oracle.Transitions
	for _, tc := range []struct {
		name       string
		start, end time.Time
		want       time.Duration
	}{
		{"entire-fold", points[0].Instant, points[3].Instant, 90 * time.Minute},
		{"closed-repeated-period", points[1].Instant, points[2].Instant, 0},
		{"clipped-components", points[0].Instant.Add(15 * time.Minute), points[2].Instant.Add(15 * time.Minute), 30 * time.Minute},
	} {
		t.Run(tc.name, func(t *testing.T) {
			duration, err := schedule.OpenDuration(tc.start, tc.end)
			if err != nil || duration != tc.want {
				t.Fatalf("actual-open duration=%s error=%v; want %s", duration, err, tc.want)
			}
		})
	}
}
