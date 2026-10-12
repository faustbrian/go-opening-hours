// Package postgres provides nullable JSONB persistence values. The root
// Schedule directly implements database/sql and the interfaces selected by
// pgx JSON/JSONB codecs, so no connection-global type registration is needed.
// Use New for a valid nullable value, or zero JSONB for SQL NULL. Callers
// own DB connections, transactions and receiver synchronization; this package
// opens no service and performs no server migration or type-map registration.
package postgres

import (
	"database/sql/driver"

	openinghours "github.com/faustbrian/go-opening-hours/v4"
)

// JSONB is a nullable database wrapper around an immutable Schedule. Zero means
// SQL NULL, which is distinct from New(Schedule{}), a valid closed schedule.
// Scan mutates caller-owned wrapper storage; concurrent mutation needs locking.
type JSONB struct {
	schedule openinghours.Schedule
	valid    bool
}

// New returns a valid nullable JSONB value.
func New(schedule openinghours.Schedule) JSONB {
	return JSONB{schedule: schedule, valid: true}
}

// Get returns the stored schedule and validity flag.
func (j JSONB) Get() (openinghours.Schedule, bool) {
	return j.schedule, j.valid
}

// Value implements driver.Valuer.
func (j JSONB) Value() (driver.Value, error) {
	if !j.valid {
		return nil, nil
	}

	return j.schedule.Value()
}

// Scan implements sql.Scanner/pgx JSONB scanning, copying/decoding before
// assignment. NULL resets validity; malformed input leaves the prior wrapper
// unchanged. Nil receivers return the root CodeInvalidState error.
func (j *JSONB) Scan(source any) error {
	if j == nil {
		var schedule *openinghours.Schedule
		return schedule.Scan(source)
	}
	if source == nil {
		j.schedule = openinghours.Schedule{}
		j.valid = false
		return nil
	}
	var schedule openinghours.Schedule
	if err := schedule.Scan(source); err != nil {
		return err
	}
	j.schedule = schedule
	j.valid = true

	return nil
}
