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

func (m *mockClientTasksClient) GetTaskTypes(ctx sirius.Context, taskTypeSelected []string) ([]model.TaskType, error) {
	if m.count == nil {
		m.count = make(map[string]int)
	}
	m.count["GetTaskTypes"] += 1
	m.lastCtx = ctx

	return m.taskTypeData, m.err
}

func (m *mockClientTasksClient) GetTaskList(ctx sirius.Context, search int, displayTaskLimit int, selectedTeamId model.Team, taskTypeSelected []string, LoadTasks []model.TaskType, assigneeSelected []string, dueDateFrom *time.Time, dueDateTo *time.Time) (sirius.TaskList, error) {
	if m.count == nil {
		m.count = make(map[string]int)
	}
	m.count["GetTaskList"] += 1
	m.lastCtx = ctx

	return m.taskListData, m.err
}

func (m *mockClientTasksClient) AssignTasksToCaseManager(ctx sirius.Context, newAssigneeIdForTask int, selectedTask []string, prioritySelected string) (string, error) {
	if m.count == nil {
		m.count = make(map[string]int)
	}
	m.count["AssignTasksToCaseManager"] += 1
	m.lastCtx = ctx

	return "", m.err
}

var mockTaskTypeData = []model.TaskType{
	{
		Handle:     "CDFC",
		Incomplete: "Correspondence - Review failed draft",
		Category:   "supervision",
		Complete:   "Correspondence - Reviewed draft failure",
		User:       true,
	},
}

var mockTaskListData = sirius.TaskList{
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
		taskTypeData: mockTaskTypeData,
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

	want := ClientTasksVars{App: app, TasksPerPage: 25, TaskTypes: mockTaskTypeData}

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

	client := &mockClientTasksClient{taskTypeData: mockTaskTypeData, taskListData: mockTaskListData}
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
	client := &mockClientTasksClient{taskTypeData: mockTaskTypeData, taskListData: mockTaskListData}
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

func TestSetTaskCount_WithMatchingTaskType(t *testing.T) {
	var mockTaskListData = sirius.TaskList{
		MetaData: sirius.MetaData{
			TaskTypeCount: []sirius.TypeAndCount{
				{Type: "ORAL", Count: 25},
			},
		},
	}

	assert.Equal(t, 25, setTaskCount("ORAL", mockTaskListData))
}

func TestSetTaskCount_NoMatchingTaskTypeWillReturnZero(t *testing.T) {
	var mockTaskListData = sirius.TaskList{
		MetaData: sirius.MetaData{
			TaskTypeCount: []sirius.TypeAndCount{
				{Type: "ORAL", Count: 25},
			},
		},
	}

	assert.Equal(t, 0, setTaskCount("FREA", mockTaskListData))
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

func TestCalculateTaskCounts(t *testing.T) {
	taskTypes := []model.TaskType{
		{
			Handle: "ECM_TASKS",
		},
		{
			Handle:  "CDFC",
			EcmTask: false,
		},
		{
			Handle:  "NONO",
			EcmTask: false,
		},
		{
			Handle:  "ECM_1",
			EcmTask: true,
		},
		{
			Handle:  "ECM_2",
			EcmTask: true,
		},
	}
	tasks := sirius.TaskList{
		MetaData: sirius.MetaData{
			TaskTypeCount: []sirius.TypeAndCount{
				{Type: "CDFC", Count: 25},
				{Type: "ECM_1", Count: 33},
				{Type: "ECM_2", Count: 44},
			},
		},
	}

	expected := []model.TaskType{
		{
			Handle:    "ECM_TASKS",
			TaskCount: 77,
		}, {
			Handle:    "CDFC",
			TaskCount: 25,
		},
		{
			Handle:    "NONO",
			TaskCount: 0,
		},
		{
			Handle:    "ECM_1",
			TaskCount: 33,
		},
		{
			Handle:    "ECM_2",
			TaskCount: 44,
		},
	}

	assert.Equal(t, expected, calculateTaskCounts(taskTypes, tasks))
}

func TestSuccessMessageForReassignAndPrioritiseTasks(t *testing.T) {
	assert.Equal(t, "You have assigned 1 task(s) to assignee name as a priority", successMessageForReassignAndPrioritiseTasks("2", "yes", []string{"1"}, "assignee name"))
	assert.Equal(t, "You have assigned 1 task(s) to assignee name and removed priority", successMessageForReassignAndPrioritiseTasks("2", "no", []string{"1"}, "assignee name"))
	assert.Equal(t, "1 task(s) have been reassigned", successMessageForReassignAndPrioritiseTasks("2", "", []string{"1"}, "assignee name"))
	assert.Equal(t, "You have assigned 1 task(s) as a priority", successMessageForReassignAndPrioritiseTasks("0", "yes", []string{"1"}, "assignee name"))
	assert.Equal(t, "You have removed 1 task(s) as a priority", successMessageForReassignAndPrioritiseTasks("0", "no", []string{"1"}, "assignee name"))
}

func TestClientTasksVars_CreateUrlBuilder(t *testing.T) {
	wantFilters := []urlbuilder.Filter{
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
	}

	tests := []struct {
		clientTasksVars ClientTasksVars
		wantBuilder     urlbuilder.UrlBuilder
		wantFilters     []urlbuilder.Filter
	}{
		{
			clientTasksVars: ClientTasksVars{},
			wantBuilder:     urlbuilder.UrlBuilder{Path: "client-tasks"},
			wantFilters:     wantFilters,
		},
		{
			clientTasksVars: ClientTasksVars{App: WorkflowVars{SelectedTeam: model.Team{Selector: "test-team"}}},
			wantBuilder:     urlbuilder.UrlBuilder{Path: "client-tasks", SelectedTeam: "test-team", SelectedFilters: wantFilters},
			wantFilters:     wantFilters,
		},
		{
			clientTasksVars: ClientTasksVars{App: WorkflowVars{SelectedTeam: model.Team{Selector: "test-team"}}, TasksPerPage: 55},
			wantBuilder:     urlbuilder.UrlBuilder{Path: "client-tasks", SelectedTeam: "test-team", SelectedPerPage: 55, SelectedFilters: wantFilters},
			wantFilters:     wantFilters,
		},
		{
			clientTasksVars: ClientTasksVars{
				App:                 WorkflowVars{SelectedTeam: model.Team{Selector: "test-team"}},
				TasksPerPage:        55,
				SelectedTaskTypes:   []string{"type1", "type2"},
				SelectedAssignees:   []string{"user1", "user2"},
				SelectedUnassigned:  "test-unassigned",
				SelectedDueDateFrom: "2010-10-10",
				SelectedDueDateTo:   "2020-10-10",
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
			assert.Equal(t, test.wantBuilder, test.clientTasksVars.CreateUrlBuilder())
		})
	}
}
