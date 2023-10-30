package sirius

import (
	"bytes"
	"github.com/ministryofjustice/opg-sirius-workflow/internal/mocks"
	"github.com/ministryofjustice/opg-sirius-workflow/internal/model"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestGetTeamsForSelection(t *testing.T) {
	logger, mockClient := SetUpTest()
	client, _ := NewApiClient(mockClient, "http://localhost:3000", logger)

	json := `[
		{
			"id":21,
			"displayName":"Allocations - (Supervision)",
			"email":"allocations.team@opgtest.com",
			"phoneNumber":"0123456789",
			"members":[
				{
					"id":71,
					"displayName":"Allocations User1",
					"email":"allocations@opgtest.com"
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
			"email":"lay.team.1@opgtest.com",
			"phoneNumber":"0123456789",
			"members":[],
			"teamType":{
				"handle":"LAY",
				"label":"Lay Team"
			}
		},
		{
			"id":23,
			"displayName":"Pro Team 1",
			"email":"pro.team.1@opgtest.com",
			"phoneNumber":"0123456789",
			"members":[],
			"teamType":{
				"handle":"PRO",
				"label":"Pro Team"
			}
		},
		{
			"id":24,
			"displayName":"LPA Team",
			"email":"lpa.team@opgtest.com",
			"phoneNumber":"0987654321",
			"members":[
				{
					"id":72,
					"displayName":"LPA User1",
					"email":"lpa.user@opgtest.com"
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

	teams, err := client.GetTeamsForSelection(getContext(nil), []string{})
	assert.Equal(t, expectedResponse, teams)
	assert.Equal(t, nil, err)
}

func TestGetTeamsForSelectionCanReturn500(t *testing.T) {
	logger, _ := SetUpTest()
	svr := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusInternalServerError)
	}))
	defer svr.Close()

	client, _ := NewApiClient(http.DefaultClient, svr.URL, logger)

	_, err := client.GetTeamsForSelection(getContext(nil), []string{})

	assert.Equal(t, StatusError{
		Code:   http.StatusInternalServerError,
		URL:    svr.URL + "/api/v1/teams",
		Method: http.MethodGet,
	}, err)
}
