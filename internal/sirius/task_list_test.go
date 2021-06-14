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
						Body:    dsl.Like(map[string]interface{}{}),
					})
			},
			cookies: []*http.Cookie{
				{Name: "XSRF-TOKEN", Value: "abcde"},
				{Name: "Other", Value: "other"},
			},
			expectedResponse: TaskList{},
		},
	}
	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			tc.setup()
			assert.Nil(t, pact.Verify(func() error {
				client, _ := NewClient(http.DefaultClient, fmt.Sprintf("http://localhost:%d", pact.Server.Port))
				taskList, taskDetails, err := client.GetTaskList(getContext(tc.cookies), 1, 25, 13, 13, []string{}, []ApiTaskTypes{})
				assert.Equal(t, tc.expectedResponse.WholeTaskList, taskList.WholeTaskList, taskDetails)
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

func SetUpTaskTypeWithACase(ApiTaskHandleInput string, ApiTaskTypeInput string, TaskTypeNameInput string, AssigneeDisplayNameInput string, AssigneeIdInput int) ApiTask {
	v := ApiTask{
		ApiTaskAssignee: CaseManagement{
			CaseManagerName: AssigneeDisplayNameInput,
			Id:              AssigneeIdInput,
		},
		ApiTaskCaseItems: []CaseItemsDetails{{
			CaseItemClient: Clients{
				ClientCaseRecNumber: "13636617",
				ClientFirstName:     "Pamela",
				ClientId:            37259351,
				ClientSupervisionCaseOwner: CaseManagement{
					Id:              4321,
					CaseManagerName: "Richard Fox",
				},
				ClientSurname: "Pragnell",
			},
		}},
		ApiTaskDueDate: "01/06/2021",
		ApiTaskId:      40904862,
		ApiTaskHandle:  ApiTaskHandleInput,
		ApiTaskType:    ApiTaskTypeInput,
		TaskTypeName:   TaskTypeNameInput,
	}
	return v
}

func SetUpTaskTypeWithoutACase(ApiTaskHandleInput string, ApiTaskTypeInput string, TaskTypeNameInput string, AssigneeDisplayNameInput string, AssigneeIdInput int) ApiTask {
	v := ApiTask{
		ApiTaskAssignee: CaseManagement{
			CaseManagerName: AssigneeDisplayNameInput,
			Id:              AssigneeIdInput,
		},
		ApiClients: []Clients{
			{
				ClientCaseRecNumber: "13636617",
				ClientFirstName:     "WithoutACase",
				ClientId:            37259351,
				ClientSupervisionCaseOwner: CaseManagement{
					Id:              1234,
					CaseManagerName: "Richard Fox",
					Team: []UserTeam{
						{
							Name: "Go TaskForce Team",
							Id:   999,
						},
					},
				},
				ClientSurname: "WithoutACase",
			},
		},
		ApiTaskDueDate: "01/06/2021",
		ApiTaskId:      40904862,
		ApiTaskHandle:  ApiTaskHandleInput,
		ApiTaskType:    ApiTaskTypeInput,
		TaskTypeName:   TaskTypeNameInput,
	}
	return v
}

func SetUpTaskTypeWithoutAClient() ApiTask {
	v := ApiTask{
		ApiTaskCaseItems: []CaseItemsDetails{
			{
				CaseItemClient: Clients{
					ClientSupervisionCaseOwner: CaseManagement{
						Team: []UserTeam{
							{
								Name: "Go TaskForce Team",
								Id:   888,
							},
						},
					},
				},
			},
		},
	}
	return v
}

func SetUpLoadTasks() []ApiTaskTypes {
	loadTasks := []ApiTaskTypes{
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
	}
	return loadTasks
}

func TestGetTaskTypeNameWillReturnIncompleteNameAsTaskTypeName(t *testing.T) {

	taskType := SetUpTaskTypeWithACase("CWGN", "", "", "", 0)
	loadTasks := SetUpLoadTasks()

	assert.Equal(t, getTaskName(taskType, loadTasks), "Casework - General")
}

func TestGetTaskTypeNameWillReturnOrginalTaskNameIfNoMatchToHandle(t *testing.T) {
	taskType := SetUpTaskTypeWithACase("FAKE", "Fake type", "", "", 0)
	loadTasks := SetUpLoadTasks()

	assert.Equal(t, getTaskName(taskType, loadTasks), "Fake type")
}

func TestGetTaskTypeNameWillOverwriteAnIncorrectNameWithHandleName(t *testing.T) {
	taskType := SetUpTaskTypeWithACase("CWGN", "Fake name that doesnt match handle", "", "", 0)
	loadTasks := SetUpLoadTasks()
	expectedResult := "Casework - General"

	assert.Equal(t, getTaskName(taskType, loadTasks), expectedResult)
}

func TestGetAssigneeDisplayNameIfTaskIsAssignedToCaseOwnerWillTakeTheCaseItems(t *testing.T) {
	taskType := SetUpTaskTypeWithACase("", "", "", "Unassigned", 0)
	expectedResult := "Richard Fox"

	assert.Equal(t, getAssigneeDisplayName(taskType), expectedResult)
}
func TestGetAssigneeDisplayNameIfTaskIsAssignedToCaseOwnerWillTakeTheClients(t *testing.T) {
	taskType := SetUpTaskTypeWithoutACase("", "", "", "Unassigned", 0)
	expectedResult := "Richard Fox"

	assert.Equal(t, getAssigneeDisplayName(taskType), expectedResult)
}

func TestGetAssigneeDisplayNameIfTaskIsNotAssignedToCaseOwnerWillTakeTheClients(t *testing.T) {
	taskType := SetUpTaskTypeWithoutACase("", "", "", "Go Taskforce", 0)
	expectedResult := "Go Taskforce"

	assert.Equal(t, getAssigneeDisplayName(taskType), expectedResult)
}

func TestGetAssigneeIdWithOutACase(t *testing.T) {
	taskType := SetUpTaskTypeWithoutACase("", "", "", "Go Taskforce", 0)
	expectedResult := 1234

	assert.Equal(t, getAssigneeId(taskType), expectedResult)
}

func TestGetAssigneeIdWithACase(t *testing.T) {
	taskType := SetUpTaskTypeWithACase("", "", "", "Go Taskforce", 0)
	expectedResult := 4321

	assert.Equal(t, getAssigneeId(taskType), expectedResult)
}

func TestGetAssigneeIdWithACaseAndAssignneNotToCaseOwner(t *testing.T) {
	taskType := SetUpTaskTypeWithACase("", "", "", "Go Taskforce", 1122)
	expectedResult := 1122

	assert.Equal(t, getAssigneeId(taskType), expectedResult)
}

func TestGetClientInformationWithACase(t *testing.T) {
	taskType := SetUpTaskTypeWithACase("", "", "", "Go Taskforce", 1122)
	expectedResult := Clients{
		ClientId:            37259351,
		ClientCaseRecNumber: "13636617",
		ClientFirstName:     "Pamela",
		ClientSurname:       "Pragnell",
		ClientSupervisionCaseOwner: CaseManagement{
			CaseManagerName: "Richard Fox",
			Id:              4321,
		},
	}

	assert.Equal(t, getClientInformation(taskType), expectedResult)
}

func SetUpUserTeamStruct(TeamName string, TeamId int) ApiTask {
	v := ApiTask{
		ApiTaskAssignee: CaseManagement{
			Team: []UserTeam{
				{
					Name: TeamName,
					Id:   TeamId,
				},
			},
		},
	}

	return v
}

func TestGetClientInformationWithoutACase(t *testing.T) {
	taskType := SetUpTaskTypeWithoutACase("", "", "", "Go Taskforce", 1122)
	expectedResult := Clients{
		ClientId:            37259351,
		ClientCaseRecNumber: "13636617",
		ClientFirstName:     "WithoutACase",
		ClientSurname:       "WithoutACase",
		ClientSupervisionCaseOwner: CaseManagement{
			CaseManagerName: "Richard Fox",
			Id:              1234,
		},
	}

	assert.Equal(t, getClientInformation(taskType), expectedResult)
}

func TestGetAssigneeTeamsReturnsOriginalContentIfGivenATeam(t *testing.T) {
	taskType := SetUpUserTeamStruct("", 0)
	expectedResult := []UserTeam{
		{
			Name: "Test Team Name",
			Id:   11,
		},
	}

	assert.Equal(t, getAssigneeTeams(taskType), expectedResult)
}

func TestGetAssigneeTeamsReplacesContentWithAPIClientsInfoIfNoTeam(t *testing.T) {
	taskType := SetUpTaskTypeWithoutACase("", "", "", "", 0)
	expectedResult := []UserTeam{
		{
			Name: "Go TaskForce Team",
			Id:   999,
		},
	}

	assert.Equal(t, getAssigneeTeams(taskType), expectedResult)
}

func TestGetAssigneeTeamsReplacesContentWithAPICaseitemsInfoIfNoTeamOrClients(t *testing.T) {
	taskType := SetUpTaskTypeWithoutAClient()
	expectedResult := []UserTeam{
		{
			Name: "Go TaskForce Team",
			Id:   888,
		},
	}

	assert.Equal(t, getAssigneeTeams(taskType), expectedResult)
}
