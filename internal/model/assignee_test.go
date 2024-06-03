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

func TestAssignee_IsOnlyCaseManager(t *testing.T) {
	assignee1 := Assignee{Id: 10, Roles: []string{"Opg User", "Case Manager"}}
	assignee2 := Assignee{Id: 11, Roles: []string{"Opg User", "Case Manager", "System Admin"}}
	assignee3 := Assignee{Id: 12, Roles: []string{"Opg User", "System Admin"}}
	assignee4 := Assignee{Id: 13, Roles: []string{""}}

	assert.Truef(t, assignee1.IsOnlyCaseManager(), "Assignee ID 10 is case manager")
	assert.False(t, assignee2.IsOnlyCaseManager(), "Assignee ID 11 has many roles")
	assert.False(t, assignee3.IsOnlyCaseManager(), "Assignee ID 12 is not case manager")
	assert.False(t, assignee4.IsOnlyCaseManager(), "Assignee ID 13 is not case manager")
}

func TestGetCount(t *testing.T) {
	tests := []struct {
		testname         string
		selectedAssignee Assignee
		want             string
		url              string
	}{
		{
			testname:         "Returns count for assignee",
			selectedAssignee: Assignee{Id: 10},
			want:             "(11)",
			url:              "test",
		},
		{
			testname:         "Returns null count for assignee",
			selectedAssignee: Assignee{Id: 11},
			want:             "(0)",
			url:              "test",
		},
		{
			testname:         "Returns null if assignee not in list",
			selectedAssignee: Assignee{Id: 22},
			want:             "(0)",
			url:              "test",
		},
	}
	for _, test := range tests {
		selectedAssignees := []AssigneeAndCount{
			{AssigneeId: 10, Count: 11},
			{AssigneeId: 11, Count: 0},
			{AssigneeId: 12, Count: 1},
		}

		t.Run(test.testname, func(t *testing.T) {
			assert.Equal(t, test.selectedAssignee.GetCountAsString(selectedAssignees, test.url), test.want)
		})
	}
}
