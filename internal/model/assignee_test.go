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

//func TestAssignee_IsCaseManager(t *testing.T) {
//	assignee1 := Assignee{Id: 10, Roles: []string{"System Admin": "System Admin", "COP User": "COP User"}}
//	assignee2 := Assignee{Id: 11, Roles: make([]string, 0)}
//
//	assert.Truef(t, assignee1.IsCaseManager(), "Assignee ID 10 is case manager")
//	assert.Truef(t, assignee2.IsCaseManager(), "Assignee ID 11 is not case manager")
//}

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
