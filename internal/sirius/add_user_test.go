package sirius

import (
	"fmt"
	"net/http"
	"testing"

	"github.com/pact-foundation/pact-go/dsl"
	"github.com/stretchr/testify/assert"
)

type addUserBadRequestResponse struct {
	ErrorMessages *struct {
		Email *struct {
			EmailAddressLengthExceeded string `json:"emailAddressLengthExceeded" pact:"example=The input is more than 255 characters long"`
		} `json:"email"`
	} `json:"errorMessages"`
}

func TestAddUser(t *testing.T) {
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
		email         string
		firstName     string
		lastName      string
		organisation  string
		roles         []string
		expectedError error
	}{
		{
			name: "Created",
			setup: func() {
				pact.
					AddInteraction().
					Given("An admin user").
					UponReceiving("A request to add a new user").
					WithRequest(dsl.Request{
						Method: http.MethodPost,
						Path:   dsl.String("/auth/user"),
						Headers: dsl.MapMatcher{
							"X-XSRF-TOKEN":        dsl.String("abcde"),
							"Cookie":              dsl.String("XSRF-TOKEN=abcde; Other=other"),
							"OPG-Bypass-Membrane": dsl.String("1"),
							"Content-Type":        dsl.String("application/json"),
						},
						Body: map[string]interface{}{
							"firstname": "John",
							"surname":   "Doe",
							"email":     "john.doe@example.com",
							"roles":     []string{"COP User", "other1", "other2"},
						},
					}).
					WillRespondWith(dsl.Response{
						Status: http.StatusCreated,
					})
			},
			cookies: []*http.Cookie{
				{Name: "XSRF-TOKEN", Value: "abcde"},
				{Name: "Other", Value: "other"},
			},
			firstName:    "John",
			lastName:     "Doe",
			email:        "john.doe@example.com",
			organisation: "COP User",
			roles:        []string{"other1", "other2"},
		},

		{
			name: "Unauthorized",
			setup: func() {
				pact.
					AddInteraction().
					Given("An admin user").
					UponReceiving("A request to add a new user without cookies").
					WithRequest(dsl.Request{
						Method: http.MethodPost,
						Path:   dsl.String("/auth/user"),
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
			name: "Errors",
			setup: func() {
				pact.
					AddInteraction().
					Given("An admin user").
					UponReceiving("A request to add a new user errors").
					WithRequest(dsl.Request{
						Method: http.MethodPost,
						Path:   dsl.String("/auth/user"),
						Headers: dsl.MapMatcher{
							"X-XSRF-TOKEN":        dsl.String("abcde"),
							"Cookie":              dsl.String("XSRF-TOKEN=abcde; Other=other"),
							"OPG-Bypass-Membrane": dsl.String("1"),
							"Content-Type":        dsl.String("application/json"),
						},
						Body: map[string]interface{}{
							"firstname": "John",
							"surname":   "Doe",
							"email":     "john.doefkhjgerhergjgerjkrgejgerjgerjegrjhkgrehjergjgerhjkgerhjkegrhjkgerhjkegrhjkegrhjkegrhjkgerhjkgerhjkgerhjkgerhjkgerhjkgerhjkegrhjkgerhjkgerhjkgerhjkgerhjkerghjkgerhjkgerhjkgerhjkgrhjkgrehjgerhjkgerhjkegrhjkgerhjkgrerghger@example.com",
							"roles":     []string{"COP User", "other1", "other2"},
						},
					}).
					WillRespondWith(dsl.Response{
						Status: http.StatusBadRequest,
						Body:   dsl.Match(addUserBadRequestResponse{}),
					})
			},
			cookies: []*http.Cookie{
				{Name: "XSRF-TOKEN", Value: "abcde"},
				{Name: "Other", Value: "other"},
			},
			firstName:    "John",
			lastName:     "Doe",
			email:        "john.doefkhjgerhergjgerjkrgejgerjgerjegrjhkgrehjergjgerhjkgerhjkegrhjkgerhjkegrhjkegrhjkegrhjkgerhjkgerhjkgerhjkgerhjkgerhjkgerhjkegrhjkgerhjkgerhjkgerhjkgerhjkerghjkgerhjkgerhjkgerhjkgrhjkgrehjgerhjkgerhjkegrhjkgerhjkgrerghger@example.com",
			organisation: "COP User",
			roles:        []string{"other1", "other2"},
			expectedError: ValidationError{
				Errors: ValidationErrors{
					"email": {
						"emailAddressLengthExceeded": "The input is more than 255 characters long",
					},
				},
			},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			tc.setup()

			assert.Nil(t, pact.Verify(func() error {
				client, _ := NewClient(http.DefaultClient, fmt.Sprintf("http://localhost:%d", pact.Server.Port))

				err := client.AddUser(getContext(tc.cookies), tc.email, tc.firstName, tc.lastName, tc.organisation, tc.roles)
				assert.Equal(t, tc.expectedError, err)
				return nil
			}))
		})
	}
}

func TestAddUserStatusError(t *testing.T) {
	s := teapotServer()
	defer s.Close()

	client, _ := NewClient(http.DefaultClient, s.URL)

	err := client.AddUser(getContext(nil), "", "", "", "", nil)
	assert.Equal(t, StatusError{
		Code:   http.StatusTeapot,
		URL:    s.URL + "/auth/user",
		Method: http.MethodPost,
	}, err)
}
