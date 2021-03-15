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

type mockEditMyDetailsClient struct {
	count         int
	saveCount     int
	lastCtx       sirius.Context
	lastRequest   string
	err           error
	errSave       error
	data          sirius.MyDetails
	lastArguments struct {
		ID          int
		PhoneNumber string
	}
}

func (m *mockEditMyDetailsClient) MyDetails(ctx sirius.Context) (sirius.MyDetails, error) {
	m.count += 1
	m.lastCtx = ctx
	m.lastRequest = "MyDetails"

	return m.data, m.err
}

func (m *mockEditMyDetailsClient) EditMyDetails(ctx sirius.Context, id int, phoneNumber string) error {
	m.saveCount += 1
	m.lastCtx = ctx
	m.lastRequest = "EditMyDetails"
	m.lastArguments.ID = id
	m.lastArguments.PhoneNumber = phoneNumber

	return m.errSave
}

func (m *mockEditMyDetailsClient) requiredPermissions() sirius.PermissionSet {
	return sirius.PermissionSet{"v1-users-updatetelephonenumber": sirius.PermissionGroup{Permissions: []string{"put"}}}
}

func TestGetEditMyDetails(t *testing.T) {
	assert := assert.New(t)

	data := sirius.MyDetails{
		ID:          123,
		Firstname:   "John",
		Surname:     "Doe",
		Email:       "john@doe.com",
		PhoneNumber: "123",
		Roles:       []string{"A", "COP User", "B"},
		Teams: []sirius.MyDetailsTeam{
			{DisplayName: "A Team"},
		},
	}
	client := &mockEditMyDetailsClient{data: data}
	template := &mockTemplate{}

	w := httptest.NewRecorder()
	r, _ := http.NewRequest("GET", "/path", nil)

	handler := editMyDetails(client, template)
	err := handler(client.requiredPermissions(), w, r)

	assert.Nil(err)

	resp := w.Result()
	assert.Equal(http.StatusOK, resp.StatusCode)
	assert.Equal(getContext(r), client.lastCtx)

	assert.Equal(1, template.count)
	assert.Equal("page", template.lastName)
	assert.Equal(editMyDetailsVars{
		Path:        "/path",
		PhoneNumber: "123",
	}, template.lastVars)
}

func TestGetEditMyDetailsUnauthenticated(t *testing.T) {
	assert := assert.New(t)

	client := &mockEditMyDetailsClient{err: sirius.ErrUnauthorized}
	template := &mockTemplate{}

	w := httptest.NewRecorder()
	r, _ := http.NewRequest("GET", "", nil)

	handler := editMyDetails(client, template)
	err := handler(client.requiredPermissions(), w, r)

	assert.Equal(sirius.ErrUnauthorized, err)

	assert.Equal(0, template.count)
}

func TestGetEditMyDetailsNotPermitted(t *testing.T) {
	assert := assert.New(t)

	client := &mockEditMyDetailsClient{}
	template := &mockTemplate{}

	w := httptest.NewRecorder()
	r, _ := http.NewRequest("GET", "", nil)

	handler := editMyDetails(client, template)
	err := handler(sirius.PermissionSet{}, w, r)

	assert.Equal(StatusError(http.StatusForbidden), err)

	assert.Equal(0, template.count)
}

func TestGetEditMyDetailsSiriusErrors(t *testing.T) {
	assert := assert.New(t)

	client := &mockEditMyDetailsClient{err: errors.New("err")}
	template := &mockTemplate{}

	w := httptest.NewRecorder()
	r, _ := http.NewRequest("GET", "", nil)

	handler := editMyDetails(client, template)
	err := handler(client.requiredPermissions(), w, r)

	assert.Equal("err", err.Error())

	assert.Equal(0, template.count)
}

func TestPostEditMyDetails(t *testing.T) {
	assert := assert.New(t)

	client := &mockEditMyDetailsClient{
		data: sirius.MyDetails{
			ID: 31,
		},
	}

	template := &mockTemplate{}

	w := httptest.NewRecorder()
	r, _ := http.NewRequest("POST", "/path", strings.NewReader("phonenumber=0189202"))
	r.Header.Add("Content-type", "application/x-www-form-urlencoded")

	handler := editMyDetails(client, template)
	err := handler(client.requiredPermissions(), w, r)

	assert.Nil(err)

	assert.Equal(1, client.count)
	assert.Equal(1, client.saveCount)

	assert.Equal(getContext(r), client.lastCtx)
	assert.Equal("EditMyDetails", client.lastRequest)
	assert.Equal(31, client.lastArguments.ID)
	assert.Equal("0189202", client.lastArguments.PhoneNumber)

	assert.Equal(1, template.count)
	assert.Equal("page", template.lastName)
	assert.Equal(editMyDetailsVars{
		Path:        "/path",
		Success:     true,
		PhoneNumber: "0189202",
	}, template.lastVars)
}

func TestPostEditMyDetailsUnauthenticated(t *testing.T) {
	assert := assert.New(t)

	client := &mockEditMyDetailsClient{errSave: sirius.ErrUnauthorized}
	template := &mockTemplate{}

	w := httptest.NewRecorder()
	r, _ := http.NewRequest("POST", "/path", strings.NewReader("phonenumber=0189202"))

	handler := editMyDetails(client, template)
	err := handler(client.requiredPermissions(), w, r)

	assert.Equal(sirius.ErrUnauthorized, err)

	assert.Equal(1, client.count)
	assert.Equal(1, client.saveCount)

	assert.Equal(0, template.count)
}

func TestPostEditMyDetailsSiriusErrors(t *testing.T) {
	assert := assert.New(t)

	client := &mockEditMyDetailsClient{errSave: errors.New("err")}
	template := &mockTemplate{}

	w := httptest.NewRecorder()
	r, _ := http.NewRequest("POST", "/path", strings.NewReader("phonenumber=0189202"))
	r.Header.Add("Content-type", "application/x-www-form-urlencoded")

	handler := editMyDetails(client, template)
	err := handler(client.requiredPermissions(), w, r)

	assert.Equal("err", err.Error())

	assert.Equal(1, client.count)
	assert.Equal(1, client.saveCount)

	assert.Equal(0, template.count)
}

func TestPostEditMyDetailsInvalidRequest(t *testing.T) {
	assert := assert.New(t)

	validationError := &sirius.ValidationError{
		Errors: sirius.ValidationErrors{
			"phoneNumber": {
				"invalidNumber": "Phone number is not in valid format",
			},
		},
	}

	client := &mockEditMyDetailsClient{errSave: validationError}
	template := &mockTemplate{}

	w := httptest.NewRecorder()
	r, _ := http.NewRequest("POST", "/path", strings.NewReader("phonenumber=invalid+phone+number"))
	r.Header.Add("Content-type", "application/x-www-form-urlencoded")

	handler := editMyDetails(client, template)
	err := handler(client.requiredPermissions(), w, r)

	assert.Nil(err)

	resp := w.Result()
	assert.Equal(http.StatusBadRequest, resp.StatusCode)

	assert.Equal(1, client.count)
	assert.Equal(1, client.saveCount)

	assert.Equal(getContext(r), client.lastCtx)
	assert.Equal("EditMyDetails", client.lastRequest)

	assert.Equal(1, template.count)
	assert.Equal(editMyDetailsVars{
		Path:        "/path",
		PhoneNumber: "invalid phone number",
		Errors: map[string]map[string]string{
			"phoneNumber": {
				"invalidNumber": "Phone number is not in valid format",
			},
		},
	}, template.lastVars)
}
