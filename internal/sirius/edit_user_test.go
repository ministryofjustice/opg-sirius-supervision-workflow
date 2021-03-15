package sirius

import (
	"fmt"
	"net/http"
	"testing"

	"github.com/pact-foundation/pact-go/dsl"
	"github.com/stretchr/testify/assert"
)

type editUserErrorsResponse struct {
	Message string `json:"message" pact:"example=oops"`
}

func TestEditUser(t *testing.T) {
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
		user          AuthUser
		cookies       []*http.Cookie
		expectedError error
	}{
		{
			name: "OK",
			user: AuthUser{
				ID:           123,
				Firstname:    "a",
				Surname:      "b",
				Organisation: "d",
				Roles:        []string{"e", "f"},
				Locked:       false,
				Suspended:    true,
			},
			setup: func() {
				pact.
					AddInteraction().
					Given("A user").
					UponReceiving("A request to edit the user").
					WithRequest(dsl.Request{
						Method: http.MethodPut,
						Path:   dsl.String("/auth/user/123"),
						Headers: dsl.MapMatcher{
							"X-XSRF-TOKEN":        dsl.String("abcde"),
							"Cookie":              dsl.String("XSRF-TOKEN=abcde; Other=other"),
							"OPG-Bypass-Membrane": dsl.String("1"),
						},
						Body: map[string]interface{}{
							"id":        123,
							"firstname": "a",
							"surname":   "b",
							"roles":     []string{"e", "f", "d"},
							"locked":    false,
							"suspended": true,
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
		},

		{
			name: "Unauthorized",
			user: AuthUser{
				ID: 123,
			},
			setup: func() {
				pact.
					AddInteraction().
					Given("A user").
					UponReceiving("A request edit the user without cookies").
					WithRequest(dsl.Request{
						Method: http.MethodPut,
						Path:   dsl.String("/auth/user/123"),
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

		{
			name: "Validation Errors",
			user: AuthUser{
				ID:        123,
				Firstname: "grehjreghjerghjerghjgerhjegrhjgrehgrehjgjherbhjgergrehjreghjerghjerghjgerhjegrhjgrehgrehjgjherbhjgergrehjreghjerghjerghjgerhjegrhjgrehgrehjgjherbhjgergrehjreghjerghjerghjgerhjegrhjgrehgrehjgjherbhjgergrehjreghjerghjerghjgerhjegrhjgrehgrehjgjherbhjgergrehjreghjerghjerghjgerhjegrhjgrehgrehjgjherbhjger",
			},
			setup: func() {
				pact.
					AddInteraction().
					Given("A user").
					UponReceiving("A request to edit the user errors on validation").
					WithRequest(dsl.Request{
						Method: http.MethodPut,
						Path:   dsl.String("/auth/user/123"),
						Headers: dsl.MapMatcher{
							"X-XSRF-TOKEN":        dsl.String("abcde"),
							"Cookie":              dsl.String("XSRF-TOKEN=abcde; Other=other"),
							"OPG-Bypass-Membrane": dsl.String("1"),
						},
						Body: map[string]interface{}{
							"id":        123,
							"firstname": "grehjreghjerghjerghjgerhjegrhjgrehgrehjgjherbhjgergrehjreghjerghjerghjgerhjegrhjgrehgrehjgjherbhjgergrehjreghjerghjerghjgerhjegrhjgrehgrehjgjherbhjgergrehjreghjerghjerghjgerhjegrhjgrehgrehjgjherbhjgergrehjreghjerghjerghjgerhjegrhjgrehgrehjgjherbhjgergrehjreghjerghjerghjgerhjegrhjgrehgrehjgjherbhjger",
							"surname":   "",
							"roles":     []string{""},
							"locked":    false,
							"suspended": false,
						},
					}).
					WillRespondWith(dsl.Response{
						Status: http.StatusBadRequest,
						Body:   dsl.Match(editUserErrorsResponse{}),
					})
			},
			cookies: []*http.Cookie{
				{Name: "XSRF-TOKEN", Value: "abcde"},
				{Name: "Other", Value: "other"},
			},
			expectedError: ClientError("oops"),
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			tc.setup()

			assert.Nil(t, pact.Verify(func() error {
				client, _ := NewClient(http.DefaultClient, fmt.Sprintf("http://localhost:%d", pact.Server.Port))

				err := client.EditUser(getContext(tc.cookies), tc.user)

				assert.Equal(t, tc.expectedError, err)
				return nil
			}))
		})
	}
}

func TestEditUserStatusError(t *testing.T) {
	s := teapotServer()
	defer s.Close()

	client, _ := NewClient(http.DefaultClient, s.URL)

	err := client.EditUser(getContext(nil), AuthUser{ID: 123})
	assert.Equal(t, StatusError{
		Code:   http.StatusTeapot,
		URL:    s.URL + "/auth/user/123",
		Method: http.MethodPut,
	}, err)
}
