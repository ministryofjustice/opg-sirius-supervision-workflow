package server

import (
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
