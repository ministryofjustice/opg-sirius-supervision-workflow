package server

import (
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/ministryofjustice/opg-sirius-user-management/internal/sirius"
	"github.com/stretchr/testify/assert"
)

type mockUnlockUserClient struct {
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
}

func (m *mockUnlockUserClient) User(ctx sirius.Context, id int) (sirius.AuthUser, error) {
	m.user.count += 1
	m.user.lastCtx = ctx
	m.user.lastID = id

	return m.user.data, m.user.err
}

func (m *mockUnlockUserClient) EditUser(ctx sirius.Context, user sirius.AuthUser) error {
	m.editUser.count += 1
	m.editUser.lastCtx = ctx
	m.editUser.lastUser = user

	return m.editUser.err
}

func (m *mockUnlockUserClient) requiredPermissions() sirius.PermissionSet {
	return sirius.PermissionSet{"v1-users": sirius.PermissionGroup{Permissions: []string{"put"}}}
}

func TestGetUnlockUser(t *testing.T) {
	assert := assert.New(t)

	client := &mockUnlockUserClient{}
	client.user.data = sirius.AuthUser{Firstname: "test"}
	template := &mockTemplate{}

	w := httptest.NewRecorder()
	r, _ := http.NewRequest("GET", "/unlock-user/123", nil)

	err := unlockUser(client, template)(client.requiredPermissions(), w, r)
	assert.Nil(err)

	assert.Equal(1, client.user.count)
	assert.Equal(123, client.user.lastID)
	assert.Equal(0, client.editUser.count)

	assert.Equal(1, template.count)
	assert.Equal("page", template.lastName)
	assert.Equal(unlockUserVars{
		Path: "/unlock-user/123",
		User: client.user.data,
	}, template.lastVars)
}

func TestGetUnlockUserNoPermission(t *testing.T) {
	assert := assert.New(t)

	w := httptest.NewRecorder()
	r, _ := http.NewRequest("GET", "/path", nil)

	err := unlockUser(nil, nil)(sirius.PermissionSet{}, w, r)
	assert.Equal(StatusError(http.StatusForbidden), err)
}

func TestGetUnlockUserError(t *testing.T) {
	assert := assert.New(t)

	expectedError := errors.New("oops")

	client := &mockUnlockUserClient{}
	client.user.err = expectedError
	template := &mockTemplate{}

	w := httptest.NewRecorder()
	r, _ := http.NewRequest("GET", "/unlock-user/123", nil)

	err := unlockUser(client, template)(client.requiredPermissions(), w, r)
	assert.Equal(expectedError, err)

	assert.Equal(1, client.user.count)
	assert.Equal(123, client.user.lastID)
	assert.Equal(0, client.editUser.count)
}

func TestGetUnlockUserBadPath(t *testing.T) {
	for name, path := range map[string]string{
		"empty":       "/unlock-user/",
		"non-numeric": "/unlock-user/hello",
		"suffixed":    "/unlock-user/123/no",
	} {
		t.Run(name, func(t *testing.T) {
			assert := assert.New(t)

			client := &mockUnlockUserClient{}
			template := &mockTemplate{}

			w := httptest.NewRecorder()
			r, _ := http.NewRequest("GET", path, nil)

			err := unlockUser(client, template)(client.requiredPermissions(), w, r)
			assert.Equal(StatusError(http.StatusNotFound), err)

			assert.Equal(0, client.user.count)
			assert.Equal(0, client.editUser.count)
			assert.Equal(0, template.count)
		})
	}
}

func TestPostUnlockUser(t *testing.T) {
	assert := assert.New(t)

	client := &mockUnlockUserClient{}
	client.user.data = sirius.AuthUser{ID: 123, Email: "user@opgtest.com", Locked: true}
	template := &mockTemplate{}

	w := httptest.NewRecorder()
	r, _ := http.NewRequest("POST", "/unlock-user/123", nil)

	err := unlockUser(client, template)(client.requiredPermissions(), w, r)
	assert.Equal(RedirectError("/edit-user/123"), err)

	assert.Equal(1, client.editUser.count)
	assert.Equal(getContext(r), client.editUser.lastCtx)
	assert.Equal(sirius.AuthUser{
		ID:     123,
		Email:  "user@opgtest.com",
		Locked: false,
	}, client.editUser.lastUser)

	assert.Equal(1, client.user.count)
	assert.Equal(0, template.count)
}

func TestPostUnlockUserClientError(t *testing.T) {
	assert := assert.New(t)

	client := &mockUnlockUserClient{}
	client.editUser.err = sirius.ClientError("problem")
	template := &mockTemplate{}

	w := httptest.NewRecorder()
	r, _ := http.NewRequest("POST", "/unlock-user/123", nil)

	err := unlockUser(client, template)(client.requiredPermissions(), w, r)
	assert.Nil(err)

	assert.Equal(1, client.user.count)
	assert.Equal(1, client.editUser.count)

	assert.Equal(1, template.count)
	assert.Equal("page", template.lastName)
	assert.Equal(unlockUserVars{
		Path: "/unlock-user/123",
		User: client.user.data,
		Errors: sirius.ValidationErrors{
			"": {
				"": "problem",
			},
		},
	}, template.lastVars)
}

func TestPostUnlockUserOtherError(t *testing.T) {
	assert := assert.New(t)

	expectedErr := errors.New("oops")
	client := &mockUnlockUserClient{}
	client.editUser.err = expectedErr
	template := &mockTemplate{}

	w := httptest.NewRecorder()
	r, _ := http.NewRequest("POST", "/unlock-user/123", nil)

	err := unlockUser(client, template)(client.requiredPermissions(), w, r)
	assert.Equal(expectedErr, err)

	assert.Equal(1, client.user.count)
	assert.Equal(1, client.editUser.count)
	assert.Equal(0, template.count)
}

func TestPutUnlockUser(t *testing.T) {
	assert := assert.New(t)

	client := &mockUnlockUserClient{}
	w := httptest.NewRecorder()
	r, _ := http.NewRequest("PUT", "/unlock-user/123", nil)

	err := unlockUser(nil, nil)(client.requiredPermissions(), w, r)
	assert.Equal(StatusError(http.StatusMethodNotAllowed), err)
}
