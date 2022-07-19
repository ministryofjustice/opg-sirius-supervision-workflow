package sirius

import (
	"fmt"
	"net/http"
	"testing"

	"github.com/pact-foundation/pact-go/dsl"
	"github.com/stretchr/testify/assert"
)

func TestGetTeamsForSelection(t *testing.T) {
	pact := &dsl.Pact{
		Consumer:                 "sirius-workflow",
		Provider:                 "sirius",
		Host:                     "localhost",
		PactFileWriteMode:        "merge",
		LogDir:                   "../../logs",
		PactDir:                  "../../pacts",
		DisableToolValidityCheck: true,
	}
	defer pact.Teardown()

	testCases := []struct {
		name             string
		setup            func()
		cookies          []*http.Cookie
		expectedResponse []ReturnedTeamCollection
		expectedError    error
	}{
		{
			name: "Test Team Selection",
			setup: func() {
				pact.
					AddInteraction().
					Given("I am a Lay Team user").
					UponReceiving("A request to get all teams for dropdown").
					WithRequest(dsl.Request{
						Method: http.MethodGet,
						Path:   dsl.String("/api/v1/teams"),
						Headers: dsl.MapMatcher{
							"X-XSRF-TOKEN":        dsl.String("abcde"),
							"Cookie":              dsl.String("XSRF-TOKEN=abcde; Other=other"),
							"OPG-Bypass-Membrane": dsl.String("1"),
						},
					}).
					WillRespondWith(dsl.Response{
						Status:  http.StatusOK,
						Headers: dsl.MapMatcher{"Content-Type": dsl.String("application/json")},
						Body: dsl.EachLike(map[string]interface{}{
							"id":   dsl.Like(13),
							"name": dsl.Like("Lay Team 1 - (Supervision)"),
						}, 1),
					})
			},
			cookies: []*http.Cookie{
				{Name: "XSRF-TOKEN", Value: "abcde"},
				{Name: "Other", Value: "other"},
			},
			expectedResponse: []ReturnedTeamCollection(nil),
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			tc.setup()

			assert.Nil(t, pact.Verify(func() error {
				client, _ := NewClient(http.DefaultClient, fmt.Sprintf("http://localhost:%d", pact.Server.Port))

				myTeamCollection, err := client.GetTeamsForSelection(getContext(tc.cookies), 13, []string{})
				assert.Equal(t, tc.expectedResponse, myTeamCollection)
				assert.Equal(t, tc.expectedError, err)
				return nil
			}))
		})
	}
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
