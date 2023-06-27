package sirius

import (
	"bytes"
	"github.com/ministryofjustice/opg-sirius-workflow/internal/mocks"
	"github.com/stretchr/testify/assert"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestGetCaseloadListCanReturn200(t *testing.T) {
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
                        "dueDate": "21/12/2023"
                    },
					"orderStatus": {
						"handle": "CLOSED",
						"label": "Closed",
						"deprecated": false
					}
                }
            ],
            "supervisionLevel": "Minimal"
        }
    ]
}
`

	r := io.NopCloser(bytes.NewReader([]byte(json)))

	mocks.GetDoFunc = func(*http.Request) (*http.Response, error) {
		return &http.Response{
			StatusCode: 200,
			Body:       r,
		}, nil
	}

	expectedResponse := ClientList{
		Clients: []Client{
			{
				Id:            63,
				CaseRecNumber: "42687883",
				FirstName:     "Ro",
				Surname:       "Bot",
				Orders: []Order{
					{
						Id: 92,
						Status: RefData{
							Handle: "CLOSED",
							Label:  "Closed",
						},
						LatestAnnualReport: AnnualReport{
							DueDate: "21/12/2023",
						},
					},
				},
				SupervisionLevel: "Minimal",
			},
		},
		Pages: PageInformation{
			PageCurrent: 1,
			PageTotal:   1,
		},
		TotalClients: 1,
	}

	selectedTeam := Team{Id: 13}

	clientList, err := client.GetClientList(getContext(nil), selectedTeam.Id)

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

	clientList, err := client.GetClientList(getContext(nil), 13)

	expectedResponse := ClientList{
		Clients:      nil,
		Pages:        PageInformation{},
		TotalClients: 0,
	}

	assert.Equal(t, expectedResponse, clientList)

	assert.Equal(t, StatusError{
		Code:   http.StatusInternalServerError,
		URL:    svr.URL + "/api/v1/assignees/13/clients",
		Method: http.MethodGet,
	}, err)
}
