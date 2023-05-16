package server

import (
	"github.com/ministryofjustice/opg-sirius-workflow/internal/sirius"
	"github.com/stretchr/testify/assert"
	"testing"
)

type urlFields struct {
	SelectedTeam        string
	CurrentPage         int
	PerPageLimit        int
	SelectedAssignees   []string
	SelectedUnassigned  string
	SelectedTaskTypes   []string
	SelectedDueDateFrom string
	SelectedDueDateTo   string
}

func createWorkflowVars(fields urlFields) WorkflowVars {
	return WorkflowVars{
		SelectedTeam: sirius.ReturnedTeamCollection{Selector: fields.SelectedTeam},
		PageDetails: sirius.PageDetails{
			CurrentPage:     fields.CurrentPage,
			StoredTaskLimit: fields.PerPageLimit,
		},
		SelectedAssignees:   fields.SelectedAssignees,
		SelectedUnassigned:  fields.SelectedUnassigned,
		SelectedTaskTypes:   fields.SelectedTaskTypes,
		SelectedDueDateFrom: fields.SelectedDueDateFrom,
		SelectedDueDateTo:   fields.SelectedDueDateTo,
	}
}

func TestWorkflowVars_GetClearFiltersUrl(t *testing.T) {
	tests := []struct {
		name   string
		fields urlFields
		want   string
	}{
		{
			name:   "Per page limit is retained",
			fields: urlFields{SelectedTeam: "lay", CurrentPage: 1, PerPageLimit: 50},
			want:   "?team=lay&page=1&per-page=50",
		},
		{
			name:   "Assignees are removed",
			fields: urlFields{SelectedTeam: "lay", CurrentPage: 1, PerPageLimit: 25, SelectedAssignees: []string{"1", "2"}},
			want:   "?team=lay&page=1&per-page=25",
		},
		{
			name:   "Unassigned is removed",
			fields: urlFields{SelectedTeam: "lay", CurrentPage: 1, PerPageLimit: 25, SelectedUnassigned: "1"},
			want:   "?team=lay&page=1&per-page=25",
		},
		{
			name:   "Task types are removed",
			fields: urlFields{SelectedTeam: "lay", CurrentPage: 1, PerPageLimit: 25, SelectedTaskTypes: []string{"1", "2"}},
			want:   "?team=lay&page=1&per-page=25",
		},
		{
			name:   "Due date filters are removed",
			fields: urlFields{SelectedTeam: "lay", CurrentPage: 1, PerPageLimit: 25, SelectedDueDateFrom: "2022-12-17", SelectedDueDateTo: "2022-12-18"},
			want:   "?team=lay&page=1&per-page=25",
		},
		{
			name:   "Page is reset back to 1",
			fields: urlFields{SelectedTeam: "lay", CurrentPage: 2, PerPageLimit: 25, SelectedAssignees: []string{"1", "2"}, SelectedUnassigned: "1", SelectedTaskTypes: []string{"task"}},
			want:   "?team=lay&page=1&per-page=25",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			w := createWorkflowVars(tt.fields)
			assert.Equalf(t, tt.want, w.GetClearFiltersUrl(), "GetClearFiltersUrl()")
		})
	}
}

func TestWorkflowVars_GetPaginationUrl(t *testing.T) {
	type args struct {
		page    int
		perPage int
	}
	tests := []struct {
		name   string
		fields urlFields
		args   args
		want   string
	}{
		{
			name:   "Page number is updated",
			fields: urlFields{SelectedTeam: "lay", CurrentPage: 1, PerPageLimit: 25},
			args:   args{page: 2, perPage: 25},
			want:   "?team=lay&page=2&per-page=25",
		},
		{
			name:   "Per page limit is updated",
			fields: urlFields{SelectedTeam: "lay", CurrentPage: 1, PerPageLimit: 25},
			args:   args{page: 1, perPage: 50},
			want:   "?team=lay&page=1&per-page=50",
		},
		{
			name:   "Per page limit is retained",
			fields: urlFields{SelectedTeam: "lay", CurrentPage: 1, PerPageLimit: 100},
			args:   args{page: 2},
			want:   "?team=lay&page=2&per-page=100",
		},
		{
			name:   "Assignees are retained",
			fields: urlFields{SelectedTeam: "lay", CurrentPage: 1, PerPageLimit: 25, SelectedAssignees: []string{"1", "2"}},
			args:   args{page: 2, perPage: 25},
			want:   "?team=lay&page=2&per-page=25&assignee=1&assignee=2",
		},
		{
			name:   "Unassigned is retained",
			fields: urlFields{SelectedTeam: "lay", CurrentPage: 1, PerPageLimit: 25, SelectedUnassigned: "1"},
			args:   args{page: 2, perPage: 25},
			want:   "?team=lay&page=2&per-page=25&unassigned=1",
		},
		{
			name:   "Task types are retained",
			fields: urlFields{SelectedTeam: "lay", CurrentPage: 1, PerPageLimit: 25, SelectedTaskTypes: []string{"1", "2"}},
			args:   args{page: 2, perPage: 25},
			want:   "?team=lay&page=2&per-page=25&task-type=1&task-type=2",
		},
		{
			name:   "Due date filters are retained",
			fields: urlFields{SelectedTeam: "lay", CurrentPage: 1, PerPageLimit: 25, SelectedDueDateFrom: "2022-12-17", SelectedDueDateTo: "2022-12-18"},
			args:   args{page: 2, perPage: 25},
			want:   "?team=lay&page=2&per-page=25&due-date-from=2022-12-17&due-date-to=2022-12-18",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			w := createWorkflowVars(tt.fields)
			var result string
			if tt.args.perPage == 0 {
				result = w.GetPaginationUrl(tt.args.page)
			} else {
				result = w.GetPaginationUrl(tt.args.page, tt.args.perPage)
			}
			assert.Equalf(t, tt.want, result, "GetPaginationUrl(%v, %v)", tt.args.page, tt.args.perPage)
		})
	}
}

func TestWorkflowVars_GetRemoveFilterUrl(t *testing.T) {
	type args struct {
		name  string
		value interface{}
	}
	tests := []struct {
		name   string
		fields urlFields
		args   args
		want   string
	}{
		{
			name:   "Assignee filter removed",
			fields: urlFields{SelectedTeam: "lay", CurrentPage: 1, PerPageLimit: 25, SelectedAssignees: []string{"1", "2"}, SelectedUnassigned: "1", SelectedTaskTypes: []string{"3", "4"}, SelectedDueDateFrom: "2022-12-17", SelectedDueDateTo: "2022-12-18"},
			args:   args{name: "assignee", value: 2},
			want:   "?team=lay&page=1&per-page=25&task-type=3&task-type=4&assignee=1&unassigned=1&due-date-from=2022-12-17&due-date-to=2022-12-18",
		},
		{
			name:   "Unassigned filter removed",
			fields: urlFields{SelectedTeam: "lay", CurrentPage: 1, PerPageLimit: 25, SelectedAssignees: []string{"1", "2"}, SelectedUnassigned: "1", SelectedTaskTypes: []string{"3", "4"}, SelectedDueDateFrom: "2022-12-17", SelectedDueDateTo: "2022-12-18"},
			args:   args{name: "unassigned", value: 1},
			want:   "?team=lay&page=1&per-page=25&task-type=3&task-type=4&assignee=1&assignee=2&due-date-from=2022-12-17&due-date-to=2022-12-18",
		},
		{
			name:   "Task type filter removed",
			fields: urlFields{SelectedTeam: "lay", CurrentPage: 1, PerPageLimit: 25, SelectedAssignees: []string{"1", "2"}, SelectedUnassigned: "1", SelectedTaskTypes: []string{"3", "4"}, SelectedDueDateFrom: "2022-12-17", SelectedDueDateTo: "2022-12-18"},
			args:   args{name: "task-type", value: 3},
			want:   "?team=lay&page=1&per-page=25&task-type=4&assignee=1&assignee=2&unassigned=1&due-date-from=2022-12-17&due-date-to=2022-12-18",
		},
		{
			name:   "Due date from filter removed",
			fields: urlFields{SelectedTeam: "lay", CurrentPage: 1, PerPageLimit: 25, SelectedAssignees: []string{"1", "2"}, SelectedUnassigned: "1", SelectedTaskTypes: []string{"3", "4"}, SelectedDueDateFrom: "2022-12-17", SelectedDueDateTo: "2022-12-18"},
			args:   args{name: "due-date-from", value: "2022-12-17"},
			want:   "?team=lay&page=1&per-page=25&task-type=3&task-type=4&assignee=1&assignee=2&unassigned=1&due-date-to=2022-12-18",
		},
		{
			name:   "Due date to filter removed",
			fields: urlFields{SelectedTeam: "lay", CurrentPage: 1, PerPageLimit: 25, SelectedAssignees: []string{"1", "2"}, SelectedUnassigned: "1", SelectedTaskTypes: []string{"3", "4"}, SelectedDueDateFrom: "2022-12-17", SelectedDueDateTo: "2022-12-18"},
			args:   args{name: "due-date-to", value: "2022-12-18"},
			want:   "?team=lay&page=1&per-page=25&task-type=3&task-type=4&assignee=1&assignee=2&unassigned=1&due-date-from=2022-12-17",
		},
		{
			name:   "Page is reset back to 1 on removing a filter",
			fields: urlFields{SelectedTeam: "lay", CurrentPage: 3, PerPageLimit: 25, SelectedAssignees: []string{"1", "2"}, SelectedUnassigned: "1", SelectedTaskTypes: []string{"3", "4"}},
			args:   args{name: "task-type", value: 3},
			want:   "?team=lay&page=1&per-page=25&task-type=4&assignee=1&assignee=2&unassigned=1",
		},
		{
			name:   "All filters retained if filter not found",
			fields: urlFields{SelectedTeam: "lay", CurrentPage: 1, PerPageLimit: 25, SelectedAssignees: []string{"1", "2"}, SelectedUnassigned: "1", SelectedTaskTypes: []string{"3", "4"}, SelectedDueDateFrom: "2022-12-17", SelectedDueDateTo: "2022-12-18"},
			args:   args{name: "non-existent", value: 3},
			want:   "?team=lay&page=1&per-page=25&task-type=3&task-type=4&assignee=1&assignee=2&unassigned=1&due-date-from=2022-12-17&due-date-to=2022-12-18",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			w := createWorkflowVars(tt.fields)
			assert.Equalf(t, tt.want, w.GetRemoveFilterUrl(tt.args.name, tt.args.value), "GetRemoveFilterUrl(%v, %v)", tt.args.name, tt.args.value)
		})
	}
}

func TestWorkflowVars_GetTeamUrl(t *testing.T) {
	tests := []struct {
		name   string
		fields urlFields
		team   string
		want   string
	}{
		{
			name:   "Team is retained",
			fields: urlFields{SelectedTeam: "lay", CurrentPage: 1, PerPageLimit: 25},
			team:   "lay",
			want:   "?team=lay&page=1&per-page=25",
		},
		{
			name:   "Per page limit is retained",
			fields: urlFields{SelectedTeam: "lay", CurrentPage: 1, PerPageLimit: 50},
			team:   "pro",
			want:   "?team=pro&page=1&per-page=50",
		},
		{
			name:   "Assignees are removed",
			fields: urlFields{SelectedTeam: "lay", CurrentPage: 1, PerPageLimit: 25, SelectedAssignees: []string{"1", "2"}},
			team:   "pro",
			want:   "?team=pro&page=1&per-page=25",
		},
		{
			name:   "Unassigned is removed",
			fields: urlFields{SelectedTeam: "lay", CurrentPage: 1, PerPageLimit: 25, SelectedUnassigned: "1"},
			team:   "pro",
			want:   "?team=pro&page=1&per-page=25",
		},
		{
			name:   "Task types are retained",
			fields: urlFields{SelectedTeam: "lay", CurrentPage: 1, PerPageLimit: 25, SelectedTaskTypes: []string{"1", "2"}},
			team:   "pro",
			want:   "?team=pro&page=1&per-page=25&task-type=1&task-type=2",
		},
		{
			name:   "Due date filters are retained",
			fields: urlFields{SelectedTeam: "lay", CurrentPage: 1, PerPageLimit: 25, SelectedDueDateFrom: "2022-12-17", SelectedDueDateTo: "2022-12-18"},
			team:   "pro",
			want:   "?team=pro&page=1&per-page=25&due-date-from=2022-12-17&due-date-to=2022-12-18",
		},
		{
			name:   "Page is reset back to 1",
			fields: urlFields{SelectedTeam: "lay", CurrentPage: 2, PerPageLimit: 25, SelectedAssignees: []string{"1", "2"}, SelectedUnassigned: "1", SelectedTaskTypes: []string{"task"}},
			team:   "pro",
			want:   "?team=pro&page=1&per-page=25&task-type=task",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			w := createWorkflowVars(tt.fields)
			assert.Equalf(t, tt.want, w.GetTeamUrl(tt.team), "GetTeamUrl(%v)", tt.team)
		})
	}
}
