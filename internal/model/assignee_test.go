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
