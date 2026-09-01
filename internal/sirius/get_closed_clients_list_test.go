package sirius

import (
	"bytes"
	"fmt"
	"io"
	"net/http"
	"net/http/httptest"
	"strconv"
	"testing"
	"time"

	"github.com/ministryofjustice/opg-go-common/telemetry"
	"github.com/ministryofjustice/opg-sirius-workflow/internal/mocks"
	"github.com/ministryofjustice/opg-sirius-workflow/internal/model"
	"github.com/pact-foundation/pact-go/v2/consumer"
	"github.com/pact-foundation/pact-go/v2/matchers"
	"github.com/stretchr/testify/assert"
)

func TestGetClosedCaseloadListCanReturn200(t *testing.T) {
	logger, mockClient := SetUpTest()
	client := NewApiClient(mockClient, "http://localhost:3000", logger)

	json := `
{
    "limit": 15,
    "metadata":{"assigneeClientCount": [{"assignee" : 1, "count": 14}]},
    "pages": {
        "current": 1,
        "total": 1
    },
    "total": 1,
    "clients": [
        {
            "id": 63,
            "caseRecNumber": "42687883",
            "firstname": "Ro",
            "surname": "Bot",
            "cases": [
                {
                    "id": 92,
                    "caseRecNumber": "33594483",
                    "latestAnnualReport": {
                        "dueDate": "21\/12\/2023"
                    },
					"orderStatus": {
						"handle": "CLOSED",
						"label": "Closed",
						"deprecated": false
					},
					"madeActiveDate": "01\/06\/2023",
					"introductoryTargetDate": "20\/06\/2023",
					"howDeputyAppointed": {
						"handle": "SOLE",
						"label": "Sole"
					}
                }
            ],
            "supervisionLevel": {
				"handle": "MINIMAL",
				"label": "Minimal"
			},
			"cachedDebtTotal": 10010,
			"lastActionDate": "2023-12-12T12:35:56+00:00",
			"closedOnDate": "2022-02-02T12:35:56+00:00"
        }
    ]
}
`

	r := io.NopCloser(bytes.NewReader([]byte(json)))

	mocks.GetDoFunc = func(rq *http.Request) (*http.Response, error) {
		return &http.Response{
			StatusCode: 200,
			Body:       r,
		}, nil
	}

	lastActionDate := time.Date(2023, time.Month(12), 12, 12, 35, 56, 0, time.UTC)
	closedOnDate := time.Date(2022, time.Month(2), 2, 12, 35, 56, 0, time.UTC)

	expectedResponse := ClientList{
		Clients: []model.Client{
			{
				Id:            63,
				CaseRecNumber: "42687883",
				FirstName:     "Ro",
				Surname:       "Bot",
				Orders: []model.Order{
					{
						Id: 92,
						Status: model.RefData{
							Handle: "CLOSED",
							Label:  "Closed",
						},
						LatestAnnualReport: model.AnnualReport{
							DueDate: "21/12/2023",
						},
						MadeActiveDate:         model.NewDate("01/06/2023"),
						IntroductoryTargetDate: model.NewDate("20/06/2023"),
						HowDeputyAppointed: model.RefData{
							Handle: "SOLE",
							Label:  "Sole",
						},
					},
				},
				SupervisionLevel: model.RefData{
					Handle: "MINIMAL",
					Label:  "Minimal",
				},
				CachedDebtTotal: 10010,
				LastActionDate:  model.Date{Time: lastActionDate},
				ClosedOnDate:    model.Date{Time: closedOnDate},
			},
		},
		Pages: model.PageInformation{
			PageCurrent: 1,
			PageTotal:   1,
		},
		TotalClients: 1,
		MetaData: ClientMetaData{
			[]model.AssigneeAndCount{
				{AssigneeId: 1, Count: 14},
			},
		},
	}

	clientList, err := client.GetClosedClientList(getContext(nil), ClientListParams{
		Team:    model.Team{Id: 40, Name: "Supervision closed cases"},
		Page:    1,
		PerPage: 25,
	})

	assert.Equal(t, nil, err)
	assert.Equal(t, expectedResponse, clientList)
}

func TestGetClosedCaseloadListCanThrow500Error(t *testing.T) {
	logger, _ := SetUpTest()
	svr := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusInternalServerError)
	}))
	defer svr.Close()

	client := NewApiClient(http.DefaultClient, svr.URL, logger)

	clientList, err := client.GetClosedClientList(getContext(nil), ClientListParams{
		Team:    model.Team{Id: 13},
		Page:    1,
		PerPage: 25,
	})

	expectedResponse := ClientList{
		Clients:      nil,
		Pages:        model.PageInformation{},
		TotalClients: 0,
	}

	assert.Equal(t, expectedResponse, clientList)

	assert.Equal(t, StatusError{
		Code:   http.StatusInternalServerError,
		URL:    svr.URL + "/v1/assignees/closed-clients?limit=25&page=1&filter=",
		Method: http.MethodGet,
	}, err)
}

func TestCreateMemberIdArray(t *testing.T) {
	tests := []struct {
		params ClientListParams
		want   []string
	}{
		{
			params: ClientListParams{
				Team: model.Team{
					Id:   40,
					Name: "Closed Cases Team",
				},
			},
			want: []string{"40"},
		},
		{
			params: ClientListParams{
				Team: model.Team{
					Id:   40,
					Name: "Closed Cases Team",
					Members: []model.Assignee{
						{
							Id:   20,
							Name: "Person 1",
						},
						{
							Id:   21,
							Name: "Person 2",
						},
					},
				},
			},
			want: []string{"40", "20", "21"},
		},
	}
	for i, test := range tests {
		t.Run("Scenario "+strconv.Itoa(i+1), func(t *testing.T) {
			assert.Equal(t, test.want, CreateMemberIdArray(test.params))
		})
	}
}

func TestGetClosedClientList_contract(t *testing.T) {
	pact, err := consumer.NewV4Pact(consumer.MockHTTPProviderConfig{
		Consumer: "sirius-supervision-workflow",
		Provider: "sirius",
		LogDir:   "../../../logs",
		PactDir:  "../../../pacts",
	})
	assert.NoError(t, err)

	err = pact.
		AddInteraction().
		Given("Closed clients exist for requested teams").
		UponReceiving("A request for the closed client list").
		WithRequest("GET", "/supervision-api/v1/assignees/closed-clients", func(b *consumer.V4RequestBuilder) {
			b.Header("Accept", matchers.S("application/json"))
			b.Query("limit", matchers.S("25"))
			b.Query("page", matchers.S("1"))
			b.Query("filter", matchers.S(""))
			b.JSONBody(matchers.StructMatcher{
				"teamIds": matchers.EachLike("40", 1),
			})
		}).
		WillRespondWith(200, func(b *consumer.V4ResponseBuilder) {
			b.Header("Content-Type", matchers.S("application/json"))
			b.JSONBody(matchers.StructMatcher{
				"pages": matchers.StructMatcher{
					"current": matchers.Like(1),
					"total":   matchers.Like(1),
				},
				"total": matchers.Like(1),
				"metadata": matchers.StructMatcher{
					"assigneeClientCount": matchers.EachLike(matchers.StructMatcher{
						"assignee": matchers.Like(1),
						"count":    matchers.Like(14),
					}, 1),
				},
				"clients": matchers.EachLike(matchers.StructMatcher{
					"id":            matchers.Like(63),
					"caseRecNumber": matchers.Like("42687883"),
					"firstname":     matchers.Like("Ro"),
					"surname":       matchers.Like("Bot"),
					"cases": matchers.EachLike(matchers.StructMatcher{
						"id":            matchers.Like(92),
						"caseRecNumber": matchers.Like("33594483"),
						"latestAnnualReport": matchers.StructMatcher{
							"dueDate": matchers.Like("21/12/2023"),
						},
						"orderStatus": matchers.StructMatcher{
							"handle": matchers.Like("CLOSED"),
							"label":  matchers.Like("Closed"),
						},
						"madeActiveDate":         matchers.Like("01/06/2023"),
						"introductoryTargetDate": matchers.Like("20/06/2023"),
						"howDeputyAppointed": matchers.StructMatcher{
							"handle": matchers.Like("SOLE"),
							"label":  matchers.Like("Sole"),
						},
					}, 1),
					"supervisionLevel": matchers.StructMatcher{
						"handle": matchers.Like("MINIMAL"),
						"label":  matchers.Like("Minimal"),
					},
					"cachedDebtTotal": matchers.Like(10010),
					"lastActionDate":  matchers.Like("2023-12-12T12:35:56+00:00"),
					"closedOnDate":    matchers.Like("2022-02-02T12:35:56+00:00"),
				}, 1),
			})
		}).
		ExecuteTest(t, func(config consumer.MockServerConfig) error {
			client := NewApiClient(http.DefaultClient, fmt.Sprintf("http://%s:%d/supervision-api", config.Host, config.Port), telemetry.NewLogger("test"))

			clientList, err := client.GetClosedClientList(getContext(nil), ClientListParams{
				Team:    model.Team{Id: 40, Name: "Supervision closed cases"},
				Page:    1,
				PerPage: 25,
			})
			assert.NoError(t, err)

			assert.EqualValues(t, 1, clientList.TotalClients)
			assert.EqualValues(t, 1, len(clientList.Clients))
			assert.EqualValues(t, "Bot", clientList.Clients[0].Surname)
			return nil
		})

	assert.NoError(t, err)
}
