package sirius

import (
	"bytes"
	"github.com/ministryofjustice/opg-sirius-workflow/internal/mocks"
	"io"
	"net/http"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestGetTeamsForSelection(t *testing.T) {
	mockClient := &mocks.MockClient{}
	client, _ := NewClient(mockClient, "http://localhost:3000")

	json := `[{
			"id":21,"displayName":"Allocations - (Supervision)", "email":"allocations.team@opgtest.com", "phoneNumber":"0123456789",
			"members":[
				{
					"id":71,"displayName":"Allocations User1", "email":"allocations@opgtest.com"
				}
			],
			"teamType":{
				"handle":"ALLOCATIONS","label":"Allocations"
			}
		}]`

	r := io.NopCloser(bytes.NewReader([]byte(json)))

	mocks.GetDoFunc = func(*http.Request) (*http.Response, error) {
		return &http.Response{
			StatusCode: 200,
			Body:       r,
		}, nil
	}

	expectedResponse := []ReturnedTeamCollection{
		{
			Id: 21,
			Members: []TeamMembers{
				{
					TeamMembersId: 71, TeamMembersName: "Allocations User1", TeamMembersDisplayName: "",
				},
			},
			Name:             "Allocations - (Supervision)",
			UserSelectedTeam: 21,
			SelectedTeamId:   0,
			Type:             "ALLOCATIONS",
			TypeLabel:        "Allocations",
			IsTeamSelected:   false,
		},
	}

	teams, err := client.GetTeamsForSelection(getContext(nil), 21, []string{""})
	assert.Equal(t, expectedResponse, teams)
	assert.Equal(t, nil, err)
}

func TestFilterOutNonLayTeamsReturnsOnlySupervisionTeams(t *testing.T) {
	teamCollection := []ReturnedTeamCollection{
		{
			Name:      "Supervision Team",
			Type:      "Supervision",
			TypeLabel: "Only",
		},
		{
			Name:      "LPA Team",
			Type:      "",
			TypeLabel: "",
		},
	}

	expectedTeamCollection := []ReturnedTeamCollection{
		{
			Name:      "Supervision Team",
			Type:      "Supervision",
			TypeLabel: "Only",
		},
	}

	notExpectedTeamCollection := []ReturnedTeamCollection{
		{
			Name:      "LPA Team",
			Type:      "",
			TypeLabel: "",
		},
	}

	assert.Equal(t, FilterOutNonLayTeams(teamCollection), expectedTeamCollection)
	assert.NotEqual(t, FilterOutNonLayTeams(teamCollection), notExpectedTeamCollection)
}

func TestGetIsTeamSelectedReturnsTrueIfTeamIdIsInAssigneeFiltersArrayAndAlsoEqualToMyTeamId(t *testing.T) {
	assigneeSelectedWithTeam := []string{"15", "88", "89"}
	assert.Equal(t, IsTeamSelected(15, assigneeSelectedWithTeam, 15), true)
}
func TestGetIsTeamSelectedReturnsFalseIfTeamIdIsNotEqualToMyTeamId(t *testing.T) {
	assigneeSelectedWithTeam := []string{"15", "88", "89"}
	assert.Equal(t, IsTeamSelected(15, assigneeSelectedWithTeam, 25), false)
}
func TestGetIsTeamSelectedReturnsFalseIfTeamIdIsNotInAssigneeFiltersArray(t *testing.T) {
	assigneeSelectedWithoutTeam := []string{"99", "88", "89"}
	assert.Equal(t, IsTeamSelected(15, assigneeSelectedWithoutTeam, 25), false)
}
