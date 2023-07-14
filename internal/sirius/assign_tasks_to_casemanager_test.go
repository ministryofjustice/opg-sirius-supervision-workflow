package sirius

import (
	"bytes"
	"encoding/json"
	"github.com/ministryofjustice/opg-sirius-workflow/internal/mocks"
	"github.com/stretchr/testify/assert"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestUpdateAssignTasksToCaseManager(t *testing.T) {
	logger, mockClient := SetUpTest()
	client, _ := NewApiClient(mockClient, "http://localhost:3000", logger)

	json := `{	
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

	r := io.NopCloser(bytes.NewReader([]byte(json)))

	mocks.GetDoFunc = func(*http.Request) (*http.Response, error) {
		return &http.Response{
			StatusCode: 200,
			Body:       r,
		}, nil
	}

	expectedResponse := "Allocations - (Supervision)"
	assigneeDisplayName, err := client.AssignTasksToCaseManager(getContext(nil), 1, []string{"76"}, "")
	assert.Equal(t, expectedResponse, assigneeDisplayName)
	assert.Equal(t, nil, err)
}

func TestAssignTasksToCaseManagerReturnsNewStatusError(t *testing.T) {
	logger, _ := SetUpTest()
	svr := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusMethodNotAllowed)
	}))
	defer svr.Close()

	client, _ := NewApiClient(http.DefaultClient, svr.URL, logger)

	_, err := client.AssignTasksToCaseManager(getContext(nil), 53, []string{"76"}, "")

	assert.Equal(t, StatusError{
		Code:   http.StatusMethodNotAllowed,
		URL:    svr.URL + "/api/v1/reassign-tasks",
		Method: http.MethodPut,
	}, err)
}

func TestAssignTasksToCaseManagerReturnsUnauthorisedClientError(t *testing.T) {
	logger, _ := SetUpTest()
	svr := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusUnauthorized)
	}))
	defer svr.Close()

	client, _ := NewApiClient(http.DefaultClient, svr.URL, logger)
	_, err := client.AssignTasksToCaseManager(getContext(nil), 53, []string{"76"}, "")
	assert.Equal(t, ErrUnauthorized, err)
}

func TestAssignTasksToCaseManagerReturnsForbiddenClientError(t *testing.T) {
	logger, _ := SetUpTest()
	svr := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusForbidden)
	}))
	defer svr.Close()

	client, _ := NewApiClient(http.DefaultClient, svr.URL, logger)
	_, err := client.AssignTasksToCaseManager(getContext(nil), 53, []string{"76"}, "")
	assert.Equal(t, "Only managers can set priority on tasks", err.Error())
}

func TestAssignTasksToCaseManagerReturnsInternalServerError(t *testing.T) {
	logger, _ := SetUpTest()

	svr := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusInternalServerError)
	}))
	defer svr.Close()

	client, _ := NewApiClient(http.DefaultClient, svr.URL, logger)
	_, err := client.AssignTasksToCaseManager(getContext(nil), 53, []string{"76"}, "")

	expectedResponse := StatusError{
		Code:   http.StatusInternalServerError,
		URL:    svr.URL + "/api/v1/reassign-tasks",
		Method: http.MethodPut,
	}

	assert.Equal(t, expectedResponse, err)
}

func TestAssignTasksToCaseManager_IsPriority(t *testing.T) {
	cases := map[string]*bool{
		"Yes": boolPointer(true),
		"No":  boolPointer(false),
		"":    nil,
	}

	logger, mockClient := SetUpTest()
	client, _ := NewApiClient(mockClient, "http://localhost:3000", logger)

	r := io.NopCloser(bytes.NewReader([]byte{}))

	mocks.GetDoFunc = func(*http.Request) (*http.Response, error) {
		return &http.Response{
			StatusCode: 200,
			Body:       r,
		}, nil
	}

	for s, e := range cases {
		_, _ = client.AssignTasksToCaseManager(getContext(nil), 1, []string{"1"}, s)

		var r ReassignTaskDetails
		_ = json.NewDecoder(mockClient.GetRequestBody()).Decode(&r)
		assert.Equal(t, r.IsPriority, e)
	}
}

func boolPointer(b bool) *bool {
	return &b
}
