package server

import (
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/ministryofjustice/opg-sirius-user-management/internal/sirius"
	"github.com/stretchr/testify/assert"
)

type mockListTeamsClient struct {
	count   int
	lastCtx sirius.Context
	err     error
	data    []sirius.Team
}

func (m *mockListTeamsClient) Teams(ctx sirius.Context) ([]sirius.Team, error) {
	m.count += 1
	m.lastCtx = ctx

	return m.data, m.err
}

func (m *mockListTeamsClient) requiredPermissions() sirius.PermissionSet {
	return sirius.PermissionSet{"team": sirius.PermissionGroup{Permissions: []string{"put"}}}
}

func TestListTeams(t *testing.T) {
	assert := assert.New(t)

	data := []sirius.Team{
		{
			ID:          29,
			DisplayName: "Milo Nihei",
			Members:     make([]sirius.TeamMember, 10),
			Type:        "Top Notch",
		},
	}
	client := &mockListTeamsClient{
		data: data,
	}
	template := &mockTemplate{}

	w := httptest.NewRecorder()
	r, _ := http.NewRequest("GET", "/path", nil)

	err := listTeams(client, template)(client.requiredPermissions(), w, r)
	assert.Nil(err)

	resp := w.Result()
	assert.Equal(http.StatusOK, resp.StatusCode)
	assert.Equal(getContext(r), client.lastCtx)

	assert.Equal(1, client.count)

	assert.Equal(1, template.count)
	assert.Equal("page", template.lastName)
	assert.Equal(listTeamsVars{
		Path:  "/path",
		Teams: data,
	}, template.lastVars)
}

func TestListTeamsNoPermission(t *testing.T) {
	assert := assert.New(t)

	w := httptest.NewRecorder()
	r, _ := http.NewRequest("GET", "/path", nil)

	err := listTeams(nil, nil)(sirius.PermissionSet{}, w, r)
	assert.Equal(StatusError(http.StatusForbidden), err)
}

func TestListTeamsSearch(t *testing.T) {
	assert := assert.New(t)

	data := []sirius.Team{
		{
			ID:          29,
			DisplayName: "Milo Nihei",
			Members:     make([]sirius.TeamMember, 10),
			Type:        "Top Notch",
		},
		{
			ID:          3,
			DisplayName: "Who",
			Members:     make([]sirius.TeamMember, 5),
			Type:        "Terrible",
		},
	}
	client := &mockListTeamsClient{
		data: data,
	}
	template := &mockTemplate{}

	w := httptest.NewRecorder()
	r, _ := http.NewRequest("GET", "/path?search=milo", nil)

	err := listTeams(client, template)(client.requiredPermissions(), w, r)
	assert.Nil(err)

	resp := w.Result()
	assert.Equal(http.StatusOK, resp.StatusCode)
	assert.Equal(getContext(r), client.lastCtx)

	assert.Equal(1, client.count)

	assert.Equal(1, template.count)
	assert.Equal("page", template.lastName)
	assert.Equal(listTeamsVars{
		Path:   "/path",
		Search: "milo",
		Teams: []sirius.Team{
			{
				ID:          29,
				DisplayName: "Milo Nihei",
				Members:     make([]sirius.TeamMember, 10),
				Type:        "Top Notch",
			},
		},
	}, template.lastVars)
}

func TestListTeamsError(t *testing.T) {
	assert := assert.New(t)

	expectedErr := errors.New("err")
	client := &mockListTeamsClient{err: expectedErr}
	template := &mockTemplate{}

	w := httptest.NewRecorder()
	r, _ := http.NewRequest("GET", "/?search=long", nil)

	err := listTeams(client, template)(client.requiredPermissions(), w, r)

	assert.Equal(expectedErr, err)
	assert.Equal(0, template.count)
}

func TestPostListTeams(t *testing.T) {
	assert := assert.New(t)

	client := &mockListTeamsClient{}
	w := httptest.NewRecorder()
	r, _ := http.NewRequest("POST", "", nil)

	err := listTeams(nil, nil)(client.requiredPermissions(), w, r)
	assert.Equal(StatusError(http.StatusMethodNotAllowed), err)
}
