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

func TestUpdateReassignClients(t *testing.T) {
	jsonResponse := `{"successful":[63],"error":[],"reassignName":"LayTeam1 User2"}`

	tests := []struct {
		params             ReassignClientsParams
		wantAssigneeId     int
		wantSuccessMessage string
	}{
		{
			params:             ReassignClientsParams{AssignTeam: "10", ClientIds: []string{"1", "2"}},
			wantAssigneeId:     10,
			wantSuccessMessage: "You have reassigned 2 client(s) to LayTeam1 User2",
		},
		{
			params:             ReassignClientsParams{AssignTeam: "10", AssignCM: "20", ClientIds: []string{"1"}},
			wantAssigneeId:     20,
			wantSuccessMessage: "You have reassigned 1 client(s) to LayTeam1 User2",
		},
	}
	for i, test := range tests {
		t.Run("Scenario "+strconv.Itoa(i), func(t *testing.T) {
			logger, mockClient := SetUpTest()
			client, _ := NewApiClient(mockClient, "http://localhost:3000", logger)

			r := io.NopCloser(bytes.NewReader([]byte(jsonResponse)))

			mocks.GetDoFunc = func(rq *http.Request) (*http.Response, error) {
				var params ReassignClientsParams
				err := json.NewDecoder(rq.Body).Decode(&params)
				assert.Nil(t, err)
				assert.Equal(t, test.wantAssigneeId, params.AssigneeId)
				assert.True(t, params.IsWorkflow)
				assert.Equal(t, test.params.ClientIds, params.ClientIds)

				return &http.Response{
					StatusCode: 200,
					Body:       r,
				}, nil
			}

			successMessage, err := client.ReassignClients(getContext(nil), test.params)
			assert.Equal(t, test.wantSuccessMessage, successMessage)
			assert.Equal(t, nil, err)
		})
	}
}

func TestReassignClientsReturnsNewStatusError(t *testing.T) {
	logger, _ := SetUpTest()
	svr := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusMethodNotAllowed)
	}))
	defer svr.Close()

	client, _ := NewApiClient(http.DefaultClient, svr.URL, logger)

	_, err := client.ReassignClients(getContext(nil), ReassignClientsParams{AssignTeam: "10"})

	assert.Equal(t, StatusError{
		Code:   http.StatusMethodNotAllowed,
		URL:    svr.URL + "/api/v1/clients/edit/reassign",
		Method: http.MethodPut,
	}, err)
}

func TestReassignClientsReturnsUnauthorisedClientError(t *testing.T) {
	logger, _ := SetUpTest()
	svr := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusUnauthorized)
	}))
	defer svr.Close()

	client, _ := NewApiClient(http.DefaultClient, svr.URL, logger)
	_, err := client.ReassignClients(getContext(nil), ReassignClientsParams{AssignTeam: "10"})
	assert.Equal(t, ErrUnauthorized, err)
}

func TestReassignClientsReturnsForbiddenClientError(t *testing.T) {
	logger, _ := SetUpTest()
	svr := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusForbidden)
	}))
	defer svr.Close()

	client, _ := NewApiClient(http.DefaultClient, svr.URL, logger)
	_, err := client.ReassignClients(getContext(nil), ReassignClientsParams{AssignTeam: "10"})
	assert.Equal(t, "Only managers can reassign client cases", err.Error())
}

func TestReassignClientsReturnsInternalServerError(t *testing.T) {
	logger, _ := SetUpTest()

	svr := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusInternalServerError)
	}))
	defer svr.Close()

	client, _ := NewApiClient(http.DefaultClient, svr.URL, logger)
	_, err := client.ReassignClients(getContext(nil), ReassignClientsParams{AssignTeam: "10"})

	expectedResponse := StatusError{
		Code:   http.StatusInternalServerError,
		URL:    svr.URL + "/api/v1/clients/edit/reassign",
		Method: http.MethodPut,
	}

	assert.Equal(t, expectedResponse, err)
}
