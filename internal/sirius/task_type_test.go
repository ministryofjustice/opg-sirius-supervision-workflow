package sirius

import (
	"fmt"
	"net/http"
	"testing"

	"github.com/pact-foundation/pact-go/dsl"
	"github.com/stretchr/testify/assert"
)

func TestTaskTypes(t *testing.T) {
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
		expectedResponse WholeTaskList
		expectedError    error
	}{
		{
			name: "Test Types",
			setup: func() {
				pact.
					AddInteraction().
					Given("User logged in").
					UponReceiving("A request to get task types").
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
								"handle":     dsl.Like("CDFC"),
								"incomplete": dsl.Like("Correspondence - Review failed draft"),
								"category":   dsl.Like("supervision"),
								"complete":   dsl.Like("Correspondence - Reviewed draft failure"),
								"user":       dsl.Like(true),
							}),
						}),
					})
			},
			cookies: []*http.Cookie{
				{Name: "XSRF-TOKEN", Value: "abcde"},
				{Name: "Other", Value: "other"},
			},
			expectedResponse: WholeTaskList{
				AllTaskList: ApiTaskTypes{
					Handle:     "CDFC",
					Incomplete: "Correspondence - Review failed draft",
					Category:   "supervision",
					Complete:   "Correspondence - Reviewed draft failure",
					User:       true,
				},
			},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			tc.setup()

			assert.Nil(t, pact.Verify(func() error {
				client, _ := NewClient(http.DefaultClient, fmt.Sprintf("http://localhost:%d", pact.Server.Port))

				myTaskTypes, err := client.GetTaskDetails(getContext(tc.cookies))
				assert.Equal(t, tc.expectedResponse, myTaskTypes)
				assert.Equal(t, tc.expectedError, err)
				return nil
			}))
		})
	}
}
