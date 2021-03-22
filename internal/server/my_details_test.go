package server

import (
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/ministryofjustice/opg-sirius-workflow/internal/sirius"
	"github.com/stretchr/testify/assert"
)

// type mockAddTeamMemberClient struct {
// 	team struct {
// 		count   int
// 		lastCtx sirius.Context
// 		lastID  int
// 		data    sirius.Team
// 		err     error
// 	}
// 	editTeam struct {
// 		count    int
// 		lastCtx  sirius.Context
// 		lastTeam sirius.Team
// 		err      error
// 	}
// 	searchUsers struct {
// 		count      int
// 		lastCtx    sirius.Context
// 		lastSearch string
// 		data       []sirius.User
// 		err        error
// 	}
// }

// func (c *mockAddTeamMemberClient) Team(ctx sirius.Context, id int) (sirius.Team, error) {
// 	c.team.count += 1
// 	c.team.lastCtx = ctx
// 	c.team.lastID = id

// 	return c.team.data, c.team.err
// }

// func (c *mockAddTeamMemberClient) EditTeam(ctx sirius.Context, team sirius.Team) error {
// 	c.editTeam.count += 1
// 	c.editTeam.lastCtx = ctx
// 	c.editTeam.lastTeam = team

// 	return c.editTeam.err
// }

// func (c *mockAddTeamMemberClient) SearchUsers(ctx sirius.Context, search string) ([]sirius.User, error) {
// 	c.searchUsers.count += 1
// 	c.searchUsers.lastCtx = ctx
// 	c.searchUsers.lastSearch = search

// 	return c.searchUsers.data, c.searchUsers.err
// }

type mockUserDetailsClient struct {
	count           int
	lastCtx         sirius.Context
	err             error
	userdetailsdata sirius.UserDetails
	taskdetailsdata []sirius.ApiTaskTypes
}

func (m *mockUserDetailsClient) SiriusUserDetails(ctx sirius.Context) (sirius.UserDetails, error) {
	m.count += 1
	m.lastCtx = ctx

	return m.userdetailsdata, m.err
}

func (c *mockUserDetailsClient) GetTaskDetails(ctx sirius.Context) ([]sirius.ApiTaskTypes, error) {
	c.count += 1
	c.lastCtx = ctx

	return c.taskdetailsdata, c.err
}

func TestGetMyDetails(t *testing.T) {
	assert := assert.New(t)

	data := sirius.UserDetails{
		Firstname: "John",
		Surname:   "Doe",
	}
	client := &mockUserDetailsClient{userdetailsdata: data}
	template := &mockTemplates{}

	w := httptest.NewRecorder()
	r, _ := http.NewRequest("GET", "", nil)

	handler := loggingInfoForWorflow(client, template)
	err := handler(sirius.PermissionSet{}, w, r)

	assert.Nil(err)

	resp := w.Result()
	assert.Equal(http.StatusOK, resp.StatusCode)
	assert.Equal(getContext(r), client.lastCtx)

	assert.Equal(1, template.count)
	assert.Equal("page", template.lastName)
	assert.Equal(userDetailsVars{
		Path:      "",
		Firstname: "John",
		Surname:   "Doe",
	}, template.lastVars)
}

func TestGetTaskTypes(t *testing.T) {
	assert := assert.New(t)

	data := []sirius.ApiTaskTypes{
		{
			Handle:     "TestHandle",
			Incomplete: "TestIncomplete",
			Category:   "TestCategory",
			Complete:   "TestComplete",
			User:       true,
		},
	}
	client := &mockUserDetailsClient{taskdetailsdata: data}
	template := &mockTemplates{}

	w := httptest.NewRecorder()
	r, _ := http.NewRequest("GET", "", nil)

	handler := listTaskTypes(client, template)
	err := handler(sirius.PermissionSet{}, w, r)
	assert.Nil(err)

	resp := w.Result()
	assert.Equal(http.StatusOK, resp.StatusCode)
	assert.Equal(getContext(r), client.lastCtx)

	assert.Equal(1, template.count)
	assert.Equal("page", template.lastName)
	assert.Equal(listTaskTypeVars{
		Path:      "",
		LoadTasks: data,
	}, template.lastVars)
}

func TestGetMyDetailsUnauthenticated(t *testing.T) {
	assert := assert.New(t)

	client := &mockUserDetailsClient{err: sirius.ErrUnauthorized}
	templates := &mockTemplates{}

	w := httptest.NewRecorder()
	r, _ := http.NewRequest("GET", "", nil)

	handler := loggingInfoForWorflow(client, templates)
	err := handler(sirius.PermissionSet{}, w, r)

	assert.Equal(sirius.ErrUnauthorized, err)

	assert.Equal(0, templates.count)
}

func TestGetMyDetailsSiriusErrors(t *testing.T) {
	assert := assert.New(t)

	client := &mockUserDetailsClient{err: errors.New("err")}
	template := &mockTemplates{}

	w := httptest.NewRecorder()
	r, _ := http.NewRequest("GET", "", nil)

	handler := loggingInfoForWorflow(client, template)
	err := handler(sirius.PermissionSet{}, w, r)

	assert.Equal("err", err.Error())

	assert.Equal(0, template.count)
}

func TestPostMyDetails(t *testing.T) {
	assert := assert.New(t)
	templates := &mockTemplates{}

	w := httptest.NewRecorder()
	r, _ := http.NewRequest("POST", "", nil)

	handler := loggingInfoForWorflow(nil, templates)
	err := handler(sirius.PermissionSet{}, w, r)

	assert.Equal(StatusError(http.StatusMethodNotAllowed), err)

	assert.Equal(0, templates.count)
}
