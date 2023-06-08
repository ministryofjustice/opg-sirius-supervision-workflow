package server

import (
	"errors"
	"github.com/ministryofjustice/opg-sirius-workflow/internal/sirius"
	"github.com/stretchr/testify/assert"
	"net/http"
	"testing"
)

type mockWorkflowVarsClient struct {
	count             map[string]int
	lastCtx           sirius.Context
	err               error
	userData          sirius.UserDetails
	teamSelectionData []sirius.ReturnedTeamCollection
}

func (m *mockWorkflowVarsClient) GetCurrentUserDetails(ctx sirius.Context) (sirius.UserDetails, error) {
	if m.count == nil {
		m.count = make(map[string]int)
	}
	m.count["GetCurrentUserDetails"] += 1
	m.lastCtx = ctx

	return m.userData, m.err
}

func (m *mockWorkflowVarsClient) GetTeamsForSelection(ctx sirius.Context) ([]sirius.ReturnedTeamCollection, error) {
	if m.count == nil {
		m.count = make(map[string]int)
	}
	m.count["GetTeamsForSelection"] += 1
	m.lastCtx = ctx

	return m.teamSelectionData, m.err
}

var mockUserDetailsData = sirius.UserDetails{
	ID:        123,
	Firstname: "John",
	Surname:   "Doe",
	Teams: []sirius.MyDetailsTeam{
		{
			TeamId:      13,
			DisplayName: "Lay Team 1 - (Supervision)",
		},
	},
}

var mockTeamSelectionData = []sirius.ReturnedTeamCollection{
	{
		Id: 13,
		Members: []sirius.TeamMember{
			{
				ID:   86,
				Name: "LayTeam1 User11",
			},
		},
		Name: "Lay Team 1 - (Supervision)",
	},
}

func TestNewWorkflowVars(t *testing.T) {
	client := &mockWorkflowVarsClient{userData: mockUserDetailsData, teamSelectionData: mockTeamSelectionData}
	r, _ := http.NewRequest("GET", "/path", nil)

	defaultTeamId := 19
	vars, err := NewWorkflowVars(client, r, defaultTeamId)

	assert.Nil(t, err)
	assert.Equal(t, WorkflowVars{
		Path:           "/path",
		XSRFToken:      "",
		MyDetails:      mockUserDetailsData,
		TeamSelection:  mockTeamSelectionData,
		SelectedTeam:   mockTeamSelectionData[0],
		SuccessMessage: "",
		Errors:         nil,
	}, *vars)
}

func TestGetLoggedInTeamId(t *testing.T) {
	assert.Equal(t, 13, getLoggedInTeamId(mockUserDetailsData, 25))
	assert.Equal(t, 25, getLoggedInTeamId(sirius.UserDetails{
		ID:          65,
		Name:        "case",
		DisplayName: "case manager",
	}, 25))
}

func TestGetSelectedTeam(t *testing.T) {
	teams := []sirius.ReturnedTeamCollection{
		{Selector: "1"},
		{Selector: "13"},
		{Selector: "2"},
	}

	tests := []struct {
		name           string
		url            string
		loggedInTeamId int
		defaultTeamId  int
		expectedTeam   sirius.ReturnedTeamCollection
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
			expectedTeam:   sirius.ReturnedTeamCollection{},
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
