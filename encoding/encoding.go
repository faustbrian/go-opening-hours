// Package encoding exposes the canonical schedule wire contract without
// shadowing human-readable formatting concerns in the root package.
// Use Marshal/Unmarshal for wire-edition 1 documents, or ImportLocation and
// ImportSpatie for bounded weekday structures. The parser normalizes accepted
// documents; this package owns no provider transport, authorization or file IO.
package encoding

import openinghours "github.com/faustbrian/go-opening-hours/v5"

// Marshal returns canonical versioned JSON.
func Marshal(schedule openinghours.Schedule) ([]byte, error) {
	return schedule.CanonicalJSON()
}

// Unmarshal strictly parses canonical JSON into an immutable schedule.
func Unmarshal(data []byte) (openinghours.Schedule, error) {
	return openinghours.ParseJSON(data)
}
