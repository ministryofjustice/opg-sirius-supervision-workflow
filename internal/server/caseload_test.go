package server

import (
	"github.com/ministryofjustice/opg-sirius-workflow/internal/model"
	"github.com/ministryofjustice/opg-sirius-workflow/internal/sirius"
	"github.com/stretchr/testify/assert"
	"net/http"
	"net/http/httptest"
	"testing"
)

type mockCaseloadClient struct {
	count      map[string]int
	lastCtx    sirius.Context
	err        error
	clientList sirius.ClientList
}

type caseloadURLFields struct {
	SelectedTeam string
	Status       string
	DueDate      string
}

func createCaseloadVars(fields caseloadURLFields) CaseloadVars {
	return CaseloadVars{
		App: WorkflowVars{
			SelectedTeam: model.Team{Selector: fields.SelectedTeam},
		},
		ClientList: sirius.ClientList{
			Clients: []model.Client{
				{
					Orders: []model.Order{
						{
							LatestAnnualReport: model.AnnualReport{DueDate: fields.DueDate},
							Status:             model.RefData{Label: fields.Status},
						},
					},
				},
			},
		},
	}
}

func (m *mockCaseloadClient) GetClientList(ctx sirius.Context, team model.Team) (sirius.ClientList, error) {
	if m.count == nil {
		m.count = make(map[string]int)
	}
	m.count["GetTaskList"] += 1
	m.lastCtx = ctx

	return m.clientList, m.err
}

func TestCaseload(t *testing.T) {
	client := &mockCaseloadClient{}
	template := &mockTemplates{}

	w := httptest.NewRecorder()
	r, _ := http.NewRequest(http.MethodGet, "", nil)

	app := WorkflowVars{
		Path:            "test-path",
		SelectedTeam:    model.Team{Type: "LAY"},
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
		SelectedTeam:    model.Team{Type: "PRO", Selector: "19"},
		EnvironmentVars: EnvironmentVars{ShowCaseload: true},
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
			team := model.Team{Selector: tt.team}
			assert.Equalf(t, "caseload"+tt.want, w.GetTeamUrl(team), "GetTeamUrl(%v)", tt.team)
		})
	}
}
