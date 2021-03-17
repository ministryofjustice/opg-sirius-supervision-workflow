package server

import (
	"io"
	"net/http"
)

type mockAuthenticateClient struct {
	authenticated bool
}

func (m *mockAuthenticateClient) Authenticate(w http.ResponseWriter, r *http.Request) {
	m.authenticated = true
}

type mockTemplates struct {
	count    int
	lastName string
	lastVars interface{}
}

func (m *mockTemplates) ExecuteTemplate(w io.Writer, name string, vars interface{}) error {
	m.count += 1
	m.lastName = name
	m.lastVars = vars
	return nil
}
