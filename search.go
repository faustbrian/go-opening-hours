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

// NextTransition finds the first actual point-availability change strictly after
// instant. Candidates include civil endpoints in both fold occurrences and
// timezone-offset discontinuities. The positive horizon is at most 366 elapsed
// days; exhaustion returns CodeSearchExhausted. Interval expansion and duration
// use the same civil point membership; see docs/timezones.md.
func (s Schedule) NextTransition(instant time.Time, horizon time.Duration) (Transition, error) {
	if horizon <= 0 || horizon > maxSearchHorizon {
		return Transition{}, newError("next transition", CodeInvalidHorizon)
	}
	if s.data == nil {
		return Transition{}, newError("next transition", CodeSearchExhausted)
	}
	deadline := instant.Add(horizon)

	date, err := s.instantDate(instant)
	if err != nil {
		return Transition{}, err
	}
	localizedDeadline := deadline.In(s.data.location)
	lastDate, _ := NewDate(localizedDeadline.Year(), localizedDeadline.Month(), localizedDeadline.Day())
	transitions, err := s.transitionCandidates(date, lastDate, instant, deadline)
	if err != nil {
		return Transition{}, err
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

// NextOpening finds the first actual opening strictly after instant within
// NextTransition's bounded horizon.
func (s Schedule) NextOpening(instant time.Time, horizon time.Duration) (Transition, error) {
	return s.nextTransitionOfKind(instant, horizon, TransitionOpen)
}

// NextClosing finds the first actual closing strictly after instant within
// NextTransition's bounded horizon.
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

// PreviousTransition finds the nearest actual point-availability change strictly
// before instant within a positive horizon of at most 366 elapsed days. Like
// NextTransition, it includes offset discontinuities and both fold occurrences.
// Closing explanations use the pre-boundary result; opening explanations use
// the post-boundary result.
func (s Schedule) PreviousTransition(instant time.Time, horizon time.Duration) (Transition, error) {
	if horizon <= 0 || horizon > maxSearchHorizon {
		return Transition{}, newError("previous transition", CodeInvalidHorizon)
	}
	if s.data == nil {
		return Transition{}, newError("previous transition", CodeSearchExhausted)
	}
	oldest := instant.Add(-horizon)
	date, err := s.instantDate(instant)
	if err != nil {
		return Transition{}, err
	}
	localizedOldest := oldest.In(s.data.location)
	firstDate, firstDateErr := NewDate(localizedOldest.Year(), localizedOldest.Month(), localizedOldest.Day())
	if firstDateErr != nil {
		firstDate = MustDate(1, time.January, 1)
	}
	candidates, err := s.transitionCandidates(firstDate, date, oldest.Add(-time.Nanosecond), instant)
	if err != nil {
		return Transition{}, err
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
