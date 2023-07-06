package server

import (
	"fmt"
	"github.com/ministryofjustice/opg-sirius-workflow/internal/model"
	"github.com/ministryofjustice/opg-sirius-workflow/internal/paginate"
	"github.com/ministryofjustice/opg-sirius-workflow/internal/sirius"
	"net/http"
)

type CaseloadClient interface {
	GetClientList(sirius.Context, model.Team, int, int) (sirius.ClientList, error)
}

type CaseloadVars struct {
	App            WorkflowVars
	ClientList     sirius.ClientList
	Pagination     paginate.Pagination
	ClientsPerPage int
}

func caseload(client CaseloadClient, tmpl Template) Handler {
	return func(app WorkflowVars, w http.ResponseWriter, r *http.Request) error {
		if r.Method != http.MethodGet {
			return StatusError(http.StatusMethodNotAllowed)
		}

		if !app.SelectedTeam.IsLay() {
			return RedirectError(ClientTasksVars{}.GetTeamUrl(app.SelectedTeam))
		}

		params := r.URL.Query()
		page := paginate.GetRequestedPage(params.Get("page"))

		perPageOptions := []int{25, 50, 100}
		clientsPerPage := paginate.GetRequestedElementsPerPage(params.Get("per-page"), perPageOptions)

		ctx := getContext(r)
		clientList, err := client.GetClientList(ctx, app.SelectedTeam, clientsPerPage, page)
		if err != nil {
			return err
		}

		vars := CaseloadVars{
			App:            app,
			ClientList:     clientList,
			ClientsPerPage: clientsPerPage,
		}

		vars.Pagination = paginate.Pagination{
			CurrentPage:     clientList.Pages.PageCurrent,
			TotalPages:      clientList.Pages.PageTotal,
			TotalElements:   clientList.TotalClients,
			ElementsPerPage: vars.ClientsPerPage,
			ElementName:     "clients",
			PerPageOptions:  perPageOptions,
			UrlBuilder:      vars,
		}

		return tmpl.Execute(w, vars)
	}
}

func (cv CaseloadVars) GetPaginationUrl(page int, clientsPerPage ...int) string {
	perPage := cv.ClientsPerPage
	if len(clientsPerPage) > 0 {
		perPage = clientsPerPage[0]
	}
	return cv.buildUrl(cv.App.SelectedTeam.Selector, page, perPage)
}

func (cv CaseloadVars) buildUrl(team string, page int, clientsPerPage int) string {
	return fmt.Sprintf("caseload?team=%s&page=%d&per-page=%d", team, page, clientsPerPage)
}

func (cv CaseloadVars) GetTeamUrl(team model.Team) string {
	perPage := cv.ClientsPerPage
	if perPage == 0 {
		perPage = 25
	}
	return cv.buildUrl(team.Selector, 1, perPage)
}
