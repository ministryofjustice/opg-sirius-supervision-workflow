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
	count           int
	lastCtx         sirius.Context
	err             error
	userdetailsdata sirius.UserDetails
	//taskdetailsdata []sirius.ApiTaskTypes
	tasklistdetailsdata []sirius.ApiTask
}

func (m *mockUserDetailsClient) SiriusUserDetails(ctx sirius.Context) (sirius.UserDetails, error) {
	m.count += 1
	m.lastCtx = ctx

	return m.userdetailsdata, m.err
}

// func (c *mockUserDetailsClient) GetTaskDetails(ctx sirius.Context) ([]sirius.ApiTaskTypes, error) {
// 	c.count += 1
// 	c.lastCtx = ctx

// 	return c.taskdetailsdata, c.err
// }

func (d *mockUserDetailsClient) GetTaskList(ctx sirius.Context) ([]sirius.ApiTask, error) {
	d.count += 1
	d.lastCtx = ctx

	return d.tasklistdetailsdata, d.err
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

// func TestGetTaskTypes(t *testing.T) {
// 	assert := assert.New(t)

// 	data := []sirius.ApiTaskTypes{
// 		{
// 			Handle:     "TestHandle",
// 			Incomplete: "TestIncomplete",
// 			Category:   "TestCategory",
// 			Complete:   "TestComplete",
// 			User:       true,
// 		},
// 	}
// 	client := &mockUserDetailsClient{taskdetailsdata: data}
// 	template := &mockTemplates{}

// 	w := httptest.NewRecorder()
// 	r, _ := http.NewRequest("GET", "", nil)

// 	handler := loggingInfoForWorflow(client, template)
// 	err := handler(sirius.PermissionSet{}, w, r)
// 	assert.Nil(err)

// 	resp := w.Result()
// 	assert.Equal(http.StatusOK, resp.StatusCode)
// 	assert.Equal(getContext(r), client.lastCtx)

// 	assert.Equal(1, template.count)
// 	assert.Equal("page", template.lastName)
// 	assert.Equal(userDetailsVars{
// 		Path: "",
// 		// LoadTasks: data,
// 	}, template.lastVars)
// }

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
