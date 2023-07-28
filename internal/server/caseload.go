package server

import (
	"fmt"
	"github.com/ministryofjustice/opg-sirius-workflow/internal/model"
	"github.com/ministryofjustice/opg-sirius-workflow/internal/paginate"
	"github.com/ministryofjustice/opg-sirius-workflow/internal/sirius"
	"github.com/ministryofjustice/opg-sirius-workflow/internal/urlbuilder"
	"net/http"
	"strconv"
)

type CaseloadClient interface {
	GetClientList(sirius.Context, sirius.ClientListParams) (sirius.ClientList, error)
	ReassignClientToCaseManager(sirius.Context, int, []string) (string, error)
}

type CaseloadPage struct {
	ListPage
	FilterByAssignee
	FilterByStatus
	WorkflowVars
	ClientList sirius.ClientList
}

func (cv CaseloadPage) CreateUrlBuilder() urlbuilder.UrlBuilder {
	return urlbuilder.UrlBuilder{
		Path:            "caseload",
		SelectedTeam:    cv.App.SelectedTeam.Selector,
		SelectedPerPage: cv.PerPage,
		SelectedFilters: []urlbuilder.Filter{
			urlbuilder.CreateFilter("assignee", cv.SelectedAssignees, true),
			urlbuilder.CreateFilter("unassigned", cv.SelectedUnassigned, true),
			urlbuilder.CreateFilter("status", cv.SelectedStatuses, true),
		},
	}
}

func (cv CaseloadPage) GetAppliedFilters() []string {
	var appliedFilters []string
	if cv.App.SelectedTeam.Selector == cv.SelectedUnassigned {
		appliedFilters = append(appliedFilters, cv.App.SelectedTeam.Name)
	}
	for _, u := range cv.App.SelectedTeam.GetAssigneesForFilter() {
		if u.IsSelected(cv.SelectedAssignees) {
			appliedFilters = append(appliedFilters, u.Name)
		}
	}
	for _, s := range cv.StatusOptions {
		if s.IsIn(cv.SelectedStatuses) {
			appliedFilters = append(appliedFilters, s.Label)
		}
	}
	return appliedFilters
}

func caseload(client CaseloadClient, tmpl Template) Handler {
	return func(app WorkflowVars, w http.ResponseWriter, r *http.Request) error {
		ctx := getContext(r)

		if r.Method != http.MethodGet && r.Method != http.MethodPost {
			return StatusError(http.StatusMethodNotAllowed)
		}

		if !app.SelectedTeam.IsLay() {
			page := ClientTasksPage{ListPage: ListPage{PerPage: 25}}
			return RedirectError(page.CreateUrlBuilder().GetTeamUrl(app.SelectedTeam))
		}

		if r.Method == http.MethodPost {
			err := r.ParseForm()
			if err != nil {
				return err
			}

			assignTeam := r.FormValue("assignTeam")
			newAssigneeId, err := getAssigneeIdForTask(assignTeam, r.FormValue("assignCM"))
			if err != nil {
				return err
			}

			selectedClients := r.Form["selected-clients"]

			assigneeDisplayName, err := client.ReassignClientToCaseManager(ctx, newAssigneeId, selectedClients)
			if err != nil {
				return err
			}

			app.SuccessMessage = successMessageForReassignClient(selectedClients, assigneeDisplayName)
		}

		params := r.URL.Query()
		page := paginate.GetRequestedPage(params.Get("page"))

		perPageOptions := []int{25, 50, 100}
		clientsPerPage := paginate.GetRequestedElementsPerPage(params.Get("per-page"), perPageOptions)

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

		var selectedStatuses []string
		if params.Has("status") {
			selectedStatuses = params["status"]
		}

		clientList, err := client.GetClientList(ctx, sirius.ClientListParams{
			Team:          app.SelectedTeam,
			Page:          page,
			PerPage:       clientsPerPage,
			CaseOwners:    selectedAssignees,
			OrderStatuses: selectedStatuses,
		})
		if err != nil {
			return err
		}

		vars := CaseloadPage{ClientList: clientList}

		vars.PerPage = clientsPerPage
		vars.AssigneeFilterName = "Case owner"
		vars.SelectedAssignees = userSelectedAssignees
		vars.SelectedUnassigned = selectedUnassigned
		vars.SelectedStatuses = selectedStatuses
		vars.StatusOptions = []model.RefData{
			{
				Handle: "active",
				Label:  "Active",
			},
			{
				Handle: "closed",
				Label:  "Closed",
			},
		}

		vars.App = app
		vars.UrlBuilder = vars.CreateUrlBuilder()
		vars.Pagination = paginate.Pagination{
			CurrentPage:     clientList.Pages.PageCurrent,
			TotalPages:      clientList.Pages.PageTotal,
			TotalElements:   clientList.TotalClients,
			ElementsPerPage: vars.PerPage,
			ElementName:     "clients",
			PerPageOptions:  perPageOptions,
			UrlBuilder:      vars.UrlBuilder,
		}
		vars.AppliedFilters = vars.GetAppliedFilters()

		return tmpl.Execute(w, vars)
	}
}

func successMessageForReassignClient(selectedTasks []string, assigneeDisplayName string) string {
	return fmt.Sprintf("You have reassigned %d client(s) to %s", len(selectedTasks), assigneeDisplayName)
}
