package model

import (
	"strings"
	"time"
)

type Date struct {
	Time time.Time
}

func (d Date) Before(d2 Date) bool {
	return d.Time.Before(d2.Time)
}

func (d Date) String() string {
	return d.Time.Format("02/01/2006")
}

func (d *Date) UnmarshalJSON(b []byte) error {
	value := strings.Trim(string(b), `"`)
	if value == "" || value == "null" {
		return nil
	}

	t, err := time.Parse("02\\/01\\/2006", value)
	if err != nil {
		return err
	}
	*d = Date{Time: t}
	return nil
}

func (d Date) MarshalJSON() ([]byte, error) {
	return []byte(`"` + d.Time.Format("02\\/01\\/2006") + `"`), nil
}
