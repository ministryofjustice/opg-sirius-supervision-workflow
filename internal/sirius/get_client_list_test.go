package sirius

import (
	"bytes"
	"io"
	"net/http"
	"net/http/httptest"
	"strconv"
	"testing"

	"github.com/ministryofjustice/opg-sirius-workflow/internal/mocks"
	"github.com/ministryofjustice/opg-sirius-workflow/internal/model"
	"github.com/stretchr/testify/assert"
)

func TestGetCaseloadListCanReturn200(t *testing.T) {
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

	client, _ := NewApiClient(http.DefaultClient, svr.URL, logger)

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
	client, _ := NewApiClient(mockClient, "", logger)

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
	client, _ := NewApiClient(mockClient, "", logger)

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
