package model

import (
	"github.com/stretchr/testify/assert"
	"testing"
	"time"
)

func TestTask_GetName(t *testing.T) {
	taskTypes := []TaskType{
		{
			Handle:     "Type1",
			Incomplete: "Label1",
		},
		{
			Handle:     "Type2",
			Incomplete: "Label2",
		},
	}

	tests := []struct {
		name string
		task Task
		want string
	}{
		{
			name: "Incomplete name used as task name",
			task: Task{Type: "Type2", Name: "TaskName"},
			want: "Label2",
		},
		{
			name: "Original task name used when cannot be matched to a task type",
			task: Task{Type: "Type3", Name: "TaskName"},
			want: "TaskName",
		},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			assert.Equal(t, test.want, test.task.GetName(taskTypes))
		})
	}
}

func TestTask_GetAssignee(t *testing.T) {
	tests := []struct {
		name string
		task Task
		want string
	}{
		{
			name: "Unassigned task gets Assignee from Clients",
			task: Task{
				Assignee: Assignee{Name: "Unassigned"},
				Clients:  []Client{{SupervisionCaseOwner: Assignee{Name: "Johnny"}}},
				Orders:   []Order{{Client: Client{SupervisionCaseOwner: Assignee{Name: "Pamela"}}}},
			},
			want: "Johnny",
		},
		{
			name: "Unassigned task gets Assignee from Orders",
			task: Task{
				Assignee: Assignee{Name: "Unassigned"},
				Orders:   []Order{{Client: Client{SupervisionCaseOwner: Assignee{Name: "Pamela"}}}},
			},
			want: "Pamela",
		},
		{
			name: "Assigned task get Assignee from task",
			task: Task{
				Assignee: Assignee{Name: "Bob"},
				Clients:  []Client{{SupervisionCaseOwner: Assignee{Name: "Johnny"}}},
				Orders:   []Order{{Client: Client{SupervisionCaseOwner: Assignee{Name: "Pamela"}}}},
			},
			want: "Bob",
		},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			assert.Equal(t, test.want, test.task.GetAssignee().Name)
		})
	}
}

func TestTask_GetClient(t *testing.T) {
	tests := []struct {
		name string
		task Task
		want string
	}{
		{
			name: "Get client from task case",
			task: Task{
				Clients: []Client{{FirstName: "Johnny"}},
				Orders:  []Order{{Client: Client{FirstName: "Pamela"}}},
			},
			want: "Pamela",
		},
		{
			name: "Get client from task without cases",
			task: Task{
				Clients: []Client{{FirstName: "Johnny"}},
			},
			want: "Johnny",
		},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			assert.Equal(t, test.want, test.task.GetClient().FirstName)
		})
	}
}

func TestTask_GetDueDateStatus(t *testing.T) {
	tests := []struct {
		name      string
		mockToday string
		dueDate   string
		want      DueDateStatus
	}{
		{
			name:      "Due date is in the past",
			mockToday: "06/06/2023",
			dueDate:   "05/06/2023",
			want:      DueDateStatus{Name: "Overdue", Colour: "red"},
		},
		{
			name:      "Due date is today",
			mockToday: "11/06/2023",
			dueDate:   "11/06/2023",
			want:      DueDateStatus{Name: "Due Today", Colour: "red"},
		},
		{
			name:      "Due date is tomorrow",
			mockToday: "06/06/2023",
			dueDate:   "07/06/2023",
			want:      DueDateStatus{Name: "Due Tomorrow", Colour: "orange"},
		},
		{
			name:      "Due date is this week but not tomorrow",
			mockToday: "06/06/2023",
			dueDate:   "08/06/2023",
			want:      DueDateStatus{Name: "Due This Week", Colour: "orange"},
		},
		{
			name:      "Due date is Monday next week",
			mockToday: "06/06/2023",
			dueDate:   "12/06/2023",
			want:      DueDateStatus{Name: "Due Next Week", Colour: "green"},
		},
		{
			name:      "Due date on same week day as today but next week",
			mockToday: "06/06/2023",
			dueDate:   "13/06/2023",
			want:      DueDateStatus{Name: "Due Next Week", Colour: "green"},
		},
		{
			name:      "Sunday today due date Monday",
			mockToday: "11/06/2023",
			dueDate:   "12/06/2023",
			want:      DueDateStatus{Name: "Due Next Week", Colour: "green"},
		},
		{
			name:      "Due date next week",
			mockToday: "06/06/2023",
			dueDate:   "12/06/2023",
			want:      DueDateStatus{Name: "Due Next Week", Colour: "green"},
		},
		{
			name:      "Due date on same week day as today but in future",
			mockToday: "06/06/2023",
			dueDate:   "23/06/2023",
			want:      DueDateStatus{},
		},
		{
			name:      "Due date that is not next week but after",
			mockToday: "06/06/2023",
			dueDate:   "19/06/2023",
			want:      DueDateStatus{},
		},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			task := Task{DueDate: test.dueDate}
			mockNow, _ := time.Parse("02/01/2006", test.mockToday)
			assert.Equal(t, test.want, task.GetDueDateStatus(mockNow))
		})
	}
}
