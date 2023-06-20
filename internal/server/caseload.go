package server

import (
	"fmt"
	"github.com/ministryofjustice/opg-sirius-workflow/internal/sirius"
	"net/http"
)

type CaseloadClient interface {
}

type CaseloadVars struct {
	App WorkflowVars
}

func caseload(client CaseloadClient, tmpl Template) Handler {
	return func(app WorkflowVars, w http.ResponseWriter, r *http.Request) error {
		if r.Method != http.MethodGet {
			return StatusError(http.StatusMethodNotAllowed)
		}

		vars := CaseloadVars{
			App: app,
		}

		if !app.SelectedTeam.IsLay() || !app.EnvironmentVars.ShowCaseload {
			return RedirectError(ClientTasksVars{}.GetTeamUrl(app.SelectedTeam))
		}

		return tmpl.Execute(w, vars)
	}
}

func (cv CaseloadVars) buildUrl(team string) string {
	return fmt.Sprintf("caseload?team=%s", team)
}

func (cv CaseloadVars) GetTeamUrl(team sirius.Team) string {
	return cv.buildUrl(team.Selector)
}
