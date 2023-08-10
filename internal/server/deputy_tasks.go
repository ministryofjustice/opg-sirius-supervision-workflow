package server

import (
	"github.com/ministryofjustice/opg-sirius-workflow/internal/paginate"
	"github.com/ministryofjustice/opg-sirius-workflow/internal/urlbuilder"
	"net/http"
)

type DeputyTasksClient interface {
}

type DeputyTasksPage struct {
	ListPage
}

func (dt DeputyTasksPage) CreateUrlBuilder() urlbuilder.UrlBuilder {
	return urlbuilder.UrlBuilder{
		Path:            "deputy-tasks",
		SelectedTeam:    dt.App.SelectedTeam.Selector,
		SelectedPerPage: dt.PerPage,
	}
}

func deputyTasks(client DeputyTasksClient, tmpl Template) Handler {
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
		tasksPerPage := paginate.GetRequestedElementsPerPage(params.Get("per-page"), perPageOptions)

		vars := DeputyTasksPage{}
		vars.PerPage = tasksPerPage
		vars.App = app
		vars.UrlBuilder = vars.CreateUrlBuilder()
		vars.Pagination = paginate.Pagination{
			CurrentPage:     1,
			ElementsPerPage: vars.PerPage,
			ElementName:     "tasks",
			PerPageOptions:  perPageOptions,
			UrlBuilder:      vars.UrlBuilder,
		}

		return tmpl.Execute(w, vars)
	}
}
