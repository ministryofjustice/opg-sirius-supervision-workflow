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
		expectedResponse TaskList
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
							"tasks": dsl.EachLike(map[string]interface{}{
								"assignee": dsl.Like(map[string]interface{}{
									"displayName": "DisplayName",
									// "id":          1111,
								}),
								"name":    dsl.Like("Case work - General"),
								"dueDate": dsl.Like("01/02/2021"),
								"caseItems": dsl.EachLike(map[string]interface{}{
									"client": dsl.Like(map[string]interface{}{
										"caseRecNumber": "caseRecNumber",
										"firstname":     "ClientFirstname",
										"id":            3333,
										// "middlenames":   "ClientMiddlenames",
										// "salutation":    "ClientSalutation",
										"supervisionCaseOwner": dsl.Like(map[string]interface{}{
											"displayName": "supervisionDisplayName",
											// "id":          4444,
										}),
										"surname": "ClientSurname",
										// "uId":     "ClientUId",
									}),
								}, 1),
							}, 1),
						}),
					})
			},
			cookies: []*http.Cookie{
				{Name: "XSRF-TOKEN", Value: "abcde"},
				{Name: "Other", Value: "other"},
			},

			expectedResponse: TaskList{

				AllTaskList: []ApiTask{
					{
						ApiTaskAssignee: AssigneeDetails{
							AssigneeDetailsDisplayName: "DisplayName",
							//AssigneeDetailsId:  1111,
						},
						ApiTaskType:    "Case work - General",
						ApiTaskDueDate: "01/02/2021",
						ApiTaskCaseItems: []CaseItemsDetails{
							{
								CaseItemClient: ClientDetails{
									ClientDetailsCaseRecNumber: "caseRecNumber",
									ClientDetailsFirstName:     "ClientFirstname",
									ClientDetailsId:            3333,
									//ClientDetailsMiddlenames: "ClientMiddlenames",
									//ClientDetailsSalutation:  "ClientSalutation",
									ClientDetailsSupervisionCaseOwner: SupervisionCaseOwnerDetail{
										SupervisionCaseOwnerName: "supervisionDisplayName",
										//SupervisionCaseOwnerId: 4444,
									},
									ClientDetailsSurname: "ClientSurname",
									//ClientDetailsUId:   "ClientUId",
								},
							},
						},
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

				taskList, err := client.GetTaskList(getContext(tc.cookies))
				assert.Equal(t, tc.expectedResponse, taskList)
				assert.Equal(t, tc.expectedError, err)
				return nil
			}))
		})
	}
}
