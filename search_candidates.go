package openinghours

import "time"

// Search follows point membership, not resolveBoundary's expansion envelope.
// Include every zone discontinuity and both realizations of a folded endpoint;
// callers retain only candidates whose before/after point states differ.
func (s Schedule) transitionCandidates(first, last Date, start, end time.Time) ([]Transition, error) {
	candidates := make([]Transition, 0, 16)
	for cursor := start; cursor.Before(end); {
		_, boundary := cursor.In(s.data.location).ZoneBounds()
		if boundary.IsZero() || !boundary.After(cursor) || boundary.After(end) {
			break
		}
		candidates = append(candidates, Transition{Instant: boundary})
		// A fold can revisit a civil date absent from both outer bounds.
		// Include dates reached on both sides of every offset discontinuity.
		for _, instant := range []time.Time{boundary.Add(-time.Nanosecond), boundary} {
			local := instant.In(s.data.location)
			day, _ := NewDate(local.Year(), local.Month(), local.Day())
			// DaysUntil returns zero for an unrepresentable Date, so dates
			// outside the civil domain do not extend the search's date bounds.
			if day.DaysUntil(first) > 0 {
				first = day
			}
			if last.DaysUntil(day) > 0 {
				last = day
			}
		}
		cursor = boundary
	}
	for step := range first.DaysUntil(last) + 1 {
		day, _ := addDate(first, step)
		segments, _, err := s.effectiveSegments(day)
		if err != nil {
			return nil, err
		}
		for _, item := range segments {
			candidates = append(candidates, s.searchBoundaryCandidates(day, item.start)...)
			candidates = append(candidates, s.searchBoundaryCandidates(day, item.end)...)
		}
	}
	return candidates, nil
}

func (s Schedule) searchBoundaryCandidates(date Date, nanosecond int64) []Transition {
	if nanosecond == nanosecondsPerDay {
		var err error
		date, err = addDate(date, 1)
		if err != nil {
			return nil
		}
		nanosecond = 0
	}
	local := LocalTime{nanosecond: nanosecond}
	earlier, err := s.ResolveLocal(date, local, PreferEarlier)
	if err != nil {
		return nil
	} // A gap has no endpoint occurrence; its offset jump is a candidate.
	result := []Transition{{Instant: earlier.Instant}}
	if earlier.Kind == LocalFold {
		// The same immutable zone/local input was just classified as a fold.
		later, _ := s.ResolveLocal(date, local, PreferLater)
		result = append(result, Transition{Instant: later.Instant})
	}
	return result
}
