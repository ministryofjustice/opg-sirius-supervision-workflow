package server

import (
	"github.com/ministryofjustice/opg-sirius-workflow/internal/model"
	"github.com/ministryofjustice/opg-sirius-workflow/internal/sirius"
	"github.com/ministryofjustice/opg-sirius-workflow/internal/urlbuilder"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"
)

type mockBondsClient struct {
	mock.Mock
}

func (m *mockBondsClient) GetBondList(ctx sirius.Context, params sirius.BondListParams) (sirius.BondList, error) {
	args := m.Called(ctx)
	return args.Get(0).(sirius.BondList), args.Error(1)
}

var testBondList = sirius.BondList{
	Bonds: []model.Bond{
		{
			Id:                  1,
			CourtRef:            "12345678",
			FirstName:           "Joseph",
			LastName:            "Smith",
			CompanyName:         "Howden",
			BondReferenceNumber: "BOND1",
			BondAmount:          101,
			BondIssuedDate:      model.Date{Time: time.Now()},
		},
	},
}

func TestGetBonds(t *testing.T) {
	client := &mockBondsClient{}
	template := &mockTemplate{}

	client.On("GetBondList", mock.Anything).Return(testBondList, nil)

	w := httptest.NewRecorder()
	r, _ := http.NewRequest(http.MethodGet, "/bonds", nil)

	app := WorkflowVars{
		Path:         "test-path",
		SelectedTeam: model.Team{Id: 123, Type: "ALLOCATIONS", Selector: "1"},
		Teams: []model.Team{
			{Id: 123, Type: "ALLOCATIONS", Selector: "1"},
			{Id: 222, Type: "PA", Selector: "1"},
			{Id: 333, Type: "LAY", Selector: "1"},
			{Id: 444, Type: "PRO", Selector: "1"},
		},
		EnvironmentVars: EnvironmentVars{},
	}
	err := bonds(client, template)(app, w, r)

	assert.Nil(t, err)
	assert.Equal(t, 1, template.count)

	var want BondsPage
	want.App = app
	want.BondList = testBondList

	want.UrlBuilder = urlbuilder.UrlBuilder{
		Path:         "bonds",
		SelectedTeam: app.SelectedTeam.Selector,
	}

	assert.Equal(t, want, template.lastVars)
}

func TestBonds_RedirectsToClientTasksForNonAllocationsTeam(t *testing.T) {
	client := &mockBondsClient{}
	template := &mockTemplate{}

	client.On("GetBondList", mock.Anything).Return(testBondList, nil)

	w := httptest.NewRecorder()
	r, _ := http.NewRequest(http.MethodGet, "/bonds", nil)

	app := WorkflowVars{
		Path:            "test-path",
		SelectedTeam:    model.Team{Type: "LAY", Selector: "19"},
		EnvironmentVars: EnvironmentVars{},
	}
	err := bonds(client, template)(app, w, r)

	assert.Equal(t, Redirect{Path: "client-tasks?team=19&page=1&per-page=25"}, err)
	assert.Equal(t, 0, template.count)
}

func TestBonds_MethodNotAllowed(t *testing.T) {
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
			client := &mockBondsClient{}
			template := &mockTemplate{}

			w := httptest.NewRecorder()
			r, _ := http.NewRequest(method, "", nil)

			app := WorkflowVars{}
			err := bonds(client, template)(app, w, r)

			assert.Equal(t, StatusError(http.StatusMethodNotAllowed), err)
			assert.Equal(t, 0, template.count)
		})
	}
}
