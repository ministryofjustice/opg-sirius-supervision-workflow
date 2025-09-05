package server

import (
	"fmt"
	"github.com/gorilla/sessions"
	"github.com/ministryofjustice/opg-go-common/paginate"
	"github.com/ministryofjustice/opg-sirius-workflow/internal/model"
	"github.com/ministryofjustice/opg-sirius-workflow/internal/sirius"
	"github.com/ministryofjustice/opg-sirius-workflow/internal/urlbuilder"
	"net/http"
	"strconv"
)

type DeputyTasksClient interface {
	GetTaskTypes(sirius.Context, sirius.TaskTypesParams) ([]model.TaskType, error)
	GetTaskList(sirius.Context, sirius.TaskListParams) (sirius.TaskList, error)
	ReassignTasks(sirius.Context, sirius.ReassignTasksParams) (string, error)
}

type DeputyTasksPage struct {
	ListPage
	FilterByAssignee
	FilterByTaskType
	TaskList sirius.TaskList
}

func (dt DeputyTasksPage) CreateUrlBuilder() urlbuilder.UrlBuilder {
	return urlbuilder.UrlBuilder{
		Path:            "deputy-tasks",
		SelectedTeam:    dt.App.SelectedTeam.Selector,
		SelectedPerPage: dt.PerPage,
		SelectedFilters: []urlbuilder.Filter{
			urlbuilder.CreateFilter("task-type", dt.SelectedTaskTypes),
			urlbuilder.CreateFilter("assignee", dt.SelectedAssignees, true),
			urlbuilder.CreateFilter("unassigned", dt.SelectedUnassigned, true),
		},
	}
}

func (dt DeputyTasksPage) GetAppliedFilters() []string {
	var appliedFilters []string
	for _, u := range dt.TaskTypes {
		if u.IsSelected(dt.SelectedTaskTypes) {
			appliedFilters = append(appliedFilters, u.Incomplete)
		}
	}
	if dt.App.SelectedTeam.Selector == dt.SelectedUnassigned {
		appliedFilters = append(appliedFilters, dt.App.SelectedTeam.Name)
	}
	for _, u := range dt.App.SelectedTeam.GetAssigneesForFilter() {
		if u.IsSelected(dt.SelectedAssignees) {
			appliedFilters = append(appliedFilters, u.Name)
		}
	}
	return appliedFilters
}

func deputyTasks(client DeputyTasksClient, tmpl Template, cookieStore sessions.CookieStore) Handler {
	return func(app WorkflowVars, w http.ResponseWriter, r *http.Request) error {
		ctx := getContext(r)

		fmt.Println(r.FormValue("successMessage"))
		fmt.Println(r.Form["successMessage"])

		if r.Method != http.MethodGet && r.Method != http.MethodPost {
			return StatusError(http.StatusMethodNotAllowed)
		}

		if !app.SelectedTeam.IsPro() && !app.SelectedTeam.IsPA() {
			page := ClientTasksPage{ListPage: ListPage{PerPage: 25}}
			return Redirect{
				Path: page.CreateUrlBuilder().GetTeamUrl(app.SelectedTeam),
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

		var vars DeputyTasksPage
		vars.PerPage = tasksPerPage
		vars.SelectedTaskTypes = selectedTaskTypes
		vars.SelectedAssignees = userSelectedAssignees
		vars.SelectedUnassigned = selectedUnassigned
		vars.App = app

		switch r.Method {
		case http.MethodPost:
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
			vars.UrlBuilder = vars.CreateUrlBuilder()

			return Redirect{
				Path:           vars.UrlBuilder.GetPaginationUrl(page, tasksPerPage),
				SuccessMessage: app.SuccessMessage,
			}

		case http.MethodGet:

			taskTypesParams := sirius.TaskTypesParams{
				Category:  sirius.TaskTypeCategoryDeputy,
				ProDeputy: app.SelectedTeam.IsPro(),
				PADeputy:  app.SelectedTeam.IsPA(),
			}
			taskTypes, err := client.GetTaskTypes(ctx, taskTypesParams)
			if err != nil {
				return err
			}
			selectedTaskTypes = vars.ValidateSelectedTaskTypes(selectedTaskTypes, taskTypes)
			vars.SelectedTaskTypes = selectedTaskTypes

			taskList, err := client.GetTaskList(ctx, sirius.TaskListParams{
				Team:              app.SelectedTeam,
				Page:              page,
				PerPage:           tasksPerPage,
				TaskTypes:         taskTypes,
				TaskTypeCategory:  "deputy",
				SelectedTaskTypes: selectedTaskTypes,
				Assignees:         selectedAssignees,
			})
			if err != nil {
				return err
			}

			vars.TaskList = taskList
			vars.UrlBuilder = vars.CreateUrlBuilder()

			if page > taskList.Pages.PageTotal && taskList.Pages.PageTotal > 0 {
				return Redirect{
					Path: vars.UrlBuilder.GetPaginationUrl(taskList.Pages.PageTotal, tasksPerPage),
				}
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
			vars.AppliedFilters = vars.GetAppliedFilters()
			vars.AssigneeCount = vars.TaskList.MetaData.AssigneeCount

			return tmpl.Execute(w, vars)

		default:
			return StatusError(http.StatusMethodNotAllowed)
		}
	}
}
