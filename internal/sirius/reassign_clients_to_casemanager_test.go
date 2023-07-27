package sirius

import (
	"bytes"
	"github.com/ministryofjustice/opg-sirius-workflow/internal/mocks"
	"github.com/stretchr/testify/assert"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestUpdateReassignClientToCaseManager(t *testing.T) {
	logger, mockClient := SetUpTest()
	client, _ := NewApiClient(mockClient, "http://localhost:3000", logger)

	json := `{"successful":[63],"error":[],"reassignName":"LayTeam1 User2"}`

	r := io.NopCloser(bytes.NewReader([]byte(json)))

	mocks.GetDoFunc = func(*http.Request) (*http.Response, error) {
		return &http.Response{
			StatusCode: 200,
			Body:       r,
		}, nil
	}

	expectedResponse := "LayTeam1 User2"
	assigneeDisplayName, err := client.ReassignClientToCaseManager(getContext(nil), 1, []string{"76"})
	assert.Equal(t, expectedResponse, assigneeDisplayName)
	assert.Equal(t, nil, err)
}

func TestReassignClientToCaseManagerReturnsNewStatusError(t *testing.T) {
	logger, _ := SetUpTest()
	svr := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusMethodNotAllowed)
	}))
	defer svr.Close()

	client, _ := NewApiClient(http.DefaultClient, svr.URL, logger)

	_, err := client.ReassignClientToCaseManager(getContext(nil), 53, []string{"76"})

	assert.Equal(t, StatusError{
		Code:   http.StatusMethodNotAllowed,
		URL:    svr.URL + "/api/v1/clients/edit/reassign",
		Method: http.MethodPut,
	}, err)
}

func TestReassignClientToCaseManagerReturnsUnauthorisedClientError(t *testing.T) {
	logger, _ := SetUpTest()
	svr := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusUnauthorized)
	}))
	defer svr.Close()

	client, _ := NewApiClient(http.DefaultClient, svr.URL, logger)
	_, err := client.ReassignClientToCaseManager(getContext(nil), 53, []string{"76"})
	assert.Equal(t, ErrUnauthorized, err)
}

func TestReassignClientToCaseManagerReturnsForbiddenClientError(t *testing.T) {
	logger, _ := SetUpTest()
	svr := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusForbidden)
	}))
	defer svr.Close()

	client, _ := NewApiClient(http.DefaultClient, svr.URL, logger)
	_, err := client.ReassignClientToCaseManager(getContext(nil), 53, []string{"76"})
	assert.Equal(t, "Only managers can reassign client cases", err.Error())
}

func TestReassignClientToCaseManagerReturnsInternalServerError(t *testing.T) {
	logger, _ := SetUpTest()

	svr := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusInternalServerError)
	}))
	defer svr.Close()

	client, _ := NewApiClient(http.DefaultClient, svr.URL, logger)
	_, err := client.ReassignClientToCaseManager(getContext(nil), 53, []string{"76"})

	expectedResponse := StatusError{
		Code:   http.StatusInternalServerError,
		URL:    svr.URL + "/api/v1/clients/edit/reassign",
		Method: http.MethodPut,
	}

	assert.Equal(t, expectedResponse, err)
}
