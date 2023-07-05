package server

import (
	"fmt"
	"github.com/ministryofjustice/opg-sirius-workflow/internal/model"
	"github.com/ministryofjustice/opg-sirius-workflow/internal/paginate"
	"github.com/ministryofjustice/opg-sirius-workflow/internal/sirius"
	"net/http"
	"strconv"
	"time"
)

type ClientTasksClient interface {
	GetTaskTypes(sirius.Context, []string) ([]model.TaskType, error)
	GetTaskList(sirius.Context, int, int, model.Team, []string, []model.TaskType, []string, *time.Time, *time.Time) (sirius.TaskList, error)
	AssignTasksToCaseManager(sirius.Context, int, []string, string) (string, error)
}

type ClientTasksVars struct {
	App                 WorkflowVars
	TaskList            sirius.TaskList
	TaskTypes           []model.TaskType
	SelectedAssignees   []string
	SelectedUnassigned  string
	SelectedTaskTypes   []string
	SelectedDueDateFrom string
	SelectedDueDateTo   string
	AppliedFilters      []string
	TasksPerPage        int
	Pagination          paginate.Pagination
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
			newAssigneeIdForTask, err := getAssigneeIdForTask(assignTeam, r.FormValue("assignCM"))
			if err != nil {
				return err
			}

			selectedTasks := r.Form["selected-tasks"]
			prioritySelected := r.FormValue("priority")

			// Attempt to save
			assigneeDisplayName, err := client.AssignTasksToCaseManager(ctx, newAssigneeIdForTask, selectedTasks, prioritySelected)
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

		taskTypes, err := client.GetTaskTypes(ctx, selectedTaskTypes)
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

		taskList, err := client.GetTaskList(ctx, page, tasksPerPage, app.SelectedTeam, selectedTaskTypes, taskTypes, selectedAssignees, selectedDueDateFrom, selectedDueDateTo)
		if err != nil {
			return err
		}

		vars := ClientTasksVars{
			App:                app,
			TaskList:           taskList,
			SelectedAssignees:  userSelectedAssignees,
			SelectedUnassigned: selectedUnassigned,
			SelectedTaskTypes:  selectedTaskTypes,
			TasksPerPage:       tasksPerPage,
		}

		if selectedDueDateFrom != nil {
			vars.SelectedDueDateFrom = selectedDueDateFrom.Format("2006-01-02")
		}
		if selectedDueDateTo != nil {
			vars.SelectedDueDateTo = selectedDueDateTo.Format("2006-01-02")
		}

		if page > taskList.Pages.PageTotal && taskList.Pages.PageTotal > 0 {
			return RedirectError(vars.GetPaginationUrl(taskList.Pages.PageTotal, tasksPerPage))
		}

		vars.Pagination = paginate.Pagination{
			CurrentPage:     taskList.Pages.PageCurrent,
			TotalPages:      taskList.Pages.PageTotal,
			TotalElements:   taskList.TotalTasks,
			ElementsPerPage: vars.TasksPerPage,
			ElementName:     "tasks",
			PerPageOptions:  perPageOptions,
			UrlBuilder:      vars,
		}

		vars.AppliedFilters = sirius.GetAppliedFilters(app.SelectedTeam, selectedAssignees, selectedUnassigned, taskTypes, selectedDueDateFrom, selectedDueDateTo)
		vars.TaskTypes = calculateTaskCounts(taskTypes, taskList)

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

func setTaskCount(handle string, metaData sirius.TaskList) int {
	for _, q := range metaData.MetaData.TaskTypeCount {
		if handle == q.Type {
			return q.Count
		}
	}
	return 0
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

func calculateTaskCounts(taskTypes []model.TaskType, tasks sirius.TaskList) []model.TaskType {
	var taskTypeList []model.TaskType
	ecmTasksCount := 0

	for _, t := range taskTypes {
		tasksWithCounts := model.TaskType{
			Handle:     t.Handle,
			Incomplete: t.Incomplete,
			Category:   t.Category,
			Complete:   t.Complete,
			User:       t.User,
			IsSelected: t.IsSelected,
			TaskCount:  setTaskCount(t.Handle, tasks),
		}
		if t.EcmTask {
			ecmTasksCount += tasksWithCounts.TaskCount
		}
		taskTypeList = append(taskTypeList, tasksWithCounts)
	}

	taskTypeList[0].TaskCount = ecmTasksCount
	return taskTypeList
}

func successMessageForReassignAndPrioritiseTasks(assignTeam string, prioritySelected string, selectedTasks []string, assigneeDisplayName string) string {
	if assignTeam != "0" && prioritySelected == "yes" {
		return fmt.Sprintf("You have assigned %d task(s) to %s as a priority", len(selectedTasks), assigneeDisplayName)
	} else if assignTeam != "0" && prioritySelected == "no" {
		return fmt.Sprintf("You have assigned %d task(s) to %s and removed priority", len(selectedTasks), assigneeDisplayName)
	} else if assignTeam != "0" {
		return fmt.Sprintf("%d task(s) have been reassigned", len(selectedTasks))
	} else if assignTeam == "0" && prioritySelected == "yes" {
		return fmt.Sprintf("You have assigned %d task(s) as a priority", len(selectedTasks))
	} else if assignTeam == "0" && prioritySelected == "no" {
		return fmt.Sprintf("You have removed %d task(s) as a priority", len(selectedTasks))
	}
	return ""
}

func (ctv ClientTasksVars) buildUrl(team string, page int, tasksPerPage int, selectedTaskTypes []string, selectedAssignees []string, selectedUnassigned string, dueDateFrom string, dueDateTo string) string {
	url := fmt.Sprintf("client-tasks?team=%s&page=%d&per-page=%d", team, page, tasksPerPage)
	for _, taskType := range selectedTaskTypes {
		url += "&task-type=" + taskType
	}
	for _, assignee := range selectedAssignees {
		url += "&assignee=" + assignee
	}
	if selectedUnassigned != "" {
		url += "&unassigned=" + selectedUnassigned
	}
	if dueDateFrom != "" {
		url += "&due-date-from=" + dueDateFrom
	}
	if dueDateTo != "" {
		url += "&due-date-to=" + dueDateTo
	}
	return url
}

func (ctv ClientTasksVars) GetTeamUrl(team model.Team) string {
	perPage := ctv.TasksPerPage
	if perPage == 0 {
		perPage = 25
	}
	return ctv.buildUrl(team.Selector, 1, perPage, ctv.SelectedTaskTypes, []string{}, "", ctv.SelectedDueDateFrom, ctv.SelectedDueDateTo)
}

func (ctv ClientTasksVars) GetPaginationUrl(page int, tasksPerPage ...int) string {
	perPage := ctv.TasksPerPage
	if len(tasksPerPage) > 0 {
		perPage = tasksPerPage[0]
	}
	return ctv.buildUrl(ctv.App.SelectedTeam.Selector, page, perPage, ctv.SelectedTaskTypes, ctv.SelectedAssignees, ctv.SelectedUnassigned, ctv.SelectedDueDateFrom, ctv.SelectedDueDateTo)
}

func (ctv ClientTasksVars) GetClearFiltersUrl() string {
	return ctv.buildUrl(ctv.App.SelectedTeam.Selector, 1, ctv.TasksPerPage, []string{}, []string{}, "", "", "")
}

func (ctv ClientTasksVars) GetRemoveFilterUrl(name string, value interface{}) string {
	taskTypes := ctv.SelectedTaskTypes
	assignees := ctv.SelectedAssignees
	unassigned := ctv.SelectedUnassigned
	dueDateFrom := ctv.SelectedDueDateFrom
	dueDateTo := ctv.SelectedDueDateTo

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
	case "due-date-from":
		dueDateFrom = ""
	case "due-date-to":
		dueDateTo = ""
	}

	return ctv.buildUrl(ctv.App.SelectedTeam.Selector, 1, ctv.TasksPerPage, taskTypes, assignees, unassigned, dueDateFrom, dueDateTo)
}
