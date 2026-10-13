package openinghours

import "time"

// nextOffsetBoundary finds an actual offset change within (start, end]. Native
// ZoneBounds also emits POSIX-extension year boundaries with unchanged offsets.
// Go 1.27's late-year extension can report a stale end on a leap year's last day;
// advancing to the next UTC year resumes native lookup without treating that
// bookkeeping boundary as availability or an extra collected fragment.
func (s Schedule) nextOffsetBoundary(start, end time.Time) (time.Time, bool) {
	_, offset := start.In(s.data.location).Zone()
	for cursor := start; cursor.Before(end); {
		_, boundary := cursor.In(s.data.location).ZoneBounds()
		if boundary.IsZero() {
			return time.Time{}, false
		}
		if !boundary.After(cursor) {
			// The POSIX lookup's annual bounds use UTC years, not local dates.
			boundary = time.Date(cursor.UTC().Year()+1, time.January, 1, 0, 0, 0, 0, time.UTC)
		}
		if boundary.After(end) {
			return time.Time{}, false
		}
		_, followingOffset := boundary.In(s.data.location).Zone()
		if followingOffset != offset {
			return boundary, true
		}
		cursor = boundary
	}
	return time.Time{}, false
}

// instantDate admits the existing civil-year domain before presenting a year as
// int. Unix seconds and native zone offsets retain their values on 32-bit hosts;
// time.Time.Year alone can truncate a very large caller-supplied instant there.
func (s Schedule) instantDate(instant time.Time) (Date, error) {
	localized := instant.In(s.data.location)
	_, offset := localized.Zone()
	seconds := instant.Unix()
	minimum := time.Date(1, time.January, 1, 0, 0, 0, 0, time.UTC).Unix() - int64(offset)
	maximum := time.Date(10000, time.January, 1, 0, 0, 0, 0, time.UTC).Unix() - int64(offset)
	if seconds < minimum || seconds >= maximum {
		return Date{}, newError("new date", CodeInvalidDate)
	}
	return NewDate(localized.Year(), localized.Month(), localized.Day())
}
