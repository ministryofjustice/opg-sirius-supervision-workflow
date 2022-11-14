package sirius

import (
	"bytes"
	"github.com/ministryofjustice/opg-sirius-workflow/internal/mocks"
	"github.com/stretchr/testify/assert"
	"io"
	"net/http"
	"testing"
)

func TestGetTaskTypes(t *testing.T) {
	logger, mockClient := SetUpTest()
	client, _ := NewClient(mockClient, "http://localhost:3000", logger)

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

	expectedResponse := []ApiTaskTypes{
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
	client, _ := NewClient(mockClient, "http://localhost:3000", logger)

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

	expectedResponse := []ApiTaskTypes{
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

func TestIsSelected(t *testing.T) {
	assert.Equal(t, IsSelected("ORAL", []string{"ORAL"}), true)
	assert.Equal(t, IsSelected("CWGN", []string{"CWGN", "ORAL"}), true)
	assert.Equal(t, IsSelected("TEST", []string{"CWGN", "ORAL"}), false)
}
