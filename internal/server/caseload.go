package server

import (
	"github.com/ministryofjustice/opg-go-common/paginate"
	"github.com/ministryofjustice/opg-sirius-workflow/internal/model"
	"github.com/ministryofjustice/opg-sirius-workflow/internal/sirius"
	"github.com/ministryofjustice/opg-sirius-workflow/internal/urlbuilder"
	"net/http"
	"net/url"
	"strconv"
)

type CaseloadClient interface {
	GetClientList(sirius.Context, sirius.ClientListParams) (sirius.ClientList, error)
	ReassignClients(sirius.Context, sirius.ReassignClientsParams) (string, error)
}

type CaseloadPage struct {
	ListPage
	FilterByAssignee
	FilterByStatus
	FilterByDeputyType
	FilterByCaseType
	FilterByDebt
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
			urlbuilder.CreateFilter("deputy-type", cv.SelectedDeputyTypes, true),
			urlbuilder.CreateFilter("case-type", cv.SelectedCaseTypes, true),
			urlbuilder.CreateFilter("debt", cv.SelectedDebtTypes, true),
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
	for _, dt := range cv.DeputyTypes {
		if dt.IsIn(cv.SelectedDeputyTypes) {
			appliedFilters = append(appliedFilters, dt.Label)
		}
	}
	for _, ct := range cv.CaseTypes {
		if ct.IsIn(cv.SelectedCaseTypes) {
			appliedFilters = append(appliedFilters, ct.Label)
		}
	}
	for _, k := range cv.DebtTypes {
		if k.IsIn(cv.SelectedDebtTypes) {
			appliedFilters = append(appliedFilters, k.Label)
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

		if !app.SelectedTeam.IsLay() && !app.SelectedTeam.IsHW() && !app.SelectedTeam.IsClosedCases() {
			page := ClientTasksPage{ListPage: ListPage{PerPage: 25}}
			return RedirectError(page.CreateUrlBuilder().GetTeamUrl(app.SelectedTeam))
		}

		if r.Method == http.MethodPost {
			err := r.ParseForm()
			if err != nil {
				return err
			}

			app.SuccessMessage, err = client.ReassignClients(ctx, sirius.ReassignClientsParams{
				AssignTeam: r.FormValue("assignTeam"),
				AssignCM:   r.FormValue("assignCM"),
				ClientIds:  r.Form["selected-clients"],
			})
			if err != nil {
				return err
			}
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

		selectedStatuses, selectedDeputyTypes, selectedCaseTypes, selectedDebtTypes := getParams(r.URL.Query())

		clientListParams := sirius.ClientListParams{
			Team:          app.SelectedTeam,
			Page:          page,
			PerPage:       clientsPerPage,
			CaseOwners:    selectedAssignees,
			OrderStatuses: selectedStatuses,
		}

		if app.SelectedTeam.IsHW() {
			clientListParams.SubType = "hw"
			clientListParams.DeputyTypes = selectedDeputyTypes
			clientListParams.CaseTypes = selectedCaseTypes
		}

		if app.SelectedTeam.IsClosedCases() {
			clientListParams.DebtTypes = selectedDebtTypes
		}

		clientList, err := client.GetClientList(ctx, clientListParams)
		if err != nil {
			return err
		}

		vars := CaseloadPage{ClientList: clientList}

		vars.PerPage = clientsPerPage
		vars.AssigneeFilterName = "Case owner"
		vars.SelectedAssignees = userSelectedAssignees
		vars.SelectedUnassigned = selectedUnassigned
		vars.SelectedStatuses = selectedStatuses
		vars.StatusOptions = getOrderStatusOptions()
		vars.FilterByAssignee.Required = true

		if app.SelectedTeam.IsHW() {
			vars.SelectedDeputyTypes = selectedDeputyTypes
			vars.DeputyTypes = getDeputyTypes()
			vars.SelectedCaseTypes = selectedCaseTypes
			vars.CaseTypes = getCaseTypes()
		}

		if app.SelectedTeam.IsClosedCases() {
			vars.SelectedDebtTypes = selectedDebtTypes
			vars.DebtTypes = getDebtTypes()
			vars.FilterByAssignee.Required = false
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

func getCaseTypes() []model.RefData {
	return []model.RefData{
		{
			Handle: "HYBRID",
			Label:  "Hybrid",
		},
		{
			Handle: "DUAL",
			Label:  "Dual",
		},
		{
			Handle: "HW",
			Label:  "Health and welfare",
		},
		{
			Handle: "PFA",
			Label:  "Property and financial affairs",
		},
	}
}

func getDeputyTypes() []model.RefData {
	return []model.RefData{
		{
			Handle: "LAY",
			Label:  "Lay",
		},
		{
			Handle: "PRO",
			Label:  "Professional",
		},
		{
			Handle: "PA",
			Label:  "Public Authority",
		},
	}
}

func getOrderStatusOptions() []model.RefData {
	return []model.RefData{
		{
			Handle: "active",
			Label:  "Active",
		},
		{
			Handle: "closed",
			Label:  "Closed",
		},
	}
}

func getDebtTypes() []model.RefData {
	return []model.RefData{
		{
			Handle: "yes",
			Label:  "Yes",
		},
		{
			Handle: "no",
			Label:  "No",
		},
	}
}

func getParams(params url.Values) ([]string, []string, []string, []string) {
	var selectedStatuses []string
	if params.Has("status") {
		selectedStatuses = params["status"]
	}

	var selectedDeputyTypes []string
	if params.Has("deputy-type") {
		selectedDeputyTypes = params["deputy-type"]
	}

	var selectedCaseTypes []string
	if params.Has("case-type") {
		selectedCaseTypes = params["case-type"]
	}

	var selectedDebtTypes []string
	if params.Has("debt") {
		selectedDebtTypes = params["debt"]
	}

	return selectedStatuses, selectedDeputyTypes, selectedCaseTypes, selectedDebtTypes
}
