package sirius

import (
	"fmt"
	"net/http"
	"testing"

	"github.com/pact-foundation/pact-go/dsl"
	"github.com/stretchr/testify/assert"
)

func TestTeamSelection(t *testing.T) {
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
		name                string
		setup               func()
		selectedTeamMembers TeamSelected
		cookies             []*http.Cookie
		expectedResponse    []TeamCollection
		expectedError       error
	}{
		{
			name: "Test Team Selection",
			setup: func() {
				pact.
					AddInteraction().
					Given("User logged in").
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
							"name": dsl.Like("Go TaskForce"),
							"members": dsl.EachLike(map[string]interface{}{
								"id":   dsl.Like(96),
								"name": dsl.Like("LayTeam1 - (Supervision)"),
							}, 1),
						}, 1),
					})
			},
			cookies: []*http.Cookie{
				{Name: "XSRF-TOKEN", Value: "abcde"},
				{Name: "Other", Value: "other"},
			},
			expectedResponse: []TeamCollection{
				{
					Id:   13,
					Name: "Go TaskForce",
					Members: []TeamMembers{
						{
							TeamMembersId:   96,
							TeamMembersName: "LayTeam1 - (Supervision)",
						},
					},
					UserSelectedTeam: 13,
					SelectedTeamId:   13,
				},
			},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			tc.setup()

			assert.Nil(t, pact.Verify(func() error {
				client, _ := NewClient(http.DefaultClient, fmt.Sprintf("http://localhost:%d", pact.Server.Port))

				myTeamCollection, err := client.GetTeamSelection(getContext(tc.cookies), 13, 13, tc.selectedTeamMembers)
				assert.Equal(t, tc.expectedResponse, myTeamCollection)
				assert.Equal(t, tc.expectedError, err)
				return nil
			}))
		})
	}
}
