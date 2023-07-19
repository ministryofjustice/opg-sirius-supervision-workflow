package server

import (
	"github.com/ministryofjustice/opg-sirius-workflow/internal/model"
	"github.com/ministryofjustice/opg-sirius-workflow/internal/paginate"
	"github.com/ministryofjustice/opg-sirius-workflow/internal/sirius"
	"github.com/ministryofjustice/opg-sirius-workflow/internal/urlbuilder"
	"github.com/stretchr/testify/assert"
	"net/http"
	"net/http/httptest"
	"strconv"
	"testing"
)

type mockCaseloadClient struct {
	count                map[string]int
	lastCtx              sirius.Context
	lastClientListParams sirius.ClientListParams
	err                  error
	clientList           sirius.ClientList
}

func (m *mockCaseloadClient) GetClientList(ctx sirius.Context, params sirius.ClientListParams) (sirius.ClientList, error) {
	if m.count == nil {
		m.count = make(map[string]int)
	}
	m.count["GetClientList"] += 1
	m.lastCtx = ctx
	m.lastClientListParams = params

	return m.clientList, m.err
}

func TestCaseload(t *testing.T) {
	client := &mockCaseloadClient{}
	template := &mockTemplate{}

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

	expectedClientListParams := sirius.ClientListParams{
		Team:    app.SelectedTeam,
		Page:    1,
		PerPage: 25,
	}
	assert.Equal(t, expectedClientListParams, client.lastClientListParams)

	var want CaseloadPage
	want.App = app
	want.PerPage = 25

	want.UrlBuilder = urlbuilder.UrlBuilder{
		Path:            "caseload",
		SelectedTeam:    app.SelectedTeam.Selector,
		SelectedPerPage: 25,
		SelectedFilters: []urlbuilder.Filter{
			{
				Name:                  "assignee",
				ClearBetweenTeamViews: true,
			},
			{
				Name:                  "unassigned",
				ClearBetweenTeamViews: true,
			},
		},
	}

	want.Pagination = paginate.Pagination{
		CurrentPage:     0,
		TotalPages:      0,
		TotalElements:   0,
		ElementsPerPage: 25,
		ElementName:     "clients",
		PerPageOptions:  []int{25, 50, 100},
		UrlBuilder:      want.UrlBuilder,
	}

	assert.Equal(t, want, template.lastVars)
}

func TestCaseload_RedirectsToClientTasksForNonLayDeputies(t *testing.T) {
	client := &mockCaseloadClient{}
	template := &mockTemplate{}

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
			template := &mockTemplate{}

			w := httptest.NewRecorder()
			r, _ := http.NewRequest(method, "", nil)

			app := WorkflowVars{}
			err := caseload(client, template)(app, w, r)

			assert.Equal(t, StatusError(http.StatusMethodNotAllowed), err)
			assert.Equal(t, 0, template.count)
		})
	}
}

func TestCaseloadVars_CreateUrlBuilder(t *testing.T) {
	expectedFilters := []urlbuilder.Filter{
		{
			Name:                  "assignee",
			ClearBetweenTeamViews: true,
		},
		{
			Name:                  "unassigned",
			ClearBetweenTeamViews: true,
		},
	}

	tests := []struct {
		page CaseloadPage
		want urlbuilder.UrlBuilder
	}{
		{
			page: CaseloadPage{},
			want: urlbuilder.UrlBuilder{Path: "caseload"},
		},
		{
			page: CaseloadPage{
				ListPage: ListPage{
					App: WorkflowVars{SelectedTeam: model.Team{Selector: "test-team"}},
				},
			},
			want: urlbuilder.UrlBuilder{Path: "caseload", SelectedTeam: "test-team"},
		},
		{
			page: CaseloadPage{
				ListPage: ListPage{
					App:     WorkflowVars{SelectedTeam: model.Team{Selector: "test-team"}},
					PerPage: 55,
				},
			},
			want: urlbuilder.UrlBuilder{Path: "caseload", SelectedTeam: "test-team", SelectedPerPage: 55},
		},
	}
	for i, test := range tests {
		t.Run("Scenario "+strconv.Itoa(i+1), func(t *testing.T) {
			test.want.SelectedFilters = expectedFilters
			assert.Equal(t, test.want, test.page.CreateUrlBuilder())
		})
	}
}
