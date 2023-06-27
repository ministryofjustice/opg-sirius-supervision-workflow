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
)

func TestGetTaskTypes(t *testing.T) {
	logger, mockClient := SetUpTest()
	client, _ := NewApiClient(mockClient, "http://localhost:3000", logger)

	json := `{
		"task_types":{
			"CWGN":{"handle":"CWGN","incomplete":"Casework - General","complete":"Casework - General","user":true,"category":"supervision"},
			"CNC":{"handle":"CNC","incomplete":"Casework - Non-compliant","complete":"Casework - compliant","user":true,"category":"supervision"}
		}
    }`

	r := io.NopCloser(bytes.NewReader([]byte(json)))

	mocks.GetDoFunc = func(*http.Request) (*http.Response, error) {
		return &http.Response{
			StatusCode: 200,
			Body:       r,
		}, nil
	}

	expectedResponse := []model.TaskType{
		{
			Handle:     "ECM_TASKS",
			Incomplete: "ECM Tasks",
			IsSelected: false,
		},
		{
			Handle:     "CWGN",
			Incomplete: "Casework - General",
			Category:   "supervision",
			Complete:   "Casework - General",
			User:       true,
			IsSelected: false,
		},
		{
			Handle:     "CNC",
			Incomplete: "Casework - Non-compliant",
			Category:   "supervision",
			Complete:   "Casework - compliant",
			User:       true,
			IsSelected: false,
		},
	}

	taskTypes, err := client.GetTaskTypes(getContext(nil), []string{""})

	assert.Equal(t, expectedResponse, taskTypes)
	assert.Equal(t, nil, err)

}

func TestGetTaskTypesCanMarkSelected(t *testing.T) {
	logger, mockClient := SetUpTest()
	client, _ := NewApiClient(mockClient, "http://localhost:3000", logger)

	json := `{
		"task_types":{
			"CWGN":{"handle":"CWGN","incomplete":"Casework - General","complete":"Casework - General","user":true,"category":"supervision"},
			"CNC":{"handle":"CNC","incomplete":"Casework - Non-compliant","complete":"Casework - compliant","user":true,"category":"supervision"},
			"FCIC":{"handle": "FCIC","incomplete": "First Contact - Introductory call","complete": "First Contact - Introductory called","user": true,"category": "supervision"}
		}
    }`

	r := io.NopCloser(bytes.NewReader([]byte(json)))

	mocks.GetDoFunc = func(*http.Request) (*http.Response, error) {
		return &http.Response{
			StatusCode: 200,
			Body:       r,
		}, nil
	}

	expectedResponse := []model.TaskType{
		{
			Handle:     "ECM_TASKS",
			Incomplete: "ECM Tasks",
			IsSelected: false,
		},
		{
			Handle:     "CWGN",
			Incomplete: "Casework - General",
			Category:   "supervision",
			Complete:   "Casework - General",
			User:       true,
			IsSelected: true,
		},
		{
			Handle:     "CNC",
			Incomplete: "Casework - Non-compliant",
			Category:   "supervision",
			Complete:   "Casework - compliant",
			User:       true,
			IsSelected: true,
		},
		{
			Handle:     "FCIC",
			Incomplete: "First Contact - Introductory call",
			Category:   "supervision",
			Complete:   "First Contact - Introductory called",
			User:       true,
			IsSelected: false,
		},
	}

	taskTypes, err := client.GetTaskTypes(getContext(nil), []string{"CWGN", "CNC"})
	assert.Equal(t, expectedResponse, taskTypes)
	assert.Equal(t, nil, err)
}

func TestGetTaskTypesCanMarkSelectedForEcmTasks(t *testing.T) {
	logger, mockClient := SetUpTest()
	client, _ := NewApiClient(mockClient, "http://localhost:3000", logger)

	json := `{
		"task_types":{
			"CWGN":{"handle":"CWGN","incomplete":"Casework - General","complete":"Casework - General","user":true,"category":"supervision"},
			"CNC":{"handle":"CNC","incomplete":"Casework - Non-compliant","complete":"Casework - compliant","user":true,"category":"supervision"},
			"FCIC":{"handle": "FCIC","incomplete": "First Contact - Introductory call","complete": "First Contact - Introductory called","user": true,"category": "supervision"}
		}
    }`

	r := io.NopCloser(bytes.NewReader([]byte(json)))

	mocks.GetDoFunc = func(*http.Request) (*http.Response, error) {
		return &http.Response{
			StatusCode: 200,
			Body:       r,
		}, nil
	}

	expectedResponse := []model.TaskType{
		{
			Handle:     "ECM_TASKS",
			Incomplete: "ECM Tasks",
			IsSelected: true,
		},
		{
			Handle:     "CWGN",
			Incomplete: "Casework - General",
			Category:   "supervision",
			Complete:   "Casework - General",
			User:       true,
			IsSelected: false,
		},
		{
			Handle:     "CNC",
			Incomplete: "Casework - Non-compliant",
			Category:   "supervision",
			Complete:   "Casework - compliant",
			User:       true,
			IsSelected: false,
		},
		{
			Handle:     "FCIC",
			Incomplete: "First Contact - Introductory call",
			Category:   "supervision",
			Complete:   "First Contact - Introductory called",
			User:       true,
			IsSelected: false,
		},
	}

	taskTypes, err := client.GetTaskTypes(getContext(nil), []string{"ECM_TASKS"})
	assert.Equal(t, expectedResponse, taskTypes)
	assert.Equal(t, nil, err)
}

func TestGetTaskTypesReturns500Error(t *testing.T) {
	logger, _ := SetUpTest()
	svr := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusInternalServerError)
	}))
	defer svr.Close()

	client, _ := NewApiClient(http.DefaultClient, svr.URL, logger)

	_, err := client.GetTaskTypes(getContext(nil), []string{"CWGN", "CNC"})

	assert.Equal(t, StatusError{
		Code:   http.StatusInternalServerError,
		URL:    svr.URL + "/api/v1/tasktypes/supervision",
		Method: http.MethodGet,
	}, err)
}

func TestIsSelected(t *testing.T) {
	assert.Equal(t, IsSelected("ORAL", []string{"ORAL"}), true)
	assert.Equal(t, IsSelected("CWGN", []string{"CWGN", "ORAL"}), true)
	assert.Equal(t, IsSelected("TEST", []string{"CWGN", "ORAL"}), false)
}
