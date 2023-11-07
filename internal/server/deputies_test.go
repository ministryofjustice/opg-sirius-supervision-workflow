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

type mockDeputiesClient struct {
	lastCtx               sirius.Context
	deputyList            sirius.DeputyList
	getDeputyListError    error
	reassignDeputiesError error
	count                 map[string]int
	lastDeputyListParams  sirius.DeputyListParams
}

func (m *mockDeputiesClient) GetDeputyList(ctx sirius.Context, params sirius.DeputyListParams) (sirius.DeputyList, error) {
	if m.count == nil {
		m.count = make(map[string]int)
	}
	m.count["GetClientList"] += 1
	m.lastCtx = ctx
	m.lastDeputyListParams = params
	return m.deputyList, m.getDeputyListError
}

func (m *mockDeputiesClient) ReassignDeputies(context sirius.Context, params sirius.ReassignDeputiesParams) (string, error) {
	if m.count == nil {
		m.count = make(map[string]int)
	}
	m.count["ReassignDeputies"] += 1
	m.lastCtx = context

	return "", m.reassignDeputiesError
}

var testDeputyList = sirius.DeputyList{
	Deputies: []model.Deputy{
		{
			Id:          1,
			DisplayName: "Test Deputy",
			Type:        model.RefData{Handle: "PRO"},
			Number:      14,
			Address:     model.Address{Town: "Derby"},
		},
	},
	Pages:              model.PageInformation{PageCurrent: 1, PageTotal: 1},
	TotalDeputies:      1,
	PaProTeamSelection: []model.Team{},
}

func TestDeputies(t *testing.T) {
	client := &mockDeputiesClient{
		deputyList: testDeputyList,
	}
	template := &mockTemplate{}

	w := httptest.NewRecorder()
	r, _ := http.NewRequest(http.MethodGet, "", nil)

	app := WorkflowVars{
		Path:         "test-path",
		SelectedTeam: model.Team{Id: 123, Type: "PRO", Selector: "1"},
		TeamSelection: []model.Team{
			{Id: 123, Type: "PRO", Selector: "1"},
			{Id: 222, Type: "PA", Selector: "1"},
			{Id: 333, Type: "LAY", Selector: "1"},
			{Id: 444, Type: "PRO", Selector: "1"},
		},
		EnvironmentVars: EnvironmentVars{ShowDeputies: true},
	}
	err := deputies(client, template)(app, w, r)

	assert.Nil(t, err)
	assert.Equal(t, 1, template.count)

	var want DeputiesPage
	want.App = app
	want.PerPage = 25
	want.DeputyList = testDeputyList
	want.DeputyList.PaProTeamSelection = []model.Team{
		{
			Id:       123,
			Type:     "PRO",
			Selector: "1",
		},
		{
			Id:       222,
			Type:     "PA",
			Selector: "1",
		},
		{
			Id:       444,
			Type:     "PRO",
			Selector: "1",
		},
	}
	want.Sort = urlbuilder.Sort{OrderBy: "deputy"}

	want.UrlBuilder = urlbuilder.UrlBuilder{
		Path:            "deputies",
		SelectedTeam:    app.SelectedTeam.Selector,
		SelectedPerPage: 25,
		SelectedSort:    want.Sort,
	}

	want.Pagination = paginate.Pagination{
		CurrentPage:     1,
		TotalPages:      1,
		TotalElements:   1,
		ElementsPerPage: 25,
		ElementName:     "deputies",
		PerPageOptions:  []int{25, 50, 100},
		UrlBuilder:      want.UrlBuilder,
	}

	assert.Equal(t, want, template.lastVars)
}

func TestDeputies_RedirectsToClientTasksForLayDeputies(t *testing.T) {
	client := &mockDeputiesClient{}
	template := &mockTemplate{}

	w := httptest.NewRecorder()
	r, _ := http.NewRequest(http.MethodGet, "", nil)

	app := WorkflowVars{
		Path:            "test-path",
		SelectedTeam:    model.Team{Type: "LAY", Selector: "19"},
		EnvironmentVars: EnvironmentVars{ShowDeputies: true},
	}
	err := deputies(client, template)(app, w, r)

	assert.Equal(t, RedirectError("client-tasks?team=19&page=1&per-page=25"), err)
	assert.Equal(t, 0, template.count)
}

func TestDeputies_MethodNotAllowed(t *testing.T) {
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
			client := &mockDeputiesClient{}
			template := &mockTemplate{}

			w := httptest.NewRecorder()
			r, _ := http.NewRequest(method, "", nil)

			app := WorkflowVars{}
			err := deputies(client, template)(app, w, r)

			assert.Equal(t, StatusError(http.StatusMethodNotAllowed), err)
			assert.Equal(t, 0, template.count)
		})
	}
}

func TestDeputies_NonExistentPageNumberWillRedirectToTheHighestExistingPageNumber(t *testing.T) {
	var testDeputyList = sirius.DeputyList{
		Deputies: []model.Deputy{{}},
		Pages: model.PageInformation{
			PageCurrent: 10,
			PageTotal:   2,
		},
	}

	client := &mockDeputiesClient{deputyList: testDeputyList}
	template := &mockTemplate{}

	w := httptest.NewRecorder()
	r, _ := http.NewRequest("GET", "/deputies?team=&page=10&per-page=25", nil)

	app := WorkflowVars{
		SelectedTeam: model.Team{Type: "PRO", Selector: "1"},
	}
	err := deputies(client, template)(app, w, r)

	assert.Equal(t, RedirectError("deputies?team=1&page=2&per-page=25&order-by=deputy&sort=asc"), err)
	assert.Equal(t, getContext(r), client.lastCtx)
	assert.Equal(t, 10, client.lastDeputyListParams.Page)
}

func TestDeputiesPage_CreateUrlBuilder(t *testing.T) {
	tests := []struct {
		page DeputiesPage
		want urlbuilder.UrlBuilder
	}{
		{
			page: DeputiesPage{},
			want: urlbuilder.UrlBuilder{Path: "deputies"},
		},
		{
			page: DeputiesPage{
				ListPage: ListPage{
					App: WorkflowVars{SelectedTeam: model.Team{Type: "PRO", Selector: "test-team"}},
				},
			},
			want: urlbuilder.UrlBuilder{Path: "deputies", SelectedTeam: "test-team"},
		},
		{
			page: DeputiesPage{
				ListPage: ListPage{
					App:     WorkflowVars{SelectedTeam: model.Team{Type: "PRO", Selector: "test-team"}},
					PerPage: 55,
					Sort:    urlbuilder.Sort{OrderBy: "test", Descending: true},
				},
			},
			want: urlbuilder.UrlBuilder{
				Path:            "deputies",
				SelectedTeam:    "test-team",
				SelectedPerPage: 55,
				SelectedSort:    urlbuilder.Sort{OrderBy: "test", Descending: true},
			},
		},
	}
	for i, test := range tests {
		t.Run("Scenario "+strconv.Itoa(i+1), func(t *testing.T) {
			assert.Equal(t, test.want, test.page.CreateUrlBuilder())
		})
	}
}

func TestListPaAndProDeputyTeams(t *testing.T) {
	allTeams := []model.Team{
		{Id: 29, Type: "PA"},
		{Id: 13, Type: "PRO"},
		{Id: 30, Type: "PA"},
		{Id: 5, Type: "PA"},
	}

	tests := []struct {
		name                string
		requiredTeamTypes   []string
		currentSelectedTeam model.Team
		expectedResponse    []model.Team
	}{
		{
			name:                "Can filter on multiple team types",
			requiredTeamTypes:   []string{"PA", "PRO"},
			currentSelectedTeam: model.Team{Type: "PRO", Id: 55},
			expectedResponse: []model.Team{
				{Type: "PRO", Id: 55},
				{Type: "PA", Id: 29},
				{Type: "PA", Id: 30},
				{Type: "PA", Id: 5},
				{Type: "PRO", Id: 13},
			},
		},
		{
			name:                "Can filter on multiple team types",
			requiredTeamTypes:   []string{"PRO"},
			currentSelectedTeam: model.Team{Type: "PRO", Id: 55},
			expectedResponse: []model.Team{
				{Type: "PRO", Id: 55},
				{Type: "PRO", Id: 13},
			},
		},
		{
			name:                "Can filter on single team types",
			requiredTeamTypes:   []string{"PA"},
			currentSelectedTeam: model.Team{Type: "PRO", Id: 55},
			expectedResponse: []model.Team{
				{Type: "PA", Id: 29},
				{Type: "PA", Id: 30},
				{Type: "PA", Id: 5},
			},
		},
		{
			name:                "Will not return current selected team if doesnt match type",
			requiredTeamTypes:   []string{""},
			currentSelectedTeam: model.Team{Type: "PRO", Id: 55},
			expectedResponse:    []model.Team{},
		},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			response := listPaAndProDeputyTeams(allTeams, test.requiredTeamTypes, test.currentSelectedTeam)
			assert.Equal(t, test.expectedResponse, response)
		})
	}
}
