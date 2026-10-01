package openinghours

import (
	"reflect"
	"time"

	clock "github.com/faustbrian/go-clock"
)

// Clock is the clock current-time capability. Core schedule queries never
// read a process-global clock.
type Clock = clock.Clock

// ElapsedClock is the clock monotonic elapsed-time capability. Observation
// helpers accept it separately so measuring a query never reads wall time.
type ElapsedClock = clock.ElapsedClock

// IsOpenNow evaluates an injected clock's current instant.
func (s Schedule) IsOpenNow(clock Clock) (Availability, error) {
	if isNilCapability(clock) {
		return Availability{}, newError("is open now", CodeInvalidClock)
	}
	instant, ok := clockInstant(clock)
	if !ok {
		return Availability{}, newError("is open now", CodeInvalidClock)
	}

	return s.IsOpen(instant)
}

func clockInstant(clock Clock) (instant time.Time, ok bool) {
	defer func() {
		if recover() != nil {
			instant, ok = time.Time{}, false
		}
	}()

	return clock.Now(), true
}

func isNilCapability(value any) bool {
	if value == nil {
		return true
	}
	reflected := reflect.ValueOf(value)
	kind := reflected.Kind()
	if kind == reflect.Chan || kind == reflect.Func || kind == reflect.Interface ||
		kind == reflect.Map || kind == reflect.Pointer || kind == reflect.Slice {
		return reflected.IsNil()
	}
	return false
}
