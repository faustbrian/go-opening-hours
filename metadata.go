package openinghours

import (
	"bytes"
	"encoding/json"
	"unicode/utf8"
)

const maxMetadataBytes = 256

// Metadata carries bounded provenance and has no interval semantics.
type Metadata struct {
	// Label is caller-supplied UTF-8 text of at most 256 bytes, not automatically redacted.
	Label string
	// Source is bounded UTF-8 provenance, not authenticated authority.
	Source string
	// Revision is an opaque bounded UTF-8 identity; it is not a content hash or freshness proof.
	Revision string
}

// OutsideEffectivePolicy controls queries outside inclusive effective dates.
type OutsideEffectivePolicy uint8

const (
	// OutsideClosed treats dates outside the effective range as closed.
	OutsideClosed OutsideEffectivePolicy = iota
	// OutsideError rejects queries outside the effective range.
	OutsideError
)

func validateMetadata(metadata Metadata) (Metadata, error) {
	for _, value := range []string{metadata.Label, metadata.Source, metadata.Revision} {
		if len(value) > maxMetadataBytes || !utf8.ValidString(value) {
			return Metadata{}, newError("validate metadata", CodeLimitExceeded)
		}
	}

	return metadata, nil
}

func (s Schedule) withinEffective(date Date) (bool, error) {
	if s.data == nil {
		return false, nil
	}
	outside := s.data.hasEffectiveStart && compareDate(date, s.data.effectiveStart) < 0 ||
		s.data.hasEffectiveEnd && compareDate(date, s.data.effectiveEnd) > 0
	if outside && s.data.outsideEffective == OutsideError {
		return false, newError("query", CodeOutsideEffectiveRange)
	}

	return !outside, nil
}

// Metadata returns a detached copy of schedule provenance.
func (s Schedule) Metadata() Metadata {
	if s.data == nil {
		return Metadata{}
	}

	return s.data.metadata
}

// Revision returns the opaque bounded revision identifier.
func (s Schedule) Revision() string { return s.Metadata().Revision }

// SemanticallyEqual compares canonical wire representations after clearing
// schedule metadata recursively. Exception provenance, timezone, effective
// policies and composition shape remain significant; this is not a proof of
// arbitrary algebraic/point-set equivalence.
func (s Schedule) SemanticallyEqual(other Schedule) bool {
	left := s.toWire()
	right := other.toWire()
	clearWireMetadata(&left)
	clearWireMetadata(&right)
	leftJSON, _ := json.Marshal(left) // wireSchedule contains no values rejected by encoding/json.
	rightJSON, _ := json.Marshal(right)

	return bytes.Equal(leftJSON, rightJSON)
}

func clearWireMetadata(wire *wireSchedule) {
	wire.Metadata = wireMetadata{}
	if wire.Composition != nil {
		clearWireMetadata(&wire.Composition.Left)
		clearWireMetadata(&wire.Composition.Right)
	}
}
