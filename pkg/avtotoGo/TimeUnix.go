package avtotoGo

import (
	"fmt"
	"strconv"
	"strings"
	"time"
)

// First create a type alias
//type TimeUnix time.Time

type TimeUnix struct {
	timeUnixValue time.Time
}

/*
	// MarshalJSON is used to convert the timestamp to JSON
	func (t TimeUnix) MarshalJSON() ([]byte, error) {
		return []byte(strconv.FormatInt(time.Time(t.timeUnixValue).Unix(), 10)), nil
	}
*/

// UnmarshalJSON is used to convert the timestamp from JSON
func (t *TimeUnix) UnmarshalJSON(s []byte) (err error) {
	r := strings.ReplaceAll(string(s), "\"", "")
	fmt.Println("r ->", r)
	q, err := strconv.ParseInt(r, 10, 64)
	if err != nil {
		return err
	}
	//*(*time.Time)(t) = time.Unix(q, 0)
	t.timeUnixValue = time.Unix(q, 0)
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

/*
	// String returns t as a formatted string
	func (t TimeUnix) String() string {
		return t.Time().String()
	}
*/

// Вернуть дату в формате string в одной строке
func (t TimeUnix) String() string {
	return t.Day() + "." + t.Month() + "." + t.Year() + " " + t.Hour() + ":" + t.Minute() // Время
}

// Вернуть дату в формате string в одной строке
func (t TimeUnix) Strings() (string, string, string) {
	return t.Day() + "." + t.Month(), // День и месяц
		t.Year(), // Год
		t.Hour() + ":" + t.Minute() // Время
}

// Составные части даты и времени
func (t TimeUnix) Day() string {
	return t.timeUnixValue.Format("02")
}
func (t TimeUnix) Month() string {
	return t.timeUnixValue.Format("01")
}
func (t TimeUnix) Year() string {
	return t.timeUnixValue.Format("2006")
}
func (t TimeUnix) Hour() string {
	return t.timeUnixValue.Format("15")
}
func (t TimeUnix) Minute() string {
	return t.timeUnixValue.Format("04")
}
