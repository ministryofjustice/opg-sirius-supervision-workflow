package sirius

import (
	"bytes"
	"github.com/ministryofjustice/opg-sirius-workflow/internal/mocks"
	"github.com/ministryofjustice/opg-sirius-workflow/internal/model"
	"github.com/stretchr/testify/assert"
	"io"
	"net/http"
	"net/http/httptest"
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
		MetaData:   MetaData{[]TypeAndCount{{Type: "FCC", Count: 14}}},
	}

	selectedTeam := model.Team{Id: 13}

	assigneeTeams, err := client.GetTaskList(getContext(nil), 1, 25, selectedTeam, []string{""}, []model.TaskType{}, []string{""}, nil, nil)

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
			expectedURL:  "/api/v1/assignees/teams/tasks?teamIds[]=13&filter=status:Not+started&limit=25&page=1&sort=dueDate:asc",
		},
		{
			name:         "Multiple Team IDs requested",
			selectedTeam: model.Team{Id: 0, Teams: []model.Team{{Id: 12}, {Id: 13}}},
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

			assigneeTeams, err := client.GetTaskList(getContext(nil), 1, 25, test.selectedTeam, []string{}, []model.TaskType{}, []string{}, nil, nil)

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

func SetUpTaskTypes() []model.TaskType {
	return []model.TaskType{
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
}
