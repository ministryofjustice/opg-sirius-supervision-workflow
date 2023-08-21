package server

import (
	"fmt"
	"github.com/ministryofjustice/opg-sirius-workflow/internal/model"
	"github.com/ministryofjustice/opg-sirius-workflow/internal/paginate"
	"github.com/ministryofjustice/opg-sirius-workflow/internal/sirius"
	"github.com/ministryofjustice/opg-sirius-workflow/internal/urlbuilder"
	"net/http"
	"strconv"
	"time"
)

type ClientTasksClient interface {
	GetTaskTypes(sirius.Context, sirius.TaskTypesParams) ([]model.TaskType, error)
	GetTaskList(sirius.Context, sirius.TaskListParams) (sirius.TaskList, error)
	AssignTasksToCaseManager(sirius.Context, int, []string, string) (string, error)
}

type ClientTasksPage struct {
	ListPage
	FilterByAssignee
	FilterByDueDate
	FilterByTaskType
	TaskList sirius.TaskList
}

func (ctp ClientTasksPage) CreateUrlBuilder() urlbuilder.UrlBuilder {
	return urlbuilder.UrlBuilder{
		Path:            "client-tasks",
		SelectedTeam:    ctp.App.SelectedTeam.Selector,
		SelectedPerPage: ctp.PerPage,
		SelectedFilters: []urlbuilder.Filter{
			urlbuilder.CreateFilter("task-type", ctp.SelectedTaskTypes),
			urlbuilder.CreateFilter("assignee", ctp.SelectedAssignees, true),
			urlbuilder.CreateFilter("unassigned", ctp.SelectedUnassigned, true),
			urlbuilder.CreateFilter("due-date-from", ctp.SelectedDueDateFrom),
			urlbuilder.CreateFilter("due-date-to", ctp.SelectedDueDateTo),
		},
	}
}

func (ctp ClientTasksPage) GetAppliedFilters(dueDateFrom *time.Time, dueDateTo *time.Time) []string {
	var appliedFilters []string
	for _, u := range ctp.TaskTypes {
		if u.IsSelected(ctp.SelectedTaskTypes) {
			appliedFilters = append(appliedFilters, u.Incomplete)
		}
	}
	if ctp.App.SelectedTeam.Selector == ctp.SelectedUnassigned {
		appliedFilters = append(appliedFilters, ctp.App.SelectedTeam.Name)
	}
	for _, u := range ctp.App.SelectedTeam.GetAssigneesForFilter() {
		if u.IsSelected(ctp.SelectedAssignees) {
			appliedFilters = append(appliedFilters, u.Name)
		}
	}
	if dueDateFrom != nil {
		appliedFilters = append(appliedFilters, "Due date from "+dueDateFrom.Format("02/01/2006")+" (inclusive)")
	}
	if dueDateTo != nil {
		appliedFilters = append(appliedFilters, "Due date to "+dueDateTo.Format("02/01/2006")+" (inclusive)")
	}
	return appliedFilters
}

func clientTasks(client ClientTasksClient, tmpl Template) Handler {
	return func(app WorkflowVars, w http.ResponseWriter, r *http.Request) error {
		ctx := getContext(r)

		if r.Method != http.MethodGet && r.Method != http.MethodPost {
			return StatusError(http.StatusMethodNotAllowed)
		}

		if r.Method == http.MethodPost {
			err := r.ParseForm()
			if err != nil {
				return err
			}

			assignTeam := r.FormValue("assignTeam")
			//this is where it picks up the new user to assign task to
			newAssigneeId, err := getAssigneeIdForTask(assignTeam, r.FormValue("assignCM"))
			if err != nil {
				return err
			}

			selectedTasks := r.Form["selected-tasks"]
			prioritySelected := r.FormValue("priority")

			// Attempt to save
			assigneeDisplayName, err := client.AssignTasksToCaseManager(ctx, newAssigneeId, selectedTasks, prioritySelected)
			if err != nil {
				return err
			}

			app.SuccessMessage = successMessageForReassignAndPrioritiseTasks(assignTeam, prioritySelected, selectedTasks, assigneeDisplayName)
		}

		params := r.URL.Query()
		page := paginate.GetRequestedPage(params.Get("page"))
		perPageOptions := []int{25, 50, 100}
		tasksPerPage := paginate.GetRequestedElementsPerPage(params.Get("per-page"), perPageOptions)

		var userSelectedAssignees []string
		if params.Has("assignee") {
			userSelectedAssignees = params["assignee"]
		}
		selectedAssignees := userSelectedAssignees
		selectedUnassigned := params.Get("unassigned")

		if selectedUnassigned == app.SelectedTeam.Selector {
			selectedAssignees = append(selectedAssignees, strconv.Itoa(app.SelectedTeam.Id))
			for _, t := range app.SelectedTeam.Teams {
				selectedAssignees = append(selectedAssignees, strconv.Itoa(t.Id))
			}
		}

		var selectedTaskTypes []string
		if params.Has("task-type") {
			selectedTaskTypes = params["task-type"]
		}

		taskTypes, err := client.GetTaskTypes(ctx, sirius.TaskTypesParams{Category: sirius.TaskTypeCategorySupervision})
		if err != nil {
			return err
		}

		selectedDueDateFrom, err := getSelectedDateFilter(params.Get("due-date-from"))
		if err != nil {
			return err
		}

		selectedDueDateTo, err := getSelectedDateFilter(params.Get("due-date-to"))
		if err != nil {
			return err
		}

		var vars ClientTasksPage

		selectedTaskTypes = vars.ValidateSelectedTaskTypes(selectedTaskTypes, taskTypes)

		taskList, err := client.GetTaskList(ctx, sirius.TaskListParams{
			Team:              app.SelectedTeam,
			Page:              page,
			PerPage:           tasksPerPage,
			TaskTypes:         taskTypes,
			SelectedTaskTypes: selectedTaskTypes,
			Assignees:         selectedAssignees,
			DueDateFrom:       selectedDueDateFrom,
			DueDateTo:         selectedDueDateTo,
		})
		if err != nil {
			return err
		}

		vars.TaskList = taskList
		vars.PerPage = tasksPerPage
		vars.SelectedTaskTypes = selectedTaskTypes
		vars.SelectedAssignees = userSelectedAssignees
		vars.SelectedUnassigned = selectedUnassigned

		if selectedDueDateFrom != nil {
			vars.SelectedDueDateFrom = selectedDueDateFrom.Format("2006-01-02")
		}
		if selectedDueDateTo != nil {
			vars.SelectedDueDateTo = selectedDueDateTo.Format("2006-01-02")
		}

		vars.App = app
		vars.UrlBuilder = vars.CreateUrlBuilder()

		if page > taskList.Pages.PageTotal && taskList.Pages.PageTotal > 0 {
			return RedirectError(vars.UrlBuilder.GetPaginationUrl(taskList.Pages.PageTotal, tasksPerPage))
		}

		vars.Pagination = paginate.Pagination{
			CurrentPage:     taskList.Pages.PageCurrent,
			TotalPages:      taskList.Pages.PageTotal,
			TotalElements:   taskList.TotalTasks,
			ElementsPerPage: vars.PerPage,
			ElementName:     "tasks",
			PerPageOptions:  perPageOptions,
			UrlBuilder:      vars.UrlBuilder,
		}

		vars.TaskTypes = taskList.CalculateTaskTypeCounts(taskTypes)
		vars.AppliedFilters = vars.GetAppliedFilters(selectedDueDateFrom, selectedDueDateTo)

		return tmpl.Execute(w, vars)
	}
}

func getAssigneeIdForTask(teamId, assigneeId string) (int, error) {
	var assigneeIdForTask int
	var err error

	if assigneeId != "" {
		assigneeIdForTask, err = strconv.Atoi(assigneeId)
	} else if teamId != "" {
		assigneeIdForTask, err = strconv.Atoi(teamId)
	}
	if err != nil {
		return 0, err
	}
	return assigneeIdForTask, nil
}

func getSelectedDateFilter(value string) (*time.Time, error) {
	if value == "" {
		return nil, nil
	}
	parsed, err := time.Parse("2006-01-02", value)
	if err != nil {
		return nil, err
	}
	return &parsed, nil
}

func successMessageForReassignAndPrioritiseTasks(assignTeam string, prioritySelected string, selectedTasks []string, assigneeDisplayName string) string {
	if assignTeam != "0" && prioritySelected == "true" {
		return fmt.Sprintf("You have assigned %d task(s) to %s as a priority", len(selectedTasks), assigneeDisplayName)
	} else if assignTeam != "0" && prioritySelected == "false" {
		return fmt.Sprintf("You have assigned %d task(s) to %s and removed priority", len(selectedTasks), assigneeDisplayName)
	} else if assignTeam != "0" {
		return fmt.Sprintf("%d task(s) have been reassigned", len(selectedTasks))
	} else if assignTeam == "0" && prioritySelected == "true" {
		return fmt.Sprintf("You have assigned %d task(s) as a priority", len(selectedTasks))
	} else if assignTeam == "0" && prioritySelected == "false" {
		return fmt.Sprintf("You have removed %d task(s) as a priority", len(selectedTasks))
	}
	return ""
}
