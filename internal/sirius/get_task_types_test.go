package sirius

import (
	"fmt"
	"net/http"
	"testing"

	"github.com/pact-foundation/pact-go/dsl"
	"github.com/stretchr/testify/assert"
)

func TestGetTaskTypes(t *testing.T) {
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
		expectedResponse []ApiTaskTypes
		expectedError    error
	}{
		{
			name: "Test Task Types",
			setup: func() {
				pact.
					AddInteraction().
					Given("User exists").
					UponReceiving("A request to get all task types").
					WithRequest(dsl.Request{
						Method: http.MethodGet,
						Path:   dsl.String("/api/v1/tasktypes/supervision"),
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
							"task_types": dsl.Like(map[string]interface{}{
								"CWGN": dsl.Like(map[string]interface{}{
									"handle":     "CWGN",
									"incomplete": "Casework - General",
									"complete":   "Casework - General",
									"user":       true,
									"category":   "supervision",
								}),
								"ORAL": dsl.Like(map[string]interface{}{
									"handle":     "ORAL",
									"incomplete": "Order - Allocate to team",
									"complete":   "Order - Allocate to team",
									"user":       true,
									"category":   "supervision",
								}),
							}),
						}),
					})
			},
			cookies: []*http.Cookie{
				{Name: "XSRF-TOKEN", Value: "abcde"},
				{Name: "Other", Value: "other"},
			},
			expectedResponse: []ApiTaskTypes{
				{
					Handle:     "CWGN",
					Incomplete: "Casework - General",
					Complete:   "Casework - General",
					User:       true,
					Category:   "supervision",
					IsSelected: true,
				},
				{
					Handle:     "ORAL",
					Incomplete: "Order - Allocate to team",
					Complete:   "Order - Allocate to team",
					User:       true,
					Category:   "supervision",
					IsSelected: false,
				},
			},
		},
	}
	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			tc.setup()
			assert.Nil(t, pact.Verify(func() error {
				client, _ := NewClient(http.DefaultClient, fmt.Sprintf("http://localhost:%d", pact.Server.Port))
				taskTypeList, _ := client.GetTaskTypes(getContext(tc.cookies), []string{"CWGN"})
				assert.Equal(t, tc.expectedResponse, taskTypeList)
				assert.Equal(t, tc.expectedError, nil)
				return nil
			}))
		})
	}
}

func TestIsSelected(t *testing.T) {
	assert.Equal(t, IsSelected("ORAL", []string{"ORAL"}), true)
	assert.Equal(t, IsSelected("CWGN", []string{"CWGN", "ORAL"}), true)
	assert.Equal(t, IsSelected("TEST", []string{"CWGN", "ORAL"}), false)
}
