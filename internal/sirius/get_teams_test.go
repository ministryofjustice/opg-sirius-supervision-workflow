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

func TestGetTeams(t *testing.T) {
	logger, mockClient := SetUpTest()
	client := NewApiClient(mockClient, "http://localhost:3000", logger)

	json := `[
		{
			"id":21,
			"displayName":"Allocations - (Supervision)",
			"members":[
				{
					"id":71,
					"displayName":"Allocations User1"
				}
			],
			"teamType":{
				"handle":"ALLOCATIONS",
				"label":"Allocations"
			}
		},
		{
			"id":22,
			"displayName":"Lay Team 1",
			"members":[],
			"teamType":{
				"handle":"LAY",
				"label":"Lay Team"
			}
		},
		{
			"id":23,
			"displayName":"Pro Team 1",
			"members":[],
			"teamType":{
				"handle":"PRO",
				"label":"Pro Team"
			}
		},
		{
			"id":24,
			"displayName":"LPA Team",
			"members":[
				{
					"id":72,
					"displayName":"LPA User1"
				}
			]
		}
	]`

	r := io.NopCloser(bytes.NewReader([]byte(json)))

	mocks.GetDoFunc = func(*http.Request) (*http.Response, error) {
		return &http.Response{
			StatusCode: 200,
			Body:       r,
		}, nil
	}

	expectedResponse := []model.Team{
		{
			Id:        21,
			Name:      "Allocations - (Supervision)",
			Type:      "ALLOCATIONS",
			TypeLabel: "Allocations",
			Selector:  "21",
			Members: []model.Assignee{
				{
					Id:   71,
					Name: "Allocations User1",
				},
			},
			Teams: []model.Team{},
		},
		{
			Name:     "Lay Deputy Team",
			Selector: "lay-team",
			Members:  []model.Assignee{},
			Teams: []model.Team{
				{
					Id:        22,
					Name:      "Lay Team 1",
					Type:      "LAY",
					TypeLabel: "Lay Team",
					Selector:  "22",
					Teams:     []model.Team{},
				},
			},
		},
		{
			Id:        22,
			Name:      "Lay Team 1",
			Type:      "LAY",
			TypeLabel: "Lay Team",
			Selector:  "22",
			Teams:     []model.Team{},
		},
		{
			Id:        23,
			Name:      "Pro Team 1",
			Type:      "PRO",
			TypeLabel: "Pro Team",
			Selector:  "23",
			Teams:     []model.Team{},
		},
		{
			Name:     "Professional Deputy Team",
			Selector: "pro-team",
			Members:  []model.Assignee{},
			Teams: []model.Team{
				{
					Id:        23,
					Name:      "Pro Team 1",
					Type:      "PRO",
					TypeLabel: "Pro Team",
					Selector:  "23",
					Teams:     []model.Team{},
				},
			},
		},
	}

	teams, err := client.GetTeams(getContext(nil))
	assert.Equal(t, expectedResponse, teams)
	assert.Equal(t, nil, err)
}

func TestGetTeamsCanReturn500(t *testing.T) {
	logger, _ := SetUpTest()
	svr := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusInternalServerError)
	}))
	defer svr.Close()

	client := NewApiClient(http.DefaultClient, svr.URL, logger)

	_, err := client.GetTeams(getContext(nil))

	assert.Equal(t, StatusError{
		Code:   http.StatusInternalServerError,
		URL:    svr.URL + "/v1/teams",
		Method: http.MethodGet,
	}, err)
}

func TestGetTeams_CachesResponse(t *testing.T) {
	logger, mockClient := SetUpTest()
	client := NewApiClient(mockClient, "http://localhost:3000", logger)

	requests := 0
	mocks.GetDoFunc = func(*http.Request) (*http.Response, error) {
		requests++
		body := io.NopCloser(bytes.NewReader([]byte(`[
			{
				"id":22,
				"displayName":"Lay Team 1",
				"members":[],
				"teamType":{
					"handle":"LAY",
					"label":"Lay Team"
				}
			}
		]`)))
		return &http.Response{
			StatusCode: http.StatusOK,
			Body:       body,
		}, nil
	}

	first, err := client.GetTeams(getContext(nil))
	assert.Nil(t, err)

	second, err := client.GetTeams(getContext(nil))
	assert.Nil(t, err)

	assert.Equal(t, first, second)
	assert.Equal(t, 1, requests)
}

func TestGetTeams_contract(t *testing.T) {
	pact, err := consumer.NewV4Pact(consumer.MockHTTPProviderConfig{
		Consumer: "sirius-supervision-workflow",
		Provider: "sirius",
		LogDir:   "../../../logs",
		PactDir:  "../../../pacts",
	})
	assert.NoError(t, err)

	err = pact.
		AddInteraction().
		Given("Teams exist").
		UponReceiving("A request for teams").
		WithRequest("GET", "/supervision-api/v1/teams", func(b *consumer.V4RequestBuilder) {
			b.Header("Accept", matchers.S("application/json"))
		}).
		WillRespondWith(200, func(b *consumer.V4ResponseBuilder) {
			b.Header("Content-Type", matchers.S("application/json"))
			b.JSONBody(matchers.EachLike(matchers.StructMatcher{
				"id":          matchers.Like(21),
				"displayName": matchers.Like("Allocations - (Supervision)"),
				"members": matchers.EachLike(matchers.StructMatcher{
					"id":          matchers.Like(71),
					"displayName": matchers.Like("Allocations User1"),
				}, 1),
				"teamType": matchers.StructMatcher{
					"handle": matchers.Like("ALLOCATIONS"),
					"label":  matchers.Like("Allocations"),
				},
			}, 1))
		}).
		ExecuteTest(t, func(config consumer.MockServerConfig) error {
			client := NewApiClient(http.DefaultClient, fmt.Sprintf("http://%s:%d/supervision-api", config.Host, config.Port), telemetry.NewLogger("test"))

			teams, err := client.GetTeams(getContext(nil))
			assert.NoError(t, err)

			assert.EqualValues(t, 3, len(teams))
			assert.EqualValues(t, "Allocations - (Supervision)", teams[0].Name)
			return nil
		})

	assert.NoError(t, err)
}
