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
		query := rq.URL.Query()
		assert.Equal(t, []string{"13"}, query["teamIds[]"])
		assert.Equal(t, "25", query.Get("limit"))
		assert.Equal(t, "1", query.Get("page"))
		assert.Equal(t, "ecm:1,ecm:2", query.Get("filter"))
		assert.Equal(t, "field:direction", query.Get("sort"))
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
		URL:    svr.URL + "/v1/assignees/teams/deputies?filter=&limit=25&page=1&sort=&teamIds%5B%5D=13",
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
		Given("A deputy exists").
		UponReceiving("A request for the deputy list").
		WithRequest("GET", "/supervision-api/v1/assignees/teams/deputies", func(b *consumer.V4RequestBuilder) {
			b.Query("teamIds[]", matchers.S("123"))
			b.Query("limit", matchers.S("25"))
			b.Query("page", matchers.S("1"))
			b.Query("filter", matchers.S("ecm:123"))
			b.Query("sort", matchers.S("field:direction"))
		}).
		WillRespondWith(200, func(b *consumer.V4ResponseBuilder) {
			b.Header("Content-Type", matchers.S("application/json"))
			// BodyMatch generates matchers purely from the DTO's Go type via
			// reflection - it ignores the literal field values below, so an
			// empty struct is sufficient and avoids implying specific example
			// values are being asserted on.
			b.BodyMatch(deputyListResponse{})
		}).
		ExecuteTest(t, func(config consumer.MockServerConfig) error {
			client := NewApiClient(http.DefaultClient, fmt.Sprintf("http://%s:%d/supervision-api", config.Host, config.Port), telemetry.NewLogger("test"))

			deputyList, err := client.GetDeputyList(getContext(nil), DeputyListParams{
				Team:         model.Team{Id: 123},
				Page:         1,
				PerPage:      25,
				Sort:         "field:direction",
				SelectedECMs: []string{"123"},
			})
			assert.NoError(t, err)

			assert.EqualValues(t, 1, deputyList.TotalDeputies)
			assert.EqualValues(t, 1, len(deputyList.Deputies))
			assert.EqualValues(t, "string", deputyList.Deputies[0].DisplayName)
			return nil
		})

	assert.NoError(t, err)
}
