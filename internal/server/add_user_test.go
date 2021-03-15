package server

import (
	"errors"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/ministryofjustice/opg-sirius-user-management/internal/sirius"
	"github.com/stretchr/testify/assert"
)

type mockAddUserClient struct {
	addUser struct {
		count            int
		lastCtx          sirius.Context
		lastEmail        string
		lastFirstname    string
		lastSurname      string
		lastOrganisation string
		lastRoles        []string
		err              error
	}

	roles struct {
		count   int
		lastCtx sirius.Context
		err     error
	}
}

func (m *mockAddUserClient) AddUser(ctx sirius.Context, email, firstname, surname, organisation string, roles []string) error {
	m.addUser.count += 1
	m.addUser.lastCtx = ctx
	m.addUser.lastEmail = email
	m.addUser.lastFirstname = firstname
	m.addUser.lastSurname = surname
	m.addUser.lastOrganisation = organisation
	m.addUser.lastRoles = roles

	return m.addUser.err
}

func (m *mockAddUserClient) Roles(ctx sirius.Context) ([]string, error) {
	m.roles.count += 1
	m.roles.lastCtx = ctx

	return []string{"System Admin", "Manager"}, m.roles.err
}

func (m *mockAddUserClient) requiredPermissions() sirius.PermissionSet {
	return sirius.PermissionSet{"v1-users": sirius.PermissionGroup{Permissions: []string{"post"}}}
}

func TestGetAddUser(t *testing.T) {
	assert := assert.New(t)

	client := &mockAddUserClient{}
	template := &mockTemplate{}

	w := httptest.NewRecorder()
	r, _ := http.NewRequest("GET", "/path", nil)

	err := addUser(client, template)(client.requiredPermissions(), w, r)
	assert.Nil(err)

	assert.Equal(1, client.roles.count)
	assert.Equal(getContext(r), client.roles.lastCtx)

	assert.Equal(0, client.addUser.count)

	assert.Equal(1, template.count)
	assert.Equal("page", template.lastName)
	assert.Equal(addUserVars{
		Path:  "/path",
		Roles: []string{"System Admin", "Manager"},
	}, template.lastVars)
}

func TestGetAddUserNoPermission(t *testing.T) {
	assert := assert.New(t)

	w := httptest.NewRecorder()
	r, _ := http.NewRequest("GET", "/path", nil)

	err := addUser(nil, nil)(sirius.PermissionSet{}, w, r)
	assert.Equal(StatusError(http.StatusForbidden), err)
}

func TestPostAddUser(t *testing.T) {
	assert := assert.New(t)

	client := &mockAddUserClient{}
	template := &mockTemplate{}

	w := httptest.NewRecorder()
	r, _ := http.NewRequest("POST", "/path", strings.NewReader("email=a&firstname=b&surname=c&organisation=d&roles=e&roles=f"))
	r.Header.Add("Content-Type", "application/x-www-form-urlencoded")

	err := addUser(client, template)(client.requiredPermissions(), w, r)
	assert.Nil(err)

	assert.Equal(1, client.roles.count)

	assert.Equal(1, client.addUser.count)
	assert.Equal(getContext(r), client.addUser.lastCtx)
	assert.Equal("a", client.addUser.lastEmail)
	assert.Equal("b", client.addUser.lastFirstname)
	assert.Equal("c", client.addUser.lastSurname)
	assert.Equal("d", client.addUser.lastOrganisation)
	assert.Equal([]string{"e", "f"}, client.addUser.lastRoles)

	assert.Equal(1, template.count)
	assert.Equal("page", template.lastName)
	assert.Equal(addUserVars{
		Path:    "/path",
		Success: true,
		Roles:   []string{"System Admin", "Manager"},
	}, template.lastVars)
}

func TestPostAddUserValidationError(t *testing.T) {
	assert := assert.New(t)

	errors := sirius.ValidationErrors{
		"x": {
			"y": "z",
		},
	}
	client := &mockAddUserClient{}
	client.addUser.err = sirius.ValidationError{
		Errors: errors,
	}
	template := &mockTemplate{}

	w := httptest.NewRecorder()
	r, _ := http.NewRequest("POST", "/path", nil)

	err := addUser(client, template)(client.requiredPermissions(), w, r)
	assert.Nil(err)

	resp := w.Result()
	assert.Equal(http.StatusBadRequest, resp.StatusCode)

	assert.Equal(1, client.roles.count)
	assert.Equal(1, client.addUser.count)

	assert.Equal(1, template.count)
	assert.Equal("page", template.lastName)
	assert.Equal(addUserVars{
		Path:   "/path",
		Roles:  []string{"System Admin", "Manager"},
		Errors: errors,
	}, template.lastVars)
}

func TestPostAddUserOtherError(t *testing.T) {
	assert := assert.New(t)

	expectedErr := errors.New("oops")
	client := &mockAddUserClient{}
	client.addUser.err = expectedErr
	template := &mockTemplate{}

	w := httptest.NewRecorder()
	r, _ := http.NewRequest("POST", "/path", nil)

	err := addUser(client, template)(client.requiredPermissions(), w, r)
	assert.Equal(expectedErr, err)

	assert.Equal(1, client.roles.count)
	assert.Equal(1, client.addUser.count)
	assert.Equal(0, template.count)
}

func TestPostAddUserRolesError(t *testing.T) {
	assert := assert.New(t)

	expectedErr := errors.New("oops")
	client := &mockAddUserClient{}
	client.roles.err = expectedErr
	template := &mockTemplate{}

	w := httptest.NewRecorder()
	r, _ := http.NewRequest("POST", "/path", nil)

	err := addUser(client, template)(client.requiredPermissions(), w, r)
	assert.Equal(expectedErr, err)

	assert.Equal(1, client.roles.count)
	assert.Equal(0, client.addUser.count)
	assert.Equal(0, template.count)
}
