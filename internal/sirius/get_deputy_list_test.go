package sirius

import (
	"bytes"
	"github.com/ministryofjustice/opg-sirius-workflow/internal/mocks"
	"github.com/ministryofjustice/opg-sirius-workflow/internal/model"
	"github.com/stretchr/testify/assert"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestApiClient_GetDeputyList_Returns200(t *testing.T) {
	logger, mockClient := SetUpTest()
	client, _ := NewApiClient(mockClient, "http://localhost:3000", logger)

	json := `
{
    "limit": 15,
    "metadata": [],
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
		Team:    model.Team{Id: 13},
		Page:    1,
		PerPage: 25,
		Sort:    "field:direction",
	}

	r := io.NopCloser(bytes.NewReader([]byte(json)))

	mocks.GetDoFunc = func(rq *http.Request) (*http.Response, error) {
		assert.Contains(t, rq.URL.RawQuery, "teamIds[]=13&limit=25&page=1&sort=field:direction")
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

	client, _ := NewApiClient(http.DefaultClient, svr.URL, logger)

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
		URL:    svr.URL + "/api/v1/assignees/teams/deputies?teamIds[]=13&limit=25&page=1&sort=",
		Method: http.MethodGet,
	}, err)
}
