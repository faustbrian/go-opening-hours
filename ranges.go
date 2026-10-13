package openinghours

import (
	"slices"
	"time"
)

const maxOutputRanges = 8192

// DailyRange is a start-inclusive, end-exclusive interval within one civil
// date. EndAtDayBoundary distinguishes midnight at the end from midnight at
// the start of the date.
type DailyRange struct {
	// Start is the inclusive local/absolute endpoint.
	Start LocalTime
	// End is the exclusive endpoint; DailyRange midnight requires EndAtDayBoundary to
	// distinguish the end of the date.
	End LocalTime
	// EndAtDayBoundary marks the date-end midnight represented by zero End, not a zero-length
	// interval.
	EndAtDayBoundary bool
}

// InstantRange is a start-inclusive, end-exclusive absolute interval.
type InstantRange struct {
	// Start is the inclusive local/absolute endpoint.
	Start time.Time
	// End is the exclusive absolute endpoint.
	End time.Time
}

// EffectiveRanges returns normalized availability fragments for one civil date.
func (s Schedule) EffectiveRanges(date Date) ([]DailyRange, error) {
	if !validDate(date) {
		return nil, newError("effective ranges", CodeInvalidDate)
	}
	if s.data == nil {
		return []DailyRange{}, nil
	}
	segments, _, err := s.effectiveSegments(date)
	if err != nil {
		return nil, err
	}
	if len(segments) > maxOutputRanges {
		return nil, newError("effective ranges", CodeLimitExceeded)
	}
	result := make([]DailyRange, 0, len(segments))
	for _, item := range segments {
		endAtBoundary := item.end == nanosecondsPerDay
		end := item.end
		if endAtBoundary {
			end = 0
		}
		result = append(result, DailyRange{
			Start:            LocalTime{nanosecond: item.start},
			End:              LocalTime{nanosecond: end},
			EndAtDayBoundary: endAtBoundary,
		})
	}

	return result, nil
}

// EffectiveInstantRanges returns clipped absolute availability within a bounded
// interval of at most 366 elapsed days. Ranges represent exactly IsOpen's civil
// point membership: folds may produce disjoint intervals and gaps contribute
// only existing wall times. At most 8192 collected fragments are admitted before
// normalization; errors discard partial results. Returned ranges are detached
// and sorted/merged.
func (s Schedule) EffectiveInstantRanges(start, end time.Time) ([]InstantRange, error) {
	if !end.After(start) || end.Sub(start) > maxSearchHorizon {
		return nil, newError("effective instant ranges", CodeInvalidInterval)
	}
	if s.data == nil {
		return []InstantRange{}, nil
	}
	result := make([]InstantRange, 0, 16)
	for cursor := start; cursor.Before(end); {
		_, zoneEnd := cursor.In(s.data.location).ZoneBounds()
		until := end
		if !zoneEnd.IsZero() && zoneEnd.After(cursor) && zoneEnd.Before(until) {
			until = zoneEnd
		}
		var err error
		result, err = s.appendConstantOffsetRanges(result, cursor, until)
		if err != nil {
			return nil, err
		}
		cursor = until
	}
	return normalizeInstantRanges(result), nil
}

// Within one native ZoneBounds span the wall-to-instant mapping is affine.
// Intersect civil fragments with that span instead of resolving ambiguous or
// nonexistent endpoints. A fold is visited twice; a gap's missing wall times
// are never visited. Native offset boundaries also retain revisited civil dates.
func (s Schedule) appendConstantOffsetRanges(result []InstantRange, start, end time.Time) ([]InstantRange, error) {
	localized := start.In(s.data.location)
	first, _ := NewDate(localized.Year(), localized.Month(), localized.Day())
	_, offset := localized.Zone()
	localizedEnd := end.Add(-time.Nanosecond).In(s.data.location)
	last, _ := NewDate(localizedEnd.Year(), localizedEnd.Month(), localizedEnd.Day())
	dayCount := first.DaysUntil(last) + 1
	for step := range dayCount {
		date, _ := addDate(first, step)
		segments, _, err := s.effectiveSegments(date)
		if err != nil {
			return nil, err
		}
		for _, item := range segments {
			// Preserve the civil-date domain even when a caller clips an
			// all-day fragment before its out-of-domain closing boundary.
			if item.end == nanosecondsPerDay {
				if _, err := addDate(date, 1); err != nil {
					return nil, err
				}
			}
			midnight := time.Date(date.Year(), date.Month(), date.Day(), 0, 0, 0, 0, time.UTC)
			midnight = midnight.Add(-time.Duration(offset) * time.Second)
			opening := midnight.Add(time.Duration(item.start))
			closing := midnight.Add(time.Duration(item.end))
			opening = latestInstant(opening, start)
			closing = earliestInstant(closing, end)
			if opening.Before(closing) {
				result = append(result, InstantRange{Start: opening, End: closing})
				if len(result) > maxOutputRanges {
					return nil, newError("effective instant ranges", CodeLimitExceeded)
				}
			}
		}
	}

	return result, nil
}

func normalizeInstantRanges(input []InstantRange) []InstantRange {
	if len(input) == 0 {
		return []InstantRange{}
	}
	result := make([]InstantRange, len(input))
	copy(result, input)
	slices.SortFunc(result, func(left, right InstantRange) int {
		switch comparison := left.Start.Compare(right.Start); comparison {
		case -1, 1:
			return comparison
		default:
			return left.End.Compare(right.End)
		}
	})
	output := result[:1]
	for _, item := range result[1:] {
		last := &output[len(output)-1]
		if !item.Start.After(last.End) {
			if item.End.After(last.End) {
				last.End = item.End
			}
			continue
		}
		output = append(output, item)
	}

	return output
}

func latestInstant(left, right time.Time) time.Time {
	if left.After(right) {
		return left
	}

	return right
}

func earliestInstant(left, right time.Time) time.Time {
	if left.Before(right) {
		return left
	}

	return right
}

// OpenDuration sums elapsed duration of EffectiveInstantRanges' actual open
// intervals, not nominal wall-clock hours. It inherits that method's bounds and
// agrees with IsOpen's civil point membership across gaps and folds.
func (s Schedule) OpenDuration(start, end time.Time) (time.Duration, error) {
	ranges, err := s.EffectiveInstantRanges(start, end)
	if err != nil {
		return 0, err
	}
	var total time.Duration
	for _, item := range ranges {
		total += item.End.Sub(item.Start)
	}

	return total, nil
}
