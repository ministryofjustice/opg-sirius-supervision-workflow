package model

import (
	"github.com/stretchr/testify/assert"
	"strconv"
	"testing"
)

func TestTaskType_IsSelected(t *testing.T) {
	tests := []struct {
		selectedTaskTypes []string
		taskType          TaskType
		wantIsSelected    bool
	}{
		{
			selectedTaskTypes: nil,
			taskType:          TaskType{},
			wantIsSelected:    false,
		},
		{
			selectedTaskTypes: []string{"T1"},
			taskType:          TaskType{Handle: "T1"},
			wantIsSelected:    true,
		},
		{
			selectedTaskTypes: []string{"T0", "T1"},
			taskType:          TaskType{Handle: "T1"},
			wantIsSelected:    true,
		},
		{
			selectedTaskTypes: []string{"T0", "T1"},
			taskType:          TaskType{Handle: "T2"},
			wantIsSelected:    false,
		},
	}
	for i, test := range tests {
		t.Run("Scenario "+strconv.Itoa(i), func(t *testing.T) {
			assert.Equal(t, test.wantIsSelected, test.taskType.IsSelected(test.selectedTaskTypes))
		})
	}
}
