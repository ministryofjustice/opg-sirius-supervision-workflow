package server

import (
	"github.com/ministryofjustice/opg-go-common/paginate"
	"github.com/ministryofjustice/opg-sirius-workflow/internal/urlbuilder"
	"net/http"
)

type DeputiesClient interface {
}

type DeputiesPage struct {
	ListPage
}

func (dt DeputiesPage) CreateUrlBuilder() urlbuilder.UrlBuilder {
	return urlbuilder.UrlBuilder{
		Path:            "deputies",
		SelectedTeam:    dt.App.SelectedTeam.Selector,
		SelectedPerPage: dt.PerPage,
	}
}

func deputies(client DeputiesClient, tmpl Template) Handler {
	return func(app WorkflowVars, w http.ResponseWriter, r *http.Request) error {
		if r.Method != http.MethodGet {
			return StatusError(http.StatusMethodNotAllowed)
		}

		if !app.SelectedTeam.IsPro() && !app.SelectedTeam.IsPA() {
			page := ClientTasksPage{ListPage: ListPage{PerPage: 25}}
			return RedirectError(page.CreateUrlBuilder().GetTeamUrl(app.SelectedTeam))
		}

		params := r.URL.Query()

		perPageOptions := []int{25, 50, 100}
		deputiesPerPage := paginate.GetRequestedElementsPerPage(params.Get("per-page"), perPageOptions)

		vars := DeputiesPage{}
		vars.PerPage = deputiesPerPage
		vars.App = app
		vars.UrlBuilder = vars.CreateUrlBuilder()
		vars.Pagination = paginate.Pagination{
			CurrentPage:     1,
			ElementsPerPage: vars.PerPage,
			ElementName:     "deputies",
			PerPageOptions:  perPageOptions,
			UrlBuilder:      vars.UrlBuilder,
		}

		return tmpl.Execute(w, vars)
	}
}
