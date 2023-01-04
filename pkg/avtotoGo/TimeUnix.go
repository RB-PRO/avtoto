package avtotoGo

import (
	"strconv"
	"time"
)

// First create a type alias
//type TimeUnix time.Time

type TimeUnix struct {
	timeUnixValue time.Time
}

// MarshalJSON is used to convert the timestamp to JSON
func (t TimeUnix) MarshalJSON() ([]byte, error) {
	return []byte(strconv.FormatInt(time.Time(t.timeUnixValue).Unix(), 10)), nil
}

// UnmarshalJSON is used to convert the timestamp from JSON
func (t *TimeUnix) UnmarshalJSON(s []byte) (err error) {
	r := string(s)
	q, err := strconv.ParseInt(r, 10, 64)
	if err != nil {
		return err
	}
	*(*time.Time)(t.timeUnixValue) = time.Unix(q, 0)
	return nil
}

// Unix returns t as a Unix time, the number of seconds elapsed
// since January 1, 1970 UTC. The result does not depend on the
// location associated with t.
func (t TimeUnix) Unix() int64 {
	return time.Time(t.timeUnixValue).Unix()
}

// Time returns the JSON time as a time.Time instance in UTC
func (t TimeUnix) Time() time.Time {
	return time.Time(t.timeUnixValue).UTC()
}

// String returns t as a formatted string
func (t TimeUnix) String() string {
	return t.Time().String()
}
