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
		taskDetails      TaskDetails
		expectedError    error
	}{
		{
			name: "Test Task List",
			setup: func() {
				pact.
					AddInteraction().
					Given("User is logged in").
					UponReceiving("A request to get tasks for a team which have long names").
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
									"id":          86,
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
							AssigneeId:          86,
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
				taskList, taskDetails, err := client.GetTaskList(getContext(tc.cookies), 1, 25, 13, 13, []string{"CWGN", "CNC"})
				assert.Equal(t, tc.expectedResponse.WholeTaskList[0], taskList.WholeTaskList[0], taskDetails)
				assert.Equal(t, tc.expectedError, err)
				return nil
			}))
		})
	}
}

func TestGetPreviousPageNumber(t *testing.T) {
	assert.Equal(t, getPreviousPageNumber(0), 1)
	assert.Equal(t, getPreviousPageNumber(1), 1)
	assert.Equal(t, getPreviousPageNumber(2), 1)
	assert.Equal(t, getPreviousPageNumber(3), 2)
	assert.Equal(t, getPreviousPageNumber(5), 4)
}

func setUpGetNextPageNumber(pageCurrent int, pageTotal int, totalTasks int) TaskList {
	taskList := TaskList{
		Pages: PageDetails{
			PageCurrent: pageCurrent,
			PageTotal:   pageTotal,
		},
		TotalTasks: totalTasks,
	}
	return taskList
}

func TestGetNextPageNumber(t *testing.T) {
	taskList := setUpGetNextPageNumber(1, 5, 0)

	assert.Equal(t, getNextPageNumber(taskList, 0), 2)
	assert.Equal(t, getNextPageNumber(taskList, 2), 3)
	assert.Equal(t, getNextPageNumber(taskList, 15), 5)
}

func TestGetShowingLowerLimitNumberAlwaysReturns1IfOnly1Page(t *testing.T) {
	taskList := setUpGetNextPageNumber(1, 0, 13)

	assert.Equal(t, getShowingLowerLimitNumber(taskList, 25), 1)
	assert.Equal(t, getShowingLowerLimitNumber(taskList, 50), 1)
	assert.Equal(t, getShowingLowerLimitNumber(taskList, 100), 1)
}

func TestGetShowingLowerLimitNumberAlwaysReturns0If0Tasks(t *testing.T) {
	taskList := setUpGetNextPageNumber(1, 0, 0)

	assert.Equal(t, getShowingLowerLimitNumber(taskList, 25), 0)
	assert.Equal(t, getShowingLowerLimitNumber(taskList, 50), 0)
	assert.Equal(t, getShowingLowerLimitNumber(taskList, 100), 0)
}

func TestGetShowingLowerLimitNumberCanIncrementOnPages(t *testing.T) {
	taskList := setUpGetNextPageNumber(2, 0, 100)

	assert.Equal(t, getShowingLowerLimitNumber(taskList, 25), 26)
	assert.Equal(t, getShowingLowerLimitNumber(taskList, 50), 51)
	assert.Equal(t, getShowingLowerLimitNumber(taskList, 100), 101)
}

func TestGetShowingLowerLimitNumberCanIncrementOnManyPages(t *testing.T) {
	taskList := setUpGetNextPageNumber(5, 0, 5000)

	assert.Equal(t, getShowingLowerLimitNumber(taskList, 25), 101)
	assert.Equal(t, getShowingLowerLimitNumber(taskList, 50), 201)
	assert.Equal(t, getShowingLowerLimitNumber(taskList, 100), 401)
}

func TestGetShowingUpperLimitNumberWillReturnTotalTasksIfOnFinalPage(t *testing.T) {
	taskList := setUpGetNextPageNumber(1, 0, 10)

	assert.Equal(t, getShowingUpperLimitNumber(taskList, 25), 10)
	assert.Equal(t, getShowingUpperLimitNumber(taskList, 50), 10)
	assert.Equal(t, getShowingUpperLimitNumber(taskList, 100), 10)
}

func makeListOfPagesRange(min, max int) []int {
	a := make([]int, max-min+1)
	for i := range a {
		a[i] = min + i
	}
	return a
}

func setUpPagesTests(pageCurrent int, lastPage int) (TaskList, TaskDetails) {

	ListOfPages := makeListOfPagesRange(1, lastPage)

	taskList := TaskList{
		Pages: PageDetails{
			PageCurrent: pageCurrent,
		},
	}
	taskDetails := TaskDetails{
		LastPage:    lastPage,
		ListOfPages: ListOfPages,
	}

	return taskList, taskDetails
}

func TestGetPaginationLimitsWillReturnARangeTwoBelowAndTwoAboveCurrentPage(t *testing.T) {
	taskList, taskDetails := setUpPagesTests(3, 10)

	assert.Equal(t, getPaginationLimits(taskList, taskDetails), []int{1, 2, 3, 4, 5})
}

func TestGetPaginationLimitsWillReturnARangeOnlyTwoAboveCurrentPage(t *testing.T) {
	taskList, taskDetails := setUpPagesTests(1, 10)

	assert.Equal(t, getPaginationLimits(taskList, taskDetails), []int{1, 2, 3})
}

func TestGetPaginationLimitsWillReturnARangeOneBelowAndTwoAboveCurrentPage(t *testing.T) {
	taskList, taskDetails := setUpPagesTests(2, 10)

	assert.Equal(t, getPaginationLimits(taskList, taskDetails), []int{1, 2, 3, 4})
}

func TestGetPaginationLimitsWillReturnARangeTwoBelowAndOneAboveCurrentPage(t *testing.T) {
	taskList, taskDetails := setUpPagesTests(4, 5)

	assert.Equal(t, getPaginationLimits(taskList, taskDetails), []int{2, 3, 4, 5})
}

func TestGetPaginationLimitsWillReturnARangeTwoBelowAndCurrentPage(t *testing.T) {
	taskList, taskDetails := setUpPagesTests(5, 5)

	assert.Equal(t, getPaginationLimits(taskList, taskDetails), []int{3, 4, 5})
}

func TestCreateTaskTypeFilter(t *testing.T) {
	assert.Equal(t, createTaskTypeFilter([]string{}, ""), "")
	assert.Equal(t, createTaskTypeFilter([]string{"CWGN"}, ""), "type:CWGN")
	assert.Equal(t, createTaskTypeFilter([]string{"CWGN", "ORAL"}, ""), "type:CWGN,type:ORAL")
	assert.Equal(t, createTaskTypeFilter([]string{"CWGN", "ORAL", "FAKE", "TEST"}, ""), "type:CWGN,type:ORAL,type:FAKE,type:TEST")
}

func TestGetStoredTaskFilterReturnsNilIfNoLastFilterOrIfHasNewTaskFilter(t *testing.T) {
	taskDetails := TaskDetails{
		LastFilter: "",
	}

	assert.Equal(t, getStoredTaskFilter(taskDetails, []string{}, ""), "")
	assert.Equal(t, getStoredTaskFilter(taskDetails, []string{"ORAL"}, ""), "")
}

func TestGetStoredTaskFilterReturnsLastFilter(t *testing.T) {
	taskDetails := TaskDetails{
		LastFilter: "CWGN",
	}

	assert.Equal(t, getStoredTaskFilter(taskDetails, []string{}, "type:CWGN"), "type:CWGN")
}
