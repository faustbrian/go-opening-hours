// Package openinghoursconfig adapts canonical schedule documents from
// configuration sources while leaving environment and file ownership to the
// caller.
// Use Parse for a bounded document or NewValue for the config value seam.
// Returned schedules are immutable, but decoding mutates Value and needs
// external synchronization; this package performs no environment or file IO.
package openinghoursconfig

import (
	"errors"

	openinghours "github.com/faustbrian/go-opening-hours/v6"
)

// ErrInvalidValue reports a configuration value that is not canonical JSON
// text. The rejected value is never included.
var ErrInvalidValue = errors.New("openinghoursconfig: invalid value")

// Parse strictly decodes one bounded canonical configuration value.
func Parse(value string) (openinghours.Schedule, error) {
	if len(value) > openinghours.MaxJSONBytes {
		return openinghours.Schedule{}, &openinghours.Error{
			Code: openinghours.CodeInvalidEncoding, Op: "parse json",
		}
	}
	return openinghours.ParseJSON([]byte(value))
}

// Value implements config's value-unmarshal seam without owning sources.
type Value struct{ schedule openinghours.Schedule }

// NewValue wraps an immutable schedule for configuration serialization.
func NewValue(schedule openinghours.Schedule) Value { return Value{schedule: schedule} }

// Schedule returns the decoded immutable schedule.
func (value Value) Schedule() openinghours.Schedule { return value.schedule }

// UnmarshalConfigValue accepts only a string and assigns after successful Parse.
// Nil receiver/non-string inputs return ErrInvalidValue; parse errors propagate
// without replacing the old schedule. Callers synchronize mutable Value access.
func (value *Value) UnmarshalConfigValue(input any) error {
	if value == nil {
		return ErrInvalidValue
	}
	text, ok := input.(string)
	if !ok {
		return ErrInvalidValue
	}
	schedule, err := Parse(text)
	if err != nil {
		return err
	}
	value.schedule = schedule
	return nil
}

// MarshalText returns canonical configuration text.
func (value Value) MarshalText() ([]byte, error) { return value.schedule.CanonicalJSON() }
