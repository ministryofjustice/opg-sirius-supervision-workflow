package model

import (
	"strings"
	"time"
)

type Date struct {
	Time time.Time
}

func NewDate(date string) Date {
	t, err := time.Parse("02/01/2006", date)
	if err != nil {
		panic(err)
	}
	return Date{Time: t}
}

func (d Date) Before(d2 Date) bool {
	return d.Time.Before(d2.Time)
}

func (d Date) String() string {
	nullDate := NewDate("01/01/0001")
	if nullDate.Before(d) {
		return d.Time.Format("02/01/2006")
	}
	return ""
}

func (d *Date) UnmarshalJSON(b []byte) error {
	value := strings.Trim(string(b), `"`)
	if value == "" || value == "null" {
		return nil
	}

	value = strings.ReplaceAll(value, `\`, "")
	t, err := time.Parse("02/01/2006", value)
	if err != nil {
		return err
	}
	*d = Date{Time: t}
	return nil
}

func (d Date) MarshalJSON() ([]byte, error) {
	return []byte(`"` + d.Time.Format("02\\/01\\/2006") + `"`), nil
}
