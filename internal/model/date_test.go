package model

import (
	"encoding/json"
	"github.com/stretchr/testify/assert"
	"testing"
	"time"
)

type testJsonDateStruct struct {
	TestDate Date `json:"testDate"`
}

func TestNewDate(t *testing.T) {
	want := Date{Time: time.Date(2020, 12, 31, 0, 0, 0, 0, time.UTC)}
	assert.Equal(t, want, NewDate("31/12/2020"))

	defer func() {
		if r := recover(); r == nil {
			t.Errorf("NewDate should panic with incorrect date format")
		}
	}()
	NewDate("12/31/2020") // wrong format should trigger a panic
}

func TestDate_Before(t *testing.T) {
	tests := []struct {
		name  string
		date1 Date
		date2 Date
		want  bool
	}{
		{
			name:  "Date1 is before Date2",
			date1: NewDate("01/01/2020"),
			date2: NewDate("02/01/2020"),
			want:  true,
		},
		{
			name:  "Date1 is after Date2",
			date1: NewDate("02/01/2020"),
			date2: NewDate("01/01/2020"),
			want:  false,
		},
		{
			name:  "Date1 is the same as Date2",
			date1: NewDate("01/01/2020"),
			date2: NewDate("01/01/2020"),
			want:  false,
		},
		{
			name:  "Date1 is empty",
			date1: Date{},
			date2: NewDate("02/01/2020"),
			want:  true,
		},
		{
			name:  "Date2 is empty",
			date1: NewDate("01/01/2020"),
			date2: Date{},
			want:  false,
		},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			assert.Equal(t, test.want, test.date1.Before(test.date2))
		})
	}
}

func TestDate_MarshalJSON(t *testing.T) {
	v := testJsonDateStruct{TestDate: NewDate("01/01/2020")}
	b, err := json.Marshal(v)
	assert.Nil(t, err)
	assert.Equal(t, `{"testDate":"01\/01\/2020"}`, string(b))
}

func TestDate_String(t *testing.T) {
	assert.Equal(t, "01/01/2020", NewDate("01/01/2020").String())
}

func TestDate_UnmarshalJSON(t *testing.T) {
	var v *testJsonDateStruct
	err := json.Unmarshal([]byte(`{"testDate":"01\/01\/2020"}`), &v)
	assert.Nil(t, err)
	assert.Equal(t, "01/01/2020", v.TestDate.String())
}
