// Package openinghours provides immutable recurring opening-hours schedules,
// dated exceptions, explicit timezone evaluation, bounded search, and stable
// canonical persistence. Start with NewSchedule and IsOpen for explicit instant
// queries. LocalTime is a civil wall value, not an elapsed duration; ranges own
// their start date and are start-inclusive/end-exclusive. A zero Schedule is
// closed without a timezone, while a zero DayRule is inherited.
//
// Constructed schedules and returned detached results support concurrent reads.
// Constructors copy input collections; callers must not mutate them during a
// call. Pointer decoding/scanning replaces caller storage and needs external
// synchronization. The core owns no goroutines, clocks, network IO or database
// lifecycle. Borrowed clock/observer capabilities execute synchronously.
//
// Date and zone resolution use Calendar v2; timezone behavior depends on the
// installed Go/IANA data, not a frozen rule snapshot. Canonical wire edition 1
// is package-owned and carries provenance, not authentication. Optional adapters
// introduce their declared client/module contracts and own no external service.
package openinghours
