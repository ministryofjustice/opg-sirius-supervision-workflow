package server

import (
	"github.com/ministryofjustice/opg-go-common/paginate"
	"github.com/ministryofjustice/opg-sirius-workflow/internal/model"
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

func (m *mockCaseloadClient) ReassignClients(ctx sirius.Context, params sirius.ReassignClientsParams) (string, error) {
	if m.count == nil {
		m.count = make(map[string]int)
	}
	m.count["ReassignClients"] += 1
	m.lastCtx = ctx

	return "", m.err
}

func TestCaseload(t *testing.T) {
	tests := []struct {
		name                  string
		teamType              string
		wantDeputyTypes       []model.RefData
		wantCaseTypes         []model.RefData
		wantSupervisionLevels []model.RefData
	}{
		{
			name:     "Caseload page is viewable for Lay teams",
			teamType: "LAY",
			wantSupervisionLevels: []model.RefData{
				{
					Handle: "GENERAL",
					Label:  "General",
				},
				{
					Handle: "MINIMAL",
					Label:  "Minimal",
				},
			},
		},
		{
			name:     "Caseload page is viewable for Health & Welfare teams",
			teamType: "HW",
			wantDeputyTypes: []model.RefData{
				{
					Handle: "LAY",
					Label:  "Lay",
				},
				{
					Handle: "PRO",
					Label:  "Professional",
				},
				{
					Handle: "PA",
					Label:  "Public Authority",
				},
			},
			wantCaseTypes: []model.RefData{
				{
					Handle: "HYBRID",
					Label:  "Hybrid",
				},
				{
					Handle: "DUAL",
					Label:  "Dual",
				},
				{
					Handle: "HW",
					Label:  "Health and welfare",
				},
				{
					Handle: "PFA",
					Label:  "Property and financial affairs",
				},
			},
		},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			client := &mockCaseloadClient{}
			template := &mockTemplate{}

			w := httptest.NewRecorder()
			r, _ := http.NewRequest(http.MethodGet, "", nil)

			app := WorkflowVars{
				Path:         "test-path",
				SelectedTeam: model.Team{Type: test.teamType, Selector: "1"},
			}
			err := caseload(client, template)(app, w, r)

			assert.Nil(t, err)
			assert.Equal(t, 1, template.count)

			expectedClientListParams := sirius.ClientListParams{
				Team:    app.SelectedTeam,
				Page:    1,
				PerPage: 25,
			}
			if test.teamType == "HW" {
				expectedClientListParams.SubType = "hw"
			}

			assert.Equal(t, expectedClientListParams, client.lastClientListParams)

			var want CaseloadPage
			want.App = app
			want.PerPage = 25
			want.AssigneeFilterName = "Case owner"
			want.StatusOptions = []model.RefData{
				{
					Handle: "active",
					Label:  "Active",
				},
				{
					Handle: "closed",
					Label:  "Closed",
				},
			}
			want.DeputyTypes = test.wantDeputyTypes
			want.CaseTypes = test.wantCaseTypes
			want.SupervisionLevels = test.wantSupervisionLevels

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
					{
						Name:                  "status",
						ClearBetweenTeamViews: true,
					},
					{
						Name:                  "deputy-type",
						ClearBetweenTeamViews: true,
					},
					{
						Name:                  "case-type",
						ClearBetweenTeamViews: true,
					},
					{
						Name:                  "supervision-level",
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
		})
	}
}

func TestCaseload_RedirectsToClientTasksForNonLayNonHWTeams(t *testing.T) {
	client := &mockCaseloadClient{}
	template := &mockTemplate{}

	w := httptest.NewRecorder()
	r, _ := http.NewRequest(http.MethodGet, "", nil)

	app := WorkflowVars{
		Path:         "test-path",
		SelectedTeam: model.Team{Type: "PRO", Selector: "19"},
	}
	err := caseload(client, template)(app, w, r)

	assert.Equal(t, RedirectError("client-tasks?team=19&page=1&per-page=25"), err)
	assert.Equal(t, 0, template.count)
}

func TestCaseload_MethodNotAllowed(t *testing.T) {
	methods := []string{
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

func TestCaseloadPage_CreateUrlBuilder(t *testing.T) {
	expectedFilters := []urlbuilder.Filter{
		{
			Name:                  "assignee",
			ClearBetweenTeamViews: true,
		},
		{
			Name:                  "unassigned",
			ClearBetweenTeamViews: true,
		},
		{
			Name:                  "status",
			ClearBetweenTeamViews: true,
		},
		{
			Name:                  "deputy-type",
			ClearBetweenTeamViews: true,
		},
		{
			Name:                  "case-type",
			ClearBetweenTeamViews: true,
		},
		{
			Name:                  "supervision-level",
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

func TestCaseloadPage_GetAppliedFilters(t *testing.T) {
	tests := []struct {
		selectedAssignees   []string
		selectedUnassigned  string
		selectedStatuses    []string
		selectedDeputyTypes []string
		selectedCaseTypes   []string
		want                []string
	}{
		{
			want: nil,
		},
		{
			selectedAssignees: []string{"2"},
			want:              []string{"User 2"},
		},
		{
			selectedUnassigned: "lay-team",
			want:               []string{"Lay team"},
		},
		{
			selectedStatuses: []string{"active"},
			want:             []string{"Active"},
		},
		{
			selectedAssignees:  []string{"1", "2"},
			selectedUnassigned: "lay-team",
			selectedStatuses:   []string{"active", "closed"},
			want:               []string{"Lay team", "User 1", "User 2", "Active", "Closed"},
		},
		{
			selectedDeputyTypes: []string{"LAY", "PA"},
			want:                []string{"Lay", "Public Authority"},
		},
		{
			selectedCaseTypes: []string{"HYBRID", "DUAL"},
			want:              []string{"Hybrid", "Dual"},
		},
	}
	for i, test := range tests {
		t.Run("Scenario "+strconv.Itoa(i+1), func(t *testing.T) {
			var page CaseloadPage
			page.App.SelectedTeam = model.Team{
				Name:     "Lay team",
				Selector: "lay-team",
				Members: []model.Assignee{
					{
						Id:   1,
						Name: "User 1",
					},
					{
						Id:   2,
						Name: "User 2",
					},
				},
			}
			page.StatusOptions = []model.RefData{
				{
					Handle: "active",
					Label:  "Active",
				},
				{
					Handle: "closed",
					Label:  "Closed",
				},
			}
			page.DeputyTypes = []model.RefData{
				{
					Handle: "LAY",
					Label:  "Lay",
				},
				{
					Handle: "PRO",
					Label:  "Professional",
				},
				{
					Handle: "PA",
					Label:  "Public Authority",
				},
			}
			page.CaseTypes = []model.RefData{
				{
					Handle: "HYBRID",
					Label:  "Hybrid",
				},
				{
					Handle: "DUAL",
					Label:  "Dual",
				},
				{
					Handle: "HW",
					Label:  "Health and welfare",
				},
				{
					Handle: "PFA",
					Label:  "Property and financial affairs",
				},
			}
			page.SelectedAssignees = test.selectedAssignees
			page.SelectedUnassigned = test.selectedUnassigned
			page.SelectedStatuses = test.selectedStatuses
			page.SelectedDeputyTypes = test.selectedDeputyTypes
			page.SelectedCaseTypes = test.selectedCaseTypes

			assert.Equal(t, test.want, page.GetAppliedFilters())
		})
	}
}
