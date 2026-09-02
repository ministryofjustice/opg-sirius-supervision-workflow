package sirius

import (
	"bytes"
	"fmt"
	"io"
	"net/http"
	"net/http/httptest"
	"strconv"
	"testing"

	"github.com/ministryofjustice/opg-go-common/telemetry"
	"github.com/ministryofjustice/opg-sirius-workflow/internal/mocks"
	"github.com/ministryofjustice/opg-sirius-workflow/internal/model"
	"github.com/pact-foundation/pact-go/v2/consumer"
	"github.com/pact-foundation/pact-go/v2/matchers"
	"github.com/stretchr/testify/assert"
)

func TestGetCaseloadListCanReturn200(t *testing.T) {
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
			}
        }
    ]
}
`

	r := io.NopCloser(bytes.NewReader([]byte(json)))

	mocks.GetDoFunc = func(rq *http.Request) (*http.Response, error) {
		assert.NotContains(t, rq.URL.RawQuery, "sort=made_active_date:asc")
		assert.Contains(t, rq.URL.RawQuery, "caseowner:1")
		return &http.Response{
			StatusCode: 200,
			Body:       r,
		}, nil
	}

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

	clientList, err := client.GetClientList(getContext(nil), ClientListParams{
		Team:       model.Team{Id: 13},
		Page:       1,
		PerPage:    25,
		CaseOwners: []string{"1"},
	})

	assert.Equal(t, nil, err)
	assert.Equal(t, expectedResponse, clientList)
}

func TestGetCaseloadListCanThrow500Error(t *testing.T) {
	logger, _ := SetUpTest()
	svr := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusInternalServerError)
	}))
	defer svr.Close()

	client := NewApiClient(http.DefaultClient, svr.URL, logger)

	clientList, err := client.GetClientList(getContext(nil), ClientListParams{
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
		URL:    svr.URL + "/v1/assignees/13/clients?limit=25&page=1&filter=&sort=",
		Method: http.MethodGet,
	}, err)
}

func TestGetCaseloadListSortedByMadeActiveDateForNewDeputyOrdersTeam(t *testing.T) {
	logger, mockClient := SetUpTest()
	client := NewApiClient(mockClient, "", logger)

	mocks.GetDoFunc = func(r *http.Request) (*http.Response, error) {
		assert.Contains(t, r.URL.RawQuery, "sort=made_active_date:asc")
		assert.NotContains(t, r.URL.RawQuery, "caseowner:1")
		return &http.Response{
			StatusCode: 200,
			Body:       io.NopCloser(bytes.NewReader([]byte("{}"))),
		}, nil
	}

	team := model.Team{Id: 13, Name: "Lay Team - New Deputy Orders"}
	_, err := client.GetClientList(getContext(nil), ClientListParams{
		Team:       team,
		Page:       1,
		PerPage:    25,
		CaseOwners: []string{"1"},
	})
	assert.Nil(t, err)
}

func TestGetCaseloadListSortedByReportDueDateForLayTeam(t *testing.T) {
	logger, mockClient := SetUpTest()
	client := NewApiClient(mockClient, "", logger)

	mocks.GetDoFunc = func(r *http.Request) (*http.Response, error) {
		assert.Contains(t, r.URL.RawQuery, "sort=report_due_date:asc")
		return &http.Response{
			StatusCode: 200,
			Body:       io.NopCloser(bytes.NewReader([]byte("{}"))),
		}, nil
	}

	team := model.Team{Id: 13, Name: "Lay Team 1", Type: "LAY"}
	_, err := client.GetClientList(getContext(nil), ClientListParams{
		Team:    team,
		Page:    1,
		PerPage: 25,
	})
	assert.Nil(t, err)
}

func TestClientListParams_CreateFilter(t *testing.T) {
	tests := []struct {
		params ClientListParams
		want   string
	}{
		{
			params: ClientListParams{},
			want:   "",
		},
		{
			params: ClientListParams{CaseOwners: []string{"1"}},
			want:   "caseowner:1",
		},
		{
			params: ClientListParams{SubType: "hw"},
			want:   "subtype:hw",
		},
		{
			params: ClientListParams{CaseTypes: []string{"HYBRID"}},
			want:   "case-type:HYBRID",
		},
		{
			params: ClientListParams{OrderStatuses: []string{"active", "duplicate"}},
			want:   "order-status:active,order-status:duplicate",
		},
		{
			params: ClientListParams{
				OrderStatuses: []string{"active", "closed"},
				SubType:       "hw",
				DeputyTypes:   []string{"LAY", "PA"},
				CaseTypes:     []string{"HYBRID", "DUAL", "HW", "PFA"},
				CaseOwners:    []string{"1", "2", "3"},
			},
			want: "order-status:active,order-status:closed,subtype:hw,deputy-type:LAY,deputy-type:PA,case-type:HYBRID,case-type:DUAL,case-type:HW,case-type:PFA,caseowner:1,caseowner:2,caseowner:3",
		},
	}
	for i, test := range tests {
		t.Run("Scenario "+strconv.Itoa(i+1), func(t *testing.T) {
			assert.Equal(t, test.want, test.params.CreateFilter())
		})
	}
}

func TestGetClientList_contract(t *testing.T) {
	pact, err := consumer.NewV4Pact(consumer.MockHTTPProviderConfig{
		Consumer: "sirius-supervision-workflow",
		Provider: "sirius",
		LogDir:   "../../logs",
		PactDir:  "../../pacts",
	})
	assert.NoError(t, err)

	err = pact.
		AddInteraction().
		Given("Clients exist for an assignee").
		UponReceiving("A request for the client list").
		WithRequest("GET", "/supervision-api/v1/assignees/13/clients", func(b *consumer.V4RequestBuilder) {
			b.Query("limit", matchers.S("25"))
			b.Query("page", matchers.S("1"))
			b.Query("filter", matchers.S("caseowner:1"))
			b.Query("sort", matchers.S(""))
		}).
		WillRespondWith(200, func(b *consumer.V4ResponseBuilder) {
			b.Header("Content-Type", matchers.S("application/json"))
			b.JSONBody(matchers.StructMatcher{
				"limit": matchers.Like(15),
				"metadata": matchers.StructMatcher{
					"assigneeClientCount": matchers.EachLike(matchers.StructMatcher{
						"assignee": matchers.Like(1),
						"count":    matchers.Like(14),
					}, 1),
				},
				"pages": matchers.StructMatcher{
					"current": matchers.Like(1),
					"total":   matchers.Like(1),
				},
				"total": matchers.Like(1),
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
				}, 1),
			})
		}).
		ExecuteTest(t, func(config consumer.MockServerConfig) error {
			client := NewApiClient(http.DefaultClient, fmt.Sprintf("http://%s:%d/supervision-api", config.Host, config.Port), telemetry.NewLogger("test"))

			clientList, err := client.GetClientList(getContext(nil), ClientListParams{
				Team:       model.Team{Id: 13},
				Page:       1,
				PerPage:    25,
				CaseOwners: []string{"1"},
			})
			assert.NoError(t, err)

			assert.EqualValues(t, 1, clientList.TotalClients)
			assert.EqualValues(t, 1, clientList.Pages.PageCurrent)
			assert.EqualValues(t, 1, len(clientList.Clients))
			assert.EqualValues(t, "Ro", clientList.Clients[0].FirstName)
			return nil
		})

	assert.NoError(t, err)
}
