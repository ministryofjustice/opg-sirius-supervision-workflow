package sirius

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/http/httptest"
	"strconv"
	"testing"

	"github.com/ministryofjustice/opg-go-common/telemetry"
	"github.com/ministryofjustice/opg-sirius-workflow/internal/mocks"
	"github.com/pact-foundation/pact-go/v2/consumer"
	"github.com/pact-foundation/pact-go/v2/matchers"
	"github.com/stretchr/testify/assert"
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
			client := NewApiClient(mockClient, "http://localhost:3000", logger)

			r := io.NopCloser(bytes.NewReader([]byte(jsonResponse)))

			mocks.GetDoFunc = func(rq *http.Request) (*http.Response, error) {
				var params ReassignDeputiesParams
				err := json.NewDecoder(rq.Body).Decode(&params)
				assert.Nil(t, err)
				assert.Equal(t, test.wantAssigneeId, params.AssigneeId)
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

	client := NewApiClient(http.DefaultClient, svr.URL, logger)

	_, err := client.ReassignDeputies(getContext(nil), ReassignDeputiesParams{AssignTeam: "10"})

	assert.Equal(t, StatusError{
		Code:   http.StatusMethodNotAllowed,
		URL:    svr.URL + "/v1/deputies/reassign",
		Method: http.MethodPut,
	}, err)
}

func TestReassignDeputiesReturnsUnauthorisedClientError(t *testing.T) {
	logger, _ := SetUpTest()
	svr := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusUnauthorized)
	}))
	defer svr.Close()

	client := NewApiClient(http.DefaultClient, svr.URL, logger)
	_, err := client.ReassignDeputies(getContext(nil), ReassignDeputiesParams{AssignTeam: "10"})
	assert.Equal(t, ErrUnauthorized, err)
}

func TestReassignDeputiesReturnsForbiddenClientError(t *testing.T) {
	logger, _ := SetUpTest()
	svr := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusForbidden)
	}))
	defer svr.Close()

	client := NewApiClient(http.DefaultClient, svr.URL, logger)
	_, err := client.ReassignDeputies(getContext(nil), ReassignDeputiesParams{AssignTeam: "10"})
	assert.Equal(t, "only managers can reassign deputy cases", err.Error())
}

func TestReassignDeputiesReturnsInternalServerError(t *testing.T) {
	logger, _ := SetUpTest()

	svr := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusInternalServerError)
	}))
	defer svr.Close()

	client := NewApiClient(http.DefaultClient, svr.URL, logger)
	_, err := client.ReassignDeputies(getContext(nil), ReassignDeputiesParams{AssignTeam: "10"})

	expectedResponse := StatusError{
		Code:   http.StatusInternalServerError,
		URL:    svr.URL + "/v1/deputies/reassign",
		Method: http.MethodPut,
	}

	assert.Equal(t, expectedResponse, err)
}

func TestReassignDeputies_contract(t *testing.T) {
	pact, err := consumer.NewV4Pact(consumer.MockHTTPProviderConfig{
		Consumer: "sirius-supervision-workflow",
		Provider: "sirius",
		LogDir:   "../../../logs",
		PactDir:  "../../../pacts",
	})
	assert.NoError(t, err)

	err = pact.
		AddInteraction().
		Given("Deputies can be reassigned").
		UponReceiving("A request to reassign deputies").
		WithRequest("PUT", "/supervision-api/v1/deputies/reassign", func(b *consumer.V4RequestBuilder) {
			b.Header("Accept", matchers.S("application/json"))
			b.Header("Content-Type", matchers.S("application/json"))
			b.JSONBody(matchers.StructMatcher{
				"AssignTeam": matchers.Like("10"),
				"AssignCM":   matchers.Like(""),
				"assigneeId": matchers.Like(10),
				"deputyIds":  matchers.EachLike("1", 1),
			})
		}).
		WillRespondWith(200, func(b *consumer.V4ResponseBuilder) {
			b.Header("Content-Type", matchers.S("application/json"))
			b.JSONBody(matchers.StructMatcher{
				"successful":   matchers.EachLike(matchers.Like(63), 1),
				"error":        []interface{}{},
				"reassignName": matchers.Like("LayTeam1 User2"),
			})
		}).
		ExecuteTest(t, func(config consumer.MockServerConfig) error {
			client := NewApiClient(http.DefaultClient, fmt.Sprintf("http://%s:%d/supervision-api", config.Host, config.Port), telemetry.NewLogger("test"))

			msg, err := client.ReassignDeputies(getContext(nil), ReassignDeputiesParams{
				AssignTeam: "10",
				DeputyIds:  []string{"1", "2"},
			})
			assert.NoError(t, err)

			assert.EqualValues(t, "You have reassigned 2 deputies(s) to LayTeam1 User2", msg)
			return nil
		})

	assert.NoError(t, err)
}
