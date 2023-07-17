package server

import (
	"github.com/ministryofjustice/opg-sirius-workflow/internal/model"
	"github.com/ministryofjustice/opg-sirius-workflow/internal/paginate"
	"github.com/ministryofjustice/opg-sirius-workflow/internal/sirius"
	"github.com/ministryofjustice/opg-sirius-workflow/internal/urlbuilder"
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
	UrlBuilder     urlbuilder.UrlBuilder
}

func (cv CaseloadVars) CreateUrlBuilder() urlbuilder.UrlBuilder {
	return urlbuilder.UrlBuilder{
		Path:            "caseload",
		SelectedTeam:    cv.App.SelectedTeam.Selector,
		SelectedPerPage: cv.ClientsPerPage,
	}
}

func caseload(client CaseloadClient, tmpl Template) Handler {
	return func(app WorkflowVars, w http.ResponseWriter, r *http.Request) error {
		if r.Method != http.MethodGet {
			return StatusError(http.StatusMethodNotAllowed)
		}

		if !app.SelectedTeam.IsLay() {
			urlBuilder := ClientTasksVars{TasksPerPage: 25}.CreateUrlBuilder()
			return RedirectError(urlBuilder.GetTeamUrl(app.SelectedTeam))
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

		vars.UrlBuilder = vars.CreateUrlBuilder()

		vars.Pagination = paginate.Pagination{
			CurrentPage:     clientList.Pages.PageCurrent,
			TotalPages:      clientList.Pages.PageTotal,
			TotalElements:   clientList.TotalClients,
			ElementsPerPage: vars.ClientsPerPage,
			ElementName:     "clients",
			PerPageOptions:  perPageOptions,
			UrlBuilder:      vars.UrlBuilder,
		}

		return tmpl.Execute(w, vars)
	}
}
