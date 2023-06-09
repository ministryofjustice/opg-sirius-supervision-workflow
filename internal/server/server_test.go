package server

import (
	"io"
)

type mockTemplates struct {
	count    int
	lastVars interface{}
}

func (m *mockTemplates) Execute(w io.Writer, vars interface{}) error {
	m.count += 1
	m.lastVars = vars
	return nil
}
