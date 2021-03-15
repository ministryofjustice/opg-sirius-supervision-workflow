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

type mockChangePasswordClient struct {
	count                  int
	lastCtx                sirius.Context
	lastExistingPassword   string
	lastNewPassword        string
	lastNewPasswordConfirm string
	err                    error
}

func (m *mockChangePasswordClient) ChangePassword(ctx sirius.Context, existingPassword, newPassword, newPasswordConfirm string) error {
	m.count += 1
	m.lastCtx = ctx
	m.lastExistingPassword = existingPassword
	m.lastNewPassword = newPassword
	m.lastNewPasswordConfirm = newPasswordConfirm

	return m.err
}

func TestGetChangePassword(t *testing.T) {
	assert := assert.New(t)

	template := &mockTemplate{}

	w := httptest.NewRecorder()
	r, _ := http.NewRequest("GET", "/path", nil)

	handler := changePassword(nil, template)
	err := handler(sirius.PermissionSet{}, w, r)

	assert.Nil(err)

	assert.Equal(1, template.count)
	assert.Equal("page", template.lastName)
	assert.Equal(changePasswordVars{
		Path: "/path",
	}, template.lastVars)
}

func TestPostChangePassword(t *testing.T) {
	assert := assert.New(t)

	client := &mockChangePasswordClient{}
	template := &mockTemplate{}

	w := httptest.NewRecorder()
	r, _ := http.NewRequest("POST", "/path", strings.NewReader("currentpassword=a&password1=b&password2=c"))
	r.Header.Add("Content-Type", "application/x-www-form-urlencoded")

	handler := changePassword(client, template)
	err := handler(sirius.PermissionSet{}, w, r)

	assert.Nil(err)

	assert.Equal(getContext(r), client.lastCtx)
	assert.Equal("a", client.lastExistingPassword)
	assert.Equal("b", client.lastNewPassword)
	assert.Equal("c", client.lastNewPasswordConfirm)

	assert.Equal("page", template.lastName)
	assert.Equal(changePasswordVars{
		Path:    "/path",
		Success: true,
	}, template.lastVars)
}

func TestPostChangePasswordUnauthenticated(t *testing.T) {
	assert := assert.New(t)

	client := &mockChangePasswordClient{err: sirius.ErrUnauthorized}
	template := &mockTemplate{}

	w := httptest.NewRecorder()
	r, _ := http.NewRequest("POST", "/path", nil)

	handler := changePassword(client, template)
	err := handler(sirius.PermissionSet{}, w, r)

	assert.Equal(sirius.ErrUnauthorized, err)

	assert.Equal(0, template.count)
}

func TestPostChangePasswordSiriusError(t *testing.T) {
	assert := assert.New(t)

	client := &mockChangePasswordClient{err: sirius.ClientError("Something happened")}
	template := &mockTemplate{}

	w := httptest.NewRecorder()
	r, _ := http.NewRequest("POST", "/path", nil)

	handler := changePassword(client, template)
	err := handler(sirius.PermissionSet{}, w, r)

	assert.Nil(err)

	resp := w.Result()
	assert.Equal(http.StatusBadRequest, resp.StatusCode)

	assert.Equal(1, template.count)
	assert.Equal("page", template.lastName)
	assert.Equal(changePasswordVars{
		Path: "/path",
		Errors: sirius.ValidationErrors{
			"currentpassword": {
				"": "Something happened",
			},
		},
	}, template.lastVars)
}

func TestPostChangePasswordOtherError(t *testing.T) {
	assert := assert.New(t)

	expectedErr := errors.New("oops")
	client := &mockChangePasswordClient{err: expectedErr}
	template := &mockTemplate{}

	w := httptest.NewRecorder()
	r, _ := http.NewRequest("POST", "/path", nil)

	handler := changePassword(client, template)
	err := handler(sirius.PermissionSet{}, w, r)

	assert.Equal(expectedErr, err)

	assert.Equal(0, template.count)
}
