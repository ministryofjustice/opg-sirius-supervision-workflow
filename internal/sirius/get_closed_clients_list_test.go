package sirius

import (
	"bytes"
	"github.com/ministryofjustice/opg-sirius-workflow/internal/mocks"
	"github.com/ministryofjustice/opg-sirius-workflow/internal/model"
	"github.com/stretchr/testify/assert"
	"io"
	"net/http"
	"net/http/httptest"
	"strconv"
	"testing"
	"time"
)

func TestGetClosedCaseloadListCanReturn200(t *testing.T) {
	logger, mockClient := SetUpTest()
	client, _ := NewApiClient(mockClient, "http://localhost:3000", logger)

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

	client, _ := NewApiClient(http.DefaultClient, svr.URL, logger)

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
		URL:    svr.URL + "/api/v1/assignees/closed-clients?limit=25&page=1&filter=",
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
