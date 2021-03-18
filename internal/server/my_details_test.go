package server

import (
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/ministryofjustice/opg-sirius-workflow/internal/sirius"
	"github.com/stretchr/testify/assert"
)

type mockUserDetailsClient struct {
	count   int
	lastCtx sirius.Context
	err     error
	data    sirius.UserDetails
}

func (m *mockUserDetailsClient) SiriusUserDetails(ctx sirius.Context) (sirius.UserDetails, error) {
	m.count += 1
	m.lastCtx = ctx

	return m.data, m.err
}

func TestGetMyDetails(t *testing.T) {
	assert := assert.New(t)

	data := sirius.UserDetails{
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
	client := &mockUserDetailsClient{data: data}
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
		Path:         "",
		ID:           123,
		Firstname:    "John",
		Surname:      "Doe",
		Email:        "john@doe.com",
		PhoneNumber:  "123",
		Organisation: "COP User",
		Roles:        []string{"A", "B"},
		Teams:        []string{"A Team"},
	}, template.lastVars)
}

// func TestGetMyDetailsUnauthenticated(t *testing.T) {
// 	assert := assert.New(t)

// 	client := &mockMyDetailsClient{err: sirius.ErrUnauthorized}
// 	templates := &mockTemplates{}

// 	w := httptest.NewRecorder()
// 	r, _ := http.NewRequest("GET", "", nil)

// 	loggingInfoForWorflow(nil, client, templates).ServeHTTP(w, r)

// 	resp := w.Result()
// 	assert.Equal(http.StatusOK, resp.StatusCode)
// 	assert.Equal(0, templates.count)
// 	assert.True(client.authenticated)
// }

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

// func TestPostMyDetails(t *testing.T) {
// 	assert := assert.New(t)
// 	templates := &mockTemplates{}

// 	w := httptest.NewRecorder()
// 	r, _ := http.NewRequest("POST", "", nil)

// 	myDetails(nil, nil, templates).ServeHTTP(w, r)

// 	resp := w.Result()
// 	assert.Equal(http.StatusMethodNotAllowed, resp.StatusCode)
// 	assert.Equal(0, templates.count)
// }
