package sirius

import (
	"bytes"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/ministryofjustice/opg-sirius-workflow/internal/mocks"
	"github.com/ministryofjustice/opg-sirius-workflow/internal/model"
	"github.com/stretchr/testify/assert"
)

func TestApiClient_GetPADeputies_Returns200(t *testing.T) {
	logger, mockClient := SetUpTest()
	client, _ := NewApiClient(mockClient, "http://localhost:3000", logger)

	json := `
	[
		{
			"id": 13,
			"displayName": "Mr Fee-paying Deputy"
		},
		{
			"id": 14,
			"displayName": "Ms Another Deputy"
		}
	]`

	r := io.NopCloser(bytes.NewReader([]byte(json)))

	mocks.GetDoFunc = func(rq *http.Request) (*http.Response, error) {
		assert.Equal(t, "/v1/assignees/pa-deputies", rq.URL.Path)
		return &http.Response{
			StatusCode: 200,
			Body:       r,
		}, nil
	}

	paDeputies, err := client.GetPADeputies(getContext(nil))

	assert.Equal(t, nil, err)
	assert.Equal(t, []model.Deputy{
		{Id: 13, DisplayName: "Mr Fee-paying Deputy"},
		{Id: 14, DisplayName: "Ms Another Deputy"},
	}, paDeputies)
}

func TestApiClient_GetPADeputies_Returns500(t *testing.T) {
	logger, _ := SetUpTest()
	svr := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusInternalServerError)
	}))
	defer svr.Close()

	client, _ := NewApiClient(http.DefaultClient, svr.URL, logger)

	paDeputies, err := client.GetPADeputies(getContext(nil))

	assert.Nil(t, paDeputies)
	assert.Equal(t, StatusError{
		Code:   http.StatusInternalServerError,
		URL:    svr.URL + "/v1/assignees/pa-deputies",
		Method: http.MethodGet,
	}, err)
}
