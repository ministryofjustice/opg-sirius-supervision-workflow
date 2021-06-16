package sirius

import (
	"fmt"
	"net/http"
	"testing"

	"github.com/pact-foundation/pact-go/dsl"
	"github.com/stretchr/testify/assert"
)

func TestGetTaskList(t *testing.T) {
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
				taskList, int, err := client.GetTaskList(getContext(tc.cookies), 1, 25, 13, 13, []string{}, []ApiTaskTypes{}, []string{})
				assert.Equal(t, tc.expectedResponse.WholeTaskList, taskList.WholeTaskList, int)
				assert.Equal(t, tc.expectedError, err)
				return nil
			}))
		})
	}
}

func setUpPagesTests(pageCurrent int, lastPage int) (TaskList, TaskDetails) {

	ListOfPages := MakeListOfPagesRange(1, lastPage)

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

	assert.Equal(t, GetPaginationLimits(taskList, taskDetails), []int{1, 2, 3, 4, 5})
}

func TestGetPaginationLimitsWillReturnARangeOnlyTwoAboveCurrentPage(t *testing.T) {
	taskList, taskDetails := setUpPagesTests(1, 10)

	assert.Equal(t, GetPaginationLimits(taskList, taskDetails), []int{1, 2, 3})
}

func TestGetPaginationLimitsWillReturnARangeOneBelowAndTwoAboveCurrentPage(t *testing.T) {
	taskList, taskDetails := setUpPagesTests(2, 10)

	assert.Equal(t, GetPaginationLimits(taskList, taskDetails), []int{1, 2, 3, 4})
}

func TestGetPaginationLimitsWillReturnARangeTwoBelowAndOneAboveCurrentPage(t *testing.T) {
	taskList, taskDetails := setUpPagesTests(4, 5)

	assert.Equal(t, GetPaginationLimits(taskList, taskDetails), []int{2, 3, 4, 5})
}

func TestGetPaginationLimitsWillReturnARangeTwoBelowAndCurrentPage(t *testing.T) {
	taskList, taskDetails := setUpPagesTests(5, 5)

	assert.Equal(t, GetPaginationLimits(taskList, taskDetails), []int{3, 4, 5})
}

func TestCreateTaskTypeFilter(t *testing.T) {
	assert.Equal(t, CreateTaskTypeFilter([]string{}, ""), ",")
	assert.Equal(t, CreateTaskTypeFilter([]string{"CWGN"}, ""), "type:CWGN,")
	assert.Equal(t, CreateTaskTypeFilter([]string{"CWGN", "ORAL"}, ""), "type:CWGN,type:ORAL,")
	assert.Equal(t, CreateTaskTypeFilter([]string{"CWGN", "ORAL", "FAKE", "TEST"}, ""), "type:CWGN,type:ORAL,type:FAKE,type:TEST,")
}

func TestCreateAssigneeFilter(t *testing.T) {
	assert.Equal(t, CreateAssigneeFilter([]string{}, ""), "")
	assert.Equal(t, CreateAssigneeFilter([]string{"LayTeam1"}, ""), "assigneeid:LayTeam1")
	assert.Equal(t, CreateAssigneeFilter([]string{"LayTeam1 User2", "LayTeam1 User3"}, ""), "assigneeid:LayTeam1 User2,assigneeid:LayTeam1 User3")
	assert.Equal(t, CreateAssigneeFilter([]string{"LayTeam1 User3"}, ""), "assigneeid:LayTeam1 User3")
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
		ApiTaskDueDate: "01/06/2021",
		ApiTaskId:      40904862,
		ApiTaskHandle:  ApiTaskHandleInput,
		ApiTaskType:    ApiTaskTypeInput,
		TaskTypeName:   TaskTypeNameInput,
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

func TestGetTaskTypesNameWillReturnIncompleteNameAsTaskTypeName(t *testing.T) {

	taskType := SetUpTaskTypeWithACase("CWGN", "", "", "", 0)
	loadTasks := SetUpLoadTasks()

	assert.Equal(t, GetTaskName(taskType, loadTasks), "Casework - General")
}

func TestGetTaskTypesNameWillReturnOrginalTaskNameIfNoMatchToHandle(t *testing.T) {
	taskType := SetUpTaskTypeWithACase("FAKE", "Fake type", "", "", 0)
	loadTasks := SetUpLoadTasks()

	assert.Equal(t, GetTaskName(taskType, loadTasks), "Fake type")
}

func TestGetTaskTypesNameWillOverwriteAnIncorrectNameWithHandleName(t *testing.T) {
	taskType := SetUpTaskTypeWithACase("CWGN", "Fake name that doesnt match handle", "", "", 0)
	loadTasks := SetUpLoadTasks()
	expectedResult := "Casework - General"

	assert.Equal(t, GetTaskName(taskType, loadTasks), expectedResult)
}

func TestGetAssigneeDisplayNameIfTaskIsAssignedToCaseOwnerWillTakeTheCaseItems(t *testing.T) {
	taskType := SetUpTaskTypeWithACase("", "", "", "Unassigned", 0)
	expectedResult := "Richard Fox"

	assert.Equal(t, GetAssigneeDisplayName(taskType), expectedResult)
}
func TestGetAssigneeDisplayNameIfTaskIsAssignedToCaseOwnerWillTakeTheClients(t *testing.T) {
	taskType := SetUpTaskTypeWithoutACase("", "", "", "Unassigned", 0)
	expectedResult := "Richard Fox"

	assert.Equal(t, GetAssigneeDisplayName(taskType), expectedResult)
}

func TestGetAssigneeDisplayNameIfTaskIsNotAssignedToCaseOwnerWillTakeTheClients(t *testing.T) {
	taskType := SetUpTaskTypeWithoutACase("", "", "", "Go Taskforce", 0)
	expectedResult := "Go Taskforce"

	assert.Equal(t, GetAssigneeDisplayName(taskType), expectedResult)
}

func TestGetAssigneeIdWithOutACase(t *testing.T) {
	taskType := SetUpTaskTypeWithoutACase("", "", "", "Go Taskforce", 0)
	expectedResult := 1234

	assert.Equal(t, GetAssigneeId(taskType), expectedResult)
}

func TestGetAssigneeIdWithACase(t *testing.T) {
	taskType := SetUpTaskTypeWithACase("", "", "", "Go Taskforce", 0)
	expectedResult := 4321

	assert.Equal(t, GetAssigneeId(taskType), expectedResult)
}

func TestGetAssigneeIdWithACaseAndAssignneNotToCaseOwner(t *testing.T) {
	taskType := SetUpTaskTypeWithACase("", "", "", "Go Taskforce", 1122)
	expectedResult := 1122

	assert.Equal(t, GetAssigneeId(taskType), expectedResult)
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

	assert.Equal(t, GetClientInformation(taskType), expectedResult)
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
			Team: []UserTeam{
				{
					Name: "Go TaskForce Team",
					Id:   999,
				},
			},
		},
	}

	assert.Equal(t, GetClientInformation(taskType), expectedResult)
}

func TestGetAssigneeTeamsReturnsOriginalContentIfGivenATeam(t *testing.T) {
	taskType := SetUpUserTeamStruct("Test Team Name", 11)
	expectedResult := []UserTeam{
		{
			Name: "Test Team Name",
			Id:   11,
		},
	}

	assert.Equal(t, GetAssigneeTeams(taskType), expectedResult)
}

func TestGetAssigneeTeamsReplacesContentWithAPIClientsInfoIfNoTeam(t *testing.T) {
	taskType := SetUpTaskTypeWithoutACase("", "", "", "", 0)
	expectedResult := []UserTeam{
		{
			Name: "Go TaskForce Team",
			Id:   999,
		},
	}

	assert.Equal(t, GetAssigneeTeams(taskType), expectedResult)
}

func TestGetAssigneeTeamsReplacesContentWithAPICaseitemsInfoIfNoTeamOrClients(t *testing.T) {
	taskType := SetUpTaskTypeWithoutAClient()
	expectedResult := []UserTeam{
		{
			Name: "Go TaskForce Team",
			Id:   888,
		},
	}

	assert.Equal(t, GetAssigneeTeams(taskType), expectedResult)
}

func TestGetClientInformationPullsInfoFromCaseItemClients(t *testing.T) {
	taskType := SetUpTaskTypeWithACase("", "", "", "", 0)
	expectedResult := Clients{
		ClientCaseRecNumber: "13636617",
		ClientFirstName:     "Pamela",
		ClientId:            37259351,
		ClientSupervisionCaseOwner: CaseManagement{
			Id:              4321,
			CaseManagerName: "Richard Fox",
		},
		ClientSurname: "Pragnell",
	}

	assert.Equal(t, GetClientInformation(taskType), expectedResult)
}

func TestGetClientInformationReturnsInfoIfCaseItemClientsNull(t *testing.T) {
	taskType := SetUpTaskTypeWithoutACase("", "", "", "", 0)
	expectedResult := Clients{
		ClientCaseRecNumber: "13636617",
		ClientFirstName:     "WithoutACase",
		ClientSurname:       "WithoutACase",
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
	}

	assert.Equal(t, GetClientInformation(taskType), expectedResult)
}
