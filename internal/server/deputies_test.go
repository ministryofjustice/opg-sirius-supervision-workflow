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
	lastCtx            sirius.Context
	lastParams         sirius.DeputyListParams
	deputyList         sirius.DeputyList
	getDeputyListError error
}

func (m *mockDeputiesClient) GetDeputyList(ctx sirius.Context, params sirius.DeputyListParams) (sirius.DeputyList, error) {
	m.lastCtx = ctx
	m.lastParams = params
	return m.deputyList, m.getDeputyListError
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
	Pages:         model.PageInformation{PageCurrent: 1, PageTotal: 1},
	TotalDeputies: 1,
}

func TestDeputies(t *testing.T) {
	client := &mockDeputiesClient{
		deputyList: testDeputyList,
	}
	template := &mockTemplate{}

	w := httptest.NewRecorder()
	r, _ := http.NewRequest(http.MethodGet, "", nil)

	app := WorkflowVars{
		Path:            "test-path",
		SelectedTeam:    model.Team{Type: "PRO", Selector: "1"},
		EnvironmentVars: EnvironmentVars{ShowDeputies: true},
	}
	err := deputies(client, template)(app, w, r)

	assert.Nil(t, err)
	assert.Equal(t, 1, template.count)

	var want DeputiesPage
	want.App = app
	want.PerPage = 25
	want.DeputyList = testDeputyList
	want.Sort = urlbuilder.Sort{OrderBy: "deputy"}

	want.UrlBuilder = urlbuilder.UrlBuilder{
		Path:            "deputies",
		SelectedTeam:    app.SelectedTeam.Selector,
		SelectedPerPage: 25,
		SelectedSort:    want.Sort,
		SelectedFilters: []urlbuilder.Filter{
			{
				Name:                  "ecm",
				ClearBetweenTeamViews: true,
			},
		},
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
		http.MethodPost,
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
	assert.Equal(t, 10, client.lastParams.Page)
}

func TestDeputiesPage_CreateUrlBuilder(t *testing.T) {
	filters := []urlbuilder.Filter{
		{Name: "ecm", ClearBetweenTeamViews: true},
	}

	tests := []struct {
		page DeputiesPage
		want urlbuilder.UrlBuilder
	}{
		{
			page: DeputiesPage{},
			want: urlbuilder.UrlBuilder{Path: "deputies", SelectedFilters: filters},
		},
		{
			page: DeputiesPage{
				ListPage: ListPage{
					App: WorkflowVars{SelectedTeam: model.Team{Type: "PRO", Selector: "test-team"}},
				},
			},
			want: urlbuilder.UrlBuilder{Path: "deputies", SelectedTeam: "test-team", SelectedFilters: filters},
		},
		{
			page: DeputiesPage{
				ListPage: ListPage{
					App:     WorkflowVars{SelectedTeam: model.Team{Type: "PRO", Selector: "test-team"}},
					PerPage: 55,
					Sort:    urlbuilder.Sort{OrderBy: "test", Descending: true},
				},
				FilterByECM: FilterByECM{SelectedECMs: []string{"1", "2"}},
			},
			want: urlbuilder.UrlBuilder{
				Path:            "deputies",
				SelectedTeam:    "test-team",
				SelectedPerPage: 55,
				SelectedSort:    urlbuilder.Sort{OrderBy: "test", Descending: true},
				SelectedFilters: []urlbuilder.Filter{{Name: "ecm", SelectedValues: []string{"1", "2"}, ClearBetweenTeamViews: true}},
			},
		},
	}
	for i, test := range tests {
		t.Run("Scenario "+strconv.Itoa(i+1), func(t *testing.T) {
			assert.Equal(t, test.want, test.page.CreateUrlBuilder())
		})
	}
}
