package urlbuilder

import (
	"github.com/stretchr/testify/assert"
	"strconv"
	"testing"
)

func TestCreateFilter(t *testing.T) {
	tests := []struct {
		name                  string
		selectedValues        interface{}
		clearBetweenTeamViews []bool
		want                  Filter
	}{
		{
			name:                  "",
			selectedValues:        nil,
			clearBetweenTeamViews: nil,
			want:                  Filter{},
		},
		{
			name:                  "testFilter",
			selectedValues:        "testVal",
			clearBetweenTeamViews: []bool{true},
			want:                  Filter{Name: "testFilter", SelectedValues: []string{"testVal"}, ClearBetweenTeamViews: true},
		},
		{
			name:                  "",
			selectedValues:        []string{"testVal"},
			clearBetweenTeamViews: []bool{false},
			want:                  Filter{SelectedValues: []string{"testVal"}},
		},
		{
			name:                  "",
			selectedValues:        []string{"testVal1", "testVal2"},
			clearBetweenTeamViews: []bool{false},
			want:                  Filter{SelectedValues: []string{"testVal1", "testVal2"}},
		},
	}
	for i, test := range tests {
		t.Run("Scenario "+strconv.Itoa(i+1), func(t *testing.T) {
			var filter Filter
			if test.clearBetweenTeamViews == nil {
				filter = CreateFilter(test.name, test.selectedValues)
			} else {
				filter = CreateFilter(test.name, test.selectedValues, test.clearBetweenTeamViews[0])
			}
			assert.Equal(t, test.want, filter)
		})
	}
}
