package sirius

import (
	"fmt"
	"net/http"
	"testing"

	"github.com/pact-foundation/pact-go/dsl"
	"github.com/stretchr/testify/assert"
)

func TestGetTaskType(t *testing.T) {
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
					Given("User is logged in").
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
							"task_types": dsl.EachLike(map[string]interface{}{
								"handle":     dsl.Like("Case work - General"),
								"incomplete": dsl.Like("Case work - General"),
								"category":   dsl.Like("Case work - General"),
								"complete":   dsl.Like("Case work - General"),
							}, 1),
						}),
					})
			},
			cookies: []*http.Cookie{
				{Name: "XSRF-TOKEN", Value: "abcde"},
				{Name: "Other", Value: "other"},
			},
			expectedResponse: []ApiTaskTypes{
				{
					Handle: "Case work - General",
				},
			},
		},
	}
	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			tc.setup()
			assert.Nil(t, pact.Verify(func() error {
				client, _ := NewClient(http.DefaultClient, fmt.Sprintf("http://localhost:%d", pact.Server.Port))
				taskTypeList, _ := client.GetTaskType(getContext(tc.cookies), []string{"CWGN", "CNC"})
				assert.Equal(t, tc.expectedResponse, taskTypeList[0].Handle)
				assert.Equal(t, tc.expectedError, nil)
				return nil
			}))
		})
	}
}
