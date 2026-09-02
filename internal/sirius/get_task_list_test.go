package sirius

import (
	"bytes"
	"fmt"
	"io"
	"net/http"
	"net/http/httptest"
	"strconv"
	"testing"
	"time"

	"github.com/ministryofjustice/opg-go-common/telemetry"
	"github.com/ministryofjustice/opg-sirius-workflow/internal/mocks"
	"github.com/ministryofjustice/opg-sirius-workflow/internal/model"
	"github.com/pact-foundation/pact-go/v2/consumer"
	"github.com/pact-foundation/pact-go/v2/matchers"
	"github.com/stretchr/testify/assert"
)

func TestGetTaskListCanReturn200(t *testing.T) {
	logger, mockClient := SetUpTest()
	client := NewApiClient(mockClient, "http://localhost:3000", logger)

	json := `
	{
		"limit":25,
		"metadata":{
			"taskTypeCount": [{"type" : "FCC", "count": 14}],
			"deputyTaskCount": [{"deputy": 61, "count": 3}]
		},
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
		MetaData: TaskMetaData{
			TaskTypeCount: []TypeAndCount{{Type: "FCC", Count: 14}},
			AssigneeCount: []model.AssigneeAndCount(nil),
		},
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
			expectedURL:  "/v1/assignees/teams/tasks?filter=status%3ANot+started&limit=25&page=1&sort=ispriority%3Adesc%2Cduedate%3Aasc%2Cid%3Aasc&teamIds%5B%5D=13",
		},
		{
			name:         "Multiple Team IDs requested",
			selectedTeam: model.Team{Id: 0, Teams: []model.Team{{Id: 12}, {Id: 13}}},
			expectedURL:  "/v1/assignees/teams/tasks?filter=status%3ANot+started&limit=25&page=1&sort=ispriority%3Adesc%2Cduedate%3Aasc%2Cid%3Aasc&teamIds%5B%5D=12&teamIds%5B%5D=13",
		},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			logger, _ := SetUpTest()
			svr := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				w.WriteHeader(http.StatusInternalServerError)
			}))
			defer svr.Close()

			client := NewApiClient(http.DefaultClient, svr.URL, logger)

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
			want:   "status:Not started",
		},
		{
			params: TaskListParams{SelectedTaskTypes: []string{"CWGN"}, TaskTypes: SetUpTaskTypes(), Assignees: []string{"LayTeam1"}},
			want:   "status:Not started,type:CWGN,assigneeid_or_null:LayTeam1",
		},
		{
			params: TaskListParams{SelectedTaskTypes: []string{"CWGN", "ORAL"}, TaskTypes: SetUpTaskTypes(), Assignees: []string{"LayTeam1 User2", "LayTeam1 User3"}},
			want:   "status:Not started,type:CWGN,type:ORAL,assigneeid_or_null:LayTeam1 User2,assigneeid_or_null:LayTeam1 User3",
		},
		{
			params: TaskListParams{SelectedTaskTypes: []string{"CWGN", "ORAL", "FAKE", "TEST"}, TaskTypes: SetUpTaskTypes(), Assignees: []string{"LayTeam1 User3"}},
			want:   "status:Not started,type:CWGN,type:ORAL,type:FAKE,type:TEST,assigneeid_or_null:LayTeam1 User3",
		},
		{
			params: TaskListParams{Assignees: []string{"LayTeam1"}},
			want:   "status:Not started,assigneeid_or_null:LayTeam1",
		},
		{
			params: TaskListParams{Assignees: []string{"LayTeam1"}, DueDateFrom: &selectedDueDateFrom, DueDateTo: &selectedDueDateTo},
			want:   "status:Not started,assigneeid_or_null:LayTeam1,due_date_from:2022-12-17,due_date_to:2022-12-18",
		},
		{
			params: TaskListParams{SelectedTaskTypes: []string{TaskTypeEcmHandle}, TaskTypes: SetUpTaskTypes()},
			want:   "status:Not started,type:CWGN,type:RRRR",
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

func TestGetTaskList_contract(t *testing.T) {
	pact, err := consumer.NewV4Pact(consumer.MockHTTPProviderConfig{
		Consumer: "sirius-supervision-workflow",
		Provider: "sirius",
		LogDir:   "../../logs",
		PactDir:  "../../pacts",
	})
	assert.NoError(t, err)

	err = pact.
		AddInteraction().
		Given("Tasks exist for requested teams").
		UponReceiving("A request for the task list").
		WithRequest("GET", "/supervision-api/v1/assignees/teams/tasks", func(b *consumer.V4RequestBuilder) {
			b.Query("teamIds[]", matchers.S("13"))
			b.Query("filter", matchers.S("status:Not started"))
			b.Query("limit", matchers.S("25"))
			b.Query("page", matchers.S("1"))
			b.Query("sort", matchers.S("ispriority:desc,duedate:asc,id:asc"))
		}).
		WillRespondWith(200, func(b *consumer.V4ResponseBuilder) {
			b.Header("Content-Type", matchers.S("application/json"))
			b.JSONBody(matchers.StructMatcher{
				"limit": matchers.Like(25),
				"metadata": matchers.StructMatcher{
					"taskTypeCount": matchers.EachLike(matchers.StructMatcher{
						"type":  matchers.Like("FCC"),
						"count": matchers.Like(14),
					}, 1),
					"deputyTaskCount": matchers.EachLike(matchers.StructMatcher{
						"deputy": matchers.Like(61),
						"count":  matchers.Like(3),
					}, 1),
				},
				"pages": matchers.StructMatcher{
					"current": matchers.Like(1),
					"total":   matchers.Like(1),
				},
				"total": matchers.Like(13),
				"tasks": matchers.EachLike(matchers.StructMatcher{
					"id":      matchers.Like(119),
					"type":    matchers.Like("ORAL"),
					"status":  matchers.Like("Not started"),
					"dueDate": matchers.Like("29/11/2022"),
					"name":    matchers.Like(""),
					"assignee": matchers.StructMatcher{
						"id":          matchers.Like(0),
						"displayName": matchers.Like("Unassigned"),
					},
					"clients": matchers.EachLike(matchers.StructMatcher{
						"id":            matchers.Like(61),
						"caseRecNumber": matchers.Like("92902877"),
						"firstname":     matchers.Like("Antoine"),
						"surname":       matchers.Like("Burgundy"),
						"supervisionCaseOwner": matchers.StructMatcher{
							"id":          matchers.Like(22),
							"displayName": matchers.Like("Allocations - (Supervision)"),
							"teams":       []interface{}{},
						},
					}, 1),
					"caseOwnerTask": matchers.Like(true),
				}, 1),
			})
		}).
		ExecuteTest(t, func(config consumer.MockServerConfig) error {
			client := NewApiClient(http.DefaultClient, fmt.Sprintf("http://%s:%d/supervision-api", config.Host, config.Port), telemetry.NewLogger("test"))

			taskList, err := client.GetTaskList(getContext(nil), TaskListParams{
				Team:    model.Team{Id: 13},
				Page:    1,
				PerPage: 25,
			})
			assert.NoError(t, err)

			assert.EqualValues(t, 13, taskList.TotalTasks)
			assert.EqualValues(t, 1, len(taskList.Tasks))
			assert.EqualValues(t, 119, taskList.Tasks[0].Id)
			return nil
		})

	assert.NoError(t, err)
}
