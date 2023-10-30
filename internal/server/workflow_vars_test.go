package server

import (
	"errors"
	"github.com/ministryofjustice/opg-sirius-workflow/internal/model"
	"github.com/ministryofjustice/opg-sirius-workflow/internal/sirius"
	"github.com/stretchr/testify/assert"
	"net/http"
	"testing"
)

type mockWorkflowVarsClient struct {
	count             map[string]int
	lastCtx           sirius.Context
	err               error
	userData          model.Assignee
	teamSelectionData []model.Team
}

func (m *mockWorkflowVarsClient) GetCurrentUserDetails(ctx sirius.Context) (model.Assignee, error) {
	if m.count == nil {
		m.count = make(map[string]int)
	}
	m.count["GetCurrentUserDetails"] += 1
	m.lastCtx = ctx

	return m.userData, m.err
}

func (m *mockWorkflowVarsClient) GetTeamsForSelection(ctx sirius.Context, teamTypes []string) ([]model.Team, error) {
	if m.count == nil {
		m.count = make(map[string]int)
	}
	m.count["GetTeamsForSelection"] += 1
	m.lastCtx = ctx

	return m.teamSelectionData, m.err
}

var mockUserDetailsData = model.Assignee{
	Id:        123,
	Firstname: "John",
	Surname:   "Doe",
	Teams: []model.Team{
		{
			Id:   13,
			Name: "Lay Team 1 - (Supervision)",
		},
	},
}

var mockTeamSelectionData = []model.Team{
	{
		Id: 13,
		Members: []model.Assignee{
			{
				Id:   86,
				Name: "LayTeam1 User11",
			},
		},
		Name:     "Lay Team 1 - (Supervision)",
		Selector: "13",
	},
}

//func TestNewWorkflowVars(t *testing.T) {
//	clientTasksTab := Tab{Title: "Client tasks", basePath: "client-tasks"}
//	caseloadTab := Tab{Title: "Caseload", basePath: "caseload"}
//	deputyTasksTab := Tab{Title: "Deputy tasks", basePath: "deputy-tasks"}
//	deputiesTab := Tab{Title: "Deputies", basePath: "deputies"}
//
//	tests := []struct {
//		teamType string
//		selector string
//		wantTabs []Tab
//	}{
//		{
//			teamType: "LAY",
//			wantTabs: []Tab{clientTasksTab, caseloadTab},
//		},
//		{
//			teamType: "LAY",
//			selector: "lay-team",
//			wantTabs: []Tab{clientTasksTab},
//		},
//		{
//			teamType: "PRO",
//			wantTabs: []Tab{clientTasksTab, deputyTasksTab, deputiesTab},
//		},
//		{
//			teamType: "PA",
//			wantTabs: []Tab{clientTasksTab, deputyTasksTab, deputiesTab},
//		},
//		{
//			teamType: "HW",
//			wantTabs: []Tab{clientTasksTab, caseloadTab},
//		},
//	}
//	for _, test := range tests {
//		t.Run(test.teamType+" team", func(t *testing.T) {
//			team := mockTeamSelectionData[0]
//			team.Type = test.teamType
//			team.Selector = test.selector
//			teams := []model.Team{team}
//
//			client := &mockWorkflowVarsClient{userData: mockUserDetailsData, teamSelectionData: teams}
//			r, _ := http.NewRequest("GET", "/path?team="+team.Selector, nil)
//
//			envVars := EnvironmentVars{
//				DefaultTeamId: 19,
//				ShowDeputies:  true,
//			}
//			vars, err := NewWorkflowVars(client, r, envVars)
//
//			assert.Nil(t, err)
//			assert.Equal(t, WorkflowVars{
//				Path:            "/path",
//				XSRFToken:       "",
//				MyDetails:       mockUserDetailsData,
//				TeamSelection:   teams,
//				SelectedTeam:    team,
//				SuccessMessage:  "",
//				Errors:          nil,
//				Tabs:            test.wantTabs,
//				EnvironmentVars: envVars,
//			}, *vars)
//		})
//	}
//}

func TestGetLoggedInTeamId(t *testing.T) {
	assert.Equal(t, 13, getLoggedInTeamId(mockUserDetailsData, 25))
	assert.Equal(t, 25, getLoggedInTeamId(model.Assignee{
		Id:   65,
		Name: "case manager",
	}, 25))
}

func TestGetSelectedTeam(t *testing.T) {
	teams := []model.Team{
		{Selector: "1"},
		{Selector: "13"},
		{Selector: "2"},
	}

	tests := []struct {
		name           string
		url            string
		loggedInTeamId int
		defaultTeamId  int
		expectedTeam   model.Team
		expectedError  error
	}{
		{
			name:           "Select team from URL parameter",
			url:            "?team=13",
			loggedInTeamId: 1,
			defaultTeamId:  2,
			expectedTeam:   teams[1],
			expectedError:  nil,
		},
		{
			name:           "Select logged in team",
			url:            "",
			loggedInTeamId: 1,
			defaultTeamId:  2,
			expectedTeam:   teams[0],
			expectedError:  nil,
		},
		{
			name:           "Select default team if logged in team is not a valid team for Workflow",
			url:            "",
			loggedInTeamId: 20,
			defaultTeamId:  2,
			expectedTeam:   teams[2],
			expectedError:  nil,
		},
		{
			name:           "Return error if no valid team can be selected",
			url:            "?team=16",
			loggedInTeamId: 3,
			defaultTeamId:  5,
			expectedTeam:   model.Team{},
			expectedError:  errors.New("invalid team selection"),
		},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			r, _ := http.NewRequest("GET", "client-tasks"+test.url, nil)
			selectedTeam, err := getSelectedTeam(r, test.loggedInTeamId, test.defaultTeamId, teams)
			assert.Equal(t, test.expectedTeam, selectedTeam)
			assert.Equal(t, test.expectedError, err)
		})
	}
}

func TestTab_GetURL(t *testing.T) {
	tab := Tab{basePath: "test-path"}
	team := mockTeamSelectionData[0]

	assert.Equal(t, "test-path?team=13", tab.GetURL(team))
}

func TestTab_IsSelected(t *testing.T) {
	app := WorkflowVars{Path: "test-path"}
	selectedTab := Tab{basePath: "test-path"}
	unselectedTab := Tab{basePath: "other-path"}

	assert.True(t, selectedTab.IsSelected(app))
	assert.False(t, unselectedTab.IsSelected(app))
}
