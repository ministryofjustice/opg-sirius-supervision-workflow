package sirius

import (
	"fmt"
	"net/http"
	"testing"

	"github.com/pact-foundation/pact-go/dsl"
	"github.com/stretchr/testify/assert"
)

func TestTaskList(t *testing.T) {
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
		expectedResponse ApiTask
		expectedError    error
	}{
		{
			name: "OK",
			setup: func() {
				pact.
					AddInteraction().
					Given("User logged in").
					UponReceiving("A request to get task list").
					WithRequest(dsl.Request{
						Method: http.MethodGet,
						Path:   dsl.String("/api/v1/assignees/65/tasks"),
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
							"name": dsl.Like("Case work - General"),
						}),
					})
			},
			cookies: []*http.Cookie{
				{Name: "XSRF-TOKEN", Value: "abcde"},
				{Name: "Other", Value: "other"},
			},
			// expectedTaskList: []TaskList{
			// 	ApiTaskId:   "123",
			// 	Status:      "Not started",
			// 	DueDate:     "01/10/2019",
			// 	Name:        "Case work - General",
			// 	Description: "Case work - General",
			// },
			expectedResponse: ApiTask{
				Name: "Case work - General",
			},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			tc.setup()

			assert.Nil(t, pact.Verify(func() error {
				client, _ := NewClient(http.DefaultClient, fmt.Sprintf("http://localhost:%d", pact.Server.Port))

				taskList, err := client.GetTaskList(getContext(tc.cookies))
				assert.Equal(t, tc.expectedResponse, taskList)
				assert.Equal(t, tc.expectedError, err)
				return nil
			}))
		})
	}
}
