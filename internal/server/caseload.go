package server

import (
	"fmt"
	"github.com/ministryofjustice/opg-sirius-workflow/internal/sirius"
	"net/http"
)

type CaseloadClient interface {
	GetClientList(sirius.Context, int) (sirius.ClientList, error)
}

type CaseloadVars struct {
	App        WorkflowVars
	ClientList sirius.ClientList
}

func caseload(client CaseloadClient, tmpl Template) Handler {
	return func(app WorkflowVars, w http.ResponseWriter, r *http.Request) error {
		if r.Method != http.MethodGet {
			return StatusError(http.StatusMethodNotAllowed)
		}
		ctx := getContext(r)
		teamSelected := app.SelectedTeam.Id

		clientList, err := client.GetClientList(ctx, teamSelected)
		if err != nil {
			return err
		}

		vars := CaseloadVars{
			App:        app,
			ClientList: clientList,
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
