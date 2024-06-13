package sirius

import (
	"bytes"
	"github.com/ministryofjustice/opg-sirius-workflow/internal/mocks"
	"github.com/ministryofjustice/opg-sirius-workflow/internal/model"
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
		Tasks: []model.Task{
			{
				Assignee: model.Assignee{
					Name: "Unassigned",
				},
				Orders: []model.Order{},
				Clients: []model.Client{
					{
						Id:            61,
						CaseRecNumber: "92902877",
						FirstName:     "Antoine",
						Surname:       "Burgundy",
						SupervisionCaseOwner: model.Assignee{
							Name:  "Allocations - (Supervision)",
							Id:    22,
							Teams: []model.Team{},
						},
						FeePayer: model.Deputy{},
					},
				},
				DueDate:       "29/11/2022",
				Id:            119,
				Type:          "ORAL",
				Name:          "",
				CaseOwnerTask: true,
			},
		},
		Pages: model.PageInformation{
			PageCurrent: 1,
			PageTotal:   1,
		},
		TotalTasks: 13,
		MetaData:   TaskMetaData{[]TypeAndCount{{Type: "FCC", Count: 14}}, []model.AssigneeAndCount(nil)},
	}

	selectedTeam := model.Team{Id: 13}

	assigneeTeams, err := client.GetTaskList(getContext(nil), TaskListParams{
		Team:    selectedTeam,
		Page:    1,
		PerPage: 25,
	})

	assert.Equal(t, expectedResponse, assigneeTeams)
	assert.Equal(t, nil, err)
}

func TestGetTaskListCanThrow500Error(t *testing.T) {
	tests := []struct {
		name         string
		selectedTeam model.Team
		expectedURL  string
	}{
		{
			name:         "Single Team ID requested",
			selectedTeam: model.Team{Id: 13},
			expectedURL:  "/api/v1/assignees/teams/tasks?teamIds[]=13&filter=status:Not+started&limit=25&page=1&sort=ispriority:desc,duedate:asc,id:asc",
		},
		{
			name:         "Multiple Team IDs requested",
			selectedTeam: model.Team{Id: 0, Teams: []model.Team{{Id: 12}, {Id: 13}}},
			expectedURL:  "/api/v1/assignees/teams/tasks?teamIds[]=12&teamIds[]=13&filter=status:Not+started&limit=25&page=1&sort=ispriority:desc,duedate:asc,id:asc",
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

			assigneeTeams, err := client.GetTaskList(getContext(nil), TaskListParams{
				Team:    test.selectedTeam,
				Page:    1,
				PerPage: 25,
			})

			expectedResponse := TaskList{
				Tasks:      nil,
				Pages:      model.PageInformation{},
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

func TestTaskListParams_CreateFilter(t *testing.T) {
	selectedDueDateFrom := time.Date(2022, 12, 17, 0, 0, 0, 0, time.Local)
	selectedDueDateTo := time.Date(2022, 12, 18, 0, 0, 0, 0, time.Local)

	tests := []struct {
		params TaskListParams
		want   string
	}{
		{
			params: TaskListParams{},
			want:   "status:Not+started",
		},
		{
			params: TaskListParams{SelectedTaskTypes: []string{"CWGN"}, TaskTypes: SetUpTaskTypes(), Assignees: []string{"LayTeam1"}},
			want:   "status:Not+started,type:CWGN,assigneeid_or_null:LayTeam1",
		},
		{
			params: TaskListParams{SelectedTaskTypes: []string{"CWGN", "ORAL"}, TaskTypes: SetUpTaskTypes(), Assignees: []string{"LayTeam1 User2", "LayTeam1 User3"}},
			want:   "status:Not+started,type:CWGN,type:ORAL,assigneeid_or_null:LayTeam1 User2,assigneeid_or_null:LayTeam1 User3",
		},
		{
			params: TaskListParams{SelectedTaskTypes: []string{"CWGN", "ORAL", "FAKE", "TEST"}, TaskTypes: SetUpTaskTypes(), Assignees: []string{"LayTeam1 User3"}},
			want:   "status:Not+started,type:CWGN,type:ORAL,type:FAKE,type:TEST,assigneeid_or_null:LayTeam1 User3",
		},
		{
			params: TaskListParams{Assignees: []string{"LayTeam1"}},
			want:   "status:Not+started,assigneeid_or_null:LayTeam1",
		},
		{
			params: TaskListParams{Assignees: []string{"LayTeam1"}, DueDateFrom: &selectedDueDateFrom, DueDateTo: &selectedDueDateTo},
			want:   "status:Not+started,assigneeid_or_null:LayTeam1,due_date_from:2022-12-17,due_date_to:2022-12-18",
		},
		{
			params: TaskListParams{SelectedTaskTypes: []string{TaskTypeEcmHandle}, TaskTypes: SetUpTaskTypes()},
			want:   "status:Not+started,type:CWGN,type:RRRR",
		},
	}
	for i, test := range tests {
		t.Run("Scenario "+strconv.Itoa(i+1), func(t *testing.T) {
			assert.Equal(t, test.want, test.params.CreateFilter())
		})
	}
}

func SetUpTaskTypes() []model.TaskType {
	return []model.TaskType{
		{
			Handle:     "CWGN",
			Incomplete: "Casework - General",
			Complete:   "Casework - General",
			User:       true,
			Category:   "supervision",
			EcmTask:    true,
		},
		{
			Handle:     "ORAL",
			Incomplete: "Order - Allocate to team",
			Complete:   "Order - Allocate to team",
			User:       true,
			Category:   "supervision",
			EcmTask:    false,
		},
		{
			Handle:     "RRRR",
			Incomplete: "Visit - Review red report",
			Complete:   "Visit - Review red report",
			User:       true,
			Category:   "supervision",
			EcmTask:    true,
		},
	}
}

func TestTaskList_CalculateTaskTypeCounts(t *testing.T) {
	taskTypes := []model.TaskType{
		{
			Handle: TaskTypeEcmHandle,
		},
		{
			Handle:  "CDFC",
			EcmTask: false,
		},
		{
			Handle:  "NONO",
			EcmTask: false,
		},
		{
			Handle:  "ECM_1",
			EcmTask: true,
		},
		{
			Handle:  "ECM_2",
			EcmTask: true,
		},
	}
	tasks := TaskList{
		MetaData: TaskMetaData{
			TaskTypeCount: []TypeAndCount{
				{Type: "CDFC", Count: 25},
				{Type: "ECM_1", Count: 33},
				{Type: "ECM_2", Count: 44},
			},
		},
	}

	expected := []model.TaskType{
		{
			Handle:    TaskTypeEcmHandle,
			TaskCount: 77,
		}, {
			Handle:    "CDFC",
			TaskCount: 25,
		},
		{
			Handle:    "NONO",
			TaskCount: 0,
		},
		{
			Handle:    "ECM_1",
			EcmTask:   true,
			TaskCount: 33,
		},
		{
			Handle:    "ECM_2",
			EcmTask:   true,
			TaskCount: 44,
		},
	}

	assert.Equal(t, expected, tasks.CalculateTaskTypeCounts(taskTypes))
}
