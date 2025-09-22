package server

import (
	"github.com/ministryofjustice/opg-go-common/paginate"
	"github.com/ministryofjustice/opg-sirius-workflow/internal/model"
	"github.com/ministryofjustice/opg-sirius-workflow/internal/sirius"
	"github.com/ministryofjustice/opg-sirius-workflow/internal/urlbuilder"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"net/http"
	"net/http/httptest"
	"net/url"
	"strconv"
	"testing"
)

type mockDeputyTasksClient struct {
	mock.Mock
}

func (m *mockDeputyTasksClient) GetTaskTypes(ctx sirius.Context, params sirius.TaskTypesParams) ([]model.TaskType, error) {
	args := m.Called(ctx)
	return args.Get(0).([]model.TaskType), args.Error(1)
}

func (m *mockDeputyTasksClient) GetTaskList(ctx sirius.Context, params sirius.TaskListParams) (sirius.TaskList, error) {
	args := m.Called(ctx)
	return args.Get(0).(sirius.TaskList), args.Error(1)
}

func (m *mockDeputyTasksClient) ReassignTasks(ctx sirius.Context, params sirius.ReassignTasksParams) (string, error) {
	args := m.Called(ctx)
	return args.Get(0).(string), args.Error(1)
}

var workflowVars = WorkflowVars{
	MyDetails: model.Assignee{
		Id: 123,
	},
	Path:         "deputy-tasks",
	SelectedTeam: model.Team{Type: "PRO", Selector: "1"},
}

var taskType = []model.TaskType{
	{
		Handle:     "DEPT",
		Incomplete: "Test incomplete name",
		Category:   sirius.TaskTypeCategoryDeputy,
		Complete:   "Test complete name",
		User:       true,
	},
}
var taskList = sirius.TaskList{
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
	Pages: model.PageInformation{
		PageCurrent: 1,
		PageTotal:   2,
	},
}

func TestDeputyTasks(t *testing.T) {
	client := &mockDeputyTasksClient{}
	template := &mockTemplate{}

	client.On("GetTaskTypes", mock.Anything).Return(taskType, nil)
	client.On("GetTaskList", mock.Anything).Return(taskList, nil)

	w := httptest.NewRecorder()
	r, _ := http.NewRequest(http.MethodGet, "/deputy-tasks", nil)

	handler := deputyTasks(client, template)
	err := handler(workflowVars, w, r)

	assert.Nil(t, err)
	resp := w.Result()
	assert.Equal(t, http.StatusOK, resp.StatusCode)
	assert.Equal(t, 1, template.count)

	want := DeputyTasksPage{
		TaskList: taskList,
		ListPage: ListPage{
			App:     workflowVars,
			PerPage: 25,
			Pagination: paginate.Pagination{
				CurrentPage:     1,
				TotalPages:      2,
				TotalElements:   26,
				ElementsPerPage: 25,
				ElementName:     "deputy-tasks",
				PerPageOptions:  []int{25, 50, 100},
				UrlBuilder:      nil,
			},
		},
	}
	want.TaskTypes = taskType
	want.UrlBuilder = urlbuilder.UrlBuilder{
		Path:            "deputy-tasks",
		SelectedTeam:    workflowVars.SelectedTeam.Selector,
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
		CurrentPage:     1,
		TotalPages:      2,
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

	client.On("GetTaskTypes", mock.Anything).Return(taskType, nil)
	client.On("GetTaskList", mock.Anything).Return(taskList, nil)

	w := httptest.NewRecorder()
	r, _ := http.NewRequest(http.MethodGet, "/deputy-tasks?team=15&page=10&per-page=25", nil)

	app := WorkflowVars{
		Path:         "test-path",
		SelectedTeam: model.Team{Type: "LAY", Selector: "19"},
	}
	err := deputyTasks(client, template)(app, w, r)

	assert.Equal(t, Redirect{Path: "client-tasks?team=19&page=1&per-page=25"}, err)
	assert.Equal(t, 0, template.count)
}

func TestDeputyTasks_NonExistentPageNumberWillRedirectToTheHighestExistingPageNumber(t *testing.T) {
	client := &mockDeputyTasksClient{}
	template := &mockTemplate{}

	client.On("GetTaskTypes", mock.Anything).Return(taskType, nil)
	client.On("GetTaskList", mock.Anything).Return(taskList, nil)

	w := httptest.NewRecorder()
	r, _ := http.NewRequest("GET", "/deputy-tasks?team=&page=10&per-page=25", nil)

	err := deputyTasks(client, template)(workflowVars, w, r)

	assert.Equal(t, Redirect{Path: "deputy-tasks?team=1&page=2&per-page=25"}, err)
	assert.Equal(t, 0, template.count)
}

func TestDeputyTasks_ReassignTasks(t *testing.T) {
	client := &mockDeputyTasksClient{}
	template := &mockTemplate{}

	client.On("ReassignTasks", mock.Anything).Return("reassign success", nil)

	expectedParams := sirius.ReassignTasksParams{
		AssignTeam: "10",
		AssignCM:   "20",
		TaskIds:    []string{"1", "2"},
		IsPriority: "true",
	}

	w := httptest.NewRecorder()
	r, _ := http.NewRequest(http.MethodPost, "/deputy-tasks?team=1&page=2&per-page=25", nil)
	r.PostForm = url.Values{
		"assignTeam":     {expectedParams.AssignTeam},
		"assignCM":       {expectedParams.AssignCM},
		"selected-tasks": expectedParams.TaskIds,
		"priority":       {expectedParams.IsPriority},
	}

	err := deputyTasks(client, template)(workflowVars, w, r)
	assert.Equal(t, Redirect{
		Path:           "deputy-tasks?team=1&page=2&per-page=25",
		SuccessMessage: "reassign success",
	}, err)
	assert.Equal(t, 0, template.count)
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
