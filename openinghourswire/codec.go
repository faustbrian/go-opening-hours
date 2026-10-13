// Package openinghourswire is the compatibility path for
// [github.com/faustbrian/go-opening-hours/v5/adapters/wire].
// Calls delegate to the canonical adapter with the same ownership/resource
// limits; path-specific named types and shared sentinel identities are retained.
// No external source or background lifecycle is added here.
//
// Deprecated: use github.com/faustbrian/go-opening-hours/v5/adapters/wire. This
// package remains supported for the longer of 180 days after successor public
// availability and two subsequently published stable minor releases.
package openinghourswire

import (
	openinghours "github.com/faustbrian/go-opening-hours/v5"
	canonical "github.com/faustbrian/go-opening-hours/v5/adapters/wire"
	wire "github.com/faustbrian/go-wire/v3"
)

// Format is the stable registry name for canonical opening-hours JSON.
const Format = canonical.Format

// WireFormat is the typed wire registry identity.
const WireFormat wire.Format = canonical.WireFormat

// Codec is stateless and safe for concurrent use.
type Codec struct{}

// Encode returns canonical bytes.
func (Codec) Encode(schedule openinghours.Schedule) ([]byte, error) {
	return canonical.Codec{}.Encode(schedule)
}

// Decode strictly parses canonical bytes.
func (Codec) Decode(data []byte) (openinghours.Schedule, error) {
	return canonical.Codec{}.Decode(data)
}
