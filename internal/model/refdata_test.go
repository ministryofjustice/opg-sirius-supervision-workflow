package model

import (
	"github.com/stretchr/testify/assert"
	"strconv"
	"testing"
)

func TestRefData_Is(t *testing.T) {
	tests := []struct {
		refData RefData
		handle  string
		want    bool
	}{
		{
			refData: RefData{},
			handle:  "",
			want:    true,
		},
		{
			refData: RefData{Handle: "test"},
			handle:  "test",
			want:    true,
		},
		{
			refData: RefData{Handle: "test"},
			handle:  "no-test",
			want:    false,
		},
	}
	for i, test := range tests {
		t.Run("Scenario "+strconv.Itoa(i+1), func(t *testing.T) {
			assert.Equal(t, test.want, test.refData.Is(test.handle))
		})
	}
}

func TestRefData_IsIn(t *testing.T) {
	tests := []struct {
		refData RefData
		handles []string
		want    bool
	}{
		{
			refData: RefData{},
			handles: nil,
			want:    false,
		},
		{
			refData: RefData{Handle: "test"},
			handles: []string{"foo", "test"},
			want:    true,
		},
		{
			refData: RefData{Handle: "test"},
			handles: []string{"no-test", "foo"},
			want:    false,
		},
	}
	for i, test := range tests {
		t.Run("Scenario "+strconv.Itoa(i+1), func(t *testing.T) {
			assert.Equal(t, test.want, test.refData.IsIn(test.handles))
		})
	}
}
