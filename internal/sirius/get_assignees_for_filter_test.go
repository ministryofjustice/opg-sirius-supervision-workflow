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

func TestGetMembersForTeamReturned(t *testing.T) {

	mockClient := &mocks.MockClient{}
	client, _ := NewClient(mockClient, "http://localhost:3000")

	json := `    {
      	"id": 13,
      	"name": "Lay Team 1 - (Supervision)"
    }`

	r := io.NopCloser(bytes.NewReader([]byte(json)))

	mocks.GetDoFunc = func(*http.Request) (*http.Response, error) {
		return &http.Response{
			StatusCode: 200,
			Body:       r,
		}, nil
	}

	expectedResponse := AssigneesTeam{
		Id:      13,
		Name:    "Lay Team 1 - (Supervision)",
		Members: []AssigneeTeamMembers{},
	}

	assigneeTeams, err := client.GetAssigneesForFilter(getContext(nil), 13, []string{""})

	assert.Equal(t, expectedResponse, assigneeTeams)
	assert.Equal(t, nil, err)
}

func TestIsAssigneeSelected(t *testing.T) {
	teamMembersSelected := []string{"15", "88", "89", "90"}
	assert.Equal(t, IsAssigneeSelected(90, teamMembersSelected), true)

	teamMembersNotSelected := []string{"99", "88", "89"}
	assert.Equal(t, IsAssigneeSelected(90, teamMembersNotSelected), false)
}

func TestSortMembersAlphabetically(t *testing.T) {
	expected := []AssigneeTeamMembers{
		{93, "Andrews Andrews", "Apple", false},
		{91, "Anthony Anthony", "Baldwin", false},
		{92, "Ben Ben", "Benjamin", false},
		{88, "Cat Cat", "Cathy", false},
		{89, "LayTeam1 LayTeam1", "User10", false},
	}

	unsortedMembers := []AssigneeTeamMembers{
		{88, "Cat Cat", "Cathy", false},
		{89, "LayTeam1 LayTeam1", "User10", false},
		{91, "Anthony Anthony", "Baldwin", false},
		{92, "Ben Ben", "Benjamin", false},
		{93, "Andrews Andrews", "Apple", false},
	}

	assert.Equal(t, SortMembersAlphabetically(unsortedMembers), expected)
}

func TestAssigneesForFilterReturnsNewStatusError(t *testing.T) {
	svr := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusMethodNotAllowed)
	}))
	defer svr.Close()

	client, _ := NewClient(http.DefaultClient, svr.URL)

	assigneeTeams, err := client.GetAssigneesForFilter(getContext(nil), 13, []string{""})

	expectedResponse := AssigneesTeam{
		Id:      0,
		Name:    "",
		Members: []AssigneeTeamMembers(nil),
	}

	assert.Equal(t, expectedResponse, assigneeTeams)

	assert.Equal(t, StatusError{
		Code:   http.StatusMethodNotAllowed,
		URL:    svr.URL + "/api/v1/teams/13",
		Method: http.MethodGet,
	}, err)
}

func TestAssigneesForFilterReturnsUnauthorisedClientError(t *testing.T) {
	svr := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusUnauthorized)
	}))
	defer svr.Close()

	client, _ := NewClient(http.DefaultClient, svr.URL)

	assigneeTeams, err := client.GetAssigneesForFilter(getContext(nil), 13, []string{""})

	expectedResponse := AssigneesTeam{
		Id:      0,
		Name:    "",
		Members: []AssigneeTeamMembers(nil),
	}

	assert.Equal(t, ErrUnauthorized, err)
	assert.Equal(t, expectedResponse, assigneeTeams)
}
