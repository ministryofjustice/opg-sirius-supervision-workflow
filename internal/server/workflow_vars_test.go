package server

import (
	"encoding/base64"
	"errors"
	"github.com/ministryofjustice/opg-sirius-workflow/internal/model"
	"github.com/ministryofjustice/opg-sirius-workflow/internal/sirius"
	"github.com/stretchr/testify/assert"
	"net/http"
	"net/http/httptest"
	"testing"
)

type mockWorkflowVarsClient struct {
	count     map[string]int
	lastCtx   sirius.Context
	err       error
	userData  model.Assignee
	teamsData []model.Team
}

func (m *mockWorkflowVarsClient) GetCurrentUserDetails(ctx sirius.Context) (model.Assignee, error) {
	if m.count == nil {
		m.count = make(map[string]int)
	}
	m.count["GetCurrentUserDetails"] += 1
	m.lastCtx = ctx

	return m.userData, m.err
}

func (m *mockWorkflowVarsClient) GetTeams(ctx sirius.Context) ([]model.Team, error) {
	if m.count == nil {
		m.count = make(map[string]int)
	}
	m.count["GetTeams"] += 1
	m.lastCtx = ctx

	return m.teamsData, m.err
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

var mockTeamsData = []model.Team{
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

func TestNewWorkflowVars(t *testing.T) {
	clientTasksTab := Tab{Title: "Client tasks", basePath: "client-tasks", IsMyTeamPage: true}
	caseloadTab := Tab{Title: "Caseload", basePath: "caseload"}
	deputyTasksTab := Tab{Title: "Deputy tasks", basePath: "deputy-tasks"}
	deputiesTab := Tab{Title: "Deputies", basePath: "deputies"}

	tests := []struct {
		teamType      string
		selector      string
		wantTabs      []Tab
		showProPaTeam bool
	}{
		{
			teamType:      "LAY",
			wantTabs:      []Tab{clientTasksTab, caseloadTab},
			showProPaTeam: false,
		},
		{
			teamType:      "LAY",
			selector:      "lay-team",
			wantTabs:      []Tab{clientTasksTab},
			showProPaTeam: false,
		},
		{
			teamType:      "PRO",
			wantTabs:      []Tab{clientTasksTab, deputyTasksTab, deputiesTab},
			showProPaTeam: true,
		},
		{
			teamType:      "PA",
			wantTabs:      []Tab{clientTasksTab, deputyTasksTab, deputiesTab},
			showProPaTeam: true,
		},
		{
			teamType:      "HW",
			wantTabs:      []Tab{clientTasksTab, caseloadTab},
			showProPaTeam: false,
		},
	}
	for _, test := range tests {
		t.Run(test.teamType+" team", func(t *testing.T) {
			team := mockTeamsData[0]
			team.Type = test.teamType
			team.Selector = test.selector
			teams := []model.Team{team}

			client := &mockWorkflowVarsClient{userData: mockUserDetailsData, teamsData: teams}
			r, _ := http.NewRequest("GET", "/path?team="+team.Selector, nil)

			envVars := EnvironmentVars{
				DefaultWorkflowTeamID: 19,
			}
			vars, err := NewWorkflowVars(client, r, envVars)

			assert.Nil(t, err)
			assert.Equal(t, WorkflowVars{
				Path:            "/path",
				XSRFToken:       "",
				MyDetails:       mockUserDetailsData,
				Teams:           teams,
				SelectedTeam:    team,
				SuccessMessage:  "",
				Errors:          nil,
				Tabs:            test.wantTabs,
				EnvironmentVars: envVars,
			}, *vars)
		})
	}
}

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
	team := mockTeamsData[0]

	assert.Equal(t, "test-path?team=13", tab.GetURL(team))
}

func TestTab_IsSelected(t *testing.T) {
	app := WorkflowVars{Path: "test-path"}
	selectedTab := Tab{basePath: "test-path"}
	unselectedTab := Tab{basePath: "other-path"}

	assert.True(t, selectedTab.IsSelected(app))
	assert.False(t, unselectedTab.IsSelected(app))
}

func TestTab_GetURLReturnsPreselected(t *testing.T) {
	tab := Tab{basePath: "test-path", IsMyTeamPage: true}
	tab2 := Tab{basePath: "test-path", IsMyTeamPage: false}

	team := mockTeamsData[0]

	assert.Equal(t, "test-path?team=13&preselect", tab.GetURL(team))
	assert.Equal(t, "test-path?team=13", tab2.GetURL(team))
}

func TestCheckIfOnMyTeamPage(t *testing.T) {
	assert.True(t, checkIfOnMyTeamPage(15, 15))
	assert.False(t, checkIfOnMyTeamPage(15, 98))
	assert.False(t, checkIfOnMyTeamPage(15, 0))
}

func TestGetSuccessMessage(t *testing.T) {
	w := httptest.NewRecorder()
	valueAsByte := []byte("test success message")
	c := &http.Cookie{Name: "success-message", Value: base64.URLEncoding.EncodeToString(valueAsByte)}
	r, _ := http.NewRequest(http.MethodGet, "test-url", nil)
	r.AddCookie(c)

	successMessage, err := getSuccessMessage(r, w, "success-message")
	assert.Nil(t, err)
	assert.Equal(t, "test success message", successMessage)

	//check once the success message is read its reset to null
	r, _ = http.NewRequest(http.MethodGet, "test-url", nil)
	successMessage, err = getSuccessMessage(r, w, "success-message")
	assert.Nil(t, err)
	assert.Equal(t, "", successMessage)
}
