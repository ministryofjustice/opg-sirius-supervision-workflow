package sirius

import (
	"bytes"
	"github.com/ministryofjustice/opg-sirius-workflow/internal/mocks"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestGetCurrentUserDetails(t *testing.T) {
	logger, mockClient := SetUpTest()
	client, _ := NewClient(mockClient, "http://localhost:3000", logger)

	json := `{
			   "id":65,
			   "name":"case",
			   "phoneNumber":"12345678",
			   "teams":[{
				  "displayName":"Lay Team 1 - (Supervision)",
				  "id":13
			   }],
			   "displayName":"case manager",
			   "deleted":false,
			   "email":"case.manager@opgtest.com",
			   "firstname":"case",
			   "surname":"manager",
			   "roles":[
				  "Case Manager"
			   ],
			   "locked":false,
			   "suspended":false
			}`

	r := io.NopCloser(bytes.NewReader([]byte(json)))

	mocks.GetDoFunc = func(*http.Request) (*http.Response, error) {
		return &http.Response{
			StatusCode: 200,
			Body:       r,
		}, nil
	}

	expectedResponse := UserDetails{
		ID:          65,
		Name:        "case",
		PhoneNumber: "12345678",
		Teams: []MyDetailsTeam{
			{
				DisplayName: "Lay Team 1 - (Supervision)",
				TeamId:      13,
			},
		},
		DisplayName: "case manager",
		Deleted:     false,
		Email:       "case.manager@opgtest.com",
		Firstname:   "case",
		Surname:     "manager",
		Roles:       []string{"Case Manager"},
		Locked:      false,
		Suspended:   false,
	}

	teams, err := client.GetCurrentUserDetails(getContext(nil))
	assert.Equal(t, expectedResponse, teams)
	assert.Equal(t, nil, err)
}

func TestGetCurrentUserDetailsReturnsUnauthorisedClientError(t *testing.T) {
	logger, _ := SetUpTest()
	svr := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusUnauthorized)
	}))
	defer svr.Close()

	client, _ := NewClient(http.DefaultClient, svr.URL, logger)
	_, err := client.GetCurrentUserDetails(getContext(nil))
	assert.Equal(t, ErrUnauthorized, err)
}

func TestMyDetailsReturns500Error(t *testing.T) {
	logger, _ := SetUpTest()
	svr := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusInternalServerError)
	}))
	defer svr.Close()

	client, _ := NewClient(http.DefaultClient, svr.URL, logger)

	_, err := client.GetCurrentUserDetails(getContext(nil))
	assert.Equal(t, StatusError{
		Code:   http.StatusInternalServerError,
		URL:    svr.URL + "/api/v1/users/current",
		Method: http.MethodGet,
	}, err)
}

func TestMyDetailsReturns200(t *testing.T) {
	logger, mockClient := SetUpTest()
	client, _ := NewClient(mockClient, "http://localhost:3000", logger)

	json := `{
		"id": 55,
		"name": "case",
		"phoneNumber": "12345678",
		"teams": [],
		"displayName": "case manager",
		"deleted": false,
		"email": "case.manager@opgtest.com",
		"firstname": "case",
		"surname": "manager",
		"roles": [
			"OPG User",
			"Case Manager"
		],
		"locked": false,
		"suspended": false
    }`

	r := io.NopCloser(bytes.NewReader([]byte(json)))

	mocks.GetDoFunc = func(*http.Request) (*http.Response, error) {
		return &http.Response{
			StatusCode: 200,
			Body:       r,
		}, nil
	}

	expectedResponse := UserDetails{
		ID:          55,
		Name:        "case",
		PhoneNumber: "12345678",
		Teams:       []MyDetailsTeam{},
		DisplayName: "case manager",
		Deleted:     false,
		Email:       "case.manager@opgtest.com",
		Firstname:   "case",
		Surname:     "manager",
		Roles:       []string{"OPG User", "Case Manager"},
		Locked:      false,
		Suspended:   false,
	}

	user, err := client.GetCurrentUserDetails(getContext(nil))
	assert.Equal(t, err, nil)
	assert.Equal(t, user, expectedResponse)
}
