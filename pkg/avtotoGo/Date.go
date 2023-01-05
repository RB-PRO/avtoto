package avtotoGo

import (
	"strings"
	"time"
)

type Date struct {
	timeDateValue time.Time
}

// UnmarshalJSON is used to convert the timestamp from JSON
func (t *Date) UnmarshalJSON(s []byte) error {
	var errorTime error
	stringTime := strings.ReplaceAll(string(s), "\"", "")
	if len(stringTime) > 10 { // 04.01.2023 16:55
		t.timeDateValue, errorTime = time.Parse("02.01.2006 15:04", stringTime)
	} else if len(stringTime) == 6 { // 05\/01
		t.timeDateValue, errorTime = time.Parse("02\\/01", stringTime)
	} else {
		t.timeDateValue, errorTime = time.Parse("02.01.2006", stringTime)
	}
	return errorTime
}

// Вернуть дату в формате string в одной строке
func (t Date) String() string {
	return t.Day() + "." + t.Month() + "." + t.Year() + " " + t.Hour() + ":" + t.Minute() // Время
}

// Вернуть дату в формате string по частям
func (t Date) Strings() (string, string, string) {
	return t.Day() + "." + t.Month(), // День и месяц
		t.Year(), // Год
		t.Hour() + ":" + t.Minute() // Время
}

// Составные части даты и времени
func (t Date) Day() string {
	return t.timeDateValue.Format("02")
}
func (t Date) Month() string {
	return t.timeDateValue.Format("01")
}
func (t Date) Year() string {
	return t.timeDateValue.Format("2006")
}
func (t Date) Hour() string {
	return t.timeDateValue.Format("15")
}
func (t Date) Minute() string {
	return t.timeDateValue.Format("04")
}
