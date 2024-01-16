package util

import (
	"github.com/ministryofjustice/opg-sirius-workflow/internal/model"
	"testing"
)

func TestIsLast(t *testing.T) {
	tests := []struct {
		name string
		i    int
		a    interface{}
		want bool
	}{
		{
			"Empty int array",
			0,
			[]int{},
			false,
		},
		{
			"Empty sirius typed array",
			0,
			[]model.Assurance{},
			false,
		},
		{
			"First of many",
			0,
			[]int{1, 2, 3},
			false,
		},
		{
			"Out of bounds",
			3,
			[]int{1, 2, 3},
			false,
		},
		{
			"Last index",
			2,
			[]int{1, 2, 3},
			true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := IsLast(tt.i, tt.a); got != tt.want {
				t.Errorf("IsLast() = %v, want %v", got, tt.want)
			}
		})
	}
}
