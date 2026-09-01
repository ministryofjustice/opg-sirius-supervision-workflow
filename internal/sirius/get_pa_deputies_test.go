package sirius

import (
	"bytes"
	"fmt"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/ministryofjustice/opg-go-common/telemetry"
	"github.com/ministryofjustice/opg-sirius-workflow/internal/mocks"
	"github.com/ministryofjustice/opg-sirius-workflow/internal/model"
	"github.com/pact-foundation/pact-go/v2/consumer"
	"github.com/pact-foundation/pact-go/v2/matchers"
	"github.com/stretchr/testify/assert"
)

func TestApiClient_GetPADeputies_Returns200(t *testing.T) {
	logger, mockClient := SetUpTest()
	client := NewApiClient(mockClient, "http://localhost:3000", logger)

	json := `
	[
		{
			"id": 13,
			"displayName": "Plompton County Council"
		},
		{
			"id": 14,
			"displayName": "Balamory County Council"
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
		{Id: 13, DisplayName: "Plompton County Council"},
		{Id: 14, DisplayName: "Balamory County Council"},
	}, paDeputies)
}

func TestApiClient_GetPADeputies_Returns500(t *testing.T) {
	logger, _ := SetUpTest()
	svr := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusInternalServerError)
	}))
	defer svr.Close()

	client := NewApiClient(http.DefaultClient, svr.URL, logger)

	paDeputies, err := client.GetPADeputies(getContext(nil))

	assert.Nil(t, paDeputies)
	assert.Equal(t, StatusError{
		Code:   http.StatusInternalServerError,
		URL:    svr.URL + "/v1/assignees/pa-deputies",
		Method: http.MethodGet,
	}, err)
}

func TestGetPADeputies_contract(t *testing.T) {
	pact, err := consumer.NewV4Pact(consumer.MockHTTPProviderConfig{
		Consumer: "sirius-supervision-workflow",
		Provider: "sirius",
		LogDir:   "../../../logs",
		PactDir:  "../../../pacts",
	})
	assert.NoError(t, err)

	err = pact.
		AddInteraction().
		Given("Public authority deputies exist").
		UponReceiving("A request for PA deputies").
		WithRequest("GET", "/supervision-api/v1/assignees/pa-deputies", func(b *consumer.V4RequestBuilder) {
			b.Header("Accept", matchers.S("application/json"))
		}).
		WillRespondWith(200, func(b *consumer.V4ResponseBuilder) {
			b.Header("Content-Type", matchers.S("application/json"))
			b.JSONBody([]interface{}{
				matchers.StructMatcher{
					"id":          matchers.Like(13),
					"displayName": matchers.Like("Plompton County Council"),
				},
				matchers.StructMatcher{
					"id":          matchers.Like(14),
					"displayName": matchers.Like("Balamory County Council"),
				},
			})
		}).
		ExecuteTest(t, func(config consumer.MockServerConfig) error {
			client := NewApiClient(http.DefaultClient, fmt.Sprintf("http://%s:%d/supervision-api", config.Host, config.Port), telemetry.NewLogger("test"))

			deputies, err := client.GetPADeputies(getContext(nil))
			assert.NoError(t, err)

			assert.EqualValues(t, []model.Deputy{
				{Id: 13, DisplayName: "Plompton County Council"},
				{Id: 14, DisplayName: "Balamory County Council"},
			}, deputies)
			return nil
		})

	assert.NoError(t, err)
}
