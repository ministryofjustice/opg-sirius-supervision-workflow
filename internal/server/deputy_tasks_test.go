package server

import (
	"github.com/ministryofjustice/opg-go-common/paginate"
	"github.com/ministryofjustice/opg-sirius-workflow/internal/model"
	"github.com/ministryofjustice/opg-sirius-workflow/internal/sirius"
	"github.com/ministryofjustice/opg-sirius-workflow/internal/urlbuilder"
	"github.com/stretchr/testify/assert"
	"net/http"
	"net/http/httptest"
	"net/url"
	"strconv"
	"testing"
)

type mockDeputyTasksClient struct {
	count                   map[string]int
	lastCtx                 sirius.Context
	lastReassignTasksParams sirius.ReassignTasksParams
	err                     error
	taskTypeData            []model.TaskType
	taskListData            sirius.TaskList
}

func (m *mockDeputyTasksClient) GetTaskTypes(ctx sirius.Context, params sirius.TaskTypesParams) ([]model.TaskType, error) {
	if m.count == nil {
		m.count = make(map[string]int)
	}
	m.count["GetTaskTypes"] += 1
	m.lastCtx = ctx

	return m.taskTypeData, m.err
}

func (m *mockDeputyTasksClient) GetTaskList(ctx sirius.Context, params sirius.TaskListParams) (sirius.TaskList, error) {
	if m.count == nil {
		m.count = make(map[string]int)
	}
	m.count["GetTaskList"] += 1
	m.lastCtx = ctx

	return m.taskListData, m.err
}

func (m *mockDeputyTasksClient) ReassignTasks(ctx sirius.Context, params sirius.ReassignTasksParams) (string, error) {
	if m.count == nil {
		m.count = make(map[string]int)
	}
	m.count["ReassignTasks"] += 1
	m.lastReassignTasksParams = params
	m.lastCtx = ctx

	return "reassign success", m.err
}

var testDeputyTaskType = []model.TaskType{
	{
		Handle:     "DEPT",
		Incomplete: "Test incomplete name",
		Category:   sirius.TaskTypeCategoryDeputy,
		Complete:   "Test complete name",
		User:       true,
	},
}

var testDeputyTaskList = sirius.TaskList{
	Tasks: []model.Task{
		{
			Assignee: model.Assignee{
				Name: "Assignee Test",
			},
			Name:    "Test task",
			DueDate: "01/02/2021",
			Deputies: []model.Deputy{
				{
					Id:          1,
					DisplayName: "Test Deputy",
					Type:        model.RefData{Handle: "PRO"},
					Number:      14,
					Address:     model.Address{Town: "Derby"},
				},
			},
		},
	},
}

func TestDeputyTasks(t *testing.T) {
	client := &mockDeputyTasksClient{
		taskTypeData: testDeputyTaskType,
		taskListData: testDeputyTaskList,
	}
	template := &mockTemplate{}

	w := httptest.NewRecorder()
	r, _ := http.NewRequest(http.MethodGet, "", nil)

	app := WorkflowVars{
		Path:            "test-path",
		SelectedTeam:    model.Team{Type: "PRO", Selector: "1"},
		EnvironmentVars: EnvironmentVars{},
	}
	err := deputyTasks(client, template)(app, w, r)

	assert.Nil(t, err)
	assert.Equal(t, 1, template.count)

	var want DeputyTasksPage
	want.App = app
	want.PerPage = 25
	want.TaskTypes = testDeputyTaskType
	want.TaskList = testDeputyTaskList

	want.UrlBuilder = urlbuilder.UrlBuilder{
		Path:            "deputy-tasks",
		SelectedTeam:    app.SelectedTeam.Selector,
		SelectedPerPage: 25,
		SelectedFilters: []urlbuilder.Filter{
			{
				Name: "task-type",
			},
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
		ElementName:     "tasks",
		PerPageOptions:  []int{25, 50, 100},
		UrlBuilder:      want.UrlBuilder,
	}

	assert.Equal(t, want, template.lastVars)
}

func TestDeputyTasks_RedirectsToClientTasksForLayDeputies(t *testing.T) {
	client := &mockDeputyTasksClient{}
	template := &mockTemplate{}

	w := httptest.NewRecorder()
	r, _ := http.NewRequest(http.MethodGet, "", nil)

	app := WorkflowVars{
		Path:            "test-path",
		SelectedTeam:    model.Team{Type: "LAY", Selector: "19"},
		EnvironmentVars: EnvironmentVars{},
	}
	err := deputyTasks(client, template)(app, w, r)

	assert.Equal(t, RedirectError("client-tasks?team=19&page=1&per-page=25"), err)
	assert.Equal(t, 0, template.count)
}

func TestDeputyTasks_NonExistentPageNumberWillRedirectToTheHighestExistingPageNumber(t *testing.T) {
	var testDeputyTaskList = sirius.TaskList{
		Tasks: []model.Task{{}},
		Pages: model.PageInformation{
			PageCurrent: 10,
			PageTotal:   2,
		},
	}

	client := &mockDeputyTasksClient{taskTypeData: testDeputyTaskType, taskListData: testDeputyTaskList}
	template := &mockTemplate{}

	w := httptest.NewRecorder()
	r, _ := http.NewRequest("GET", "/deputy-tasks?team=&page=10&per-page=25", nil)

	app := WorkflowVars{
		SelectedTeam: model.Team{Type: "PRO", Selector: "1"},
	}
	err := deputyTasks(client, template)(app, w, r)

	assert.Equal(t, RedirectError("deputy-tasks?team=1&page=2&per-page=25"), err)
	assert.Equal(t, getContext(r), client.lastCtx)
	assert.Equal(t, 2, len(client.count))
	assert.Equal(t, 1, client.count["GetTaskList"])
}

func TestDeputyTasks_ReassignTasks(t *testing.T) {
	client := &mockDeputyTasksClient{taskTypeData: testDeputyTaskType, taskListData: testDeputyTaskList}
	template := &mockTemplate{}

	expectedParams := sirius.ReassignTasksParams{
		AssignTeam: "10",
		AssignCM:   "20",
		TaskIds:    []string{"1", "2"},
		IsPriority: "true",
	}

	w := httptest.NewRecorder()
	r, _ := http.NewRequest("POST", "", nil)
	r.PostForm = url.Values{
		"assignTeam":     {expectedParams.AssignTeam},
		"assignCM":       {expectedParams.AssignCM},
		"selected-tasks": expectedParams.TaskIds,
		"priority":       {expectedParams.IsPriority},
	}

	app := WorkflowVars{
		SelectedTeam: model.Team{Type: "PRO", Selector: "1"},
	}
	err := deputyTasks(client, template)(app, w, r)

	assert.Nil(t, err)
	assert.Equal(t, 1, client.count["ReassignTasks"])
	assert.Equal(t, expectedParams, client.lastReassignTasksParams)
	assert.Equal(t, "reassign success", template.lastVars.(DeputyTasksPage).App.SuccessMessage)
}

func TestDeputyTasks_MethodNotAllowed(t *testing.T) {
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
			client := &mockDeputyTasksClient{}
			template := &mockTemplate{}

			w := httptest.NewRecorder()
			r, _ := http.NewRequest(method, "", nil)

			app := WorkflowVars{}
			err := deputyTasks(client, template)(app, w, r)

			assert.Equal(t, StatusError(http.StatusMethodNotAllowed), err)
			assert.Equal(t, 0, template.count)
		})
	}
}

func TestDeputyTasksPage_CreateUrlBuilder(t *testing.T) {
	wantFilters := []urlbuilder.Filter{
		{Name: "task-type"},
		{Name: "assignee", ClearBetweenTeamViews: true},
		{Name: "unassigned", ClearBetweenTeamViews: true},
	}

	tests := []struct {
		page        DeputyTasksPage
		wantBuilder urlbuilder.UrlBuilder
		wantFilters []urlbuilder.Filter
	}{
		{
			page:        DeputyTasksPage{},
			wantBuilder: urlbuilder.UrlBuilder{Path: "deputy-tasks"},
			wantFilters: wantFilters,
		},
		{
			page: DeputyTasksPage{
				ListPage: ListPage{
					App: WorkflowVars{SelectedTeam: model.Team{Type: "PRO", Selector: "test-team"}},
				},
			},
			wantBuilder: urlbuilder.UrlBuilder{Path: "deputy-tasks", SelectedTeam: "test-team", SelectedFilters: wantFilters},
			wantFilters: wantFilters,
		},
		{
			page: DeputyTasksPage{
				ListPage: ListPage{
					App:     WorkflowVars{SelectedTeam: model.Team{Type: "PRO", Selector: "test-team"}},
					PerPage: 55,
				},
			},
			wantBuilder: urlbuilder.UrlBuilder{Path: "deputy-tasks", SelectedTeam: "test-team", SelectedPerPage: 55, SelectedFilters: wantFilters},
			wantFilters: wantFilters,
		},
		{
			page: DeputyTasksPage{
				ListPage: ListPage{
					App:     WorkflowVars{SelectedTeam: model.Team{Selector: "test-team"}},
					PerPage: 55,
				},
				FilterByTaskType: FilterByTaskType{
					SelectedTaskTypes: []string{"type1", "type2"},
				},
				FilterByAssignee: FilterByAssignee{
					SelectedAssignees:  []string{"user1", "user2"},
					SelectedUnassigned: "test-unassigned",
				},
			},
			wantBuilder: urlbuilder.UrlBuilder{Path: "deputy-tasks", SelectedTeam: "test-team", SelectedPerPage: 55, SelectedFilters: wantFilters},
			wantFilters: []urlbuilder.Filter{
				{
					Name:           "task-type",
					SelectedValues: []string{"type1", "type2"},
				},
				{
					Name:                  "assignee",
					SelectedValues:        []string{"user1", "user2"},
					ClearBetweenTeamViews: true,
				},
				{
					Name:                  "unassigned",
					SelectedValues:        []string{"test-unassigned"},
					ClearBetweenTeamViews: true,
				},
			},
		},
	}
	for i, test := range tests {
		t.Run("Scenario "+strconv.Itoa(i+1), func(t *testing.T) {
			test.wantBuilder.SelectedFilters = test.wantFilters
			assert.Equal(t, test.wantBuilder, test.page.CreateUrlBuilder())
		})
	}
}

func TestDeputyTasksPage_GetAppliedFilters(t *testing.T) {
	tests := []struct {
		taskTypes          []model.TaskType
		selectedTaskTypes  []string
		selectedAssignees  []string
		selectedUnassigned string
		want               []string
	}{
		{
			want: nil,
		},
		{
			taskTypes: []model.TaskType{
				{Incomplete: "TaskType1", Handle: "TT1"},
				{Incomplete: "TaskType2", Handle: "TT2"},
				{Incomplete: "TaskType3", Handle: "TT3"},
			},
			selectedTaskTypes: []string{"TT1", "TT3"},
			want:              []string{"TaskType1", "TaskType3"},
		},
		{
			selectedAssignees: []string{"2"},
			want:              []string{"User 2"},
		},
		{
			selectedUnassigned: "pro-team",
			want:               []string{"Pro team"},
		},
		{
			taskTypes:          []model.TaskType{{Incomplete: "TaskType1", Handle: "TT1"}},
			selectedTaskTypes:  []string{"TT1"},
			selectedAssignees:  []string{"1"},
			selectedUnassigned: "pro-team",
			want:               []string{"TaskType1", "Pro team", "User 1"},
		},
	}
	for i, test := range tests {
		t.Run("Scenario "+strconv.Itoa(i+1), func(t *testing.T) {
			var page DeputyTasksPage
			page.App.SelectedTeam = model.Team{
				Name:     "Pro team",
				Selector: "pro-team",
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
			page.TaskTypes = test.taskTypes
			page.SelectedTaskTypes = test.selectedTaskTypes
			page.SelectedAssignees = test.selectedAssignees
			page.SelectedUnassigned = test.selectedUnassigned

			assert.Equal(t, test.want, page.GetAppliedFilters())
		})
	}
}
