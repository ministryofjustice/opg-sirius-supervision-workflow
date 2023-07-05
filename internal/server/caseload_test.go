package server

import (
	"github.com/ministryofjustice/opg-sirius-workflow/internal/model"
	"github.com/ministryofjustice/opg-sirius-workflow/internal/paginate"
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
	SelectedTeam   string
	Status         string
	DueDate        string
	ClientsPerPage int
}

func createCaseloadVars(fields caseloadURLFields) CaseloadVars {
	return CaseloadVars{
		App: WorkflowVars{
			SelectedTeam: model.Team{Selector: fields.SelectedTeam},
		},
		ClientsPerPage: fields.ClientsPerPage,
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

func (m *mockCaseloadClient) GetClientList(ctx sirius.Context, team model.Team, clientsPerPage int, page int) (sirius.ClientList, error) {
	if m.count == nil {
		m.count = make(map[string]int)
	}
	m.count["GetClientList"] += 1
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

	want := CaseloadVars{App: app, ClientsPerPage: 25}

	want.Pagination = paginate.Pagination{
		CurrentPage:     0,
		TotalPages:      0,
		TotalElements:   0,
		ElementsPerPage: 25,
		ElementName:     "clients",
		PerPageOptions:  []int{25, 50, 100},
		UrlBuilder:      want,
	}

	assert.Equal(t, want, template.lastVars)
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
			fields: caseloadURLFields{SelectedTeam: "lay", ClientsPerPage: 25},
			team:   "lay",
			want:   "?team=lay&page=1&per-page=25",
		},
		{
			name:   "Per page limit is retained",
			fields: caseloadURLFields{SelectedTeam: "lay", ClientsPerPage: 50},
			team:   "lay",
			want:   "?team=lay&page=1&per-page=50",
		},
		{
			name:   "Per page limit defaults to 25",
			fields: caseloadURLFields{SelectedTeam: "lay", ClientsPerPage: 0},
			team:   "lay",
			want:   "?team=lay&page=1&per-page=25",
		},
		{
			name:   "Page is reset back to 1",
			fields: caseloadURLFields{SelectedTeam: "lay", ClientsPerPage: 25},
			team:   "pro",
			want:   "?team=pro&page=1&per-page=25",
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

func createCaseloadVars(fields caseloadURLFields) CaseloadVars {
	return CaseloadVars{
		App: WorkflowVars{
			SelectedTeam: model.Team{Selector: "lay"},
		},
		ClientList:     sirius.ClientList{},
		Pagination:     paginate.Pagination{},
		ClientsPerPage: fields.ClientsPerPage,
	}
}

func TestCaseloadVars_GetPaginationUrl(t *testing.T) {
	type args struct {
		page    int
		perPage int
	}
	tests := []struct {
		name   string
		args   args
		fields caseloadURLFields
		want   string
	}{
		{
			name: "Page number is updated",
			args: args{page: 2, perPage: 25},
			fields: caseloadURLFields{
				SelectedTeam:   "lay",
				ClientsPerPage: 25,
			},
			want: "?team=lay&page=2&per-page=25",
		},
		{
			name: "Per page limit is updated",
			args: args{page: 1, perPage: 50},
			fields: caseloadURLFields{
				SelectedTeam:   "lay",
				ClientsPerPage: 25,
			},
			want: "?team=lay&page=1&per-page=50",
		},
		{
			name: "Per page limit is retained",
			args: args{page: 2},
			fields: caseloadURLFields{
				SelectedTeam:   "lay",
				ClientsPerPage: 100,
			},
			want: "?team=lay&page=2&per-page=100",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			w := createClientVars(tt.fields)
			var result string
			if tt.args.perPage == 0 {
				result = w.GetPaginationUrl(tt.args.page)
			} else {
				result = w.GetPaginationUrl(tt.args.page, tt.args.perPage)
			}
			assert.Equalf(t, "caseload"+tt.want, result, "GetPaginationUrl(%v, %v)", tt.args.page, tt.args.perPage)
		})
	}
}
