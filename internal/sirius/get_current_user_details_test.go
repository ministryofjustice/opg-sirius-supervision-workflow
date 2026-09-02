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

func TestGetCurrentUserDetails(t *testing.T) {
	logger, mockClient := SetUpTest()
	client := NewApiClient(mockClient, "http://localhost:3000", logger)

	json := `{
			   "id":65,
			   "name":"case",
			   "phoneNumber":"12345678",
			   "teams":[{
				  "displayName":"Lay Team 1 - (Supervision)",
				  "id":13
			   }],
			   "displayName":"case manager",
			   "deleted":false,
			   "email":"case.manager@opgtest.com",
			   "firstname":"case",
			   "surname":"manager",
			   "roles":[
				  "Case Manager"
			   ],
			   "locked":false,
			   "suspended":false
			}`

	r := io.NopCloser(bytes.NewReader([]byte(json)))

	mocks.GetDoFunc = func(*http.Request) (*http.Response, error) {
		return &http.Response{
			StatusCode: 200,
			Body:       r,
		}, nil
	}

	expectedResponse := model.Assignee{
		Id:          65,
		PhoneNumber: "12345678",
		Teams: []model.Team{
			{
				Name: "Lay Team 1 - (Supervision)",
				Id:   13,
			},
		},
		Name:      "case manager",
		Deleted:   false,
		Email:     "case.manager@opgtest.com",
		Firstname: "case",
		Surname:   "manager",
		Roles:     []string{"Case Manager"},
		Locked:    false,
		Suspended: false,
	}

	teams, err := client.GetCurrentUserDetails(getContext(nil))
	assert.Equal(t, expectedResponse, teams)
	assert.Equal(t, nil, err)
}

func TestGetCurrentUserDetailsReturnsUnauthorisedClientError(t *testing.T) {
	logger, _ := SetUpTest()
	svr := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusUnauthorized)
	}))
	defer svr.Close()

	client := NewApiClient(http.DefaultClient, svr.URL, logger)
	_, err := client.GetCurrentUserDetails(getContext(nil))
	assert.Equal(t, ErrUnauthorized, err)
}

func TestGetCurrentUserDetailsReturns500Error(t *testing.T) {
	logger, _ := SetUpTest()
	svr := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusInternalServerError)
	}))
	defer svr.Close()

	client := NewApiClient(http.DefaultClient, svr.URL, logger)

	_, err := client.GetCurrentUserDetails(getContext(nil))
	assert.Equal(t, StatusError{
		Code:   http.StatusInternalServerError,
		URL:    svr.URL + "/v1/users/current",
		Method: http.MethodGet,
	}, err)
}

func TestGetCurrentUserDetailsReturns200(t *testing.T) {
	logger, mockClient := SetUpTest()
	client := NewApiClient(mockClient, "http://localhost:3000", logger)

	json := `{
		"id": 55,
		"name": "case",
		"phoneNumber": "12345678",
		"teams": [],
		"displayName": "case manager",
		"deleted": false,
		"email": "case.manager@opgtest.com",
		"firstname": "case",
		"surname": "manager",
		"roles": [
			"OPG User",
			"Case Manager"
		],
		"locked": false,
		"suspended": false
    }`

	r := io.NopCloser(bytes.NewReader([]byte(json)))

	mocks.GetDoFunc = func(*http.Request) (*http.Response, error) {
		return &http.Response{
			StatusCode: 200,
			Body:       r,
		}, nil
	}

	expectedResponse := model.Assignee{
		Id:          55,
		PhoneNumber: "12345678",
		Teams:       []model.Team{},
		Name:        "case manager",
		Deleted:     false,
		Email:       "case.manager@opgtest.com",
		Firstname:   "case",
		Surname:     "manager",
		Roles:       []string{"OPG User", "Case Manager"},
		Locked:      false,
		Suspended:   false,
	}

	user, err := client.GetCurrentUserDetails(getContext(nil))
	assert.Equal(t, err, nil)
	assert.Equal(t, user, expectedResponse)
}

func TestGetCurrentUserDetails_contract(t *testing.T) {
	pact, err := consumer.NewV4Pact(consumer.MockHTTPProviderConfig{
		Consumer: "sirius-supervision-workflow",
		Provider: "sirius",
		LogDir:   "../../logs",
		PactDir:  "../../pacts",
	})
	assert.NoError(t, err)

	err = pact.
		AddInteraction().
		Given("User exists").
		UponReceiving("A request for the current user").
		WithRequest("GET", "/supervision-api/v1/users/current", func(b *consumer.V4RequestBuilder) {
			b.Header("Accept", matchers.S("application/json"))
		}).
		WillRespondWith(200, func(b *consumer.V4ResponseBuilder) {
			b.Header("Content-Type", matchers.S("application/json"))
			// we only use the below subset of fields for user details, so we only include those in the contract;
			// if more fields are needed in the future, they should be added to the test
			b.JSONBody(matchers.MapMatcher{
				"id":    matchers.Like(1),
				"roles": matchers.EachLike("Case Manager", 1),
				"teams": matchers.EachLike(matchers.MapMatcher{"id": matchers.Like(1)}, 1),
			})
		}).
		ExecuteTest(t, func(config consumer.MockServerConfig) error {
			client := NewApiClient(http.DefaultClient, fmt.Sprintf("http://%s:%d/supervision-api", config.Host, config.Port), telemetry.NewLogger("test"))

			user, _ := client.GetCurrentUserDetails(getContext(nil))

			assert.EqualValues(t, model.Assignee{
				Id:    1,
				Roles: []string{"Case Manager"},
				Teams: []model.Team{{Id: 1}},
			}, user)
			return nil
		})

	assert.NoError(t, err)
}
