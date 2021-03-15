package server

import (
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/ministryofjustice/opg-sirius-user-management/internal/sirius"
	"github.com/stretchr/testify/assert"
)

type mockDeleteTeamClient struct {
	team struct {
		count   int
		lastCtx sirius.Context
		lastID  int
		data    sirius.Team
		err     error
	}

	deleteTeam struct {
		count      int
		lastCtx    sirius.Context
		lastTeamID int
		err        error
	}
}

func (m *mockDeleteTeamClient) Team(ctx sirius.Context, id int) (sirius.Team, error) {
	m.team.count += 1
	m.team.lastCtx = ctx
	m.team.lastID = id

	return m.team.data, m.team.err
}

func (m *mockDeleteTeamClient) DeleteTeam(ctx sirius.Context, teamID int) error {
	m.deleteTeam.count += 1
	m.deleteTeam.lastCtx = ctx
	m.deleteTeam.lastTeamID = teamID

	return m.deleteTeam.err
}

func (m *mockDeleteTeamClient) requiredPermissions() sirius.PermissionSet {
	return sirius.PermissionSet{"v1-teams": sirius.PermissionGroup{Permissions: []string{"delete"}}}
}

func TestGetDeleteTeam(t *testing.T) {
	assert := assert.New(t)

	client := &mockDeleteTeamClient{}
	client.team.data = sirius.Team{DisplayName: "Filing - Pool 5"}
	template := &mockTemplate{}

	w := httptest.NewRecorder()
	r, _ := http.NewRequest("GET", "/teams/delete/461", nil)

	err := deleteTeam(client, template)(client.requiredPermissions(), w, r)
	assert.Nil(err)

	assert.Equal(1, client.team.count)
	assert.Equal(461, client.team.lastID)
	assert.Equal(0, client.deleteTeam.count)

	assert.Equal(1, template.count)
	assert.Equal("page", template.lastName)
	assert.Equal(deleteTeamVars{
		Path: "/teams/delete/461",
		Team: client.team.data,
	}, template.lastVars)
}

func TestGetDeleteTeamNoPermission(t *testing.T) {
	assert := assert.New(t)

	w := httptest.NewRecorder()
	r, _ := http.NewRequest("GET", "/path", nil)

	err := deleteTeam(nil, nil)(sirius.PermissionSet{}, w, r)
	assert.Equal(StatusError(http.StatusForbidden), err)
}

func TestGetDeleteTeamError(t *testing.T) {
	assert := assert.New(t)

	expectedError := errors.New("oops")

	client := &mockDeleteTeamClient{}
	client.team.err = expectedError
	template := &mockTemplate{}

	w := httptest.NewRecorder()
	r, _ := http.NewRequest("GET", "/teams/delete/461", nil)

	err := deleteTeam(client, template)(client.requiredPermissions(), w, r)
	assert.Equal(expectedError, err)

	assert.Equal(1, client.team.count)
	assert.Equal(461, client.team.lastID)
	assert.Equal(0, client.deleteTeam.count)
}

func TestGetDeleteTeamBadPath(t *testing.T) {
	for name, path := range map[string]string{
		"empty":       "/teams/delete/",
		"non-numeric": "/teams/delete/hello",
		"suffixed":    "/teams/delete/461/no",
	} {
		t.Run(name, func(t *testing.T) {
			assert := assert.New(t)

			client := &mockDeleteTeamClient{}
			template := &mockTemplate{}

			w := httptest.NewRecorder()
			r, _ := http.NewRequest("GET", path, nil)

			err := deleteTeam(client, template)(client.requiredPermissions(), w, r)
			assert.Equal(StatusError(http.StatusNotFound), err)

			assert.Equal(0, client.team.count)
			assert.Equal(0, client.deleteTeam.count)
			assert.Equal(0, template.count)
		})
	}
}

func TestPostDeleteTeam(t *testing.T) {
	assert := assert.New(t)

	client := &mockDeleteTeamClient{}
	client.team.data = sirius.Team{DisplayName: "Filing - Pool 5"}
	template := &mockTemplate{}

	w := httptest.NewRecorder()
	r, _ := http.NewRequest("POST", "/teams/delete/461", nil)

	err := deleteTeam(client, template)(client.requiredPermissions(), w, r)
	assert.Nil(err)

	assert.Equal(1, client.team.count)
	assert.Equal(1, client.deleteTeam.count)
	assert.Equal(getContext(r), client.deleteTeam.lastCtx)
	assert.Equal(461, client.deleteTeam.lastTeamID)

	assert.Equal(1, template.count)
	assert.Equal(deleteTeamVars{
		Path:           "/teams/delete/461",
		Team:           client.team.data,
		SuccessMessage: "The team \"Filing - Pool 5\" was deleted.",
	}, template.lastVars)
}

func TestPostDeleteTeamClientError(t *testing.T) {
	assert := assert.New(t)

	client := &mockDeleteTeamClient{}
	client.deleteTeam.err = sirius.ClientError("problem")
	template := &mockTemplate{}

	w := httptest.NewRecorder()
	r, _ := http.NewRequest("POST", "/teams/delete/461", nil)

	err := deleteTeam(client, template)(client.requiredPermissions(), w, r)
	assert.Nil(err)

	assert.Equal(1, client.team.count)
	assert.Equal(1, client.deleteTeam.count)

	assert.Equal(1, template.count)
	assert.Equal("page", template.lastName)
	assert.Equal(deleteTeamVars{
		Path: "/teams/delete/461",
		Team: client.team.data,
		Errors: sirius.ValidationErrors{
			"": {
				"": "problem",
			},
		},
	}, template.lastVars)
}

func TestPostDeleteTeamOtherError(t *testing.T) {
	assert := assert.New(t)

	expectedErr := errors.New("oops")
	client := &mockDeleteTeamClient{}
	client.deleteTeam.err = expectedErr
	template := &mockTemplate{}

	w := httptest.NewRecorder()
	r, _ := http.NewRequest("POST", "/teams/delete/461", nil)

	err := deleteTeam(client, template)(client.requiredPermissions(), w, r)
	assert.Equal(expectedErr, err)

	assert.Equal(1, client.team.count)
	assert.Equal(1, client.deleteTeam.count)
	assert.Equal(0, template.count)
}

func TestPutDeleteTeam(t *testing.T) {
	assert := assert.New(t)

	client := &mockDeleteTeamClient{}
	w := httptest.NewRecorder()
	r, _ := http.NewRequest("PUT", "/teams/delete/461", nil)

	err := deleteTeam(nil, nil)(client.requiredPermissions(), w, r)
	assert.Equal(StatusError(http.StatusMethodNotAllowed), err)
}
