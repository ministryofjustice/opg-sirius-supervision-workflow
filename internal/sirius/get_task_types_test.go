package sirius

import (
	"bytes"
	"errors"
	"github.com/ministryofjustice/opg-sirius-workflow/internal/mocks"
	"github.com/ministryofjustice/opg-sirius-workflow/internal/model"
	"github.com/stretchr/testify/assert"
	"io"
	"net/http"
	"net/http/httptest"
	"strconv"
	"testing"
)

func TestApiClient_GetTaskTypes(t *testing.T) {
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
			Handle:     TaskTypeEcmHandle,
			Incomplete: TaskTypeEcmLabel,
		},
		{
			Handle:     "CWGN",
			Incomplete: "Casework - General",
			Category:   "supervision",
			Complete:   "Casework - General",
			User:       true,
		},
		{
			Handle:     "CNC",
			Incomplete: "Casework - Non-compliant",
			Category:   "supervision",
			Complete:   "Casework - compliant",
			User:       true,
		},
	}

	taskTypes, err := client.GetTaskTypes(getContext(nil), TaskTypesParams{Category: TaskTypeCategorySupervision})

	assert.Equal(t, expectedResponse, taskTypes)
	assert.Equal(t, nil, err)
}

func TestApiClient_GetTaskTypes_Params(t *testing.T) {
	tests := []struct {
		params       TaskTypesParams
		wantEndpoint string
		wantQuery    string
	}{
		{
			params:       TaskTypesParams{Category: "supervision"},
			wantEndpoint: "/api/v1/tasktypes/supervision",
		},
		{
			params:       TaskTypesParams{Category: "deputy"},
			wantEndpoint: "/api/v1/tasktypes/deputy",
		},
		{
			params:       TaskTypesParams{Category: "deputy", ProDeputy: true},
			wantEndpoint: "/api/v1/tasktypes/deputy",
			wantQuery:    "pro_deputy=true",
		},
		{
			params:       TaskTypesParams{Category: "deputy", PADeputy: true},
			wantEndpoint: "/api/v1/tasktypes/deputy",
			wantQuery:    "pa_deputy=true",
		},
	}
	for i, test := range tests {
		t.Run("Scenario "+strconv.Itoa(i), func(t *testing.T) {
			logger, mockClient := SetUpTest()
			client, _ := NewApiClient(mockClient, "http://localhost:3000", logger)

			mocks.GetDoFunc = func(req *http.Request) (*http.Response, error) {
				assert.Equal(t, test.wantEndpoint, req.URL.Path)
				assert.Equal(t, test.wantQuery, req.URL.RawQuery)
				return nil, errors.New("endpoint checked")
			}

			_, err := client.GetTaskTypes(getContext(nil), test.params)
			assert.Equal(t, errors.New("endpoint checked"), err)
		})
	}
}

func TestApiClient_GetTaskTypes_Returns500Error(t *testing.T) {
	logger, _ := SetUpTest()
	svr := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusInternalServerError)
	}))
	defer svr.Close()

	client, _ := NewApiClient(http.DefaultClient, svr.URL, logger)

	_, err := client.GetTaskTypes(getContext(nil), TaskTypesParams{})

	assert.Equal(t, StatusError{
		Code:   http.StatusInternalServerError,
		URL:    svr.URL + "/api/v1/tasktypes/",
		Method: http.MethodGet,
	}, err)
}
