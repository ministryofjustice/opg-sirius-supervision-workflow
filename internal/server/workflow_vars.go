package server

import (
	"fmt"
	"github.com/ministryofjustice/opg-sirius-workflow/internal/sirius"
	"strconv"
)

type WorkflowVars struct {
	Path               string
	XSRFToken          string
	MyDetails          sirius.UserDetails
	TaskList           sirius.TaskList
	PageDetails        sirius.PageDetails
	LoadTasks          []sirius.ApiTaskTypes
	TeamSelection      []sirius.ReturnedTeamCollection
	SelectedTeam       sirius.ReturnedTeamCollection
	SelectedAssignees  []string
	SelectedUnassigned string
	SelectedTaskTypes  []string
	AppliedFilters     []string
	SuccessMessage     string
	Error              string
	Errors             sirius.ValidationErrors
}

func (w WorkflowVars) buildUrl(team string, page int, tasksPerPage int, selectedTaskTypes []string, selectedAssignees []string, selectedUnassigned string) string {
	url := fmt.Sprintf("?team=%s&page=%d&per-page=%d", team, page, tasksPerPage)
	for _, taskType := range selectedTaskTypes {
		url += "&task-type=" + taskType
	}
	for _, assignee := range selectedAssignees {
		url += "&assignee=" + assignee
	}
	if selectedUnassigned != "" {
		url += "&unassigned=" + selectedUnassigned
	}
	return url
}

func (w WorkflowVars) GetTeamUrl(team string) string {
	return w.buildUrl(team, 1, w.PageDetails.StoredTaskLimit, w.SelectedTaskTypes, []string{}, "")
}

func (w WorkflowVars) GetPaginationUrl(page int, tasksPerPage ...int) string {
	perPage := w.PageDetails.StoredTaskLimit
	if len(tasksPerPage) > 0 {
		perPage = tasksPerPage[0]
	}
	return w.buildUrl(w.SelectedTeam.Selector, page, perPage, w.SelectedTaskTypes, w.SelectedAssignees, w.SelectedUnassigned)
}

func (w WorkflowVars) GetClearFiltersUrl() string {
	return w.buildUrl(w.SelectedTeam.Selector, 1, w.PageDetails.StoredTaskLimit, []string{}, []string{}, "")
}

func (w WorkflowVars) GetRemoveFilterUrl(name string, value interface{}) string {
	taskTypes := w.SelectedTaskTypes
	assignees := w.SelectedAssignees
	unassigned := w.SelectedUnassigned

	removeFilter := func(filters []string, filter string) []string {
		var newFilters []string
		for _, v := range filters {
			if v != filter {
				newFilters = append(newFilters, v)
			}
		}
		return newFilters
	}

	var stringValue string
	switch v := value.(type) {
	case int:
		stringValue = strconv.Itoa(v)
	case string:
		stringValue = v
	}

	switch name {
	case "task-type":
		taskTypes = removeFilter(taskTypes, stringValue)
	case "assignee":
		assignees = removeFilter(assignees, stringValue)
	case "unassigned":
		unassigned = ""
	}

	return w.buildUrl(w.SelectedTeam.Selector, 1, w.PageDetails.StoredTaskLimit, taskTypes, assignees, unassigned)
}
