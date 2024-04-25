package model

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestAssignee_IsSelected(t *testing.T) {
	selectedAssignee := Assignee{Id: 10}
	unselectedAssignee := Assignee{Id: 11}

	selectedAssignees := []string{"9", "10", "12", "13"}

	assert.Truef(t, selectedAssignee.IsSelected(selectedAssignees), "Assignee ID 10 is not selected")
	assert.False(t, unselectedAssignee.IsSelected(selectedAssignees), "Assignee ID 11 is selected")
}

func TestGetCount(t *testing.T) {
	tests := []struct {
		testname         string
		selectedAssignee Assignee
		want             string
	}{
		{
			testname:         "Returns count for assignee",
			selectedAssignee: Assignee{Id: 10},
			want:             "(11)",
		},
		{
			testname:         "Returns null count for assignee",
			selectedAssignee: Assignee{Id: 11},
			want:             "(0)",
		},
		{
			testname:         "Returns null if assignee not in list",
			selectedAssignee: Assignee{Id: 22},
			want:             "(0)",
		},
	}
	for _, test := range tests {
		selectedAssignees := []AssigneeAndCount{
			{AssigneeId: 10, Count: 11},
			{AssigneeId: 11, Count: 0},
			{AssigneeId: 12, Count: 1},
		}

		t.Run(test.testname, func(t *testing.T) {
			assert.Equal(t, test.selectedAssignee.GetCount(selectedAssignees), test.want)
		})
	}
}

func TestGetCount2(t *testing.T) {
	selectedAssignee := Assignee{Id: 11}
	selectedAssignees := []AssigneeAndCount{
		{
			AssigneeId: 10,
			Count:      11,
		},
		{
			AssigneeId: 11,
			Count:      0,
		},
	}

	assert.Equal(t, selectedAssignee.GetCount(selectedAssignees), "(0)")
}
