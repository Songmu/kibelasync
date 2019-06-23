package kibela

import "time"

type Time struct {
	time.Time
}

const rfc3339Milli = "2006-01-02T15:04:05.999Z07:00"

func (t *Time) UnmarshalJSON(data []byte) (err error) {
	*t, err = time.Parse(rfc3339Milli, string(data))
	return
}

func (t *Time) MarshalJSON() ([]byte, error) {
	return t.Time.Format(rfc3339Milli), nil
}
