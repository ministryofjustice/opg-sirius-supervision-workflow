package server

import (
	"errors"
	"github.com/ministryofjustice/opg-sirius-workflow/internal/sirius"
	"net/http"
	"strconv"
	"strings"
)

type WorkflowVars struct {
	Path           string
	XSRFToken      string
	MyDetails      sirius.UserDetails
	TeamSelection  []sirius.ReturnedTeamCollection
	SelectedTeam   sirius.ReturnedTeamCollection
	Tabs           []Tab
	SuccessMessage string
	Errors         sirius.ValidationErrors
}

type Tab struct {
	Title    string
	basePath string
}

type WorkflowVarsClient interface {
	GetCurrentUserDetails(sirius.Context) (sirius.UserDetails, error)
	GetTeamsForSelection(sirius.Context) ([]sirius.ReturnedTeamCollection, error)
}

func NewWorkflowVars(client WorkflowVarsClient, r *http.Request, defaultTeamId int) (*WorkflowVars, error) {
	ctx := getContext(r)

	myDetails, err := client.GetCurrentUserDetails(ctx)
	if err != nil {
		return nil, err
	}

	teamSelection, err := client.GetTeamsForSelection(ctx)
	if err != nil {
		return nil, err
	}

	loggedInTeamId := getLoggedInTeamId(myDetails, defaultTeamId)

	selectedTeam, err := getSelectedTeam(r, loggedInTeamId, defaultTeamId, teamSelection)
	if err != nil {
		return nil, err
	}

	vars := WorkflowVars{
		Path:          r.URL.Path,
		XSRFToken:     ctx.XSRFToken,
		MyDetails:     myDetails,
		TeamSelection: teamSelection,
		SelectedTeam:  selectedTeam,
		Tabs: []Tab{
			{
				Title:    "Client tasks",
				basePath: "client-tasks",
			},
			{
				Title:    "Caseload",
				basePath: "caseload",
			},
		},
	}

	return &vars, nil
}

func getLoggedInTeamId(myDetails sirius.UserDetails, defaultTeamId int) int {
	if len(myDetails.Teams) < 1 {
		return defaultTeamId
	} else {
		return myDetails.Teams[0].TeamId
	}
}

func getSelectedTeam(r *http.Request, loggedInTeamId int, defaultTeamId int, teamSelection []sirius.ReturnedTeamCollection) (sirius.ReturnedTeamCollection, error) {
	selectors := []string{
		r.URL.Query().Get("team"),
		strconv.Itoa(loggedInTeamId),
		strconv.Itoa(defaultTeamId),
	}

	for _, selector := range selectors {
		for _, team := range teamSelection {
			if team.Selector == selector {
				return team, nil
			}
		}
	}

	return sirius.ReturnedTeamCollection{}, errors.New("invalid team selection")
}

func (w WorkflowVars) IsTabSelected(tab Tab) bool {
	return strings.HasSuffix(w.Path, tab.basePath)
}

func (w WorkflowVars) GetTabURL(tab Tab) string {
	return tab.basePath + "?team=" + w.SelectedTeam.Selector
}
