package server

import (
	"errors"
	"github.com/ministryofjustice/opg-sirius-workflow/internal/model"
	"github.com/ministryofjustice/opg-sirius-workflow/internal/sirius"
	"net/http"
	"strconv"
	"strings"
)

type WorkflowVars struct {
	Path            string
	XSRFToken       string
	MyDetails       model.Assignee
	Teams           []model.Team
	SelectedTeam    model.Team
	Tabs            []Tab
	SuccessMessage  string
	Errors          sirius.ValidationErrors
	EnvironmentVars EnvironmentVars
}

type Tab struct {
	Title        string
	basePath     string
	IsMyTeamPage bool
}

type WorkflowVarsClient interface {
	GetCurrentUserDetails(sirius.Context) (model.Assignee, error)
	GetTeams(sirius.Context) ([]model.Team, error)
}

func NewWorkflowVars(client WorkflowVarsClient, r *http.Request, envVars EnvironmentVars) (*WorkflowVars, error) {
	ctx := getContext(r)

	myDetails, err := client.GetCurrentUserDetails(ctx)
	if err != nil {
		return nil, err
	}

	teams, err := client.GetTeams(ctx)
	if err != nil {
		return nil, err
	}

	loggedInTeamId := getLoggedInTeamId(myDetails, envVars.DefaultWorkflowTeamID)

	selectedTeam, err := getSelectedTeam(r, loggedInTeamId, envVars.DefaultWorkflowTeamID, teams)
	if err != nil {
		return nil, err
	}

	vars := WorkflowVars{
		Path:         r.URL.Path,
		XSRFToken:    ctx.XSRFToken,
		MyDetails:    myDetails,
		Teams:        teams,
		SelectedTeam: selectedTeam,
		Tabs: []Tab{
			{
				Title:        "Client tasks",
				basePath:     "client-tasks",
				IsMyTeamPage: checkIfOnMyTeamPage(loggedInTeamId, selectedTeam.Id),
			},
		},
		EnvironmentVars: envVars,
	}

	if (selectedTeam.IsLay() && !selectedTeam.IsFullLayTeam()) || (selectedTeam.IsHW()) {
		vars.Tabs = append(vars.Tabs,
			Tab{
				Title:    "Caseload",
				basePath: "caseload",
			})
	}

	if selectedTeam.IsPro() || selectedTeam.IsPA() {
		vars.Tabs = append(vars.Tabs,
			Tab{
				Title:    "Deputy tasks",
				basePath: "deputy-tasks",
			})
	}

	if selectedTeam.IsPro() || selectedTeam.IsPA() {
		vars.Tabs = append(vars.Tabs,
			Tab{
				Title:    "Deputies",
				basePath: "deputies",
			})
	}
	return &vars, nil
}

func getLoggedInTeamId(myDetails model.Assignee, defaultTeamId int) int {
	if len(myDetails.Teams) < 1 {
		return defaultTeamId
	} else {
		return myDetails.Teams[0].Id
	}
}

func getSelectedTeam(r *http.Request, loggedInTeamId int, defaultTeamId int, teams []model.Team) (model.Team, error) {
	selectors := []string{
		r.URL.Query().Get("team"),
		strconv.Itoa(loggedInTeamId),
		strconv.Itoa(defaultTeamId),
	}

	for _, selector := range selectors {
		for _, team := range teams {
			if team.Selector == selector {
				return team, nil
			}
		}
	}

	return model.Team{}, errors.New("invalid team selection")
}

func checkIfOnMyTeamPage(loggedInTeamId, selectedTeamId int) bool {
	return loggedInTeamId == selectedTeamId
}

func (t Tab) GetURL(team model.Team) string {
	if t.IsMyTeamPage {
		return t.basePath + "?team=" + team.Selector + "&preselect"
	}
	return t.basePath + "?team=" + team.Selector
}

func (t Tab) IsSelected(app WorkflowVars) bool {
	return strings.HasSuffix(app.Path, t.basePath)
}
