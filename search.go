package openinghours

import (
	"cmp"
	"slices"
	"time"
)

const maxSearchHorizon = 366 * 24 * time.Hour

// TransitionKind identifies an opening or closing boundary.
type TransitionKind uint8

const (
	// TransitionOpen identifies a boundary whose following instant is open.
	TransitionOpen TransitionKind = iota
	// TransitionClose identifies a boundary whose following instant is closed.
	TransitionClose
)

// Transition is an explained availability boundary.
type Transition struct {
	// Instant is the selected absolute boundary.
	Instant time.Time
	// Kind reports availability immediately after the boundary.
	Kind TransitionKind
	// Explanation is post-boundary provenance for forward searches; PreviousTransition uses
	// pre-boundary provenance for closing.
	Explanation Explanation
}

// NextTransition finds the first resolved civil-range endpoint after instant
// whose before/after IsOpen values differ. The horizon is positive and at most
// 366 elapsed days; exhaustion returns CodeSearchExhausted. It does not enumerate
// timezone-offset jumps or every fold occurrence, so it can omit wall-membership
// changes at those instants. See docs/timezones.md for the documented fold limit.
func (s Schedule) NextTransition(instant time.Time, horizon time.Duration) (Transition, error) {
	if horizon <= 0 || horizon > maxSearchHorizon {
		return Transition{}, newError("next transition", CodeInvalidHorizon)
	}
	if s.data == nil {
		return Transition{}, newError("next transition", CodeSearchExhausted)
	}
	deadline := instant.Add(horizon)

	localized := instant.In(s.data.location)
	date, _ := NewDate(localized.Year(), localized.Month(), localized.Day())
	localizedDeadline := deadline.In(s.data.location)
	lastDate, _ := NewDate(localizedDeadline.Year(), localizedDeadline.Month(), localizedDeadline.Day())
	dayCount := date.DaysUntil(lastDate) + 1
	transitions := make([]Transition, 0, 16)
	for step := range dayCount {
		day, _ := addDate(date, step)
		segments, _, err := s.effectiveSegments(day)
		if err != nil {
			return Transition{}, err
		}
		for _, item := range segments {
			opening, openingErr := s.resolveBoundary(day, item.start, true)
			if openingErr == nil {
				transitions = append(transitions, Transition{Instant: opening, Kind: TransitionOpen})
			}
			closing, closingErr := s.resolveBoundary(day, item.end, false)
			if closingErr == nil {
				transitions = append(transitions, Transition{Instant: closing, Kind: TransitionClose})
			}
		}
	}
	slices.SortFunc(transitions, func(left, right Transition) int {
		if comparison := left.Instant.Compare(right.Instant); comparison != 0 {
			return comparison
		}
		return cmp.Compare(left.Kind, right.Kind)
	})
	for _, transition := range transitions {
		if !transition.Instant.After(instant) || transition.Instant.After(deadline) {
			continue
		}
		// Candidate generation already evaluated these civil dates, so point
		// queries at their boundaries cannot introduce a new schedule error.
		before, _ := s.IsOpen(transition.Instant.Add(-time.Nanosecond))
		after, _ := s.IsOpen(transition.Instant)
		if before.Open == after.Open {
			continue
		}
		if after.Open {
			transition.Kind = TransitionOpen
		} else {
			transition.Kind = TransitionClose
		}
		transition.Explanation = after.Explanation
		return transition, nil
	}

	return Transition{}, newError("next transition", CodeSearchExhausted)
}

// NextOpening selects the next opening among NextTransition's endpoint
// candidates within the same bounded horizon; it inherits its fold-search limit.
func (s Schedule) NextOpening(instant time.Time, horizon time.Duration) (Transition, error) {
	return s.nextTransitionOfKind(instant, horizon, TransitionOpen)
}

// NextClosing selects the next closing among NextTransition's endpoint
// candidates within the same bounded horizon; it inherits its fold-search limit.
func (s Schedule) NextClosing(instant time.Time, horizon time.Duration) (Transition, error) {
	return s.nextTransitionOfKind(instant, horizon, TransitionClose)
}

func (s Schedule) nextTransitionOfKind(instant time.Time, horizon time.Duration, kind TransitionKind) (Transition, error) {
	if horizon < 1 || horizon > maxSearchHorizon {
		return Transition{}, newError("next typed transition", CodeInvalidHorizon)
	}
	deadline := instant.Add(horizon)
	cursor := instant
	for cursor.Before(deadline) {
		transition, err := s.NextTransition(cursor, deadline.Sub(cursor))
		if err != nil {
			return Transition{}, err
		}
		if transition.Kind == kind {
			return transition, nil
		}
		cursor = advanceTransitionCursor(cursor, transition, deadline)
	}

	return Transition{}, newError("next typed transition", CodeSearchExhausted)
}

func advanceTransitionCursor(cursor time.Time, transition Transition, deadline time.Time) time.Time {
	if !transition.Instant.After(cursor) {
		return deadline
	}

	return transition.Instant
}

// PreviousTransition finds the nearest resolved endpoint before instant whose
// before/after IsOpen values differ, within a positive horizon of at most 366
// elapsed days. Like NextTransition, it does not enumerate timezone-offset jumps
// or every fold occurrence and can omit changes there. Closing explanations use
// the pre-boundary result; opening explanations use the post-boundary result.
func (s Schedule) PreviousTransition(instant time.Time, horizon time.Duration) (Transition, error) {
	if horizon <= 0 || horizon > maxSearchHorizon {
		return Transition{}, newError("previous transition", CodeInvalidHorizon)
	}
	if s.data == nil {
		return Transition{}, newError("previous transition", CodeSearchExhausted)
	}
	oldest := instant.Add(-horizon)
	localized := instant.In(s.data.location)
	date, _ := NewDate(localized.Year(), localized.Month(), localized.Day())
	localizedOldest := oldest.In(s.data.location)
	firstDate, firstDateErr := NewDate(localizedOldest.Year(), localizedOldest.Month(), localizedOldest.Day())
	if firstDateErr != nil {
		firstDate = MustDate(1, time.January, 1)
	}
	dayCount := firstDate.DaysUntil(date) + 1
	candidates := make([]Transition, 0, 16)
	for step := range dayCount {
		day, _ := addDate(firstDate, step)
		segments, _, err := s.effectiveSegments(day)
		if err != nil {
			return Transition{}, err
		}
		for _, item := range segments {
			opening, openingErr := s.resolveBoundary(day, item.start, true)
			if openingErr == nil {
				candidates = append(candidates, Transition{Instant: opening, Kind: TransitionOpen})
			}
			closing, closingErr := s.resolveBoundary(day, item.end, false)
			if closingErr == nil {
				candidates = append(candidates, Transition{Instant: closing, Kind: TransitionClose})
			}
		}
	}
	slices.SortFunc(candidates, func(left, right Transition) int {
		return right.Instant.Compare(left.Instant)
	})
	for _, transition := range candidates {
		if !transition.Instant.Before(instant) || transition.Instant.Before(oldest) {
			continue
		}
		// Candidate generation already evaluated these civil dates, so point
		// queries at their boundaries cannot introduce a new schedule error.
		before, _ := s.IsOpen(transition.Instant.Add(-time.Nanosecond))
		after, _ := s.IsOpen(transition.Instant)
		switch {
		case before.Open == after.Open:
		case after.Open:
			transition.Kind = TransitionOpen
			transition.Explanation = after.Explanation
			return transition, nil
		default:
			transition.Kind = TransitionClose
			transition.Explanation = before.Explanation
			return transition, nil
		}
	}

	return Transition{}, newError("previous transition", CodeSearchExhausted)
}

// Interval expansion resolves folds to an earlier opening/later closing and
// shifts gaps forward, unlike caller-selected ResolveLocal policy. Point queries
// still evaluate wall membership independently; do not equate this boundary
// expansion with enumerating every possible timezone-offset transition.
func (s Schedule) resolveBoundary(date Date, nanosecond int64, opening bool) (time.Time, error) {
	if nanosecond == nanosecondsPerDay {
		var err error
		date, err = addDate(date, 1)
		if err != nil {
			return time.Time{}, err
		}
		nanosecond = 0
	}
	localTime := LocalTime{nanosecond: nanosecond}
	policy := PreferLater
	if opening {
		policy = PreferEarlier
	}
	resolved, err := s.ResolveLocal(date, localTime, policy)
	if IsCode(err, CodeNonexistentLocalTime) {
		resolved, err = s.ResolveLocal(date, localTime, ShiftForward)
	}
	if err != nil {
		return time.Time{}, err
	}

	return resolved.Instant, nil
}
