package server

import (
	"github.com/ministryofjustice/opg-sirius-workflow/internal/model"
	"github.com/stretchr/testify/assert"
	"strconv"
	"testing"
)

type mockFilterByAssignee struct {
	FilterByAssignee
}

type mockFilterByTaskType struct {
	FilterByTaskType
}

type mockFilterByDueDate struct {
	FilterByDueDate
}

type mockFilterByStatus struct {
	FilterByStatus
}

type mockFilterByDeputyType struct {
	FilterByDeputyType
}

func TestListPage_HasFilterBy(t *testing.T) {
	tests := []struct {
		page   interface{}
		filter string
		want   bool
	}{
		{
			page:   mockFilterByAssignee{},
			filter: "assignee",
			want:   true,
		},
		{
			page:   mockFilterByDueDate{},
			filter: "due-date",
			want:   true,
		},
		{
			page:   mockFilterByTaskType{},
			filter: "task-type",
			want:   true,
		},
		{
			page:   mockFilterByStatus{},
			filter: "status",
			want:   true,
		},
		{
			page:   mockFilterByDeputyType{},
			filter: "deputy-type",
			want:   true,
		},
		{
			page:   mockFilterByTaskType{},
			filter: "assignee",
			want:   false,
		},
		{
			page:   mockFilterByAssignee{},
			filter: "non-existent",
			want:   false,
		},
		{
			page:   struct{}{},
			filter: "non-existent",
			want:   false,
		},
	}
	for i, test := range tests {
		t.Run("Scenario "+strconv.Itoa(i+1), func(t *testing.T) {
			assert.Equal(t, test.want, ListPage{}.HasFilterBy(test.page, test.filter))
		})
	}
}

func TestFilterByTaskType_ValidateSelectedTaskTypes(t *testing.T) {
	tests := []struct {
		taskTypes         []model.TaskType
		selectedTaskTypes []string
		want              []string
	}{
		{
			taskTypes:         nil,
			selectedTaskTypes: []string{"test"},
			want:              nil,
		},
		{
			taskTypes:         []model.TaskType{{Handle: "test2"}},
			selectedTaskTypes: []string{"test1", "test2"},
			want:              []string{"test2"},
		},
		{
			taskTypes:         []model.TaskType{{Handle: "test3"}},
			selectedTaskTypes: []string{"test1", "test2"},
			want:              nil,
		},
	}
	for i, test := range tests {
		t.Run("Scenario "+strconv.Itoa(i), func(t *testing.T) {
			validatedTaskTypes := FilterByTaskType{}.ValidateSelectedTaskTypes(test.selectedTaskTypes, test.taskTypes)
			assert.Equal(t, test.want, validatedTaskTypes)
		})
	}
}
