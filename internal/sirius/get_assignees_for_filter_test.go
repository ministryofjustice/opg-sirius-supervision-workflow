package sirius

import (
	"fmt"
	"net/http"
	"testing"

	"github.com/pact-foundation/pact-go/dsl"
	"github.com/stretchr/testify/assert"
)

func TestGetMembersForTeam(t *testing.T) {
	pact := &dsl.Pact{
		Consumer:          "sirius-workflow",
		Provider:          "sirius",
		Host:              "localhost",
		PactFileWriteMode: "merge",
		LogDir:            "../../logs",
		PactDir:           "../../pacts",
	}
	defer pact.Teardown()

	testCases := []struct {
		name             string
		setup            func()
		cookies          []*http.Cookie
		expectedResponse AssigneesTeam
		expectedError    error
	}{
		{
			name: "Test Get Members for Team",
			setup: func() {
				pact.
					AddInteraction().
					Given("User logged in").
					UponReceiving("A request to get default team members for selected team").
					WithRequest(dsl.Request{
						Method: http.MethodGet,
						Path:   dsl.String("/api/v1/teams/13"),
						Headers: dsl.MapMatcher{
							"X-XSRF-TOKEN":        dsl.String("abcde"),
							"Cookie":              dsl.String("XSRF-TOKEN=abcde; Other=other"),
							"OPG-Bypass-Membrane": dsl.String("1"),
						},
					}).
					WillRespondWith(dsl.Response{
						Status:  http.StatusOK,
						Headers: dsl.MapMatcher{"Content-Type": dsl.String("application/json")},
						Body: dsl.Like(map[string]interface{}{
							"id":   dsl.Like(13),
							"name": dsl.Like("Lay Team 1 - (Supervision)"),
						}),
					})
			},
			cookies: []*http.Cookie{
				{Name: "XSRF-TOKEN", Value: "abcde"},
				{Name: "Other", Value: "other"},
			},
			expectedResponse: AssigneesTeam{
				Id:      13,
				Name:    "Lay Team 1 - (Supervision)",
				Members: []AssigneeTeamMembers{},
			},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			tc.setup()

			assert.Nil(t, pact.Verify(func() error {
				client, _ := NewClient(http.DefaultClient, fmt.Sprintf("http://localhost:%d", pact.Server.Port))

				myTeamMembers, err := client.GetAssigneesForFilter(getContext(tc.cookies), 13, []string{""})
				assert.Equal(t, tc.expectedResponse, myTeamMembers)
				assert.Equal(t, tc.expectedError, err)
				return nil
			}))
		})
	}
}

func TestIsAssigneeSelected(t *testing.T) {
	teamMembersSelected := []string{"15", "88", "89", "90"}
	assert.Equal(t, IsAssigneeSelected(90, teamMembersSelected), true)

	teamMembersNotSelected := []string{"99", "88", "89"}
	assert.Equal(t, IsAssigneeSelected(90, teamMembersNotSelected), false)
}

func TestSortMembersAlphabetically(t *testing.T) {
	expected := []AssigneeTeamMembers{
		{93, "Andrews Andrews", "Apple", false},
		{91, "Anthony Anthony", "Baldwin", false},
		{92, "Ben Ben", "Benjamin", false},
		{88, "Cat Cat", "Cathy", false},
		{89, "LayTeam1 LayTeam1", "User10", false},
	}

	unsortedMembers := []AssigneeTeamMembers{
		{88, "Cat Cat", "Cathy", false},
		{89, "LayTeam1 LayTeam1", "User10", false},
		{91, "Anthony Anthony", "Baldwin", false},
		{92, "Ben Ben", "Benjamin", false},
		{93, "Andrews Andrews", "Apple", false},
	}

	assert.Equal(t, SortMembersAlphabetically(unsortedMembers), expected)
}
