package sirius

import (
	"fmt"
	"net/http"
	"testing"

	"github.com/pact-foundation/pact-go/dsl"
	"github.com/stretchr/testify/assert"
)

type exampleUser struct {
	ID          int    `json:"id" pact:"example=47"`
	DisplayName string `json:"displayName" pact:"example=system admin"`
	Surname     string `json:"surname" pact:"example=admin"`
	Email       string `json:"email" pact:"example=system.admin@opgtest.com"`
	Locked      bool   `json:"locked" pact:"example=false"`
	Suspended   bool   `json:"suspended" pact:"example=false"`
}

type exampleUserList struct {
	Data []exampleUser `json:"data" pact:"min=1"`
}

func TestSearchUsers(t *testing.T) {
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
		name             string
		setup            func()
		cookies          []*http.Cookie
		expectedResponse []User
		expectedError    error
	}{
		{
			name: "OK",
			setup: func() {
				pact.
					AddInteraction().
					Given("User exists").
					UponReceiving("A search for admin users").
					WithRequest(dsl.Request{
						Method: http.MethodGet,
						Path:   dsl.String("/api/search/users"),
						Query: dsl.MapMatcher{
							"query": dsl.String("admin"),
						},
						Headers: dsl.MapMatcher{
							"X-XSRF-TOKEN":        dsl.String("abcde"),
							"Cookie":              dsl.String("XSRF-TOKEN=abcde; Other=other"),
							"OPG-Bypass-Membrane": dsl.String("1"),
						},
					}).
					WillRespondWith(dsl.Response{
						Status:  http.StatusOK,
						Headers: dsl.MapMatcher{"Content-Type": dsl.String("application/json")},
						Body: dsl.Match(&exampleUserList{
							Data: []exampleUser{},
						}),
					})
			},
			cookies: []*http.Cookie{
				{Name: "XSRF-TOKEN", Value: "abcde"},
				{Name: "Other", Value: "other"},
			},
			expectedResponse: []User{
				{
					ID:          47,
					DisplayName: "system admin",
					Email:       "system.admin@opgtest.com",
					Status:      "Active",
				},
			},
		},

		{
			name: "Unauthorized",
			setup: func() {
				pact.
					AddInteraction().
					Given("User exists").
					UponReceiving("A search for admin users without cookies").
					WithRequest(dsl.Request{
						Method: http.MethodGet,
						Path:   dsl.String("/api/search/users"),
						Query: dsl.MapMatcher{
							"query": dsl.String("admin"),
						},
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

				users, err := client.SearchUsers(getContext(tc.cookies), "admin")
				assert.Equal(t, tc.expectedResponse, users)
				assert.Equal(t, tc.expectedError, err)
				return nil
			}))
		})
	}
}

func TestSearchUsersStatusError(t *testing.T) {
	s := teapotServer()
	defer s.Close()

	client, _ := NewClient(http.DefaultClient, s.URL)

	_, err := client.SearchUsers(getContext(nil), "abc")
	assert.Equal(t, StatusError{
		Code:   http.StatusTeapot,
		URL:    s.URL + "/api/search/users?query=abc",
		Method: http.MethodGet,
	}, err)
}

func TestSearchUsersTooShort(t *testing.T) {
	client, _ := NewClient(http.DefaultClient, "")

	users, err := client.SearchUsers(getContext(nil), "ad")
	assert.Nil(t, users)
	assert.Equal(t, ClientError("Search term must be at least three characters"), err)
}

func TestUserStatus(t *testing.T) {
	assert.Equal(t, "string", UserStatus("string").String())

	assert.Equal(t, "", UserStatus("string").TagColour())
	assert.Equal(t, "govuk-tag--grey", UserStatus("Suspended").TagColour())
	assert.Equal(t, "govuk-tag--orange", UserStatus("Locked").TagColour())
}
