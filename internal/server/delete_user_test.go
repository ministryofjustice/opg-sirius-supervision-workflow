package server

import (
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/ministryofjustice/opg-sirius-user-management/internal/sirius"
	"github.com/stretchr/testify/assert"
)

type mockDeleteUserClient struct {
	user struct {
		count   int
		lastCtx sirius.Context
		lastID  int
		data    sirius.AuthUser
		err     error
	}

	deleteUser struct {
		count      int
		lastCtx    sirius.Context
		lastUserID int
		err        error
	}
}

func (m *mockDeleteUserClient) User(ctx sirius.Context, id int) (sirius.AuthUser, error) {
	m.user.count += 1
	m.user.lastCtx = ctx
	m.user.lastID = id

	return m.user.data, m.user.err
}

func (m *mockDeleteUserClient) DeleteUser(ctx sirius.Context, userID int) error {
	m.deleteUser.count += 1
	m.deleteUser.lastCtx = ctx
	m.deleteUser.lastUserID = userID

	return m.deleteUser.err
}

func (m *mockDeleteUserClient) requiredPermissions() sirius.PermissionSet {
	return sirius.PermissionSet{"v1-users": sirius.PermissionGroup{Permissions: []string{"delete"}}}
}

func TestGetDeleteUser(t *testing.T) {
	assert := assert.New(t)

	client := &mockDeleteUserClient{}
	client.user.data = sirius.AuthUser{Firstname: "test"}
	template := &mockTemplate{}

	w := httptest.NewRecorder()
	r, _ := http.NewRequest("GET", "/delete-user/123", nil)

	err := deleteUser(client, template)(client.requiredPermissions(), w, r)
	assert.Nil(err)

	assert.Equal(1, client.user.count)
	assert.Equal(123, client.user.lastID)
	assert.Equal(0, client.deleteUser.count)

	assert.Equal(1, template.count)
	assert.Equal("page", template.lastName)
	assert.Equal(deleteUserVars{
		Path: "/delete-user/123",
		User: client.user.data,
	}, template.lastVars)
}

func TestGetDeleteUserNoPermission(t *testing.T) {
	assert := assert.New(t)

	w := httptest.NewRecorder()
	r, _ := http.NewRequest("GET", "/path", nil)

	err := deleteUser(nil, nil)(sirius.PermissionSet{}, w, r)
	assert.Equal(StatusError(http.StatusForbidden), err)
}

func TestGetDeleteUserError(t *testing.T) {
	assert := assert.New(t)

	expectedError := errors.New("oops")

	client := &mockDeleteUserClient{}
	client.user.err = expectedError
	template := &mockTemplate{}

	w := httptest.NewRecorder()
	r, _ := http.NewRequest("GET", "/delete-user/123", nil)

	err := deleteUser(client, template)(client.requiredPermissions(), w, r)
	assert.Equal(expectedError, err)

	assert.Equal(1, client.user.count)
	assert.Equal(123, client.user.lastID)
	assert.Equal(0, client.deleteUser.count)
}

func TestGetDeleteUserBadPath(t *testing.T) {
	for name, path := range map[string]string{
		"empty":       "/delete-user/",
		"non-numeric": "/delete-user/hello",
		"suffixed":    "/delete-user/123/no",
	} {
		t.Run(name, func(t *testing.T) {
			assert := assert.New(t)

			client := &mockDeleteUserClient{}
			template := &mockTemplate{}

			w := httptest.NewRecorder()
			r, _ := http.NewRequest("GET", path, nil)

			err := deleteUser(client, template)(client.requiredPermissions(), w, r)
			assert.Equal(StatusError(http.StatusNotFound), err)

			assert.Equal(0, client.user.count)
			assert.Equal(0, client.deleteUser.count)
			assert.Equal(0, template.count)
		})
	}
}

func TestPostDeleteUser(t *testing.T) {
	assert := assert.New(t)

	client := &mockDeleteUserClient{}
	client.user.data = sirius.AuthUser{Firstname: "test", Surname: "user", Email: "user@opgtest.com"}
	template := &mockTemplate{}

	w := httptest.NewRecorder()
	r, _ := http.NewRequest("POST", "/delete-user/123", nil)

	err := deleteUser(client, template)(client.requiredPermissions(), w, r)
	assert.Nil(err)

	assert.Equal(1, client.deleteUser.count)
	assert.Equal(getContext(r), client.deleteUser.lastCtx)
	assert.Equal(123, client.deleteUser.lastUserID)

	assert.Equal(1, client.user.count)
	assert.Equal(1, template.count)

	assert.Equal(deleteUserVars{
		Path:           "/delete-user/123",
		User:           client.user.data,
		SuccessMessage: "User test user (user@opgtest.com) was deleted.",
	}, template.lastVars)
}

func TestPostDeleteUserClientError(t *testing.T) {
	assert := assert.New(t)

	client := &mockDeleteUserClient{}
	client.deleteUser.err = sirius.ClientError("problem")
	template := &mockTemplate{}

	w := httptest.NewRecorder()
	r, _ := http.NewRequest("POST", "/delete-user/123", nil)

	err := deleteUser(client, template)(client.requiredPermissions(), w, r)
	assert.Nil(err)

	assert.Equal(1, client.deleteUser.count)
	assert.Equal(1, client.user.count)

	assert.Equal(1, template.count)
	assert.Equal("page", template.lastName)
	assert.Equal(deleteUserVars{
		Path: "/delete-user/123",
		User: client.user.data,
		Errors: sirius.ValidationErrors{
			"": {
				"": "problem",
			},
		},
	}, template.lastVars)
}

func TestPostDeleteUserOtherError(t *testing.T) {
	assert := assert.New(t)

	expectedErr := errors.New("oops")
	client := &mockDeleteUserClient{}
	client.deleteUser.err = expectedErr
	template := &mockTemplate{}

	w := httptest.NewRecorder()
	r, _ := http.NewRequest("POST", "/delete-user/123", nil)

	err := deleteUser(client, template)(client.requiredPermissions(), w, r)
	assert.Equal(expectedErr, err)

	assert.Equal(1, client.deleteUser.count)
	assert.Equal(1, client.user.count)
	assert.Equal(0, template.count)
}

func TestPutDeleteUser(t *testing.T) {
	assert := assert.New(t)

	client := &mockDeleteUserClient{}
	w := httptest.NewRecorder()
	r, _ := http.NewRequest("PUT", "/delete-user/123", nil)

	err := deleteUser(nil, nil)(client.requiredPermissions(), w, r)
	assert.Equal(StatusError(http.StatusMethodNotAllowed), err)
}
