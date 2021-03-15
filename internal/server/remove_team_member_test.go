package server

import (
	"errors"
	"net/http"
	"net/http/httptest"
	"strconv"
	"strings"
	"testing"

	"github.com/ministryofjustice/opg-sirius-user-management/internal/sirius"
	"github.com/stretchr/testify/assert"
)

type mockRemoveTeamMemberClient struct {
	team struct {
		count   int
		lastCtx sirius.Context
		lastID  int
		data    sirius.Team
		err     error
	}
	editTeam struct {
		count    int
		lastCtx  sirius.Context
		lastTeam sirius.Team
		err      error
	}
}

func (c *mockRemoveTeamMemberClient) Team(ctx sirius.Context, id int) (sirius.Team, error) {
	c.team.count += 1
	c.team.lastCtx = ctx
	c.team.lastID = id

	return c.team.data, c.team.err
}

func (c *mockRemoveTeamMemberClient) EditTeam(ctx sirius.Context, team sirius.Team) error {
	c.editTeam.count += 1
	c.editTeam.lastCtx = ctx
	c.editTeam.lastTeam = team

	return c.editTeam.err
}

func (m *mockRemoveTeamMemberClient) requiredPermissions() sirius.PermissionSet {
	return sirius.PermissionSet{"team": sirius.PermissionGroup{Permissions: []string{"put"}}}
}

func generateTeamWithIds(ids ...int) sirius.Team {
	team := sirius.Team{
		ID: 123,
	}

	for _, id := range ids {
		team.Members = append(team.Members, sirius.TeamMember{
			ID:          id,
			DisplayName: "User " + strconv.Itoa(id),
		})
	}

	return team
}

func TestPostRemoveTeamMember(t *testing.T) {
	assert := assert.New(t)

	client := &mockRemoveTeamMemberClient{}
	client.team.data = generateTeamWithIds(12, 16, 45)
	template := &mockTemplate{}

	w := httptest.NewRecorder()
	r, _ := http.NewRequest("POST", "/teams/remove-member/123", strings.NewReader("selected[]=12&selected[]=45"))
	r.Header.Add("Content-Type", "application/x-www-form-urlencoded")

	err := removeTeamMember(client, template)(client.requiredPermissions(), w, r)
	assert.Nil(err)

	assert.Equal(1, client.team.count)
	assert.Equal(getContext(r), client.team.lastCtx)
	assert.Equal(123, client.team.lastID)

	assert.Equal(0, client.editTeam.count)

	assert.Equal(1, template.count)
	assert.Equal("page", template.lastName)
	assert.Equal(removeTeamMemberVars{
		Path: "/teams/remove-member/123",
		Team: client.team.data,
		Selected: map[int]string{
			12: "User 12",
			45: "User 45",
		},
	}, template.lastVars)
}

func TestPostRemoveTeamMemberNoPermission(t *testing.T) {
	assert := assert.New(t)

	w := httptest.NewRecorder()
	r, _ := http.NewRequest("POST", "/path", nil)

	err := removeTeamMember(nil, nil)(sirius.PermissionSet{}, w, r)
	assert.Equal(StatusError(http.StatusForbidden), err)
}

func TestPostRemoveTeamMemberBadPath(t *testing.T) {
	for name, path := range map[string]string{
		"empty":       "/teams/remove-member/",
		"non-numeric": "/teams/remove-member/hello",
		"suffixed":    "/teams/remove-member/123/no",
		"elsewhere":   "/teams/add-member/123",
	} {
		t.Run(name, func(t *testing.T) {
			assert := assert.New(t)

			client := &mockRemoveTeamMemberClient{}
			r, _ := http.NewRequest("POST", path, strings.NewReader("selected[]=12&selected[]=45"))
			r.Header.Add("Content-Type", "application/x-www-form-urlencoded")
			err := editTeam(nil, nil)(client.requiredPermissions(), nil, r)

			assert.Equal(StatusError(http.StatusNotFound), err)
		})
	}
}

func TestPostRemoveTeamMemberTeamError(t *testing.T) {
	assert := assert.New(t)

	expectedError := errors.New("oops")

	client := &mockRemoveTeamMemberClient{}
	client.team.err = expectedError
	template := &mockTemplate{}

	w := httptest.NewRecorder()
	r, _ := http.NewRequest("POST", "/teams/remove-member/123", strings.NewReader("selected[]=12&selected[]=45"))
	r.Header.Add("Content-Type", "application/x-www-form-urlencoded")

	err := removeTeamMember(client, template)(client.requiredPermissions(), w, r)
	assert.Equal(expectedError, err)

	assert.Equal(1, client.team.count)
	assert.Equal(0, client.editTeam.count)
}

func TestPostRemoveTeamMemberBadData(t *testing.T) {
	for name, data := range map[string]string{
		"invalid":     "%1",
		"non-numeric": "selected[]=string",
	} {
		t.Run(name, func(t *testing.T) {
			assert := assert.New(t)

			client := &mockRemoveTeamMemberClient{}
			template := &mockTemplate{}

			w := httptest.NewRecorder()
			r, _ := http.NewRequest("POST", "/teams/remove-member/123", strings.NewReader(data))
			r.Header.Add("Content-Type", "application/x-www-form-urlencoded")

			err := removeTeamMember(client, template)(client.requiredPermissions(), w, r)
			assert.Equal(StatusError(http.StatusBadRequest), err)

			assert.Equal(0, client.editTeam.count)
		})
	}
}

func TestPostRemoveTeamMemberIgnoresBadIds(t *testing.T) {
	assert := assert.New(t)

	client := &mockRemoveTeamMemberClient{}
	client.team.data = generateTeamWithIds(12, 16, 45)
	template := &mockTemplate{}

	w := httptest.NewRecorder()
	r, _ := http.NewRequest("POST", "/teams/remove-member/123", strings.NewReader("selected[]=19&selected[]=45"))
	r.Header.Add("Content-Type", "application/x-www-form-urlencoded")

	err := removeTeamMember(client, template)(client.requiredPermissions(), w, r)
	assert.Nil(err)

	assert.Equal(removeTeamMemberVars{
		Path: "/teams/remove-member/123",
		Team: client.team.data,
		Selected: map[int]string{
			45: "User 45",
		},
	}, template.lastVars)
}

func TestConfirmPostRemoveTeamMember(t *testing.T) {
	assert := assert.New(t)

	client := &mockRemoveTeamMemberClient{}
	client.team.data = generateTeamWithIds(12, 16, 45)
	template := &mockTemplate{}

	w := httptest.NewRecorder()
	r, _ := http.NewRequest("POST", "/teams/remove-member/123", strings.NewReader("selected[]=12&selected[]=45&confirm=true"))
	r.Header.Add("Content-Type", "application/x-www-form-urlencoded")

	err := removeTeamMember(client, template)(client.requiredPermissions(), w, r)
	assert.Equal(RedirectError("/teams/123"), err)

	assert.Equal(1, client.team.count)
	assert.Equal(1, client.editTeam.count)

	assert.Equal([]sirius.TeamMember{
		{
			ID:          16,
			DisplayName: "User 16",
		},
	}, client.editTeam.lastTeam.Members)
}

func TestConfirmPostRemoveTeamMemberClientError(t *testing.T) {
	assert := assert.New(t)

	client := &mockRemoveTeamMemberClient{}
	client.team.data = generateTeamWithIds(12, 16, 45)
	client.editTeam.err = sirius.ClientError("Team has been deleted")
	template := &mockTemplate{}

	w := httptest.NewRecorder()
	r, _ := http.NewRequest("POST", "/teams/remove-member/123", strings.NewReader("selected[]=12&selected[]=45&confirm=true"))
	r.Header.Add("Content-Type", "application/x-www-form-urlencoded")

	err := removeTeamMember(client, template)(client.requiredPermissions(), w, r)
	assert.Nil(err)

	assert.Equal(1, client.team.count)
	assert.Equal(1, client.editTeam.count)

	assert.Equal(removeTeamMemberVars{
		Path: "/teams/remove-member/123",
		Team: client.team.data,
		Selected: map[int]string{
			12: "User 12",
			45: "User 45",
		},
		Errors: sirius.ValidationErrors{
			"_": {
				"": "Team has been deleted",
			},
		},
	}, template.lastVars)
}

func TestConfirmPostRemoveTeamMemberOtherError(t *testing.T) {
	assert := assert.New(t)

	expectedError := errors.New("oops")

	client := &mockRemoveTeamMemberClient{}
	client.team.data = generateTeamWithIds(12, 16, 45)
	client.editTeam.err = expectedError
	template := &mockTemplate{}

	w := httptest.NewRecorder()
	r, _ := http.NewRequest("POST", "/teams/remove-member/123", strings.NewReader("selected[]=12&selected[]=45&confirm=true"))
	r.Header.Add("Content-Type", "application/x-www-form-urlencoded")

	err := removeTeamMember(client, template)(client.requiredPermissions(), w, r)
	assert.Equal(expectedError, err)

	assert.Equal(1, client.team.count)
	assert.Equal(1, client.editTeam.count)
}

func TestGetRemoveTeamMemberTeam(t *testing.T) {
	assert := assert.New(t)

	client := &mockRemoveTeamMemberClient{}
	w := httptest.NewRecorder()
	r, _ := http.NewRequest("GET", "/teams/remove-member/123", nil)

	err := removeTeamMember(nil, nil)(client.requiredPermissions(), w, r)
	assert.Equal(StatusError(http.StatusMethodNotAllowed), err)
}
