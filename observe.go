package openinghours

import "time"

// Operation identifies an observable package operation.
type Operation uint8

const (
	// OperationIsOpen identifies an observed instant availability query.
	OperationIsOpen Operation = iota
	// OperationNextTransition identifies an observed transition search.
	OperationNextTransition
)

// Outcome is a bounded, non-sensitive operation result.
type Outcome uint8

const (
	// OutcomeClosed reports a successful query whose resource is closed.
	OutcomeClosed Outcome = iota
	// OutcomeOpen reports a successful query whose resource is open.
	OutcomeOpen
	// OutcomeFound reports a transition found within the search horizon.
	OutcomeFound
	// OutcomeError reports a typed query or search failure.
	OutcomeError
)

// Observation contains only bounded operational data. It intentionally has no
// schedule label, source, revision, date, timezone, or customer fields.
type Observation struct {
	// Operation identifies the observed query or search.
	Operation Operation
	// Outcome records success classification or an error without embedding diagnostic text.
	Outcome Outcome
	// RangeCount is the number of separately evaluated effective fragments for ObserveIsOpen, or
	// zero if unavailable.
	RangeCount int
	// SearchSteps is a horizon-based estimate capped at 367, not a measured count of candidate
	// evaluations.
	SearchSteps int
	// Duration is caller-supplied elapsed measurement; absent/panicking capabilities report zero.
	Duration time.Duration
}

// Observer receives one completed observation outside any lock. Panics are
// contained and cannot alter the query result. Nil is ignored. Calls are
// synchronous: blocking code has no cancellation bound and remains caller-owned;
// concurrent observations may invoke the observer concurrently.
type Observer func(Observation)

// ObserveIsOpen runs IsOpen and reports bounded operational data. A nil clock
// disables elapsed measurement and reports a zero duration. Typed nil, panic or
// a nil measurement also means zero. Measurement covers IsOpen itself, excluding
// the separate RangeCount query and observer callback; collaborator work may block.
func (s Schedule) ObserveIsOpen(instant time.Time, elapsedClock ElapsedClock, observer Observer) (Availability, error) {
	elapsed := measureElapsed(elapsedClock)
	result, err := s.IsOpen(instant)
	observation := Observation{Operation: OperationIsOpen, Duration: elapsed()}
	if err != nil {
		observation.Outcome = OutcomeError
	} else if result.Open {
		observation.Outcome = OutcomeOpen
	} else {
		observation.Outcome = OutcomeClosed
	}
	if s.data != nil {
		localized := instant.In(s.data.location)
		date, dateErr := NewDate(localized.Year(), localized.Month(), localized.Day())
		if dateErr == nil {
			ranges, rangeErr := s.EffectiveRanges(date)
			if rangeErr == nil {
				observation.RangeCount = len(ranges)
			}
		}
	}
	callObserver(observer, observation)

	return result, err
}

// ObserveNextTransition runs a bounded transition search and reports its
// outcome. A nil/panicking elapsed capability reports zero duration. SearchSteps
// estimates horizon-day cardinality rather than actual candidate visits; observer
// callback work is outside the measured query and remains synchronous.
func (s Schedule) ObserveNextTransition(instant time.Time, horizon time.Duration, elapsedClock ElapsedClock, observer Observer) (Transition, error) {
	elapsed := measureElapsed(elapsedClock)
	result, err := s.NextTransition(instant, horizon)
	outcome := OutcomeFound
	if err != nil {
		outcome = OutcomeError
	}
	observation := Observation{
		Operation: OperationNextTransition, Outcome: outcome,
		SearchSteps: boundedSearchSteps(horizon), Duration: elapsed(),
	}
	callObserver(observer, observation)

	return result, err
}

func measureElapsed(clock ElapsedClock) (elapsed func() time.Duration) {
	zero := func() time.Duration { return 0 }
	if isNilCapability(clock) {
		return zero
	}
	elapsed = zero
	defer func() {
		if recover() != nil {
			elapsed = zero
		}
	}()
	measurement := clock.Measure()
	if measurement == nil {
		return zero
	}

	return func() (duration time.Duration) {
		defer func() {
			if recover() != nil {
				duration = 0
			}
		}()
		return measurement()
	}
}

func boundedSearchSteps(horizon time.Duration) int {
	if horizon <= 0 {
		return 0
	}
	steps := int(horizon/(24*time.Hour)) + 1

	return min(steps, 367)
}

func callObserver(observer Observer, observation Observation) {
	if observer == nil {
		return
	}
	defer func() {
		_ = recover()
	}()
	observer(observation)
}
