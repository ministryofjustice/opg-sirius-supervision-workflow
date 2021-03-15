package sirius

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/pact-foundation/pact-go/dsl"
	"github.com/stretchr/testify/assert"
)

func TestDeleteTeam(t *testing.T) {
	pact := &dsl.Pact{
		Consumer:          "sirius-user-management",
		Provider:          "sirius",
		Host:              "localhost",
		PactFileWriteMode: "merge",
		LogDir:            "../../logs",
		PactDir:           "../../pacts",
	}
	defer pact.Teardown()

	testCases := []struct {
		name          string
		setup         func()
		teamID        int
		cookies       []*http.Cookie
		expectedError error
	}{
		{
			name:   "OK",
			teamID: 461,
			setup: func() {
				pact.
					AddInteraction().
					Given("A team that can be deleted").
					UponReceiving("A request to delete the team").
					WithRequest(dsl.Request{
						Method: http.MethodDelete,
						Path:   dsl.String("/api/v1/teams/461"),
						Headers: dsl.MapMatcher{
							"X-XSRF-TOKEN":        dsl.String("abcde"),
							"Cookie":              dsl.String("XSRF-TOKEN=abcde; Other=other"),
							"OPG-Bypass-Membrane": dsl.String("1"),
						},
					}).
					WillRespondWith(dsl.Response{
						Status: http.StatusNoContent,
					})
			},
			cookies: []*http.Cookie{
				{Name: "XSRF-TOKEN", Value: "abcde"},
				{Name: "Other", Value: "other"},
			},
		},

		{
			name:   "Unauthorized",
			teamID: 461,
			setup: func() {
				pact.
					AddInteraction().
					Given("A team that can be deleted").
					UponReceiving("A request to delete the team without cookies").
					WithRequest(dsl.Request{
						Method: http.MethodDelete,
						Path:   dsl.String("/api/v1/teams/461"),
						Headers: dsl.MapMatcher{
							"OPG-Bypass-Membrane": dsl.String("1"),
						},
					}).
					WillRespondWith(dsl.Response{
						Status: http.StatusUnauthorized,
					})
			},
			expectedError: ErrUnauthorized,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			tc.setup()

			assert.Nil(t, pact.Verify(func() error {
				client, _ := NewClient(http.DefaultClient, fmt.Sprintf("http://localhost:%d", pact.Server.Port))

				err := client.DeleteTeam(getContext(tc.cookies), tc.teamID)

				assert.Equal(t, tc.expectedError, err)
				return nil
			}))
		})
	}
}

func TestDeleteTeamClientError(t *testing.T) {
	s := httptest.NewServer(
		http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			http.Error(w, `{"detail":"oops"}`, http.StatusBadRequest)
		}),
	)
	defer s.Close()

	client, _ := NewClient(http.DefaultClient, s.URL)

	err := client.DeleteTeam(getContext(nil), 461)
	assert.Equal(t, ClientError("oops"), err)
}

func TestDeleteTeamStatusError(t *testing.T) {
	s := teapotServer()
	defer s.Close()

	client, _ := NewClient(http.DefaultClient, s.URL)

	err := client.DeleteTeam(getContext(nil), 461)
	assert.Equal(t, StatusError{
		Code:   http.StatusTeapot,
		URL:    s.URL + "/api/v1/teams/461",
		Method: http.MethodDelete,
	}, err)
}
