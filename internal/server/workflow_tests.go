package server

import (
	"context"
	"errors"
	"io/ioutil"
	"log"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/ministryofjustice/opg-sirius-user-management/internal/sirius"
	"github.com/stretchr/testify/assert"
)

type mockMyDetailsClient struct {
	mockAuthenticateClient
	count       int
	lastCookies []*http.Cookie
	err         error
	data        sirius.MyDetails
}

func (m *mockMyDetailsClient) MyDetails(ctx context.Context, cookies []*http.Cookie) (sirius.MyDetails, error) {
	m.count += 1
	m.lastCookies = cookies

	return m.data, m.err
}

func TestGetMyDetails(t *testing.T) {
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
	client := &mockMyDetailsClient{data: data}
	templates := &mockTemplates{}

	w := httptest.NewRecorder()
	r, _ := http.NewRequest("GET", "", nil)
	r.AddCookie(&http.Cookie{Name: "test", Value: "val"})

	myDetails(nil, client, templates).ServeHTTP(w, r)

	resp := w.Result()
	assert.Equal(http.StatusOK, resp.StatusCode)
	assert.Equal(r.Cookies(), client.lastCookies)

	assert.Equal(1, templates.count)
	assert.Equal("my-details.gotmpl", templates.lastName)
	assert.Equal(myDetailsVars{
		Path:         "",
		ID:           123,
		Firstname:    "John",
		Surname:      "Doe",
		Email:        "john@doe.com",
		PhoneNumber:  "123",
		Organisation: "COP User",
		Roles:        []string{"A", "B"},
		Teams:        []string{"A Team"},
	}, templates.lastVars)
}

func TestGetMyDetailsUnauthenticated(t *testing.T) {
	assert := assert.New(t)

	client := &mockMyDetailsClient{err: sirius.ErrUnauthorized}
	templates := &mockTemplates{}

	w := httptest.NewRecorder()
	r, _ := http.NewRequest("GET", "", nil)

	myDetails(nil, client, templates).ServeHTTP(w, r)

	resp := w.Result()
	assert.Equal(http.StatusOK, resp.StatusCode)
	assert.Equal(0, templates.count)
	assert.True(client.authenticated)
}

func TestGetMyDetailsSiriusErrors(t *testing.T) {
	assert := assert.New(t)

	logger := log.New(ioutil.Discard, "", 0)
	client := &mockMyDetailsClient{err: errors.New("err")}
	templates := &mockTemplates{}

	w := httptest.NewRecorder()
	r, _ := http.NewRequest("GET", "", nil)

	myDetails(logger, client, templates).ServeHTTP(w, r)

	resp := w.Result()
	assert.Equal(http.StatusInternalServerError, resp.StatusCode)
	assert.Equal(0, templates.count)
}

func TestPostMyDetails(t *testing.T) {
	assert := assert.New(t)
	templates := &mockTemplates{}

	w := httptest.NewRecorder()
	r, _ := http.NewRequest("POST", "", nil)

	myDetails(nil, nil, templates).ServeHTTP(w, r)

	resp := w.Result()
	assert.Equal(http.StatusMethodNotAllowed, resp.StatusCode)
	assert.Equal(0, templates.count)
}
