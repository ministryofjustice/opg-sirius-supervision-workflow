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

func (d Date) After(d2 Date) bool {
	return d.Time.After(d2.Time)
}

func (d Date) String() string {
	if d.IsNull() {
		return ""
	}
	return d.Time.Format("02/01/2006")
}

func (d Date) IsNull() bool {
	nullDate := NewDate("01/01/0001")
	return d.Time.Equal(nullDate.Time)
}

func (d *Date) UnmarshalJSON(b []byte) error {
	value := strings.Trim(string(b), `"`)
	if value == "" || value == "null" {
		return nil
	}

	value = strings.ReplaceAll(value, `\`, "")
	supportedFormats := []string{
		"02/01/2006",
		"2006-01-02T15:04:05+00:00",
	}

	var t time.Time
	var err error

	for _, format := range supportedFormats {
		t, err = time.Parse(format, value)
		if err != nil {
			continue
		}
		break
	}
	if err != nil {
		return err
	}

	*d = Date{Time: t}
	return nil
}

func (d Date) MarshalJSON() ([]byte, error) {
	return []byte(`"` + d.Time.Format("02\\/01\\/2006") + `"`), nil
}
