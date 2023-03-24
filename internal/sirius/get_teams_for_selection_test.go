package sirius

import (
	"bytes"
	"github.com/ministryofjustice/opg-sirius-workflow/internal/mocks"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestGetTeamsForSelection(t *testing.T) {
	logger, mockClient := SetUpTest()
	client, _ := NewClient(mockClient, "http://localhost:3000", logger)

	json := `[
		{
			"id":21,
			"displayName":"Allocations - (Supervision)",
			"email":"allocations.team@opgtest.com",
			"phoneNumber":"0123456789",
			"members":[
				{
					"id":71,
					"displayName":"Allocations User1",
					"email":"allocations@opgtest.com"
				}
			],
			"teamType":{
				"handle":"ALLOCATIONS",
				"label":"Allocations"
			}
		},
		{
			"id":22,
			"displayName":"Lay Team 1",
			"email":"lay.team.1@opgtest.com",
			"phoneNumber":"0123456789",
			"members":[],
			"teamType":{
				"handle":"LAY",
				"label":"Lay Team"
			}
		},
		{
			"id":23,
			"displayName":"Pro Team 1",
			"email":"pro.team.1@opgtest.com",
			"phoneNumber":"0123456789",
			"members":[],
			"teamType":{
				"handle":"PRO",
				"label":"Pro Team"
			}
		},
		{
			"id":24,
			"displayName":"LPA Team",
			"email":"lpa.team@opgtest.com",
			"phoneNumber":"0987654321",
			"members":[
				{
					"id":72,
					"displayName":"LPA User1",
					"email":"lpa.user@opgtest.com"
				}
			]
		}
	]`

	r := io.NopCloser(bytes.NewReader([]byte(json)))

	mocks.GetDoFunc = func(*http.Request) (*http.Response, error) {
		return &http.Response{
			StatusCode: 200,
			Body:       r,
		}, nil
	}

	expectedResponse := []ReturnedTeamCollection{
		{
			Id:        21,
			Name:      "Allocations - (Supervision)",
			Type:      "ALLOCATIONS",
			TypeLabel: "Allocations",
			Selector:  "21",
			Members: []TeamMember{
				{
					ID:   71,
					Name: "Allocations User1",
				},
			},
			Teams: []ReturnedTeamCollection{},
		},
		{
			Id:        22,
			Name:      "Lay Team 1",
			Type:      "LAY",
			TypeLabel: "Lay Team",
			Selector:  "22",
			Teams:     []ReturnedTeamCollection{},
		},
		{
			Name:     "Lay deputy team",
			Selector: "lay-team",
			Members:  []TeamMember{},
			Teams: []ReturnedTeamCollection{
				{
					Id:        22,
					Name:      "Lay Team 1",
					Type:      "LAY",
					TypeLabel: "Lay Team",
					Selector:  "22",
					Teams:     []ReturnedTeamCollection{},
				},
			},
		},
		{
			Id:        23,
			Name:      "Pro Team 1",
			Type:      "PRO",
			TypeLabel: "Pro Team",
			Selector:  "23",
			Teams:     []ReturnedTeamCollection{},
		},
		{
			Name:     "Professional deputy team",
			Selector: "pro-team",
			Members:  []TeamMember{},
			Teams: []ReturnedTeamCollection{
				{
					Id:        23,
					Name:      "Pro Team 1",
					Type:      "PRO",
					TypeLabel: "Pro Team",
					Selector:  "23",
					Teams:     []ReturnedTeamCollection{},
				},
			},
		},
	}

	teams, err := client.GetTeamsForSelection(getContext(nil))
	assert.Equal(t, expectedResponse, teams)
	assert.Equal(t, nil, err)
}

func TestGetTeamsForSelectionCanReturn500(t *testing.T) {
	logger, _ := SetUpTest()
	svr := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusInternalServerError)
	}))
	defer svr.Close()

	client, _ := NewClient(http.DefaultClient, svr.URL, logger)

	_, err := client.GetTeamsForSelection(getContext(nil))

	assert.Equal(t, StatusError{
		Code:   http.StatusInternalServerError,
		URL:    svr.URL + "/api/v1/teams",
		Method: http.MethodGet,
	}, err)
}

func TestReturnedTeamCollection_GetAssigneesForFilter(t *testing.T) {
	team := ReturnedTeamCollection{
		Members: []TeamMember{
			{ID: 1, Name: "B"},
			{ID: 2, Name: "A"},
		},
		Teams: []ReturnedTeamCollection{
			{
				Members: []TeamMember{
					{ID: 4, Name: "D"},
					{ID: 2, Name: "A"},
					{ID: 3, Name: "C"},
				},
			},
			{
				Members: []TeamMember{
					{ID: 3, Name: "C"},
				},
			},
		},
	}

	expected := []TeamMember{
		{ID: 2, Name: "A"},
		{ID: 1, Name: "B"},
		{ID: 3, Name: "C"},
		{ID: 4, Name: "D"},
	}

	assert.Equal(t, expected, team.GetAssigneesForFilter())
}

func TestReturnedTeamCollection_HasTeam(t *testing.T) {
	team := ReturnedTeamCollection{
		Id: 10,
		Teams: []ReturnedTeamCollection{
			{Id: 12},
			{Id: 13},
		},
	}

	assert.Truef(t, team.HasTeam(10), "Parent team ID 10 not found")
	assert.Truef(t, team.HasTeam(12), "Check team ID 12 not found")
	assert.Truef(t, team.HasTeam(13), "Child team ID 13 not found")
	assert.False(t, team.HasTeam(11), "Child team ID 11 should not exist")
}

func TestTeamMember_IsSelected(t *testing.T) {
	selectedTeamMember := TeamMember{ID: 10}
	unselectedTeamMember := TeamMember{ID: 11}

	selectedAssignees := []string{"9", "10", "12", "13"}

	assert.Truef(t, selectedTeamMember.IsSelected(selectedAssignees), "Team ID 10 is not selected")
	assert.False(t, unselectedTeamMember.IsSelected(selectedAssignees), "Team ID 11 is selected")
}
