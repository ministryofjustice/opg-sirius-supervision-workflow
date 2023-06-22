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
	client, _ := NewClient(mockClient, "http://localhost:3000", logger)

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
		WholeClientList: []ApiClient{
			{
				Id:            63,
				CaseRecNumber: "42687883",
				FirstName:     "Ro",
				Surname:       "Bot",
				Case: []Order{
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

	assigneeTeams, err := client.GetCaseloadList(getContext(nil), selectedTeam.Id)

	assert.Equal(t, expectedResponse, assigneeTeams)
	assert.Equal(t, nil, err)
}

func TestGetCaseloadListCanThrow500Error(t *testing.T) {
	tests := []struct {
		name         string
		selectedTeam Team
		expectedURL  string
	}{
		{
			name:         "Single Team ID requested",
			selectedTeam: Team{Id: 13},
			expectedURL:  "/api/v1/assignees/teams/tasks?teamIds[]=13&filter=status:Not+started&limit=25&page=1&sort=dueDate:asc",
		},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			logger, _ := SetUpTest()
			svr := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				w.WriteHeader(http.StatusInternalServerError)
			}))
			defer svr.Close()

			client, _ := NewClient(http.DefaultClient, svr.URL, logger)

			assigneeTeams, err := client.GetTaskList(getContext(nil), 1, 25, test.selectedTeam, []string{}, []ApiTaskTypes{}, []string{}, nil, nil)

			expectedResponse := TaskList{
				WholeTaskList: nil,
				Pages:         PageInformation{},
				TotalTasks:    0,
				ActiveFilters: nil,
			}

			assert.Equal(t, expectedResponse, assigneeTeams)

			assert.Equal(t, StatusError{
				Code:   http.StatusInternalServerError,
				URL:    svr.URL + test.expectedURL,
				Method: http.MethodGet,
			}, err)
		})
	}
}
