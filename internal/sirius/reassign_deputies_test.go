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

func TestUpdateReassignDeputies(t *testing.T) {
	jsonResponse := `{"successful":[63],"error":[],"reassignName":"LayTeam1 User2"}`

	tests := []struct {
		params             ReassignDeputiesParams
		wantAssigneeId     int
		wantSuccessMessage string
	}{
		{
			params:             ReassignDeputiesParams{AssignTeam: "10", DeputyIds: []string{"1", "2"}},
			wantAssigneeId:     10,
			wantSuccessMessage: "You have reassigned 2 deputies(s) to LayTeam1 User2",
		},
		{
			params:             ReassignDeputiesParams{AssignTeam: "10", AssignCM: "20", DeputyIds: []string{"1"}},
			wantAssigneeId:     20,
			wantSuccessMessage: "You have reassigned 1 deputies(s) to LayTeam1 User2",
		},
	}
	for i, test := range tests {
		t.Run("Scenario "+strconv.Itoa(i), func(t *testing.T) {
			logger, mockClient := SetUpTest()
			client, _ := NewApiClient(mockClient, "http://localhost:3000", logger)

			r := io.NopCloser(bytes.NewReader([]byte(jsonResponse)))

			mocks.GetDoFunc = func(rq *http.Request) (*http.Response, error) {
				var params ReassignDeputiesParams
				err := json.NewDecoder(rq.Body).Decode(&params)
				assert.Nil(t, err)
				assert.Equal(t, test.wantAssigneeId, params.AssigneeId)
				assert.True(t, params.IsWorkflow)
				assert.Equal(t, test.params.DeputyIds, params.DeputyIds)

				return &http.Response{
					StatusCode: 200,
					Body:       r,
				}, nil
			}

			successMessage, err := client.ReassignDeputies(getContext(nil), test.params)
			assert.Equal(t, test.wantSuccessMessage, successMessage)
			assert.Equal(t, nil, err)
		})
	}
}

func TestReassignDeputiesReturnsNewStatusError(t *testing.T) {
	logger, _ := SetUpTest()
	svr := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusMethodNotAllowed)
	}))
	defer svr.Close()

	client, _ := NewApiClient(http.DefaultClient, svr.URL, logger)

	_, err := client.ReassignDeputies(getContext(nil), ReassignDeputiesParams{AssignTeam: "10"})

	assert.Equal(t, StatusError{
		Code:   http.StatusMethodNotAllowed,
		URL:    svr.URL + "/api/v1/deputies/reassign",
		Method: http.MethodPut,
	}, err)
}

func TestReassignDeputiesReturnsUnauthorisedClientError(t *testing.T) {
	logger, _ := SetUpTest()
	svr := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusUnauthorized)
	}))
	defer svr.Close()

	client, _ := NewApiClient(http.DefaultClient, svr.URL, logger)
	_, err := client.ReassignDeputies(getContext(nil), ReassignDeputiesParams{AssignTeam: "10"})
	assert.Equal(t, ErrUnauthorized, err)
}

func TestReassignDeputiesReturnsForbiddenClientError(t *testing.T) {
	logger, _ := SetUpTest()
	svr := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusForbidden)
	}))
	defer svr.Close()

	client, _ := NewApiClient(http.DefaultClient, svr.URL, logger)
	_, err := client.ReassignDeputies(getContext(nil), ReassignDeputiesParams{AssignTeam: "10"})
	assert.Equal(t, "Only managers can reassign deputy cases", err.Error())
}

func TestReassignDeputiesReturnsInternalServerError(t *testing.T) {
	logger, _ := SetUpTest()

	svr := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusInternalServerError)
	}))
	defer svr.Close()

	client, _ := NewApiClient(http.DefaultClient, svr.URL, logger)
	_, err := client.ReassignDeputies(getContext(nil), ReassignDeputiesParams{AssignTeam: "10"})

	expectedResponse := StatusError{
		Code:   http.StatusInternalServerError,
		URL:    svr.URL + "/api/v1/deputies/reassign",
		Method: http.MethodPut,
	}

	assert.Equal(t, expectedResponse, err)
}
