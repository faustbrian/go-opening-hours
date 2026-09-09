// Package openinghoursconfig is the compatibility path for
// [github.com/faustbrian/go-opening-hours/adapters/config].
//
// Deprecated: use github.com/faustbrian/go-opening-hours/adapters/config.
// This package remains supported for the longer of 180 days after successor
// public availability and two subsequently published stable minor releases.
package openinghoursconfig

import (
	openinghours "github.com/faustbrian/go-opening-hours"
	canonical "github.com/faustbrian/go-opening-hours/adapters/config"
)

// ErrInvalidValue reports a configuration value that is not canonical JSON
// text. The rejected value is never included.
var ErrInvalidValue = canonical.ErrInvalidValue

// Parse strictly decodes one bounded canonical configuration value.
func Parse(value string) (openinghours.Schedule, error) {
	return canonical.Parse(value)
}

// Value implements config's value-unmarshal seam without owning sources.
type Value struct{ value canonical.Value }

// NewValue wraps an immutable schedule for configuration serialization.
func NewValue(schedule openinghours.Schedule) Value {
	return Value{value: canonical.NewValue(schedule)}
}

// Schedule returns the decoded immutable schedule.
func (value Value) Schedule() openinghours.Schedule { return value.value.Schedule() }

// UnmarshalConfigValue decodes canonical JSON text supplied by config.
func (value *Value) UnmarshalConfigValue(input any) error {
	if value == nil {
		return ErrInvalidValue
	}
	next := value.value
	if err := next.UnmarshalConfigValue(input); err != nil {
		return err
	}
	value.value = next
	return nil
}

// MarshalText returns canonical configuration text.
func (value Value) MarshalText() ([]byte, error) { return value.value.MarshalText() }
