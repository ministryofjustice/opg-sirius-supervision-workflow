package server

import (
	"github.com/ministryofjustice/opg-go-common/paginate"
	"github.com/ministryofjustice/opg-sirius-workflow/internal/model"
	"github.com/ministryofjustice/opg-sirius-workflow/internal/urlbuilder"
	"github.com/stretchr/testify/assert"
	"net/http"
	"net/http/httptest"
	"strconv"
	"testing"
)

type mockDeputiesClient struct {
}

func TestDeputies(t *testing.T) {
	client := &mockDeputiesClient{}
	template := &mockTemplate{}

	w := httptest.NewRecorder()
	r, _ := http.NewRequest(http.MethodGet, "", nil)

	app := WorkflowVars{
		Path:            "test-path",
		SelectedTeam:    model.Team{Type: "PRO", Selector: "1"},
		EnvironmentVars: EnvironmentVars{},
	}
	err := deputies(client, template)(app, w, r)

	assert.Nil(t, err)
	assert.Equal(t, 1, template.count)

	var want DeputiesPage
	want.App = app
	want.PerPage = 25

	want.UrlBuilder = urlbuilder.UrlBuilder{
		Path:            "deputies",
		SelectedTeam:    app.SelectedTeam.Selector,
		SelectedPerPage: 25,
	}

	want.Pagination = paginate.Pagination{
		CurrentPage:     1,
		TotalPages:      0,
		TotalElements:   0,
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
		EnvironmentVars: EnvironmentVars{},
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
				},
			},
			want: urlbuilder.UrlBuilder{Path: "deputies", SelectedTeam: "test-team", SelectedPerPage: 55},
		},
	}
	for i, test := range tests {
		t.Run("Scenario "+strconv.Itoa(i+1), func(t *testing.T) {
			assert.Equal(t, test.want, test.page.CreateUrlBuilder())
		})
	}
}
