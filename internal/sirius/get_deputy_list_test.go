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

func TestApiClient_GetDeputyList_Returns200(t *testing.T) {
	logger, mockClient := SetUpTest()
	client := NewApiClient(mockClient, "http://localhost:3000", logger)

	json := `
{
   "limit": 15,
   "metadata":{"ecmCount": [{"assignee" : 1, "count": 14}]},
   "pages": {
       "current": 1,
       "total": 1
   },
   "total": 1,
   "persons": [
       {
           "id": 13,
           "deputyNumber": 123456,
           "displayName": "Mr Fee-paying Deputy",
           "deputyType": {
             "handle": "PRO",
             "label": "Professional"
           },
           "deputyAddress": {
             "town": "Derby"
           },
           "executiveCaseManager": {
             "displayName": "PROTeam1 User1",
             "id": 96
           },
           "mostRecentlyCompletedAssurance": {
             "reportReviewDate" : "2023-05-26T00:00:00+00:00",
             "reportMarkedAs": {
               "handle": "GREEN",
               "label": "Green"
             },
             "assuranceType": {
               "handle": "VISIT",
               "label": "Visit"
             }
           },
           "activeClientCount": 100,
           "activeNonCompliantClientCount": 10
		}
   ]
}
`

	params := DeputyListParams{
		Team:         model.Team{Id: 13},
		Page:         1,
		PerPage:      25,
		Sort:         "field:direction",
		SelectedECMs: []string{"1", "2"},
	}

	r := io.NopCloser(bytes.NewReader([]byte(json)))

	mocks.GetDoFunc = func(rq *http.Request) (*http.Response, error) {
		assert.Contains(t, rq.URL.RawQuery, "teamIds[]=13&limit=25&page=1&filter=ecm:1,ecm:2&sort=field:direction")
		return &http.Response{
			StatusCode: 200,
			Body:       r,
		}, nil
	}

	expectedResponse := DeputyList{
		Deputies: []model.Deputy{
			{
				Id:          13,
				DisplayName: "Mr Fee-paying Deputy",
				Type:        model.RefData{Handle: "PRO", Label: "Professional"},
				Number:      123456,
				Address:     model.Address{Town: "Derby"},
				ExecutiveCaseManager: model.Assignee{
					Id:   96,
					Name: "PROTeam1 User1",
				},
				Assurance: model.Assurance{
					ReportReviewDate: model.NewDate("26/05/2023"),
					ReportMarkedAs:   model.RefData{Handle: "GREEN", Label: "Green"},
					Type:             model.RefData{Handle: "VISIT", Label: "Visit"},
				},
				ActiveClientCount:             100,
				ActiveNonCompliantClientCount: 10,
			},
		},
		Pages: model.PageInformation{
			PageCurrent: 1,
			PageTotal:   1,
		},
		TotalDeputies: 1,
		MetaData: DeputyMetaData{
			[]model.AssigneeAndCount{
				{AssigneeId: 1, Count: 14},
			},
		},
	}

	deputyList, err := client.GetDeputyList(getContext(nil), params)

	assert.Equal(t, nil, err)
	assert.Equal(t, expectedResponse, deputyList)
}

func TestApiClient_GetDeputyList_Returns500(t *testing.T) {
	logger, _ := SetUpTest()
	svr := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusInternalServerError)
	}))
	defer svr.Close()

	client := NewApiClient(http.DefaultClient, svr.URL, logger)

	deputyList, err := client.GetDeputyList(getContext(nil), DeputyListParams{
		Team:    model.Team{Id: 13},
		Page:    1,
		PerPage: 25,
	})

	expectedResponse := DeputyList{
		Deputies:      nil,
		Pages:         model.PageInformation{},
		TotalDeputies: 0,
	}

	assert.Equal(t, expectedResponse, deputyList)

	assert.Equal(t, StatusError{
		Code:   http.StatusInternalServerError,
		URL:    svr.URL + "/v1/assignees/teams/deputies?teamIds[]=13&limit=25&page=1&filter=&sort=",
		Method: http.MethodGet,
	}, err)
}

func TestGetDeputyList_contract(t *testing.T) {
	pact, err := consumer.NewV4Pact(consumer.MockHTTPProviderConfig{
		Consumer: "sirius-supervision-workflow",
		Provider: "sirius",
		LogDir:   "../../logs",
		PactDir:  "../../pacts",
	})
	assert.NoError(t, err)

	err = pact.
		AddInteraction().
		Given("Deputies exist for requested teams").
		UponReceiving("A request for the deputy list").
		WithRequest("GET", "/supervision-api/v1/assignees/teams/deputies", func(b *consumer.V4RequestBuilder) {
			b.Query("teamIds[]", matchers.S("13"))
			b.Query("limit", matchers.S("25"))
			b.Query("page", matchers.S("1"))
			b.Query("filter", matchers.S("ecm:1,ecm:2"))
			b.Query("sort", matchers.S("field:direction"))
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
					"ecmCount": matchers.EachLike(matchers.StructMatcher{
						"assignee": matchers.Like(1),
						"count":    matchers.Like(14),
					}, 1),
				},
				"persons": matchers.EachLike(matchers.StructMatcher{
					"id":           matchers.Like(13),
					"deputyNumber": matchers.Like(123456),
					"displayName":  matchers.Like("Mr Fee-paying Deputy"),
					"deputyType": matchers.StructMatcher{
						"handle": matchers.Like("PRO"),
						"label":  matchers.Like("Professional"),
					},
					"deputyAddress": matchers.StructMatcher{
						"town": matchers.Like("Derby"),
					},
					"executiveCaseManager": matchers.StructMatcher{
						"displayName": matchers.Like("PROTeam1 User1"),
						"id":          matchers.Like(96),
					},
					"mostRecentlyCompletedAssurance": matchers.StructMatcher{
						"reportReviewDate": matchers.Like("2023-05-26T00:00:00+00:00"),
						"reportMarkedAs": matchers.StructMatcher{
							"handle": matchers.Like("GREEN"),
							"label":  matchers.Like("Green"),
						},
						"assuranceType": matchers.StructMatcher{
							"handle": matchers.Like("VISIT"),
							"label":  matchers.Like("Visit"),
						},
					},
					"activeClientCount":             matchers.Like(100),
					"activeNonCompliantClientCount": matchers.Like(10),
				}, 1),
			})
		}).
		ExecuteTest(t, func(config consumer.MockServerConfig) error {
			client := NewApiClient(http.DefaultClient, fmt.Sprintf("http://%s:%d/supervision-api", config.Host, config.Port), telemetry.NewLogger("test"))

			deputyList, err := client.GetDeputyList(getContext(nil), DeputyListParams{
				Team:         model.Team{Id: 13},
				Page:         1,
				PerPage:      25,
				Sort:         "field:direction",
				SelectedECMs: []string{"1", "2"},
			})
			assert.NoError(t, err)

			assert.EqualValues(t, 1, deputyList.TotalDeputies)
			assert.EqualValues(t, 1, len(deputyList.Deputies))
			assert.EqualValues(t, "Mr Fee-paying Deputy", deputyList.Deputies[0].DisplayName)
			return nil
		})

	assert.NoError(t, err)
}
