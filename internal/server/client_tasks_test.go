package server

import (
	"errors"
	"github.com/ministryofjustice/opg-sirius-workflow/internal/sirius"
	"github.com/stretchr/testify/assert"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"
)

type mockClientTasksClient struct {
	count           map[string]int
	lastCtx         sirius.Context
	err             error
	taskTypeData    []sirius.ApiTaskTypes
	taskListData    sirius.TaskList
	pageDetailsData sirius.PageDetails
}

type clientTasksURLFields struct {
	SelectedTeam        string
	CurrentPage         int
	PerPageLimit        int
	SelectedAssignees   []string
	SelectedUnassigned  string
	SelectedTaskTypes   []string
	SelectedDueDateFrom string
	SelectedDueDateTo   string
}

func (m *mockClientTasksClient) GetTaskTypes(ctx sirius.Context, taskTypeSelected []string) ([]sirius.ApiTaskTypes, error) {
	if m.count == nil {
		m.count = make(map[string]int)
	}
	m.count["GetTaskTypes"] += 1
	m.lastCtx = ctx

	return m.taskTypeData, m.err
}

func (m *mockClientTasksClient) GetTaskList(ctx sirius.Context, search int, displayTaskLimit int, selectedTeamId sirius.Team, taskTypeSelected []string, LoadTasks []sirius.ApiTaskTypes, assigneeSelected []string, dueDateFrom *time.Time, dueDateTo *time.Time) (sirius.TaskList, error) {
	if m.count == nil {
		m.count = make(map[string]int)
	}
	m.count["GetTaskList"] += 1
	m.lastCtx = ctx

	return m.taskListData, m.err
}
func (m *mockClientTasksClient) GetPageDetails(taskList sirius.TaskList, search int, displayTaskLimit int) sirius.PageDetails {
	if m.count == nil {
		m.count = make(map[string]int)
	}
	m.count["GetPageDetails"] += 1

	return m.pageDetailsData
}

func (m *mockClientTasksClient) AssignTasksToCaseManager(ctx sirius.Context, newAssigneeIdForTask int, selectedTask []string, prioritySelected string) (string, error) {
	if m.count == nil {
		m.count = make(map[string]int)
	}
	m.count["AssignTasksToCaseManager"] += 1
	m.lastCtx = ctx

	return "", m.err
}

var mockTaskTypeData = []sirius.ApiTaskTypes{
	{
		Handle:     "CDFC",
		Incomplete: "Correspondence - Review failed draft",
		Category:   "supervision",
		Complete:   "Correspondence - Reviewed draft failure",
		User:       true,
	},
}

var mockTaskListData = sirius.TaskList{
	WholeTaskList: []sirius.ApiTask{
		{
			ApiTaskAssignee: sirius.CaseManagement{
				CaseManagerName: "Assignee Duke Clive Henry Hetley Junior Jones",
			},
			ApiTaskType:    "Case work - General",
			ApiTaskDueDate: "01/02/2021",
			ApiTaskCaseItems: []sirius.CaseItemsDetails{
				{
					CaseItemClient: sirius.Clients{
						ClientCaseRecNumber: "caseRecNumber",
						ClientFirstName:     "Client Alexander Zacchaeus",
						ClientId:            3333,
						ClientSupervisionCaseOwner: sirius.CaseManagement{
							CaseManagerName: "Supervision - Team - Name",
						},
						ClientSurname: "Client Wolfeschlegelsteinhausenbergerdorff",
					},
				},
			},
		},
	},
}

func TestClientTasks_NonExistentPageNumberWillRedirectToTheHighestExistingPageNumber(t *testing.T) {
	assert := assert.New(t)

	var mockTaskListData = sirius.TaskList{
		WholeTaskList: []sirius.ApiTask{{}},
		Pages: sirius.PageInformation{
			PageCurrent: 10,
			PageTotal:   2,
		},
	}

	client := &mockClientTasksClient{taskTypeData: mockTaskTypeData, taskListData: mockTaskListData}
	template := &mockTemplates{}

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
	template := &mockTemplates{}

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
	template := &mockTemplates{}

	w := httptest.NewRecorder()
	r, _ := http.NewRequest("GET", "", nil)

	app := WorkflowVars{}
	err := clientTasks(client, template)(app, w, r)

	assert.Equal("err", err.Error())
	assert.Equal(0, template.count)
}

func TestClientTasks_PostIsPermitted(t *testing.T) {
	client := &mockClientTasksClient{taskTypeData: mockTaskTypeData, taskListData: mockTaskListData}
	template := &mockTemplates{}

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
	taskTypes := []sirius.ApiTaskTypes{
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

	expected := []sirius.ApiTaskTypes{
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

func createClientTasksVars(fields clientTasksURLFields) ClientTasksVars {
	return ClientTasksVars{
		App: WorkflowVars{
			SelectedTeam: sirius.Team{Selector: fields.SelectedTeam},
		},
		PageDetails: sirius.PageDetails{
			CurrentPage:     fields.CurrentPage,
			StoredTaskLimit: fields.PerPageLimit,
		},
		SelectedAssignees:   fields.SelectedAssignees,
		SelectedUnassigned:  fields.SelectedUnassigned,
		SelectedTaskTypes:   fields.SelectedTaskTypes,
		SelectedDueDateFrom: fields.SelectedDueDateFrom,
		SelectedDueDateTo:   fields.SelectedDueDateTo,
	}
}

func TestClientTasksVars_GetClearFiltersUrl(t *testing.T) {
	tests := []struct {
		name   string
		fields clientTasksURLFields
		want   string
	}{
		{
			name:   "Per page limit is retained",
			fields: clientTasksURLFields{SelectedTeam: "lay", CurrentPage: 1, PerPageLimit: 50},
			want:   "?team=lay&page=1&per-page=50",
		},
		{
			name:   "Assignees are removed",
			fields: clientTasksURLFields{SelectedTeam: "lay", CurrentPage: 1, PerPageLimit: 25, SelectedAssignees: []string{"1", "2"}},
			want:   "?team=lay&page=1&per-page=25",
		},
		{
			name:   "Unassigned is removed",
			fields: clientTasksURLFields{SelectedTeam: "lay", CurrentPage: 1, PerPageLimit: 25, SelectedUnassigned: "1"},
			want:   "?team=lay&page=1&per-page=25",
		},
		{
			name:   "Task types are removed",
			fields: clientTasksURLFields{SelectedTeam: "lay", CurrentPage: 1, PerPageLimit: 25, SelectedTaskTypes: []string{"1", "2"}},
			want:   "?team=lay&page=1&per-page=25",
		},
		{
			name:   "Due date filters are removed",
			fields: clientTasksURLFields{SelectedTeam: "lay", CurrentPage: 1, PerPageLimit: 25, SelectedDueDateFrom: "2022-12-17", SelectedDueDateTo: "2022-12-18"},
			want:   "?team=lay&page=1&per-page=25",
		},
		{
			name:   "Page is reset back to 1",
			fields: clientTasksURLFields{SelectedTeam: "lay", CurrentPage: 2, PerPageLimit: 25, SelectedAssignees: []string{"1", "2"}, SelectedUnassigned: "1", SelectedTaskTypes: []string{"task"}},
			want:   "?team=lay&page=1&per-page=25",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			w := createClientTasksVars(tt.fields)
			assert.Equalf(t, "client-tasks"+tt.want, w.GetClearFiltersUrl(), "GetClearFiltersUrl()")
		})
	}
}

func TestClientTasksVars_GetPaginationUrl(t *testing.T) {
	type args struct {
		page    int
		perPage int
	}
	tests := []struct {
		name   string
		fields clientTasksURLFields
		args   args
		want   string
	}{
		{
			name:   "Page number is updated",
			fields: clientTasksURLFields{SelectedTeam: "lay", CurrentPage: 1, PerPageLimit: 25},
			args:   args{page: 2, perPage: 25},
			want:   "?team=lay&page=2&per-page=25",
		},
		{
			name:   "Per page limit is updated",
			fields: clientTasksURLFields{SelectedTeam: "lay", CurrentPage: 1, PerPageLimit: 25},
			args:   args{page: 1, perPage: 50},
			want:   "?team=lay&page=1&per-page=50",
		},
		{
			name:   "Per page limit is retained",
			fields: clientTasksURLFields{SelectedTeam: "lay", CurrentPage: 1, PerPageLimit: 100},
			args:   args{page: 2},
			want:   "?team=lay&page=2&per-page=100",
		},
		{
			name:   "Assignees are retained",
			fields: clientTasksURLFields{SelectedTeam: "lay", CurrentPage: 1, PerPageLimit: 25, SelectedAssignees: []string{"1", "2"}},
			args:   args{page: 2, perPage: 25},
			want:   "?team=lay&page=2&per-page=25&assignee=1&assignee=2",
		},
		{
			name:   "Unassigned is retained",
			fields: clientTasksURLFields{SelectedTeam: "lay", CurrentPage: 1, PerPageLimit: 25, SelectedUnassigned: "1"},
			args:   args{page: 2, perPage: 25},
			want:   "?team=lay&page=2&per-page=25&unassigned=1",
		},
		{
			name:   "Task types are retained",
			fields: clientTasksURLFields{SelectedTeam: "lay", CurrentPage: 1, PerPageLimit: 25, SelectedTaskTypes: []string{"1", "2"}},
			args:   args{page: 2, perPage: 25},
			want:   "?team=lay&page=2&per-page=25&task-type=1&task-type=2",
		},
		{
			name:   "Due date filters are retained",
			fields: clientTasksURLFields{SelectedTeam: "lay", CurrentPage: 1, PerPageLimit: 25, SelectedDueDateFrom: "2022-12-17", SelectedDueDateTo: "2022-12-18"},
			args:   args{page: 2, perPage: 25},
			want:   "?team=lay&page=2&per-page=25&due-date-from=2022-12-17&due-date-to=2022-12-18",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			w := createClientTasksVars(tt.fields)
			var result string
			if tt.args.perPage == 0 {
				result = w.GetPaginationUrl(tt.args.page)
			} else {
				result = w.GetPaginationUrl(tt.args.page, tt.args.perPage)
			}
			assert.Equalf(t, "client-tasks"+tt.want, result, "GetPaginationUrl(%v, %v)", tt.args.page, tt.args.perPage)
		})
	}
}

func TestClientTasksVars_GetRemoveFilterUrl(t *testing.T) {
	type args struct {
		name  string
		value interface{}
	}
	tests := []struct {
		name   string
		fields clientTasksURLFields
		args   args
		want   string
	}{
		{
			name:   "Assignee filter removed",
			fields: clientTasksURLFields{SelectedTeam: "lay", CurrentPage: 1, PerPageLimit: 25, SelectedAssignees: []string{"1", "2"}, SelectedUnassigned: "1", SelectedTaskTypes: []string{"3", "4"}, SelectedDueDateFrom: "2022-12-17", SelectedDueDateTo: "2022-12-18"},
			args:   args{name: "assignee", value: 2},
			want:   "?team=lay&page=1&per-page=25&task-type=3&task-type=4&assignee=1&unassigned=1&due-date-from=2022-12-17&due-date-to=2022-12-18",
		},
		{
			name:   "Unassigned filter removed",
			fields: clientTasksURLFields{SelectedTeam: "lay", CurrentPage: 1, PerPageLimit: 25, SelectedAssignees: []string{"1", "2"}, SelectedUnassigned: "1", SelectedTaskTypes: []string{"3", "4"}, SelectedDueDateFrom: "2022-12-17", SelectedDueDateTo: "2022-12-18"},
			args:   args{name: "unassigned", value: 1},
			want:   "?team=lay&page=1&per-page=25&task-type=3&task-type=4&assignee=1&assignee=2&due-date-from=2022-12-17&due-date-to=2022-12-18",
		},
		{
			name:   "Task type filter removed",
			fields: clientTasksURLFields{SelectedTeam: "lay", CurrentPage: 1, PerPageLimit: 25, SelectedAssignees: []string{"1", "2"}, SelectedUnassigned: "1", SelectedTaskTypes: []string{"3", "4"}, SelectedDueDateFrom: "2022-12-17", SelectedDueDateTo: "2022-12-18"},
			args:   args{name: "task-type", value: 3},
			want:   "?team=lay&page=1&per-page=25&task-type=4&assignee=1&assignee=2&unassigned=1&due-date-from=2022-12-17&due-date-to=2022-12-18",
		},
		{
			name:   "Due date from filter removed",
			fields: clientTasksURLFields{SelectedTeam: "lay", CurrentPage: 1, PerPageLimit: 25, SelectedAssignees: []string{"1", "2"}, SelectedUnassigned: "1", SelectedTaskTypes: []string{"3", "4"}, SelectedDueDateFrom: "2022-12-17", SelectedDueDateTo: "2022-12-18"},
			args:   args{name: "due-date-from", value: "2022-12-17"},
			want:   "?team=lay&page=1&per-page=25&task-type=3&task-type=4&assignee=1&assignee=2&unassigned=1&due-date-to=2022-12-18",
		},
		{
			name:   "Due date to filter removed",
			fields: clientTasksURLFields{SelectedTeam: "lay", CurrentPage: 1, PerPageLimit: 25, SelectedAssignees: []string{"1", "2"}, SelectedUnassigned: "1", SelectedTaskTypes: []string{"3", "4"}, SelectedDueDateFrom: "2022-12-17", SelectedDueDateTo: "2022-12-18"},
			args:   args{name: "due-date-to", value: "2022-12-18"},
			want:   "?team=lay&page=1&per-page=25&task-type=3&task-type=4&assignee=1&assignee=2&unassigned=1&due-date-from=2022-12-17",
		},
		{
			name:   "Page is reset back to 1 on removing a filter",
			fields: clientTasksURLFields{SelectedTeam: "lay", CurrentPage: 3, PerPageLimit: 25, SelectedAssignees: []string{"1", "2"}, SelectedUnassigned: "1", SelectedTaskTypes: []string{"3", "4"}},
			args:   args{name: "task-type", value: 3},
			want:   "?team=lay&page=1&per-page=25&task-type=4&assignee=1&assignee=2&unassigned=1",
		},
		{
			name:   "All filters retained if filter not found",
			fields: clientTasksURLFields{SelectedTeam: "lay", CurrentPage: 1, PerPageLimit: 25, SelectedAssignees: []string{"1", "2"}, SelectedUnassigned: "1", SelectedTaskTypes: []string{"3", "4"}, SelectedDueDateFrom: "2022-12-17", SelectedDueDateTo: "2022-12-18"},
			args:   args{name: "non-existent", value: 3},
			want:   "?team=lay&page=1&per-page=25&task-type=3&task-type=4&assignee=1&assignee=2&unassigned=1&due-date-from=2022-12-17&due-date-to=2022-12-18",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			w := createClientTasksVars(tt.fields)
			assert.Equalf(t, "client-tasks"+tt.want, w.GetRemoveFilterUrl(tt.args.name, tt.args.value), "GetRemoveFilterUrl(%v, %v)", tt.args.name, tt.args.value)
		})
	}
}

func TestClientTasksVars_GetTeamUrl(t *testing.T) {
	tests := []struct {
		name   string
		fields clientTasksURLFields
		team   string
		want   string
	}{
		{
			name:   "Team is retained",
			fields: clientTasksURLFields{SelectedTeam: "lay", CurrentPage: 1, PerPageLimit: 25},
			team:   "lay",
			want:   "?team=lay&page=1&per-page=25",
		},
		{
			name:   "Per page limit is retained",
			fields: clientTasksURLFields{SelectedTeam: "lay", CurrentPage: 1, PerPageLimit: 50},
			team:   "pro",
			want:   "?team=pro&page=1&per-page=50",
		},
		{
			name:   "Assignees are removed",
			fields: clientTasksURLFields{SelectedTeam: "lay", CurrentPage: 1, PerPageLimit: 25, SelectedAssignees: []string{"1", "2"}},
			team:   "pro",
			want:   "?team=pro&page=1&per-page=25",
		},
		{
			name:   "Unassigned is removed",
			fields: clientTasksURLFields{SelectedTeam: "lay", CurrentPage: 1, PerPageLimit: 25, SelectedUnassigned: "1"},
			team:   "pro",
			want:   "?team=pro&page=1&per-page=25",
		},
		{
			name:   "Task types are retained",
			fields: clientTasksURLFields{SelectedTeam: "lay", CurrentPage: 1, PerPageLimit: 25, SelectedTaskTypes: []string{"1", "2"}},
			team:   "pro",
			want:   "?team=pro&page=1&per-page=25&task-type=1&task-type=2",
		},
		{
			name:   "Due date filters are retained",
			fields: clientTasksURLFields{SelectedTeam: "lay", CurrentPage: 1, PerPageLimit: 25, SelectedDueDateFrom: "2022-12-17", SelectedDueDateTo: "2022-12-18"},
			team:   "pro",
			want:   "?team=pro&page=1&per-page=25&due-date-from=2022-12-17&due-date-to=2022-12-18",
		},
		{
			name:   "Page is reset back to 1",
			fields: clientTasksURLFields{SelectedTeam: "lay", CurrentPage: 2, PerPageLimit: 25, SelectedAssignees: []string{"1", "2"}, SelectedUnassigned: "1", SelectedTaskTypes: []string{"task"}},
			team:   "pro",
			want:   "?team=pro&page=1&per-page=25&task-type=task",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			w := createClientTasksVars(tt.fields)
			assert.Equalf(t, "client-tasks"+tt.want, w.GetTeamUrl(tt.team), "GetTeamUrl(%v)", tt.team)
		})
	}
}
