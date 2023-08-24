package sirius

import (
	"bytes"
	"encoding/json"
	"github.com/ministryofjustice/opg-sirius-workflow/internal/mocks"
	"github.com/stretchr/testify/assert"
	"io"
	"net/http"
	"net/http/httptest"
	"strconv"
	"testing"
)

func TestUpdateReassignTasks(t *testing.T) {
	jsonResponse := `{	
			"id":98,
			"type":"ORAL",
			"status":"Not started",
			"dueDate":"25\/05\/2023",
			"name":"",
			"description":"A client has been created",
			"ragRating":2,
			"assignee":{
				"id":25,
				"teams":[],
				"displayName":"Allocations - (Supervision)"
			},
			"createdTime":"10\/05\/2023 09:27:05",
			"caseItems":[],
			"persons":[
				{
					"id":42,
					"uId":"7000-0000-0856",
					"caseRecNumber":"71115167",
					"salutation":"Mr",
					"firstname":"Luke",
					"middlenames":"",
					"surname":"Crete",
					"supervisionCaseOwner":{
						"id":25,
						"teams":[],
						"displayName":"Allocations - (Supervision)"}
						}
			],
			"clients":[
				{
					"id":42,
					"uId":"7000-0000-0856",
					"caseRecNumber":"71115167",
					"salutation":"Mr",
					"firstname":"Luke",
					"middlenames":"",
					"surname":"Crete",
					"supervisionCaseOwner":
						{
							"id":25,
							"teams":[],
							"displayName":"Allocations - (Supervision)"
						}
				}
			],
			"caseOwnerTask":false,
			"isPriority":true
		}`

	tests := []struct {
		params             ReassignTasksParams
		wantAssigneeId     int
		wantSuccessMessage string
	}{
		{
			params:             ReassignTasksParams{AssignTeam: "0"},
			wantAssigneeId:     0,
			wantSuccessMessage: "",
		},
		{
			params:             ReassignTasksParams{AssignTeam: "10", TaskIds: []string{"1", "2"}},
			wantAssigneeId:     10,
			wantSuccessMessage: "You have assigned 2 task(s) to Allocations - (Supervision)",
		},
		{
			params:             ReassignTasksParams{AssignTeam: "10", AssignCM: "20", TaskIds: []string{"1"}, IsPriority: "true"},
			wantAssigneeId:     20,
			wantSuccessMessage: "You have assigned 1 task(s) to Allocations - (Supervision) as a priority",
		},
		{
			params:             ReassignTasksParams{AssignTeam: "10", TaskIds: []string{"1"}, IsPriority: "false"},
			wantAssigneeId:     10,
			wantSuccessMessage: "You have assigned 1 task(s) to Allocations - (Supervision) and removed priority",
		},
		{
			params:             ReassignTasksParams{AssignTeam: "0", TaskIds: []string{"1"}, IsPriority: "true"},
			wantAssigneeId:     0,
			wantSuccessMessage: "You have assigned 1 task(s) as a priority",
		},
		{
			params:             ReassignTasksParams{AssignTeam: "0", TaskIds: []string{"1"}, IsPriority: "false"},
			wantAssigneeId:     0,
			wantSuccessMessage: "You have removed 1 task(s) as a priority",
		},
	}
	for i, test := range tests {
		t.Run("Scenario "+strconv.Itoa(i), func(t *testing.T) {
			logger, mockClient := SetUpTest()
			client, _ := NewApiClient(mockClient, "http://localhost:3000", logger)

			r := io.NopCloser(bytes.NewReader([]byte(jsonResponse)))

			mocks.GetDoFunc = func(rq *http.Request) (*http.Response, error) {
				var params ReassignTasksParams
				err := json.NewDecoder(rq.Body).Decode(&params)
				assert.Nil(t, err)
				assert.Equal(t, test.wantAssigneeId, params.AssigneeId)
				assert.Equal(t, test.params.IsPriority, params.IsPriority)
				assert.Equal(t, test.params.TaskIds, params.TaskIds)

				return &http.Response{
					StatusCode: 200,
					Body:       r,
				}, nil
			}

			successMessage, err := client.ReassignTasks(getContext(nil), test.params)
			assert.Equal(t, test.wantSuccessMessage, successMessage)
			assert.Equal(t, nil, err)
		})
	}
}

func TestReassignTasksReturnsNewStatusError(t *testing.T) {
	logger, _ := SetUpTest()
	svr := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusMethodNotAllowed)
	}))
	defer svr.Close()

	client, _ := NewApiClient(http.DefaultClient, svr.URL, logger)

	_, err := client.ReassignTasks(getContext(nil), ReassignTasksParams{AssignTeam: "10"})

	assert.Equal(t, StatusError{
		Code:   http.StatusMethodNotAllowed,
		URL:    svr.URL + "/api/v1/reassign-tasks",
		Method: http.MethodPut,
	}, err)
}

func TestReassignTasksReturnsUnauthorisedClientError(t *testing.T) {
	logger, _ := SetUpTest()
	svr := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusUnauthorized)
	}))
	defer svr.Close()

	client, _ := NewApiClient(http.DefaultClient, svr.URL, logger)
	_, err := client.ReassignTasks(getContext(nil), ReassignTasksParams{AssignTeam: "10"})
	assert.Equal(t, ErrUnauthorized, err)
}

func TestReassignTasksReturnsForbiddenClientError(t *testing.T) {
	logger, _ := SetUpTest()
	svr := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusForbidden)
	}))
	defer svr.Close()

	client, _ := NewApiClient(http.DefaultClient, svr.URL, logger)
	_, err := client.ReassignTasks(getContext(nil), ReassignTasksParams{AssignTeam: "10"})
	assert.Equal(t, "Only managers can set priority on tasks", err.Error())
}

func TestReassignTasksReturnsInternalServerError(t *testing.T) {
	logger, _ := SetUpTest()

	svr := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusInternalServerError)
	}))
	defer svr.Close()

	client, _ := NewApiClient(http.DefaultClient, svr.URL, logger)
	_, err := client.ReassignTasks(getContext(nil), ReassignTasksParams{AssignTeam: "10"})

	expectedResponse := StatusError{
		Code:   http.StatusInternalServerError,
		URL:    svr.URL + "/api/v1/reassign-tasks",
		Method: http.MethodPut,
	}

	assert.Equal(t, expectedResponse, err)
}
