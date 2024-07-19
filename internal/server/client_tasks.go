package server

import (
	"github.com/ministryofjustice/opg-go-common/paginate"
	"github.com/ministryofjustice/opg-sirius-workflow/internal/model"
	"github.com/ministryofjustice/opg-sirius-workflow/internal/sirius"
	"github.com/ministryofjustice/opg-sirius-workflow/internal/urlbuilder"
	"net/http"
	"strconv"
	"time"
)

type ClientTasksClient interface {
	GetTaskTypes(sirius.Context, sirius.TaskTypesParams) ([]model.TaskType, error)
	GetTaskList(sirius.Context, sirius.TaskListParams) (sirius.TaskList, error)
	ReassignTasks(sirius.Context, sirius.ReassignTasksParams) (string, error)
}

type ClientTasksPage struct {
	ListPage
	FilterByAssignee
	FilterByDueDate
	FilterByTaskType
	TaskList sirius.TaskList
	MyTeamId string
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
		MyTeamId: ctp.MyTeamId,
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

			app.SuccessMessage, err = client.ReassignTasks(ctx, sirius.ReassignTasksParams{
				AssignTeam: r.FormValue("assignTeam"),
				AssignCM:   r.FormValue("assignCM"),
				TaskIds:    r.Form["selected-tasks"],
				IsPriority: r.FormValue("priority"),
			})
			if err != nil {
				return err
			}
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

		if app.MyDetails.IsOnlyCaseManager() && (!params.Has("team") || params.Has("preselect")) {
			selectedAssignees = append(selectedAssignees, strconv.Itoa(app.MyDetails.Id))
			userSelectedAssignees = append(userSelectedAssignees, strconv.Itoa(app.MyDetails.Id))
		}

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

		if len(vars.App.MyDetails.Teams) >= 1 {
			vars.MyTeamId = strconv.Itoa(vars.App.MyDetails.Teams[0].Id)
		}

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

		if len(selectedTaskTypes) > 0 {
			//	make another call to get original task count
			taskList2, err := client.GetTaskList(ctx, sirius.TaskListParams{
				Team:              app.SelectedTeam,
				Page:              page,
				PerPage:           tasksPerPage,
				TaskTypes:         taskTypes,
				SelectedTaskTypes: []string{},
				Assignees:         selectedAssignees,
				DueDateFrom:       selectedDueDateFrom,
				DueDateTo:         selectedDueDateTo,
			})

			if err != nil {
				return err
			}
			vars.TaskList.MetaData.TaskTypeCount = taskList2.MetaData.TaskTypeCount
		}
		taskList.MetaData.TaskTypeCount = vars.TaskList.MetaData.TaskTypeCount


		vars.TaskTypes = taskList.CalculateTaskTypeCounts(taskTypes)
		vars.AppliedFilters = vars.GetAppliedFilters(selectedDueDateFrom, selectedDueDateTo)
		vars.FilterByAssignee.AssigneeCount = vars.TaskList.MetaData.AssigneeCount
		return tmpl.Execute(w, vars)
	}
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
