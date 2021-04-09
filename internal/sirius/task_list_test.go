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
					Given("User is logged in").
					UponReceiving("A request to get tasks which have long names").
					WithRequest(dsl.Request{
						Method: http.MethodGet,
						Path:   dsl.String("/api/v1/assignees/team/tasks"),
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
									"displayName": "Assignee Duke Clive Henry Hetley Junior Jones",
									// "id":          1111,
								}),
								"name":    dsl.Like("Case work - General"),
								"dueDate": dsl.Like("01/02/2021"),
								"caseItems": dsl.EachLike(map[string]interface{}{
									"client": dsl.Like(map[string]interface{}{
										"caseRecNumber": "caseRecNumber",
										"firstname":     "Client Alexander Zacchaeus",
										"id":            3333,
										// "middlenames":   "ClientMiddlenames",
										// "salutation":    "ClientSalutation",
										"supervisionCaseOwner": dsl.Like(map[string]interface{}{
											"displayName": "Supervision - Team - Name",
											// "id":          4444,
										}),
										"surname": "Client Wolfeschlegelsteinhausenbergerdorff",
										// "uId":     "ClientUId",
									}),
								}, 1),
							}, 1),
							"pages": dsl.Like(map[string]interface{}{
								"current": 1,
								"total":   1,
							}),
							"total": dsl.Like(1),
						}),
					})
			},
			cookies: []*http.Cookie{
				{Name: "XSRF-TOKEN", Value: "abcde"},
				{Name: "Other", Value: "other"},
			},
			expectedResponse: TaskList{
				WholeTaskList: []ApiTask{
					{
						ApiTaskAssignee: AssigneeDetails{
							AssigneeDisplayName: "Assignee Duke Clive Henry Hetley Junior Jones",
							//AssigneeId:  1111,
						},
						ApiTaskType:    "Case work - General",
						ApiTaskDueDate: "01/02/2021",
						ApiTaskCaseItems: []CaseItemsDetails{
							{
								CaseItemClient: ClientDetails{
									ClientCaseRecNumber: "caseRecNumber",
									ClientFirstName:     "Client Alexander Zacchaeus",
									ClientId:            3333,
									//ClientMiddlenames: "ClientMiddlenames",
									//ClientSalutation:  "ClientSalutation",
									ClientSupervisionCaseOwner: SupervisionCaseOwnerDetail{
										SupervisionCaseOwnerName: "Supervision - Team - Name",
										//SupervisionCaseOwnerId: 4444,
									},
									ClientSurname: "Client Wolfeschlegelsteinhausenbergerdorff",
									//ClientUId:   "ClientUId",
								},
							},
						},
					},
				},
				Pages: PageDetails{
					PageCurrent: 1,
					PageTotal:   1,
				},
				TotalTasks: 1,
			},
		},
	}
	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			tc.setup()
			assert.Nil(t, pact.Verify(func() error {
				client, _ := NewClient(http.DefaultClient, fmt.Sprintf("http://localhost:%d", pact.Server.Port))
				taskList, taskDetails, err := client.GetTaskList(getContext(tc.cookies), 1, 25)
				assert.Equal(t, tc.expectedResponse, taskList, taskDetails)
				assert.Equal(t, tc.expectedError, err)
				return nil
			}))
		})
	}
}
