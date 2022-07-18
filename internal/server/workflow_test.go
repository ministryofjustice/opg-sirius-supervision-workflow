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
	taskTypeData      []sirius.TaskType
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

func (m *mockWorkflowInformation) GetTaskTypes(ctx sirius.Context, taskTypeSelected []string) ([]sirius.TaskType, error) {
	m.count += 1
	m.lastCtx = ctx

	return m.taskTypeData, m.err
}

func (m *mockWorkflowInformation) GetTaskList(ctx sirius.Context, search int, displayTaskLimit int, selectedTeamId int, loggedInTeamId int, taskTypeSelected []string, LoadTasks []sirius.TaskType, assigneeSelected []string) (sirius.TaskList, int, error) {
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

func (m *mockWorkflowInformation) GetAppliedFilters(ctx sirius.Context, teamId int, loadTaskTypes []sirius.TaskType, teamSelection []sirius.ReturnedTeamCollection, assigneesForFilter sirius.AssigneesTeam) []string {
	m.count += 1
	m.lastCtx = ctx

	return m.appliedFilters
}

var mockUserDetailsData = sirius.UserDetails{
	ID: 123,
	Teams: []sirius.MyDetailsTeam{
		{
			TeamId:      13,
			DisplayName: "Lay Team 1 - (Supervision)",
		},
	},
}

var mockTaskTypeData = []sirius.TaskType{
	{
		Handle:     "CDFC",
		Incomplete: "Correspondence - Review failed draft",
		Category:   "supervision",
		Complete:   "Correspondence - Reviewed draft failure",
		User:       true,
	},
}

var mockTaskListData = sirius.TaskList{
	WholeTaskList: []sirius.Task{
		{
			Assignee: sirius.CaseManagement{
				Name: "Assignee Duke Clive Henry Hetley Junior Jones",
			},
			Type:    "Case work - General",
			DueDate: "01/02/2021",
			CaseItems: []sirius.CaseItemsDetails{
				{
					Client: sirius.SupervisionClient{
						CaseRecNumber: "caseRecNumber",
						FirstName:     "Client Alexander Zacchaeus",
						Id:            3333,
						SupervisionCaseOwner: sirius.CaseManagement{
							Name: "Supervision - Team - Name",
						},
						Surname: "Client Wolfeschlegelsteinhausenbergerdorff",
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
				Id:   86,
				Name: "LayTeam1 User11",
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
	// 		WholeTaskList: []sirius.Task{
	// 			{
	// 				Assignee: sirius.AssigneeDetails{
	// 					AssigneeDisplayName: "Assignee Duke Clive Henry Hetley Junior Jones",
	// 				},
	// 				CaseItems: []sirius.CaseItemsDetails{
	// 					{
	// 						Client: sirius.ClientDetails{
	// 							CaseRecNumber: "caseRecNumber",
	// 							FirstName:     "Client Alexander Zacchaeus",
	// 							Id:            3333,
	// 							SupervisionCaseOwner: sirius.SupervisionCaseOwnerDetail{
	// 								SupervisionCaseOwnerName: "Supervision - Team - Name",
	// 							},
	// 							Surname: "Client Wolfeschlegelsteinhausenbergerdorff",
	// 						},
	// 					},
	// 				},
	// 				DueDate: "01/02/2021",
	// 				Type:    "Case work - General",
	// 			},
	// 		},
	// 	},
	// 	LoadTasks: []sirius.TaskType{
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
	// 					Id:   86,
	// 					Name: "LayTeam1 User11",
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
		WholeTaskList: []sirius.Task{{}},
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

	assert.Equal(7, client.count)

	assert.Equal(1, template.count)
	assert.Equal("page", template.lastName)
	assert.Equal(workflowVars{
		Path: "/path",
		MyDetails: sirius.UserDetails{
			ID: 123,
			Teams: []sirius.MyDetailsTeam{
				{
					TeamId:      13,
					DisplayName: "Lay Team 1 - (Supervision)",
				},
			},
		},

		TaskList: sirius.TaskList{
			WholeTaskList: []sirius.Task{{}},
		},
		LoadTasks: []sirius.TaskType{
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
						Id:   86,
						Name: "LayTeam1 User11",
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
	r, _ := http.NewRequest("POST", "/path", nil)
	defaultWorkflowTeam := 19
	handler := loggingInfoForWorkflow(client, template, defaultWorkflowTeam)
	err := handler(w, r)

	assert.Nil(err)
}
