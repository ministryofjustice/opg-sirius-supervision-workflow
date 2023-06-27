package sirius

import (
	"bytes"
	"github.com/ministryofjustice/opg-sirius-workflow/internal/mocks"
	"github.com/stretchr/testify/assert"
	"io"
	"net/http"
	"net/http/httptest"
	"strconv"
	"testing"
	"time"
)

func TestGetTaskListCanReturn200(t *testing.T) {
	logger, mockClient := SetUpTest()
	client, _ := NewApiClient(mockClient, "http://localhost:3000", logger)

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
		Tasks: []Task{
			{
				Assignee: Assignee{
					Name: "Unassigned",
				},
				CaseItems: []CaseItem{},
				Clients: []Client{
					{
						Id:            61,
						CaseRecNumber: "92902877",
						FirstName:     "Antoine",
						Surname:       "Burgundy",
						SupervisionCaseOwner: Assignee{
							Name:  "Allocations - (Supervision)",
							Id:    22,
							Teams: []UserTeam{},
						},
						FeePayer: Deputy{},
					},
				},
				DueDate:       "29/11/2022",
				Id:            119,
				Type:          "ORAL",
				Name:          "",
				CaseOwnerTask: true,
			},
		},
		Pages: PageInformation{
			PageCurrent: 1,
			PageTotal:   1,
		},
		TotalTasks: 13,
		MetaData:   MetaData{[]TypeAndCount{{Type: "FCC", Count: 14}}},
	}

	selectedTeam := Team{Id: 13}

	assigneeTeams, err := client.GetTaskList(getContext(nil), 1, 25, selectedTeam, []string{""}, []TaskType{}, []string{""}, nil, nil)

	assert.Equal(t, expectedResponse, assigneeTeams)
	assert.Equal(t, nil, err)
}

func TestGetTaskListCanThrow500Error(t *testing.T) {
	tests := []struct {
		name         string
		selectedTeam Team
		expectedURL  string
	}{
		{
			name:         "Single Team ID requested",
			selectedTeam: Team{Id: 13},
			expectedURL:  "/api/v1/assignees/teams/tasks?teamIds[]=13&filter=status:Not+started&limit=25&page=1&sort=dueDate:asc",
		},
		{
			name:         "Multiple Team IDs requested",
			selectedTeam: Team{Id: 0, Teams: []Team{{Id: 12}, {Id: 13}}},
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

			client, _ := NewApiClient(http.DefaultClient, svr.URL, logger)

			assigneeTeams, err := client.GetTaskList(getContext(nil), 1, 25, test.selectedTeam, []string{}, []TaskType{}, []string{}, nil, nil)

			expectedResponse := TaskList{
				Tasks:      nil,
				Pages:      PageInformation{},
				TotalTasks: 0,
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

	assert.Equal(t, CreateFilter([]string{}, []string{}, nil, nil, SetUpTaskTypes()), "status:Not+started")
	assert.Equal(t, CreateFilter([]string{"CWGN"}, []string{"LayTeam1"}, nil, nil, SetUpTaskTypes()), "status:Not+started,type:CWGN,assigneeid_or_null:LayTeam1")
	assert.Equal(t, CreateFilter([]string{"CWGN", "ORAL"}, []string{"LayTeam1 User2", "LayTeam1 User3"}, nil, nil, SetUpTaskTypes()), "status:Not+started,type:CWGN,type:ORAL,assigneeid_or_null:LayTeam1 User2,assigneeid_or_null:LayTeam1 User3")
	assert.Equal(t, CreateFilter([]string{"CWGN", "ORAL", "FAKE", "TEST"}, []string{"LayTeam1 User3"}, nil, nil, SetUpTaskTypes()), "status:Not+started,type:CWGN,type:ORAL,type:FAKE,type:TEST,assigneeid_or_null:LayTeam1 User3")
	assert.Equal(t, CreateFilter([]string{}, []string{"LayTeam1"}, nil, nil, SetUpTaskTypes()), "status:Not+started,assigneeid_or_null:LayTeam1")
	assert.Equal(t, CreateFilter([]string{}, []string{"LayTeam1"}, &selectedDueDateFrom, &selectedDueDateTo, SetUpTaskTypes()), "status:Not+started,assigneeid_or_null:LayTeam1,due_date_from:2022-12-17,due_date_to:2022-12-18")
	assert.Equal(t, CreateFilter([]string{"ECM_TASKS"}, []string{}, nil, nil, SetUpTaskTypes()), "status:Not+started,type:CWGN,type:RRRR")
}

func SetUpTaskWithACase(ApiTaskHandleInput string, ApiTaskTypeInput string, AssigneeDisplayNameInput string, AssigneeIdInput int) Task {
	v := Task{
		Assignee: Assignee{
			Name: AssigneeDisplayNameInput,
			Id:   AssigneeIdInput,
		},
		CaseItems: []CaseItem{{
			Client: Client{
				CaseRecNumber: "13636617",
				FirstName:     "Pamela",
				Id:            37259351,
				SupervisionCaseOwner: Assignee{
					Id:   4321,
					Name: "Richard Fox",
				},
				Surname: "Pragnell",
			},
		}},
		DueDate: "01/06/2021",
		Id:      40904862,
		Type:    ApiTaskHandleInput,
		Name:    ApiTaskTypeInput,
	}
	return v
}

func TestTask_GetName(t *testing.T) {
	taskTypes := SetUpTaskTypes()
	tests := []struct {
		name     string
		task     Task
		wantName string
	}{
		{
			name:     "Incomplete name used as task name",
			task:     SetUpTaskWithACase("CWGN", "", "", 0),
			wantName: "Casework - General",
		},
		{
			name:     "Original task name used when cannot be matched to a task type",
			task:     SetUpTaskWithACase("FAKE", "Fake type", "", 0),
			wantName: "Fake type",
		},
		{
			name:     "Task type name overwrites an incorrect name with a matching handle",
			task:     SetUpTaskWithACase("CWGN", "Fake name that doesnt match handle", "", 0),
			wantName: "Casework - General",
		},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			assert.Equal(t, test.wantName, test.task.GetName(taskTypes))
		})
	}
}

func TestTask_GetAssignee(t *testing.T) {
	tests := []struct {
		name         string
		task         Task
		wantAssignee string
	}{
		{
			name:         "Unassigned task gets Assignee from Clients",
			task:         SetUpTaskWithACase("", "", "Unassigned", 0),
			wantAssignee: "Richard Fox",
		},
		{
			name:         "Unassigned task gets Assignee from CaseItems",
			task:         SetUpTaskWithoutACase("Unassigned", 0),
			wantAssignee: "Richard Fox",
		},
		{
			name:         "Assigned task get Assignee from task",
			task:         SetUpTaskWithoutACase("Go Taskforce", 0),
			wantAssignee: "Go Taskforce",
		},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			assert.Equal(t, test.wantAssignee, test.task.GetAssignee().Name)
		})
	}
}

func TestTask_GetClient(t *testing.T) {
	tests := []struct {
		name string
		task Task
		want string
	}{
		{
			name: "Get client from task case",
			task: SetUpTaskWithACase("", "", "Go Taskforce", 1122),
			want: "Pamela",
		},
		{
			name: "Get client from task without cases",
			task: SetUpTaskWithoutACase("Go Taskforce", 1122),
			want: "WithoutACase",
		},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			assert.Equal(t, test.want, test.task.GetClient().FirstName)
		})
	}
}

func SetUpTaskWithoutACase(AssigneeDisplayNameInput string, AssigneeIdInput int) Task {
	v := Task{
		Assignee: Assignee{
			Name: AssigneeDisplayNameInput,
			Id:   AssigneeIdInput,
		},
		Clients: []Client{
			{
				CaseRecNumber: "13636617",
				FirstName:     "WithoutACase",
				Id:            37259351,
				SupervisionCaseOwner: Assignee{
					Id:   1234,
					Name: "Richard Fox",
					Teams: []UserTeam{
						{
							Name: "Go TaskForce Team",
							Id:   999,
						},
					},
				},
				Surname: "WithoutACase",
			},
		},
		DueDate: "01/06/2021",
		Id:      40904862,
	}
	return v
}

func SetUpTaskTypes() []TaskType {
	taskTypes := []TaskType{
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
	return taskTypes
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

func TestTask_GetDueDateStatus(t *testing.T) {
	tests := []struct {
		name      string
		mockToday string
		dueDate   string
		want      DueDateStatus
	}{
		{
			name:      "Due date is in the past",
			mockToday: "06/06/2023",
			dueDate:   "05/06/2023",
			want:      DueDateStatus{"Overdue", "red"},
		},
		{
			name:      "Due date is today",
			mockToday: "11/06/2023",
			dueDate:   "11/06/2023",
			want:      DueDateStatus{"Due Today", "red"},
		},
		{
			name:      "Due date is tomorrow",
			mockToday: "06/06/2023",
			dueDate:   "07/06/2023",
			want:      DueDateStatus{"Due Tomorrow", "orange"},
		},
		{
			name:      "Due date is this week but not tomorrow",
			mockToday: "06/06/2023",
			dueDate:   "08/06/2023",
			want:      DueDateStatus{"Due This Week", "orange"},
		},
		{
			name:      "Due date is Monday next week",
			mockToday: "06/06/2023",
			dueDate:   "12/06/2023",
			want:      DueDateStatus{"Due Next Week", "green"},
		},
		{
			name:      "Due date on same week day as today but next week",
			mockToday: "06/06/2023",
			dueDate:   "13/06/2023",
			want:      DueDateStatus{"Due Next Week", "green"},
		},
		{
			name:      "Sunday today due date Monday",
			mockToday: "11/06/2023",
			dueDate:   "12/06/2023",
			want:      DueDateStatus{"Due Next Week", "green"},
		},
		{
			name:      "Due date next week",
			mockToday: "06/06/2023",
			dueDate:   "12/06/2023",
			want:      DueDateStatus{"Due Next Week", "green"},
		},
		{
			name:      "Due date on same week day as today but in future",
			mockToday: "06/06/2023",
			dueDate:   "23/06/2023",
			want:      DueDateStatus{"", ""},
		},
		{
			name:      "Due date that is not next week but after",
			mockToday: "06/06/2023",
			dueDate:   "19/06/2023",
			want:      DueDateStatus{"", ""},
		},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			task := Task{DueDate: test.dueDate}
			mockNow, _ := time.Parse("02/01/2006", test.mockToday)
			assert.Equal(t, test.want, task.GetDueDateStatus(mockNow))
		})
	}
}

func TestClient_GetStatus(t *testing.T) {
	tests := []struct {
		orderStatuses []string
		wantStatus    string
	}{
		{
			orderStatuses: []string{"Closed", "Open", "Duplicate", "Active", "Closed", "Open", "Duplicate"},
			wantStatus:    "Active",
		},
		{
			orderStatuses: []string{"Open", "Duplicate", "Closed", "Open", "Duplicate"},
			wantStatus:    "Open",
		},
		{
			orderStatuses: []string{"Duplicate", "Closed", "Duplicate"},
			wantStatus:    "Closed",
		},
		{
			orderStatuses: []string{"Duplicate"},
			wantStatus:    "Duplicate",
		},
	}
	for i, test := range tests {
		t.Run("Scenario "+strconv.Itoa(i+1), func(t *testing.T) {
			var client Client
			for _, status := range test.orderStatuses {
				client.Orders = append(client.Orders, Order{Status: RefData{Label: status}})
			}
			assert.Equal(t, test.wantStatus, client.GetStatus())
		})
	}
}

func TestClient_GetReportDueDate(t *testing.T) {
	client := Client{
		Orders: []Order{
			{
				LatestAnnualReport: AnnualReport{
					DueDate: "12/02/2020",
				},
			},
		},
	}
	assert.Equal(t, "12/02/2020", client.GetReportDueDate())
}
