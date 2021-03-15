package sirius

import (
	"fmt"
	"net/http"
	"testing"

	"github.com/pact-foundation/pact-go/dsl"
	"github.com/stretchr/testify/assert"
)

type editTeamErrorsResponse struct {
	Errors *struct {
		TeamType *struct {
			Error string `json:"error" pact:"example=Invalid team type"`
		} `json:"type"`
	} `json:"validation_errors"`
}

func TestEditTeam(t *testing.T) {
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
		cookies       []*http.Cookie
		team          Team
		expectedError func(int) error
	}{
		{
			name: "OK",
			team: Team{
				ID:          65,
				DisplayName: "Test team",
				Type:        "INVESTIGATIONS",
				PhoneNumber: "014729583920",
				Email:       "test.team@opgtest.com",
			},
			setup: func() {
				pact.
					AddInteraction().
					Given("Supervision team with members exists").
					UponReceiving("A request to edit the team").
					WithRequest(dsl.Request{
						Method: http.MethodPut,
						Path:   dsl.String("/api/v1/teams/65"),
						Headers: dsl.MapMatcher{
							"X-XSRF-TOKEN":        dsl.String("abcde"),
							"Cookie":              dsl.String("XSRF-TOKEN=abcde; Other=other"),
							"OPG-Bypass-Membrane": dsl.String("1"),
							"Content-Type":        dsl.String("application/json"),
						},
						Body: map[string]interface{}{
							"email":       "test.team@opgtest.com",
							"name":        "Test team",
							"phoneNumber": "014729583920",
							"type":        "INVESTIGATIONS",
							"memberIds":   []int{},
						},
					}).
					WillRespondWith(dsl.Response{
						Status: http.StatusOK,
					})
			},
			cookies: []*http.Cookie{
				{Name: "XSRF-TOKEN", Value: "abcde"},
				{Name: "Other", Value: "other"},
			},
			expectedError: func(port int) error { return nil },
		},

		{
			name: "OKSendsMembers",
			team: Team{
				ID:          65,
				DisplayName: "Test team with members",
				Type:        "INVESTIGATIONS",
				PhoneNumber: "014729583920",
				Email:       "test.team@opgtest.com",
				Members: []TeamMember{
					{
						ID:    23,
						Email: "someone@opgtest.com",
					},
				},
			},
			setup: func() {
				pact.
					AddInteraction().
					Given("Supervision team with members exists").
					UponReceiving("A request to edit the team with members").
					WithRequest(dsl.Request{
						Method: http.MethodPut,
						Path:   dsl.String("/api/v1/teams/65"),
						Headers: dsl.MapMatcher{
							"X-XSRF-TOKEN":        dsl.String("abcde"),
							"Cookie":              dsl.String("XSRF-TOKEN=abcde; Other=other"),
							"OPG-Bypass-Membrane": dsl.String("1"),
							"Content-Type":        dsl.String("application/json"),
						},
						Body: map[string]interface{}{
							"email":       "test.team@opgtest.com",
							"name":        "Test team with members",
							"phoneNumber": "014729583920",
							"type":        "INVESTIGATIONS",
							"memberIds":   []int{23},
						},
					}).
					WillRespondWith(dsl.Response{
						Status: http.StatusOK,
					})
			},
			cookies: []*http.Cookie{
				{Name: "XSRF-TOKEN", Value: "abcde"},
				{Name: "Other", Value: "other"},
			},
			expectedError: func(port int) error { return nil },
		},

		{
			name: "Unauthorized",
			team: Team{
				ID:          65,
				DisplayName: "Test team",
				Type:        "FINANCE",
			},
			setup: func() {
				pact.
					AddInteraction().
					Given("Supervision team with members exists").
					UponReceiving("A request to edit the team without cookies").
					WithRequest(dsl.Request{
						Method: http.MethodPut,
						Path:   dsl.String("/api/v1/teams/65"),
						Headers: dsl.MapMatcher{
							"OPG-Bypass-Membrane": dsl.String("1"),
							"Content-Type":        dsl.String("application/json"),
						},
						Body: map[string]interface{}{
							"name":        "Test team",
							"type":        "FINANCE",
							"email":       "",
							"phoneNumber": "",
							"memberIds":   []int{},
						},
					}).
					WillRespondWith(dsl.Response{
						Status: http.StatusUnauthorized,
					})
			},
			expectedError: func(port int) error { return ErrUnauthorized },
		},

		{
			name: "Validation Errors",
			team: Team{
				ID:          65,
				DisplayName: "Test duplicate finance team",
				Type:        "FINANCE",
			},
			setup: func() {
				pact.
					AddInteraction().
					Given("Supervision team with members exists").
					UponReceiving("A request to edit the team with a non-unique type").
					WithRequest(dsl.Request{
						Method: http.MethodPut,
						Path:   dsl.String("/api/v1/teams/65"),
						Headers: dsl.MapMatcher{
							"X-XSRF-TOKEN":        dsl.String("abcde"),
							"Cookie":              dsl.String("XSRF-TOKEN=abcde; Other=other"),
							"OPG-Bypass-Membrane": dsl.String("1"),
							"Content-Type":        dsl.String("application/json"),
						},
						Body: map[string]interface{}{
							"name":        "Test duplicate finance team",
							"type":        "FINANCE",
							"email":       "",
							"phoneNumber": "",
							"memberIds":   []int{},
						},
					}).
					WillRespondWith(dsl.Response{
						Status: http.StatusBadRequest,
						Body:   dsl.Match(editTeamErrorsResponse{}),
					})
			},
			cookies: []*http.Cookie{
				{Name: "XSRF-TOKEN", Value: "abcde"},
				{Name: "Other", Value: "other"},
			},
			expectedError: func(port int) error {
				return &ValidationError{
					Errors: ValidationErrors{
						"type": {
							"error": "Invalid team type",
						},
					},
				}
			},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			tc.setup()

			assert.Nil(t, pact.Verify(func() error {
				client, _ := NewClient(http.DefaultClient, fmt.Sprintf("http://localhost:%d", pact.Server.Port))

				err := client.EditTeam(getContext(tc.cookies), tc.team)

				assert.Equal(t, tc.expectedError(pact.Server.Port), err)
				return nil
			}))
		})
	}
}

func TestEditTeamStatusError(t *testing.T) {
	s := teapotServer()
	defer s.Close()

	client, _ := NewClient(http.DefaultClient, s.URL)

	err := client.EditTeam(getContext(nil), Team{ID: 65})
	assert.Equal(t, StatusError{
		Code:   http.StatusTeapot,
		URL:    s.URL + "/api/v1/teams/65",
		Method: http.MethodPut,
	}, err)
}
