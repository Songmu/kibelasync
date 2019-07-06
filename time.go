package kibelasync

import "time"

type Time struct {
	time.Time
}

const (
	rfc3339Milli       = `2006-01-02T15:04:05.999Z07:00`
	rfc3339MilliQuoted = `"` + rfc3339Milli + `"`
)

func (t *Time) UnmarshalJSON(data []byte) (err error) {
	t.Time, err = time.Parse(rfc3339MilliQuoted, string(data))
	return
}

func (t *Time) MarshalJSON() ([]byte, error) {
	return []byte(t.Format(rfc3339MilliQuoted)), nil
}
