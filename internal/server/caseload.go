package server

import (
	"fmt"
	"github.com/ministryofjustice/opg-sirius-workflow/internal/sirius"
	"net/http"
)

type CaseloadClient interface {
	GetCaseloadList(sirius.Context, int) (sirius.ClientList, error)
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

		clientList, err := client.GetCaseloadList(ctx, teamSelected)
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

func (cv CaseloadVars) GetReportDueDate(order []sirius.Order) string {
	return order[0].LatestAnnualReport.DueDate
}

func (cv CaseloadVars) GetClientStatus(orders []sirius.Order) string {
	orderStatuses := make(map[string]string)

	for _, s := range orders {
		label := s.Status.Label
		orderStatuses[label] = label
	}

	statuses := []string{"Active", "Open", "Closed", "Duplicate"}
	for _, status := range statuses {
		if _, found := orderStatuses[status]; found {
			return status
		}
	}

	return ""
}
