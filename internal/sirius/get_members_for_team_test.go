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
		expectedResponse TeamSelected
		expectedError    error
	}{
		{
			name: "Test Get Members for Team",
			setup: func() {
				pact.
					AddInteraction().
					Given("User logged in").
					UponReceiving("A request to get team members for selected team").
					WithRequest(dsl.Request{
						Method: http.MethodGet,
						Path:   dsl.String("/api/v1/teams/14"),
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
							"id":   dsl.Like(14),
							"name": dsl.Like("Lay Team 2 - (Supervision)"),
							"members": dsl.EachLike(map[string]interface{}{
								"id":          dsl.Like(106),
								"displayName": dsl.Like("LayTeam2 User1"),
							}, 1),
						}),
					})
			},
			cookies: []*http.Cookie{
				{Name: "XSRF-TOKEN", Value: "abcde"},
				{Name: "Other", Value: "other"},
			},
			expectedResponse: TeamSelected{
				Id:   14,
				Name: "Lay Team 2 - (Supervision)",
				Members: []TeamSelectedMembers{
					{
						TeamMembersId:   106,
						TeamMembersName: "LayTeam2 User1",
					},
				},
				selectedTeamToAssignTask: 14,
			},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			tc.setup()

			assert.Nil(t, pact.Verify(func() error {
				client, _ := NewClient(http.DefaultClient, fmt.Sprintf("http://localhost:%d", pact.Server.Port))

				myTeamMembers, err := client.GetMembersForTeam(getContext(tc.cookies), 13, 14)
				assert.Equal(t, tc.expectedResponse, myTeamMembers)
				assert.Equal(t, tc.expectedError, err)
				return nil
			}))
		})
	}
}
