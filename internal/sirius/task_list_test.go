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
		name                string
		setup               func()
		cookies             []*http.Cookie
		selectedTeamMembers TeamSelected
		expectedResponse    TaskList
		taskDetails         TaskDetails
		expectedError       error
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
						Path:   dsl.String("/api/v1/assignees/team/13/tasks"),
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
								}),
								"name":    dsl.Like("Case work - General"),
								"dueDate": dsl.Like("01/02/2021"),
								"caseItems": dsl.EachLike(map[string]interface{}{
									"client": dsl.Like(map[string]interface{}{
										"caseRecNumber": "caseRecNumber",
										"firstname":     "Client Alexander Zacchaeus",
										"id":            3333,
										"supervisionCaseOwner": dsl.Like(map[string]interface{}{
											"displayName": "Supervision - Team - Name",
										}),
										"surname": "Client Wolfeschlegelsteinhausenbergerdorff",
									}),
								}, 1),
							}, 100),
							"pages": dsl.Like(map[string]interface{}{
								"current": 1,
								"total":   4,
							}),
							"total": dsl.Like(100),
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
						},
						ApiTaskType:    "Case work - General",
						ApiTaskDueDate: "01/02/2021",
						ApiTaskCaseItems: []CaseItemsDetails{
							{
								CaseItemClient: ClientDetails{
									ClientCaseRecNumber: "caseRecNumber",
									ClientFirstName:     "Client Alexander Zacchaeus",
									ClientId:            3333,
									ClientSupervisionCaseOwner: SupervisionCaseOwnerDetail{
										SupervisionCaseOwnerName: "Supervision - Team - Name",
									},
									ClientSurname: "Client Wolfeschlegelsteinhausenbergerdorff",
								},
							},
						},
					},
				},
				Pages: PageDetails{
					PageCurrent: 1,
					PageTotal:   4,
				},
				TotalTasks: 100,
			},
		},
	}
	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			tc.setup()
			assert.Nil(t, pact.Verify(func() error {
				client, _ := NewClient(http.DefaultClient, fmt.Sprintf("http://localhost:%d", pact.Server.Port))
				tc.selectedTeamMembers.Id = 13
				taskList, taskDetails, err := client.GetTaskList(getContext(tc.cookies), 1, 25, tc.selectedTeamMembers)
				assert.Equal(t, tc.expectedResponse.WholeTaskList[0], taskList.WholeTaskList[0], taskDetails)
				assert.Equal(t, tc.expectedError, err)
				return nil
			}))
		})
	}
}
