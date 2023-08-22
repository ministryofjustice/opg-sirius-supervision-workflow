package server

import (
	"errors"
	"github.com/ministryofjustice/opg-sirius-workflow/internal/model"
	"github.com/ministryofjustice/opg-sirius-workflow/internal/paginate"
	"github.com/ministryofjustice/opg-sirius-workflow/internal/sirius"
	"github.com/ministryofjustice/opg-sirius-workflow/internal/urlbuilder"
	"github.com/stretchr/testify/assert"
	"net/http"
	"net/http/httptest"
	"strconv"
	"testing"
	"time"
)

type mockClientTasksClient struct {
	count        map[string]int
	lastCtx      sirius.Context
	err          error
	taskTypeData []model.TaskType
	taskListData sirius.TaskList
}

func (m *mockClientTasksClient) GetTaskTypes(ctx sirius.Context, params sirius.TaskTypesParams) ([]model.TaskType, error) {
	if m.count == nil {
		m.count = make(map[string]int)
	}
	m.count["GetTaskTypes"] += 1
	m.lastCtx = ctx

	return m.taskTypeData, m.err
}

func (m *mockClientTasksClient) GetTaskList(ctx sirius.Context, params sirius.TaskListParams) (sirius.TaskList, error) {
	if m.count == nil {
		m.count = make(map[string]int)
	}
	m.count["GetTaskList"] += 1
	m.lastCtx = ctx

	return m.taskListData, m.err
}

func (m *mockClientTasksClient) AssignTasksToCaseManager(ctx sirius.Context, newAssigneeId int, selectedTask []string, prioritySelected string) (string, error) {
	if m.count == nil {
		m.count = make(map[string]int)
	}
	m.count["AssignTasksToCaseManager"] += 1
	m.lastCtx = ctx

	return "", m.err
}

var testTaskType = []model.TaskType{
	{
		Handle:     "CDFC",
		Incomplete: "Correspondence - Review failed draft",
		Category:   sirius.TaskTypeCategorySupervision,
		Complete:   "Correspondence - Reviewed draft failure",
		User:       true,
	},
}

var testTaskList = sirius.TaskList{
	Tasks: []model.Task{
		{
			Assignee: model.Assignee{
				Name: "Assignee Duke Clive Henry Hetley Junior Jones",
			},
			Name:    "Case work - General",
			DueDate: "01/02/2021",
			Orders: []model.Order{
				{
					Client: model.Client{
						CaseRecNumber: "caseRecNumber",
						FirstName:     "Client Alexander Zacchaeus",
						Id:            3333,
						SupervisionCaseOwner: model.Assignee{
							Name: "Supervision - Team - Name",
						},
						Surname: "Client Wolfeschlegelsteinhausenbergerdorff",
					},
				},
			},
		},
	},
}

func TestClientTasks(t *testing.T) {
	client := &mockClientTasksClient{
		taskTypeData: testTaskType,
		taskListData: testTaskList,
	}
	template := &mockTemplate{}

	w := httptest.NewRecorder()
	r, _ := http.NewRequest(http.MethodGet, "", nil)

	app := WorkflowVars{
		Path:            "test-path",
		SelectedTeam:    model.Team{Type: "LAY", Selector: "test-team"},
		EnvironmentVars: EnvironmentVars{ShowCaseload: true},
	}
	err := clientTasks(client, template)(app, w, r)

	assert.Nil(t, err)
	assert.Equal(t, 1, template.count)

	var want ClientTasksPage
	want.App = app
	want.PerPage = 25
	want.TaskTypes = testTaskType
	want.TaskList = testTaskList

	want.UrlBuilder = urlbuilder.UrlBuilder{
		Path:            "client-tasks",
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
			{
				Name: "due-date-from",
			},
			{
				Name: "due-date-to",
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

func TestClientTasks_NonExistentPageNumberWillRedirectToTheHighestExistingPageNumber(t *testing.T) {
	assert := assert.New(t)

	var mockTaskListData = sirius.TaskList{
		Tasks: []model.Task{{}},
		Pages: model.PageInformation{
			PageCurrent: 10,
			PageTotal:   2,
		},
	}

	client := &mockClientTasksClient{taskTypeData: testTaskType, taskListData: mockTaskListData}
	template := &mockTemplate{}

	w := httptest.NewRecorder()
	r, _ := http.NewRequest("GET", "/client-tasks?team=&page=10&per-page=25", nil)

	app := WorkflowVars{}
	err := clientTasks(client, template)(app, w, r)

	assert.Equal(RedirectError("client-tasks?team=&page=2&per-page=25"), err)
	assert.Equal(getContext(r), client.lastCtx)
	assert.Equal(2, len(client.count))
	assert.Equal(1, client.count["GetTaskList"])
}

func TestClientTasks_Unauthorized(t *testing.T) {
	assert := assert.New(t)

	client := &mockClientTasksClient{err: sirius.ErrUnauthorized}
	template := &mockTemplate{}

	w := httptest.NewRecorder()
	r, _ := http.NewRequest("GET", "", nil)

	app := WorkflowVars{}
	err := clientTasks(client, template)(app, w, r)

	assert.Equal(sirius.ErrUnauthorized, err)
	assert.Equal(0, template.count)
}

func TestClientTasks_SiriusErrors(t *testing.T) {
	assert := assert.New(t)

	client := &mockClientTasksClient{err: errors.New("err")}
	template := &mockTemplate{}

	w := httptest.NewRecorder()
	r, _ := http.NewRequest("GET", "", nil)

	app := WorkflowVars{}
	err := clientTasks(client, template)(app, w, r)

	assert.Equal("err", err.Error())
	assert.Equal(0, template.count)
}

func TestClientTasks_PostIsPermitted(t *testing.T) {
	client := &mockClientTasksClient{taskTypeData: testTaskType, taskListData: testTaskList}
	template := &mockTemplate{}

	w := httptest.NewRecorder()
	r, _ := http.NewRequest("POST", "", nil)

	app := WorkflowVars{}
	err := clientTasks(client, template)(app, w, r)

	assert.Nil(t, err)
}

func TestGetAssigneeIdForTask(t *testing.T) {
	expectedAssigneeId, expectedError := getAssigneeIdForTask("13", "67")
	assert.Equal(t, expectedAssigneeId, 67)
	assert.Nil(t, expectedError)

	expectedAssigneeId, expectedError = getAssigneeIdForTask("13", "")
	assert.Equal(t, expectedAssigneeId, 13)
	assert.Nil(t, expectedError)

	expectedAssigneeId, expectedError = getAssigneeIdForTask("", "")
	assert.Equal(t, expectedAssigneeId, 0)
	assert.Nil(t, expectedError)
}

func TestGetSelectedDateFilter(t *testing.T) {
	testDate := time.Date(2022, 12, 17, 0, 0, 0, 0, time.Local)

	tests := []struct {
		name         string
		value        string
		expectedDate *time.Time
		expectedErr  error
	}{
		{
			name:         "Valid date",
			value:        "2022-12-17",
			expectedDate: &testDate,
			expectedErr:  nil,
		},
		{
			name:         "Blank date",
			value:        "",
			expectedDate: nil,
			expectedErr:  nil,
		},
		{
			name:         "Invalid date",
			value:        "17/12/2022",
			expectedDate: nil,
			expectedErr:  errors.New("parsing time \"17/12/2022\" as \"2006-01-02\": cannot parse \"17/12/2022\" as \"2006\""),
		},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			date, err := getSelectedDateFilter(test.value)

			if test.expectedErr == nil {
				assert.Nil(t, err)
			} else {
				assert.Equal(t, test.expectedErr.Error(), err.Error())
			}

			if test.expectedDate == nil {
				assert.Nil(t, date)
			} else {
				assert.Equal(t, test.expectedDate.Format("2006-01-02"), date.Format("2006-01-02"))
			}
		})
	}
}

func TestSuccessMessageForReassignAndPrioritiseTasks(t *testing.T) {
	assert.Equal(t, "You have assigned 1 task(s) to assignee name as a priority", successMessageForReassignAndPrioritiseTasks("2", "true", []string{"1"}, "assignee name"))
	assert.Equal(t, "You have assigned 1 task(s) to assignee name and removed priority", successMessageForReassignAndPrioritiseTasks("2", "false", []string{"1"}, "assignee name"))
	assert.Equal(t, "1 task(s) have been reassigned", successMessageForReassignAndPrioritiseTasks("2", "", []string{"1"}, "assignee name"))
	assert.Equal(t, "You have assigned 1 task(s) as a priority", successMessageForReassignAndPrioritiseTasks("0", "true", []string{"1"}, "assignee name"))
	assert.Equal(t, "You have removed 1 task(s) as a priority", successMessageForReassignAndPrioritiseTasks("0", "false", []string{"1"}, "assignee name"))
}

func TestClientTasksVars_CreateUrlBuilder(t *testing.T) {
	wantFilters := []urlbuilder.Filter{
		{Name: "task-type"},
		{Name: "assignee", ClearBetweenTeamViews: true},
		{Name: "unassigned", ClearBetweenTeamViews: true},
		{Name: "due-date-from"},
		{Name: "due-date-to"},
	}

	tests := []struct {
		page        ClientTasksPage
		wantBuilder urlbuilder.UrlBuilder
		wantFilters []urlbuilder.Filter
	}{
		{
			page:        ClientTasksPage{},
			wantBuilder: urlbuilder.UrlBuilder{Path: "client-tasks"},
			wantFilters: wantFilters,
		},
		{
			page: ClientTasksPage{
				ListPage: ListPage{
					App: WorkflowVars{SelectedTeam: model.Team{Selector: "test-team"}},
				},
			},
			wantBuilder: urlbuilder.UrlBuilder{Path: "client-tasks", SelectedTeam: "test-team", SelectedFilters: wantFilters},
			wantFilters: wantFilters,
		},
		{
			page: ClientTasksPage{
				ListPage: ListPage{
					App:     WorkflowVars{SelectedTeam: model.Team{Selector: "test-team"}},
					PerPage: 55,
				},
			},
			wantBuilder: urlbuilder.UrlBuilder{Path: "client-tasks", SelectedTeam: "test-team", SelectedPerPage: 55, SelectedFilters: wantFilters},
			wantFilters: wantFilters,
		},
		{
			page: ClientTasksPage{
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
				FilterByDueDate: FilterByDueDate{
					SelectedDueDateFrom: "2010-10-10",
					SelectedDueDateTo:   "2020-10-10",
				},
			},
			wantBuilder: urlbuilder.UrlBuilder{Path: "client-tasks", SelectedTeam: "test-team", SelectedPerPage: 55, SelectedFilters: wantFilters},
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
				{
					Name:           "due-date-from",
					SelectedValues: []string{"2010-10-10"},
				},
				{
					Name:           "due-date-to",
					SelectedValues: []string{"2020-10-10"},
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

func TestClientTasksPage_GetAppliedFilters(t *testing.T) {
	dueDateFrom := time.Date(2022, 12, 17, 0, 0, 0, 0, time.Local)
	dueDateTo := time.Date(2022, 12, 18, 0, 0, 0, 0, time.Local)

	tests := []struct {
		taskTypes          []model.TaskType
		selectedTaskTypes  []string
		selectedAssignees  []string
		selectedUnassigned string
		dueDateFrom        *time.Time
		dueDateTo          *time.Time
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
			selectedUnassigned: "lay-team",
			want:               []string{"Lay team"},
		},
		{
			dueDateFrom: &dueDateFrom,
			want:        []string{"Due date from 17/12/2022 (inclusive)"},
		},
		{
			dueDateTo: &dueDateTo,
			want:      []string{"Due date to 18/12/2022 (inclusive)"},
		},
		{
			taskTypes:          []model.TaskType{{Incomplete: "TaskType1", Handle: "TT1"}},
			selectedTaskTypes:  []string{"TT1"},
			selectedAssignees:  []string{"1"},
			selectedUnassigned: "lay-team",
			dueDateFrom:        &dueDateFrom,
			dueDateTo:          &dueDateTo,
			want:               []string{"TaskType1", "Lay team", "User 1", "Due date from 17/12/2022 (inclusive)", "Due date to 18/12/2022 (inclusive)"},
		},
	}
	for i, test := range tests {
		t.Run("Scenario "+strconv.Itoa(i+1), func(t *testing.T) {
			var page ClientTasksPage
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
			page.TaskTypes = test.taskTypes
			page.SelectedTaskTypes = test.selectedTaskTypes
			page.SelectedAssignees = test.selectedAssignees
			page.SelectedUnassigned = test.selectedUnassigned

			assert.Equal(t, test.want, page.GetAppliedFilters(test.dueDateFrom, test.dueDateTo))
		})
	}
}
