package kibelasync

import "time"

// Time (un)marshals time for kibela
type Time struct {
	time.Time
}

const (
	rfc3339Milli       = `2006-01-02T15:04:05.999Z07:00`
	rfc3339MilliQuoted = `"` + rfc3339Milli + `"`
)

// UnmarshalJSON for encoding/json
func (t *Time) UnmarshalJSON(data []byte) (err error) {
	t.Time, err = time.Parse(rfc3339MilliQuoted, string(data))
	return
}

// MarshalJSON for encoding/json
func (t *Time) MarshalJSON() ([]byte, error) {
	return []byte(t.Format(rfc3339MilliQuoted)), nil
}
