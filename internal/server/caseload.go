package server

import (
	"fmt"
	"net/http"
)

type CaseloadClient interface {
}

type CaseloadVars struct {
	App WorkflowVars
}

func caseload(client CaseloadClient, tmpl Template) Handler {
	return func(app WorkflowVars, w http.ResponseWriter, r *http.Request) error {
		if r.Method != http.MethodGet {
			return StatusError(http.StatusMethodNotAllowed)
		}

		vars := CaseloadVars{
			App: app,
		}

		return tmpl.Execute(w, vars)
	}
}

func (cv CaseloadVars) buildUrl(team string) string {
	return fmt.Sprintf("caseload?team=%s", team)
}

func (cv CaseloadVars) GetTeamUrl(team string) string {
	return cv.buildUrl(team)
}
