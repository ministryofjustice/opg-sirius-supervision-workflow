package server

import (
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/ministryofjustice/opg-sirius-workflow/internal/sirius"
	"github.com/stretchr/testify/assert"
)

type mockWorkflowInformation struct {
	count             int
	lastCtx           sirius.Context
	err               error
	userData          sirius.UserDetails
	taskTypeData      []sirius.ApiTaskTypes
	taskListData      sirius.TaskList
	taskDetailsData   sirius.TaskDetails
	teamSelectionData []sirius.ReturnedTeamCollection
	assignees         sirius.AssigneesTeam
	teamId            int
	appliedFilters    []string
}

func (m *mockWorkflowInformation) GetCurrentUserDetails(ctx sirius.Context) (sirius.UserDetails, error) {
	m.count += 1
	m.lastCtx = ctx

	return m.userData, m.err
}

func (m *mockWorkflowInformation) GetTaskTypes(ctx sirius.Context, taskTypeSelected []string) ([]sirius.ApiTaskTypes, error) {
	m.count += 1
	m.lastCtx = ctx

	return m.taskTypeData, m.err
}

func (m *mockWorkflowInformation) GetTaskList(ctx sirius.Context, search int, displayTaskLimit int, selectedTeamId int, loggedInTeamId int, taskTypeSelected []string, LoadTasks []sirius.ApiTaskTypes, assigneeSelected []string) (sirius.TaskList, int, error) {
	m.count += 1
	m.lastCtx = ctx

	return m.taskListData, m.teamId, m.err
}
func (m *mockWorkflowInformation) GetTaskDetails(ctx sirius.Context, taskList sirius.TaskList, search int, displayTaskLimit int) sirius.TaskDetails {
	m.count += 1
	m.lastCtx = ctx

	return m.taskDetailsData
}

func (m *mockWorkflowInformation) GetAssigneesForFilter(ctx sirius.Context, teamId int, assigneeSelected []string) (sirius.AssigneesTeam, error) {
	m.count += 1
	m.lastCtx = ctx

	return m.assignees, m.err
}

func (m *mockWorkflowInformation) GetTeamsForSelection(ctx sirius.Context, teamId int, assigneeSelected []string) ([]sirius.ReturnedTeamCollection, error) {
	m.count += 1
	m.lastCtx = ctx

	return m.teamSelectionData, m.err
}

func (m *mockWorkflowInformation) AssignTasksToCaseManager(ctx sirius.Context, newAssigneeIdForTask int, selectedTask string) error {
	m.count += 1
	m.lastCtx = ctx

	return m.err
}

func (m *mockWorkflowInformation) GetAppliedFilters(ctx sirius.Context, teamId int, loadTaskTypes []sirius.ApiTaskTypes, teamSelection []sirius.ReturnedTeamCollection, assigneesForFilter sirius.AssigneesTeam) []string {
	m.count += 1
	m.lastCtx = ctx

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

	handler := loggingInfoForWorkflow(client, template)
	err := handler(sirius.PermissionSet{}, w, r)

	assert.Nil(err)

	resp := w.Result()
	assert.Equal(http.StatusOK, resp.StatusCode)
	// assert.Equal(getContext(r), client.lastCtx)

	// assert.Equal(5, client.count)

	// assert.Equal(1, template.count)
	// assert.Equal("page", template.lastName)
	// assert.Equal(workflowVars{
	// 	Path: "/path",
	// 	MyDetails: sirius.UserDetails{
	// 		ID:        123,
	// 		Firstname: "John",
	// 		Surname:   "Doe",
	// 		Teams: []sirius.MyDetailsTeam{
	// 			{
	// 				TeamId:      13,
	// 				DisplayName: "Lay Team 1 - (Supervision)",
	// 			},
	// 		},
	// 	},

	// 	TaskList: sirius.TaskList{
	// 		WholeTaskList: []sirius.ApiTask{
	// 			{
	// 				ApiTaskAssignee: sirius.AssigneeDetails{
	// 					AssigneeDisplayName: "Assignee Duke Clive Henry Hetley Junior Jones",
	// 				},
	// 				ApiTaskCaseItems: []sirius.CaseItemsDetails{
	// 					{
	// 						CaseItemClient: sirius.ClientDetails{
	// 							ClientCaseRecNumber: "caseRecNumber",
	// 							ClientFirstName:     "Client Alexander Zacchaeus",
	// 							ClientId:            3333,
	// 							ClientSupervisionCaseOwner: sirius.SupervisionCaseOwnerDetail{
	// 								SupervisionCaseOwnerName: "Supervision - Team - Name",
	// 							},
	// 							ClientSurname: "Client Wolfeschlegelsteinhausenbergerdorff",
	// 						},
	// 					},
	// 				},
	// 				ApiTaskDueDate: "01/02/2021",
	// 				ApiTaskType:    "Case work - General",
	// 			},
	// 		},
	// 	},
	// 	LoadTasks: []sirius.ApiTaskTypes{
	// 		{
	// 			Handle:     "CDFC",
	// 			Incomplete: "Correspondence - Review failed draft",
	// 			Category:   "supervision",
	// 			Complete:   "Correspondence - Reviewed draft failure",
	// 			User:       true,
	// 		},
	// 	},
	// 	TeamSelection: []sirius.TeamCollection{
	// 		{
	// 			Id: 13,
	// 			Members: []sirius.TeamMembers{
	// 				{
	// 					TeamMembersId:   86,
	// 					TeamMembersName: "LayTeam1 User11",
	// 				},
	// 			},
	// 			Name: "Lay Team 1 - (Supervision)",
	// 		},
	// 	},
	// }, template.lastVars)

}

func TestGetUserDetailsWithNoTasksWillReturnWithNoErrors(t *testing.T) {
	assert := assert.New(t)

	var mockTaskListData = sirius.TaskList{
		WholeTaskList: []sirius.ApiTask{{}},
	}

	client := &mockWorkflowInformation{userData: mockUserDetailsData, taskTypeData: mockTaskTypeData, taskListData: mockTaskListData, teamSelectionData: mockTeamSelectionData}
	template := &mockTemplates{}

	w := httptest.NewRecorder()
	r, _ := http.NewRequest("GET", "/path", nil)

	handler := loggingInfoForWorkflow(client, template)
	err := handler(sirius.PermissionSet{}, w, r)

	assert.Nil(err)

	resp := w.Result()
	assert.Equal(http.StatusOK, resp.StatusCode)
	assert.Equal(getContext(r), client.lastCtx)

	assert.Equal(7, client.count)

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

	handler := loggingInfoForWorkflow(client, template)
	err := handler(sirius.PermissionSet{}, w, r)

	assert.Equal(sirius.ErrUnauthorized, err)

	assert.Equal(0, template.count)
}

func TestWorkflowSiriusErrors(t *testing.T) {
	assert := assert.New(t)

	client := &mockWorkflowInformation{err: errors.New("err")}
	template := &mockTemplates{}

	w := httptest.NewRecorder()
	r, _ := http.NewRequest("GET", "", nil)

	handler := loggingInfoForWorkflow(client, template)
	err := handler(sirius.PermissionSet{}, w, r)

	assert.Equal("err", err.Error())

	assert.Equal(0, template.count)
}

func TestPostWorkflowIsPermitted(t *testing.T) {
	assert := assert.New(t)

	client := &mockWorkflowInformation{userData: mockUserDetailsData, taskTypeData: mockTaskTypeData, taskListData: mockTaskListData, teamSelectionData: mockTeamSelectionData}

	template := &mockTemplates{}
	w := httptest.NewRecorder()
	r, _ := http.NewRequest("POST", "/path", nil)
	handler := loggingInfoForWorkflow(client, template)
	err := handler(sirius.PermissionSet{}, w, r)

	assert.Nil(err)
}
