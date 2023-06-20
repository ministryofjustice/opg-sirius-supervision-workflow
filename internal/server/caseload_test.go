package server

import (
	"github.com/ministryofjustice/opg-sirius-workflow/internal/sirius"
	"github.com/stretchr/testify/assert"
	"net/http"
	"net/http/httptest"
	"testing"
)

type mockCaseloadClient struct {
}

type caseloadURLFields struct {
	SelectedTeam string
}

func createCaseloadVars(fields caseloadURLFields) CaseloadVars {
	return CaseloadVars{
		App: WorkflowVars{
			SelectedTeam: sirius.Team{Selector: fields.SelectedTeam},
		},
	}
}

func TestCaseload(t *testing.T) {
	client := &mockCaseloadClient{}
	template := &mockTemplates{}

	w := httptest.NewRecorder()
	r, _ := http.NewRequest(http.MethodGet, "", nil)

	app := WorkflowVars{
		Path:            "test-path",
		SelectedTeam:    sirius.Team{Type: "LAY"},
		EnvironmentVars: EnvironmentVars{ShowCaseload: true},
	}
	err := caseload(client, template)(app, w, r)

	assert.Nil(t, err)
	assert.Equal(t, 1, template.count)
	assert.Equal(t, CaseloadVars{App: app}, template.lastVars)
}

func TestCaseload_RedirectsToClientTasksForNonLayDeputies(t *testing.T) {
	client := &mockCaseloadClient{}
	template := &mockTemplates{}

	w := httptest.NewRecorder()
	r, _ := http.NewRequest(http.MethodGet, "", nil)

	app := WorkflowVars{
		Path:            "test-path",
		SelectedTeam:    sirius.Team{Type: "PRO", Selector: "19"},
		EnvironmentVars: EnvironmentVars{ShowCaseload: true},
	}
	err := caseload(client, template)(app, w, r)

	assert.Equal(t, RedirectError("client-tasks?team=19&page=1&per-page=25"), err)
	assert.Equal(t, 0, template.count)
}

func TestCaseload_RedirectsToClientTasksWhenFeatureFlagIsOff(t *testing.T) {
	client := &mockCaseloadClient{}
	template := &mockTemplates{}

	w := httptest.NewRecorder()
	r, _ := http.NewRequest(http.MethodGet, "", nil)

	app := WorkflowVars{
		Path:            "test-path",
		SelectedTeam:    sirius.Team{Type: "LAY", Selector: "19"},
		EnvironmentVars: EnvironmentVars{ShowCaseload: false},
	}
	err := caseload(client, template)(app, w, r)

	assert.Equal(t, RedirectError("client-tasks?team=19&page=1&per-page=25"), err)
	assert.Equal(t, 0, template.count)
}

func TestCaseload_MethodNotAllowed(t *testing.T) {
	methods := []string{
		http.MethodPost,
		http.MethodConnect,
		http.MethodDelete,
		http.MethodHead,
		http.MethodOptions,
		http.MethodPatch,
		http.MethodPut,
		http.MethodTrace,
	}
	for _, method := range methods {
		t.Run("Test "+method, func(t *testing.T) {
			client := &mockCaseloadClient{}
			template := &mockTemplates{}

			w := httptest.NewRecorder()
			r, _ := http.NewRequest(method, "", nil)

			app := WorkflowVars{}
			err := caseload(client, template)(app, w, r)

			assert.Equal(t, StatusError(http.StatusMethodNotAllowed), err)
			assert.Equal(t, 0, template.count)
		})
	}
}

func TestCaseloadVars_GetTeamUrl(t *testing.T) {
	tests := []struct {
		name   string
		fields caseloadURLFields
		team   string
		want   string
	}{
		{
			name:   "Team is retained",
			fields: caseloadURLFields{SelectedTeam: "lay"},
			team:   "lay",
			want:   "?team=lay",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			w := createCaseloadVars(tt.fields)
			team := sirius.Team{Selector: tt.team}
			assert.Equalf(t, "caseload"+tt.want, w.GetTeamUrl(team), "GetTeamUrl(%v)", tt.team)
		})
	}
}
