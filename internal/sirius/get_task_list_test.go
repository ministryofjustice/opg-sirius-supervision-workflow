package sirius

import (
	"bytes"
	"github.com/ministryofjustice/opg-sirius-workflow/internal/mocks"
	"github.com/stretchr/testify/assert"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"
)

func TestGetTaskListCanReturn200(t *testing.T) {
	logger, mockClient := SetUpTest()
	client, _ := NewClient(mockClient, "http://localhost:3000", logger)

	json := `
	{
		"limit":25,
		"metadata":{"taskTypeCount": [{"type" : "FCC", "count": 14}]},
		"pages":{"current":1,"total":1},
		"total":13,
		"tasks":[
		{
			"id":119,
			"type":"ORAL",
			"status":"Not started",
			"dueDate":"29\/11\/2022",
			"name":"",
			"description":"A client has been created",
			"ragRating":1,
			"assignee":{"id":0,"displayName":"Unassigned"},
			"createdTime":"14\/11\/2022 12:02:01",
			"caseItems":[],
			"persons":[{"id":61,"uId":"7000-0000-1870","caseRecNumber":"92902877","salutation":"Maquis","firstname":"Antoine","middlenames":"","surname":"Burgundy","supervisionCaseOwner":{"id":22,"teams":[],"displayName":"Allocations - (Supervision)"}}],
			"clients":[{"id":61,"uId":"7000-0000-1870","caseRecNumber":"92902877","salutation":"Maquis","firstname":"Antoine","middlenames":"","surname":"Burgundy","supervisionCaseOwner":{"id":22,"teams":[],"displayName":"Allocations - (Supervision)"}}],
			"caseOwnerTask":true
    	}
		]
	}`

	r := io.NopCloser(bytes.NewReader([]byte(json)))

	mocks.GetDoFunc = func(*http.Request) (*http.Response, error) {
		return &http.Response{
			StatusCode: 200,
			Body:       r,
		}, nil
	}

	expectedResponse := TaskList{
		WholeTaskList: []ApiTask{
			{
				ApiTaskAssignee: CaseManagement{
					"Allocations - (Supervision)",
					22,
					[]UserTeam{},
				},
				ApiTaskCaseItems: nil,
				ApiClients:       nil,
				ApiTaskDueDate:   "29/11/2022",
				ApiTaskId:        119,
				ApiTaskHandle:    "ORAL",
				ApiTaskType:      "",
				ApiCaseOwnerTask: false,
				TaskTypeName:     "",
				ClientInformation: Clients{
					ClientId:            61,
					ClientCaseRecNumber: "92902877",
					ClientFirstName:     "Antoine",
					ClientSurname:       "Burgundy",
					ClientSupervisionCaseOwner: CaseManagement{
						CaseManagerName: "Allocations - (Supervision)",
						Id:              22,
						Team:            []UserTeam{},
					},
				},
				CalculatedDueDateColour: "inThePast",
			},
		},
		Pages: PageInformation{
			PageCurrent: 1,
			PageTotal:   1,
		},
		TotalTasks: 13,
		MetaData:   MetaData{[]TypeAndCount{{Type: "FCC", Count: 14}}},
	}

	selectedTeam := ReturnedTeamCollection{Id: 13}

	assigneeTeams, err := client.GetTaskList(getContext(nil), 1, 25, selectedTeam, []string{""}, []ApiTaskTypes{}, []string{""}, nil, nil)

	assert.Equal(t, expectedResponse, assigneeTeams)
	assert.Equal(t, nil, err)
}

func TestGetTaskListCanThrow500Error(t *testing.T) {
	tests := []struct {
		name         string
		selectedTeam ReturnedTeamCollection
		expectedURL  string
	}{
		{
			name:         "Single Team ID requested",
			selectedTeam: ReturnedTeamCollection{Id: 13},
			expectedURL:  "/api/v1/assignees/teams/tasks?teamIds[]=13&filter=status:Not+started&limit=25&page=1&sort=dueDate:asc",
		},
		{
			name:         "Multiple Team IDs requested",
			selectedTeam: ReturnedTeamCollection{Id: 0, Teams: []ReturnedTeamCollection{{Id: 12}, {Id: 13}}},
			expectedURL:  "/api/v1/assignees/teams/tasks?teamIds[]=12&teamIds[]=13&filter=status:Not+started&limit=25&page=1&sort=dueDate:asc",
		},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			logger, _ := SetUpTest()
			svr := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				w.WriteHeader(http.StatusInternalServerError)
			}))
			defer svr.Close()

			client, _ := NewClient(http.DefaultClient, svr.URL, logger)

			assigneeTeams, err := client.GetTaskList(getContext(nil), 1, 25, test.selectedTeam, []string{}, []ApiTaskTypes{}, []string{}, nil, nil)

			expectedResponse := TaskList{
				WholeTaskList: nil,
				Pages:         PageInformation{},
				TotalTasks:    0,
				ActiveFilters: nil,
			}

			assert.Equal(t, expectedResponse, assigneeTeams)

			assert.Equal(t, StatusError{
				Code:   http.StatusInternalServerError,
				URL:    svr.URL + test.expectedURL,
				Method: http.MethodGet,
			}, err)
		})
	}
}

func TestGetPaginationLimitsWillReturnARangeTwoBelowAndTwoAboveCurrentPage(t *testing.T) {
	taskList, pageDetails := setUpPagesTests(3, 10)
	assert.Equal(t, GetPaginationLimits(taskList, pageDetails), []int{1, 2, 3, 4, 5})
}

func TestGetPaginationLimitsWillReturnARangeOnlyTwoAboveCurrentPage(t *testing.T) {
	taskList, pageDetails := setUpPagesTests(1, 10)
	assert.Equal(t, GetPaginationLimits(taskList, pageDetails), []int{1, 2, 3})
}

func TestGetPaginationLimitsWillReturnARangeOneBelowAndTwoAboveCurrentPage(t *testing.T) {
	taskList, pageDetails := setUpPagesTests(2, 10)
	assert.Equal(t, GetPaginationLimits(taskList, pageDetails), []int{1, 2, 3, 4})
}

func TestGetPaginationLimitsWillReturnARangeTwoBelowAndOneAboveCurrentPage(t *testing.T) {
	taskList, pageDetails := setUpPagesTests(4, 5)
	assert.Equal(t, GetPaginationLimits(taskList, pageDetails), []int{2, 3, 4, 5})
}

func TestGetPaginationLimitsWillReturnARangeTwoBelowAndCurrentPage(t *testing.T) {
	taskList, pageDetails := setUpPagesTests(5, 5)
	assert.Equal(t, GetPaginationLimits(taskList, pageDetails), []int{3, 4, 5})
}

func TestCreateFilter(t *testing.T) {
	selectedDueDateFrom := time.Date(2022, 12, 17, 0, 0, 0, 0, time.Local)
	selectedDueDateTo := time.Date(2022, 12, 18, 0, 0, 0, 0, time.Local)

	assert.Equal(t, CreateFilter([]string{}, []string{}, nil, nil, SetUpLoadTasks()), "status:Not+started")
	assert.Equal(t, CreateFilter([]string{"CWGN"}, []string{"LayTeam1"}, nil, nil, SetUpLoadTasks()), "status:Not+started,type:CWGN,assigneeid_or_null:LayTeam1")
	assert.Equal(t, CreateFilter([]string{"CWGN", "ORAL"}, []string{"LayTeam1 User2", "LayTeam1 User3"}, nil, nil, SetUpLoadTasks()), "status:Not+started,type:CWGN,type:ORAL,assigneeid_or_null:LayTeam1 User2,assigneeid_or_null:LayTeam1 User3")
	assert.Equal(t, CreateFilter([]string{"CWGN", "ORAL", "FAKE", "TEST"}, []string{"LayTeam1 User3"}, nil, nil, SetUpLoadTasks()), "status:Not+started,type:CWGN,type:ORAL,type:FAKE,type:TEST,assigneeid_or_null:LayTeam1 User3")
	assert.Equal(t, CreateFilter([]string{}, []string{"LayTeam1"}, nil, nil, SetUpLoadTasks()), "status:Not+started,assigneeid_or_null:LayTeam1")
	assert.Equal(t, CreateFilter([]string{}, []string{"LayTeam1"}, &selectedDueDateFrom, &selectedDueDateTo, SetUpLoadTasks()), "status:Not+started,assigneeid_or_null:LayTeam1,due_date_from:2022-12-17,due_date_to:2022-12-18")
	assert.Equal(t, CreateFilter([]string{"ECM_TASKS"}, []string{}, nil, nil, SetUpLoadTasks()), "status:Not+started,type:CWGN,type:RRRR")
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

func TestGetTaskNameWillReturnIncompleteNameAsTaskTypeName(t *testing.T) {
	taskType := SetUpTaskTypeWithACase("CWGN", "", "", "", 0)
	loadTasks := SetUpLoadTasks()
	assert.Equal(t, GetTaskName(taskType, loadTasks), "Casework - General")
}

func TestGetTaskNameWillReturnOriginalTaskNameIfNoMatchToHandle(t *testing.T) {
	taskType := SetUpTaskTypeWithACase("FAKE", "Fake type", "", "", 0)
	loadTasks := SetUpLoadTasks()
	assert.Equal(t, GetTaskName(taskType, loadTasks), "Fake type")
}

func TestGetTaskNameWillOverwriteAnIncorrectNameWithHandleName(t *testing.T) {
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

func TestGetAssigneeIdWithACaseAndAssigneeNotToCaseOwner(t *testing.T) {
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
			EcmTask:    true,
		},
		{
			Handle:     "ORAL",
			Incomplete: "Order - Allocate to team",
			Complete:   "Order - Allocate to team",
			User:       true,
			Category:   "supervision",
			IsSelected: false,
			EcmTask:    false,
		},
		{
			Handle:     "RRRR",
			Incomplete: "Visit - Review red report",
			Complete:   "Visit - Review red report",
			User:       true,
			Category:   "supervision",
			IsSelected: false,
			EcmTask:    true,
		},
	}
	return loadTasks
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

func makeListOfPagesRange(min, max int) []int {
	a := make([]int, max-min+1)
	for i := range a {
		a[i] = min + i
	}
	return a
}

func setUpPagesTests(pageCurrent int, lastPage int) (TaskList, PageDetails) {

	ListOfPages := makeListOfPagesRange(1, lastPage)

	taskList := TaskList{
		Pages: PageInformation{
			PageCurrent: pageCurrent,
		},
	}
	pageDetails := PageDetails{
		LastPage:    lastPage,
		ListOfPages: ListOfPages,
	}

	return taskList, pageDetails
}

func TestDeputy_GetURL(t *testing.T) {
	tests := []struct {
		name        string
		deputyType  string
		expectedUrl string
	}{
		{
			name:        "Professional deputy URL",
			deputyType:  "PRO",
			expectedUrl: "/supervision/deputies/13",
		},
		{
			name:        "PA deputy URL",
			deputyType:  "PA",
			expectedUrl: "/supervision/deputies/13",
		},
		{
			name:        "Lay deputy URL",
			deputyType:  "LAY",
			expectedUrl: "/supervision/#/deputy-hub/13",
		},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			deputy := Deputy{Id: 13, Type: RefData{Handle: test.deputyType}}
			assert.Equal(t, test.expectedUrl, deputy.GetURL())
		})
	}
}

func TestGetCalculatedDueDateColour(t *testing.T) {
	tests := []struct {
		name           string
		mockToday      string
		dueDate        string
		expectedColour string
	}{
		{
			name:           "Monday next week will be green",
			mockToday:      "06/06/2023",
			dueDate:        "12/06/2023",
			expectedColour: "dueNextWeek",
		},
		{
			name:           "Due date in the past will be red",
			mockToday:      "06/06/2023",
			dueDate:        "05/06/2023",
			expectedColour: "inThePast",
		},
		{
			name:           "Due date tomorrow will return dueTomorrow",
			mockToday:      "06/06/2023",
			dueDate:        "07/06/2023",
			expectedColour: "dueTomorrow",
		},
		{
			name:           "Due date this week but not tomorrow will return amber",
			mockToday:      "06/06/2023",
			dueDate:        "08/06/2023",
			expectedColour: "dueThisWeek",
		},
		{
			name:           "Due date that is not next week but after will return none",
			mockToday:      "06/06/2023",
			dueDate:        "19/06/2023",
			expectedColour: "none",
		},
		{
			name:           "Sunday today due date Monday will return green",
			mockToday:      "11/06/2023",
			dueDate:        "12/06/2023",
			expectedColour: "dueNextWeek",
		},
		{
			name:           "Due date is today will return dueToday",
			mockToday:      "11/06/2023",
			dueDate:        "11/06/2023",
			expectedColour: "dueToday",
		},
		{
			name:           "Due date is today will return dueToday",
			mockToday:      "06/06/2023",
			dueDate:        "12/06/2023",
			expectedColour: "dueNextWeek",
		},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			mockNow := func() time.Time {
				fakeTime, _ := time.Parse("02/01/2006", test.mockToday)
				return fakeTime
			}
			assert.Equal(t, test.expectedColour, GetCalculatedDueDateStatus(test.dueDate, mockNow))
		})
	}
}

func TestFormatTheDate(t *testing.T) {
	mockDate := time.Date(2023, 06, 11, 0, 0, 0, 0, time.Local)
	expectedResponse, _ := time.Parse("02/01/2006", "11/06/2023")
	result := formatDate(mockDate)
	assert.Equal(t, expectedResponse, result)
}
