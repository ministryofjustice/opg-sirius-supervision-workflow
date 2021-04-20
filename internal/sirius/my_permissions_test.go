package sirius

import (
	"fmt"
	"net/http"
	"testing"

	"github.com/pact-foundation/pact-go/dsl"
	"github.com/stretchr/testify/assert"
)

func TestPermissions(t *testing.T) {
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
		expectedResponse PermissionSet
		expectedError    error
	}{
		{
			name: "Test My Permissions",
			setup: func() {
				pact.
					AddInteraction().
					Given("User exists").
					UponReceiving("A request to get my permissions").
					WithRequest(dsl.Request{
						Method: http.MethodGet,
						Path:   dsl.String("/api/v1/permissions"),
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
							"user": map[string]interface{}{
								"permissions": dsl.EachLike("PATCH", 1),
							},
							"team": map[string]interface{}{
								"permissions": dsl.EachLike("POST", 1),
							},
						}),
					})
			},
			cookies: []*http.Cookie{
				{Name: "XSRF-TOKEN", Value: "abcde"},
				{Name: "Other", Value: "other"},
			},
			expectedResponse: PermissionSet{
				"user": PermissionGroup{Permissions: []string{"PATCH"}},
				"team": PermissionGroup{Permissions: []string{"POST"}},
			},
		},
		{
			name: "Unauthorized",
			setup: func() {
				pact.
					AddInteraction().
					Given("User exists").
					UponReceiving("A request to get my permissions without cookies").
					WithRequest(dsl.Request{
						Method: http.MethodGet,
						Path:   dsl.String("/api/v1/permissions"),
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

				myPermissions, err := client.MyPermissions(getContext(tc.cookies))
				assert.Equal(t, tc.expectedResponse, myPermissions)
				assert.Equal(t, tc.expectedError, err)
				return nil
			}))
		})
	}
}

func TestPermissionsIgnoredPact(t *testing.T) {
	// We need this test to produce a specific array of permissions in the
	// response so that Cypress tests will pass. Since Pact won't let us return
	// multiple array entries from `dsl.EachLike` we have to write a separate
	// test with the specific output.
	pact := &dsl.Pact{
		Consumer:          "ignored",
		Provider:          "ignored",
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
		expectedResponse PermissionSet
		expectedError    error
	}{
		{
			name: "OK",
			setup: func() {
				pact.
					AddInteraction().
					Given("User exists").
					UponReceiving("A request to get all the permissions I need").
					WithRequest(dsl.Request{
						Method: http.MethodGet,
						Path:   dsl.String("/api/v1/permissions"),
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
							"v1-users-updatetelephonenumber": map[string]interface{}{
								"permissions": []string{"PUT"},
							},
							"v1-users": map[string]interface{}{
								"permissions": []string{"PUT", "POST", "DELETE"},
							},
							"team": map[string]interface{}{
								"permissions": []string{"GET", "POST", "PUT"},
							},
							"v1-teams": map[string]interface{}{
								"permissions": []string{"DELETE"},
							},
						}),
					})
			},
			cookies: []*http.Cookie{
				{Name: "XSRF-TOKEN", Value: "abcde"},
				{Name: "Other", Value: "other"},
			},
			expectedResponse: PermissionSet{
				"v1-users-updatetelephonenumber": PermissionGroup{Permissions: []string{"PUT"}},
				"v1-users":                       PermissionGroup{Permissions: []string{"PUT", "POST", "DELETE"}},
				"team":                           PermissionGroup{Permissions: []string{"GET", "POST", "PUT"}},
				"v1-teams":                       PermissionGroup{Permissions: []string{"DELETE"}},
			},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			tc.setup()

			assert.Nil(t, pact.Verify(func() error {
				client, _ := NewClient(http.DefaultClient, fmt.Sprintf("http://localhost:%d", pact.Server.Port))

				myPermissions, err := client.MyPermissions(getContext(tc.cookies))
				assert.Equal(t, tc.expectedResponse, myPermissions)
				assert.Equal(t, tc.expectedError, err)
				return nil
			}))
		})
	}
}

func TestHasPermissionStatusError(t *testing.T) {
	s := teapotServer()
	defer s.Close()

	client, _ := NewClient(http.DefaultClient, s.URL)

	_, err := client.MyPermissions(getContext(nil))
	assert.Equal(t, StatusError{
		Code:   http.StatusTeapot,
		URL:    s.URL + "/api/v1/permissions",
		Method: http.MethodGet,
	}, err)
}

func TestPermissionSetChecksPermission(t *testing.T) {
	permissions := PermissionSet{
		"user": {
			Permissions: []string{"GET", "PATCH"},
		},
		"team": {
			Permissions: []string{"GET"},
		},
	}

	assert.True(t, permissions.HasPermission("user", "PATCH"))
	assert.True(t, permissions.HasPermission("team", "GET"))
	assert.True(t, permissions.HasPermission("team", "get"))
	assert.False(t, permissions.HasPermission("team", "PATCHs"))
}
