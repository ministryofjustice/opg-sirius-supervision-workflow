package server

import (
	"errors"
	"github.com/ministryofjustice/opg-go-common/logging"
	"github.com/ministryofjustice/opg-sirius-workflow/internal/sirius"
	"github.com/stretchr/testify/assert"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"
)

type mockWorkflowInformation struct {
	count             map[string]int
	lastCtx           sirius.Context
	err               error
	userData          sirius.UserDetails
	taskTypeData      []sirius.ApiTaskTypes
	taskListData      sirius.TaskList
	pageDetailsData   sirius.PageDetails
	teamSelectionData []sirius.ReturnedTeamCollection
	assignees         sirius.AssigneesTeam
	appliedFilters    []string
}

func (m *mockWorkflowInformation) GetCurrentUserDetails(ctx sirius.Context) (sirius.UserDetails, error) {
	if m.count == nil {
		m.count = make(map[string]int)
	}
	m.count["GetCurrentUserDetails"] += 1
	m.lastCtx = ctx

	return m.userData, m.err
}

func (m *mockWorkflowInformation) GetTaskTypes(ctx sirius.Context, taskTypeSelected []string) ([]sirius.ApiTaskTypes, error) {
	if m.count == nil {
		m.count = make(map[string]int)
	}
	m.count["GetTaskTypes"] += 1
	m.lastCtx = ctx

	return m.taskTypeData, m.err
}

func (m *mockWorkflowInformation) GetTaskList(ctx sirius.Context, search int, displayTaskLimit int, selectedTeamId int, taskTypeSelected []string, LoadTasks []sirius.ApiTaskTypes, assigneeSelected []string) (sirius.TaskList, error) {
	if m.count == nil {
		m.count = make(map[string]int)
	}
	m.count["GetTaskList"] += 1
	m.lastCtx = ctx

	return m.taskListData, m.err
}
func (m *mockWorkflowInformation) GetPageDetails(taskList sirius.TaskList, search int, displayTaskLimit int) sirius.PageDetails {
	if m.count == nil {
		m.count = make(map[string]int)
	}
	m.count["GetPageDetails"] += 1

	return m.pageDetailsData
}

func (m *mockWorkflowInformation) GetAssigneesForFilter(ctx sirius.Context, teamId int, assigneeSelected []string) (sirius.AssigneesTeam, error) {
	if m.count == nil {
		m.count = make(map[string]int)
	}
	m.count["GetAssigneesForFilter"] += 1
	m.lastCtx = ctx

	return m.assignees, m.err
}

func (m *mockWorkflowInformation) GetTeamsForSelection(ctx sirius.Context, teamId int, assigneeSelected []string) ([]sirius.ReturnedTeamCollection, error) {
	if m.count == nil {
		m.count = make(map[string]int)
	}
	m.count["GetTeamsForSelection"] += 1
	m.lastCtx = ctx

	return m.teamSelectionData, m.err
}

func (m *mockWorkflowInformation) AssignTasksToCaseManager(ctx sirius.Context, newAssigneeIdForTask int, selectedTask string) error {
	if m.count == nil {
		m.count = make(map[string]int)
	}
	m.count["AssignTasksToCaseManager"] += 1
	m.lastCtx = ctx

	return m.err
}

func (m *mockWorkflowInformation) GetAppliedFilters(teamId int, loadTaskTypes []sirius.ApiTaskTypes, teamSelection []sirius.ReturnedTeamCollection, assigneesForFilter sirius.AssigneesTeam) []string {
	if m.count == nil {
		m.count = make(map[string]int)
	}
	m.count["GetAppliedFilters"] += 1

	return m.appliedFilters
}

var mockUserDetailsData = sirius.UserDetails{
	ID:        123,
	Firstname: "John",
	Surname:   "Doe",
	Teams: []sirius.MyDetailsTeam{
		{
			TeamId:      13,
			DisplayName: "Lay Team 1 - (Supervision)",
		},
	},
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

var mockTeamSelectionData = []sirius.ReturnedTeamCollection{
	{
		Id: 13,
		Members: []sirius.TeamMembers{
			{
				TeamMembersId:   86,
				TeamMembersName: "LayTeam1 User11",
			},
		},
		Name: "Lay Team 1 - (Supervision)",
	},
}

func TestGetUserDetails(t *testing.T) {
	assert := assert.New(t)

	client := &mockWorkflowInformation{userData: mockUserDetailsData, taskTypeData: mockTaskTypeData, taskListData: mockTaskListData, teamSelectionData: mockTeamSelectionData}
	template := &mockTemplates{}

	w := httptest.NewRecorder()
	r, _ := http.NewRequest("GET", "/path", nil)

	defaultWorkflowTeam := 19
	handler := loggingInfoForWorkflow(client, template, defaultWorkflowTeam)
	err := handler(w, r)

	assert.Nil(err)

	resp := w.Result()
	assert.Equal(http.StatusOK, resp.StatusCode)
}

func TestGetUserDetailsWithNoTasksWillReturnWithNoErrors(t *testing.T) {
	assert := assert.New(t)

	var mockTaskListData = sirius.TaskList{
		WholeTaskList: []sirius.ApiTask{{}},
		Pages: sirius.PageInformation{
			PageCurrent: 2,
			PageTotal:   1,
		},
	}

	client := &mockWorkflowInformation{userData: mockUserDetailsData, taskTypeData: mockTaskTypeData, taskListData: mockTaskListData, teamSelectionData: mockTeamSelectionData}
	template := &mockTemplates{}

	w := httptest.NewRecorder()
	r, _ := http.NewRequest("GET", "/path", nil)

	defaultWorkflowTeam := 19
	handler := loggingInfoForWorkflow(client, template, defaultWorkflowTeam)
	err := handler(w, r)

	assert.Nil(err)

	resp := w.Result()
	assert.Equal(http.StatusOK, resp.StatusCode)
	assert.Equal(getContext(r), client.lastCtx)

	assert.Equal(7, len(client.count), client.count)

	assert.Equal(1, template.count)
	assert.Equal("page", template.lastName)
	assert.Equal(workflowVars{
		Path: "/path",
		MyDetails: sirius.UserDetails{
			ID:        123,
			Firstname: "John",
			Surname:   "Doe",
			Teams: []sirius.MyDetailsTeam{
				{
					TeamId:      13,
					DisplayName: "Lay Team 1 - (Supervision)",
				},
			},
		},

		TaskList: sirius.TaskList{
			WholeTaskList: []sirius.ApiTask{{}},
			Pages: sirius.PageInformation{
				PageCurrent: 2,
				PageTotal:   1,
			},
		},
		LoadTasks: []sirius.ApiTaskTypes{
			{
				Handle:     "CDFC",
				Incomplete: "Correspondence - Review failed draft",
				Category:   "supervision",
				Complete:   "Correspondence - Reviewed draft failure",
				User:       true,
			},
		},
		TeamSelection: []sirius.ReturnedTeamCollection{
			{
				Id: 13,
				Members: []sirius.TeamMembers{
					{
						TeamMembersId:   86,
						TeamMembersName: "LayTeam1 User11",
					},
				},
				Name: "Lay Team 1 - (Supervision)",
			},
		},
	}, template.lastVars)

}

func TestNonExistentPageNumberWillReturnTheHighestExistingPageNumber(t *testing.T) {
	assert := assert.New(t)

	var mockTaskListData = sirius.TaskList{
		WholeTaskList: []sirius.ApiTask{{}},
		Pages: sirius.PageInformation{
			PageCurrent: 10,
			PageTotal:   2,
		},
	}

	client := &mockWorkflowInformation{userData: mockUserDetailsData, taskTypeData: mockTaskTypeData, taskListData: mockTaskListData, teamSelectionData: mockTeamSelectionData}
	template := &mockTemplates{}

	w := httptest.NewRecorder()
	r, _ := http.NewRequest("GET", "/path?page=10", nil)

	defaultWorkflowTeam := 19
	handler := loggingInfoForWorkflow(client, template, defaultWorkflowTeam)
	err := handler(w, r)

	assert.Nil(err)

	resp := w.Result()
	assert.Equal(http.StatusOK, resp.StatusCode)
	assert.Equal(getContext(r), client.lastCtx)

	assert.Equal(7, len(client.count))
	assert.Equal(2, client.count["GetTaskList"])

	assert.Equal(1, template.count)
	assert.Equal("page", template.lastName)
	assert.Equal(workflowVars{
		Path: "/path",
		MyDetails: sirius.UserDetails{
			ID:        123,
			Firstname: "John",
			Surname:   "Doe",
			Teams: []sirius.MyDetailsTeam{
				{
					TeamId:      13,
					DisplayName: "Lay Team 1 - (Supervision)",
				},
			},
		},

		TaskList: sirius.TaskList{
			WholeTaskList: []sirius.ApiTask{{}},
			Pages: sirius.PageInformation{
				PageCurrent: 10,
				PageTotal:   2,
			},
		},
		LoadTasks: []sirius.ApiTaskTypes{
			{
				Handle:     "CDFC",
				Incomplete: "Correspondence - Review failed draft",
				Category:   "supervision",
				Complete:   "Correspondence - Reviewed draft failure",
				User:       true,
			},
		},
		TeamSelection: []sirius.ReturnedTeamCollection{
			{
				Id: 13,
				Members: []sirius.TeamMembers{
					{
						TeamMembersId:   86,
						TeamMembersName: "LayTeam1 User11",
					},
				},
				Name: "Lay Team 1 - (Supervision)",
			},
		},
	}, template.lastVars)

}

func TestWorkflowUnauthenticated(t *testing.T) {
	assert := assert.New(t)

	client := &mockWorkflowInformation{err: sirius.ErrUnauthorized}
	template := &mockTemplates{}

	w := httptest.NewRecorder()
	r, _ := http.NewRequest("GET", "", nil)

	defaultWorkflowTeam := 19
	handler := loggingInfoForWorkflow(client, template, defaultWorkflowTeam)
	err := handler(w, r)

	assert.Equal(sirius.ErrUnauthorized, err)

	assert.Equal(0, template.count)
}

func TestWorkflowSiriusErrors(t *testing.T) {
	assert := assert.New(t)

	client := &mockWorkflowInformation{err: errors.New("err")}
	template := &mockTemplates{}

	w := httptest.NewRecorder()
	r, _ := http.NewRequest("GET", "", nil)

	defaultWorkflowTeam := 19
	handler := loggingInfoForWorkflow(client, template, defaultWorkflowTeam)
	err := handler(w, r)

	assert.Equal("err", err.Error())

	assert.Equal(0, template.count)
}

func TestPostWorkflowIsPermitted(t *testing.T) {
	assert := assert.New(t)

	client := &mockWorkflowInformation{userData: mockUserDetailsData, taskTypeData: mockTaskTypeData, taskListData: mockTaskListData, teamSelectionData: mockTeamSelectionData}

	template := &mockTemplates{}
	w := httptest.NewRecorder()
	r, _ := http.NewRequest("POST", "", nil)
	defaultWorkflowTeam := 19
	handler := loggingInfoForWorkflow(client, template, defaultWorkflowTeam)
	err := handler(w, r)

	assert.Nil(err)
}

func TestCheckForChangesToSelectedPagination(t *testing.T) {
	assert.Equal(t, 50, checkForChangesToSelectedPagination([]string{"25", "50"}, "25"))
	assert.Equal(t, 25, checkForChangesToSelectedPagination([]string{"25", "50"}, "50"))
	assert.Equal(t, 25, checkForChangesToSelectedPagination([]string{"25", "25"}, "25"))
	assert.Equal(t, 100, checkForChangesToSelectedPagination([]string{"50", "100"}, "50"))
	assert.Equal(t, 25, checkForChangesToSelectedPagination([]string{}, "100"))
}

func TestGetLoggedInTeam(t *testing.T) {
	assert.Equal(t, 13, getLoggedInTeam(sirius.UserDetails{
		ID:          65,
		Name:        "case",
		PhoneNumber: "12345678",
		Teams: []sirius.MyDetailsTeam{
			{
				TeamId:      13,
				DisplayName: "Lay Team 1 - (Supervision)",
			},
		},
		DisplayName: "case manager",
	}, 25))

	assert.Equal(t, 25, getLoggedInTeam(sirius.UserDetails{
		ID:          65,
		Name:        "case",
		DisplayName: "case manager",
	}, 25))
}

func TestGetAssigneeIdForTask(t *testing.T) {
	logger := logging.New(os.Stdout, "opg-sirius-workflow ")

	expectedAssigneeId, expectedError := getAssigneeIdForTask(logger, "13", "67")
	assert.Equal(t, expectedAssigneeId, 67)
	assert.Nil(t, expectedError)

	expectedAssigneeId, expectedError = getAssigneeIdForTask(logger, "13", "")
	assert.Equal(t, expectedAssigneeId, 13)
	assert.Nil(t, expectedError)

	expectedAssigneeId, expectedError = getAssigneeIdForTask(logger, "", "")
	assert.Equal(t, expectedAssigneeId, 0)
	assert.Nil(t, expectedError)
}

func TestCreateTaskIdForUrl(t *testing.T) {
	assert.Equal(t, "", createTaskIdForUrl([]string{}))
	assert.Equal(t, "15+16+17", createTaskIdForUrl([]string{"15", "16", "17"}))
	assert.Equal(t, "15", createTaskIdForUrl([]string{"15"}))
}

func TestGetSelectedTeamId(t *testing.T) {
	r, _ := http.NewRequest("GET", "?change-team=13", nil)

	actualTeamId := getSelectedTeamId(r, 20)
	assert.Equal(t, 13, actualTeamId)
}

func TestGetSelectedTeamIdNoneSelected(t *testing.T) {
	r, _ := http.NewRequest("GET", "", nil)
	actualTeamId := getSelectedTeamId(r, 20)
	assert.Equal(t, 20, actualTeamId)
}
