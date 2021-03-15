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

type mockEditUserClient struct {
	user struct {
		count   int
		lastCtx sirius.Context
		lastID  int
		data    sirius.AuthUser
		err     error
	}

	editUser struct {
		count    int
		lastCtx  sirius.Context
		lastUser sirius.AuthUser
		err      error
	}

	roles struct {
		count   int
		lastCtx sirius.Context
		err     error
	}
}

func (m *mockEditUserClient) User(ctx sirius.Context, id int) (sirius.AuthUser, error) {
	m.user.count += 1
	m.user.lastCtx = ctx
	m.user.lastID = id

	return m.user.data, m.user.err
}

func (m *mockEditUserClient) EditUser(ctx sirius.Context, user sirius.AuthUser) error {
	m.editUser.count += 1
	m.editUser.lastCtx = ctx
	m.editUser.lastUser = user

	return m.editUser.err
}

func (m *mockEditUserClient) Roles(ctx sirius.Context) ([]string, error) {
	m.roles.count += 1
	m.roles.lastCtx = ctx

	return []string{"System Admin", "Manager"}, m.roles.err
}

func (m *mockEditUserClient) requiredPermissions() sirius.PermissionSet {
	return sirius.PermissionSet{"v1-users": sirius.PermissionGroup{Permissions: []string{"put"}}}
}

func TestGetEditUser(t *testing.T) {
	assert := assert.New(t)

	client := &mockEditUserClient{}
	client.user.data = sirius.AuthUser{Firstname: "test"}
	template := &mockTemplate{}

	w := httptest.NewRecorder()
	r, _ := http.NewRequest("GET", "/edit-user/123", nil)

	err := editUser(client, template, false)(client.requiredPermissions(), w, r)
	assert.Nil(err)

	assert.Equal(1, client.roles.count)
	assert.Equal(getContext(r), client.roles.lastCtx)
	assert.Equal(1, client.user.count)
	assert.Equal(123, client.user.lastID)
	assert.Equal(0, client.editUser.count)

	assert.Equal(1, template.count)
	assert.Equal("page", template.lastName)
	assert.Equal(editUserVars{
		Path:  "/edit-user/123",
		User:  client.user.data,
		Roles: []string{"System Admin", "Manager"},
	}, template.lastVars)
}

func TestGetEditUserNoPermission(t *testing.T) {
	assert := assert.New(t)

	w := httptest.NewRecorder()
	r, _ := http.NewRequest("GET", "/path", nil)

	err := editUser(nil, nil, false)(sirius.PermissionSet{}, w, r)
	assert.Equal(StatusError(http.StatusForbidden), err)
}

func TestGetEditUserDeleteEnabled(t *testing.T) {
	assert := assert.New(t)

	client := &mockEditUserClient{}
	client.user.data = sirius.AuthUser{Firstname: "test"}
	template := &mockTemplate{}

	w := httptest.NewRecorder()
	r, _ := http.NewRequest("GET", "/edit-user/123", nil)

	err := editUser(client, template, true)(client.requiredPermissions(), w, r)
	assert.Nil(err)

	assert.Equal(1, client.roles.count)
	assert.Equal(1, client.user.count)
	assert.Equal(123, client.user.lastID)
	assert.Equal(0, client.editUser.count)

	assert.Equal(1, template.count)
	assert.Equal("page", template.lastName)
	assert.Equal(editUserVars{
		Path:              "/edit-user/123",
		User:              client.user.data,
		Roles:             []string{"System Admin", "Manager"},
		DeleteUserEnabled: true,
	}, template.lastVars)
}

func TestGetEditUserBadPath(t *testing.T) {
	for name, path := range map[string]string{
		"empty":       "/edit-user/",
		"non-numeric": "/edit-user/hello",
		"suffixed":    "/edit-user/123/no",
	} {
		t.Run(name, func(t *testing.T) {
			assert := assert.New(t)

			client := &mockEditUserClient{}
			w := httptest.NewRecorder()
			r, _ := http.NewRequest("GET", path, nil)

			err := editUser(nil, nil, false)(client.requiredPermissions(), w, r)
			assert.Equal(StatusError(http.StatusNotFound), err)
		})
	}
}

func TestPostEditUser(t *testing.T) {
	assert := assert.New(t)

	client := &mockEditUserClient{}
	client.user.data = sirius.AuthUser{Firstname: "test"}
	template := &mockTemplate{}

	w := httptest.NewRecorder()
	r, _ := http.NewRequest("POST", "/edit-user/123", strings.NewReader("email=a&firstname=b&surname=c&organisation=d&roles=e&roles=f&locked=Yes&suspended=No"))
	r.Header.Add("Content-Type", "application/x-www-form-urlencoded")

	err := editUser(client, template, false)(client.requiredPermissions(), w, r)
	assert.Nil(err)

	assert.Equal(1, client.roles.count)
	assert.Equal(getContext(r), client.roles.lastCtx)

	assert.Equal(1, client.editUser.count)
	assert.Equal(getContext(r), client.editUser.lastCtx)
	assert.Equal(sirius.AuthUser{
		ID:           123,
		Firstname:    "b",
		Surname:      "c",
		Organisation: "d",
		Roles:        []string{"e", "f"},
		Locked:       true,
		Suspended:    false,
	}, client.editUser.lastUser)

	assert.Equal(0, client.user.count)

	assert.Equal(1, template.count)
	assert.Equal("page", template.lastName)
	assert.Equal(editUserVars{
		Path:    "/edit-user/123",
		Success: true,
		Roles:   []string{"System Admin", "Manager"},
		User: sirius.AuthUser{
			ID:           123,
			Email:        "a",
			Firstname:    "b",
			Surname:      "c",
			Organisation: "d",
			Roles:        []string{"e", "f"},
			Locked:       true,
			Suspended:    false,
		},
	}, template.lastVars)
}

func TestPostEditUserClientError(t *testing.T) {
	assert := assert.New(t)

	client := &mockEditUserClient{}
	client.editUser.err = sirius.ClientError("something")
	template := &mockTemplate{}

	w := httptest.NewRecorder()
	r, _ := http.NewRequest("POST", "/edit-user/123", strings.NewReader("email=a&firstname=b&surname=c&organisation=d&roles=e&roles=f&locked=Yes&suspended=No"))
	r.Header.Add("Content-Type", "application/x-www-form-urlencoded")

	err := editUser(client, template, false)(client.requiredPermissions(), w, r)
	assert.Nil(err)

	assert.Equal(1, client.roles.count)
	assert.Equal(1, client.editUser.count)
	assert.Equal(0, client.user.count)

	assert.Equal(1, template.count)
	assert.Equal("page", template.lastName)
	assert.Equal(editUserVars{
		Path:  "/edit-user/123",
		Roles: []string{"System Admin", "Manager"},
		User: sirius.AuthUser{
			ID:           123,
			Firstname:    "b",
			Surname:      "c",
			Organisation: "d",
			Roles:        []string{"e", "f"},
			Locked:       true,
			Suspended:    false,
		},
		Errors: sirius.ValidationErrors{
			"firstname": {
				"": "something",
			},
		},
	}, template.lastVars)
}

func TestPostEditUserOtherError(t *testing.T) {
	assert := assert.New(t)

	expectedErr := errors.New("oops")
	client := &mockEditUserClient{}
	client.editUser.err = expectedErr
	template := &mockTemplate{}

	w := httptest.NewRecorder()
	r, _ := http.NewRequest("POST", "/edit-user/123", nil)

	err := editUser(client, template, false)(client.requiredPermissions(), w, r)
	assert.Equal(expectedErr, err)

	assert.Equal(1, client.roles.count)
	assert.Equal(1, client.editUser.count)
	assert.Equal(0, client.user.count)
	assert.Equal(0, template.count)
}

func TestPostEditUserRolesError(t *testing.T) {
	assert := assert.New(t)

	expectedErr := errors.New("oops")
	client := &mockEditUserClient{}
	client.roles.err = expectedErr
	template := &mockTemplate{}

	w := httptest.NewRecorder()
	r, _ := http.NewRequest("POST", "/edit-user/123", nil)

	err := editUser(client, template, false)(client.requiredPermissions(), w, r)
	assert.Equal(expectedErr, err)

	assert.Equal(1, client.roles.count)
	assert.Equal(0, client.editUser.count)
	assert.Equal(0, client.user.count)
	assert.Equal(0, template.count)
}
