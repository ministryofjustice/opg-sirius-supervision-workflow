package server

import (
	"github.com/ministryofjustice/opg-sirius-workflow/internal/paginate"
	"github.com/ministryofjustice/opg-sirius-workflow/internal/sirius"
	"github.com/ministryofjustice/opg-sirius-workflow/internal/urlbuilder"
	"net/http"
	"strconv"
)

type CaseloadClient interface {
	GetClientList(sirius.Context, sirius.ClientListParams) (sirius.ClientList, error)
}

type CaseloadPage struct {
	ListPage
	FilterByAssignee
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
	return appliedFilters
}

func caseload(client CaseloadClient, tmpl Template) Handler {
	return func(app WorkflowVars, w http.ResponseWriter, r *http.Request) error {
		if r.Method != http.MethodGet {
			return StatusError(http.StatusMethodNotAllowed)
		}

		if !app.SelectedTeam.IsLay() {
			page := ClientTasksPage{ListPage: ListPage{PerPage: 25}}
			return RedirectError(page.CreateUrlBuilder().GetTeamUrl(app.SelectedTeam))
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

		ctx := getContext(r)
		clientList, err := client.GetClientList(ctx, sirius.ClientListParams{
			Team:       app.SelectedTeam,
			Page:       page,
			PerPage:    clientsPerPage,
			CaseOwners: selectedAssignees,
		})
		if err != nil {
			return err
		}

		vars := CaseloadPage{ClientList: clientList}

		vars.PerPage = clientsPerPage
		vars.AssigneeFilterName = "Case owner"
		vars.SelectedAssignees = userSelectedAssignees
		vars.SelectedUnassigned = selectedUnassigned

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
