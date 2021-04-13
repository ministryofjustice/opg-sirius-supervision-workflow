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
	taskTypeData      sirius.TaskTypes
	taskListData      sirius.TaskList
	taskDetailsData   sirius.TaskDetails
	teamSelectionData []sirius.TeamCollection
}

func (m *mockWorkflowInformation) SiriusUserDetails(ctx sirius.Context) (sirius.UserDetails, error) {
	m.count += 1
	m.lastCtx = ctx

	return m.userData, m.err
}

func (m *mockWorkflowInformation) GetTaskType(ctx sirius.Context) (sirius.TaskTypes, error) {
	m.count += 1
	m.lastCtx = ctx

	return m.taskTypeData, m.err
}

func (m *mockWorkflowInformation) GetTaskList(ctx sirius.Context, search int, displayTaskLimit int) (sirius.TaskList, sirius.TaskDetails, error) {
	m.count += 1
	m.lastCtx = ctx

	return m.taskListData, m.taskDetailsData, m.err
}

func (m *mockWorkflowInformation) GetTeamSelection(ctx sirius.Context) ([]sirius.TeamCollection, error) {
	m.count += 1
	m.lastCtx = ctx

	return m.teamSelectionData, m.err
}

func TestGetUserDetails(t *testing.T) {
	assert := assert.New(t)

	mockUserDetailsData := sirius.UserDetails{
		ID:        123,
		Firstname: "John",
		Surname:   "Doe",
	}

	mockTaskTypeData := sirius.TaskTypes{
		TaskTypeList: sirius.ApiTaskTypes{
			Handle:     "CDFC",
			Incomplete: "Correspondence - Review failed draft",
			Category:   "supervision",
			Complete:   "Correspondence - Reviewed draft failure",
			User:       true,
		},
	}

	mockTaskListData := sirius.TaskList{
		WholeTaskList: []sirius.ApiTask{
			{
				ApiTaskAssignee: sirius.AssigneeDetails{
					AssigneeDisplayName: "Assignee Duke Clive Henry Hetley Junior Jones",
				},
				ApiTaskType:    "Case work - General",
				ApiTaskDueDate: "01/02/2021",
				ApiTaskCaseItems: []sirius.CaseItemsDetails{
					{
						CaseItemClient: sirius.ClientDetails{
							ClientCaseRecNumber: "caseRecNumber",
							ClientFirstName:     "Client Alexander Zacchaeus",
							ClientId:            3333,
							ClientSupervisionCaseOwner: sirius.SupervisionCaseOwnerDetail{
								SupervisionCaseOwnerName: "Supervision - Team - Name",
							},
							ClientSurname: "Client Wolfeschlegelsteinhausenbergerdorff",
						},
					},
				},
			},
		},
	}

	mockTeamSelectionData := []sirius.TeamCollection{
		{
			Id: 13,
			Members: []sirius.TeamMembers{
				{
					TeamMembersId:   96,
					TeamMembersName: "LayTeam1 User11",
				},
			},
			Name: "Lay Team 1 - (Supervision)",
		},
	}

	client := &mockWorkflowInformation{userData: mockUserDetailsData, taskTypeData: mockTaskTypeData, taskListData: mockTaskListData, teamSelectionData: mockTeamSelectionData}
	template := &mockTemplates{}

	w := httptest.NewRecorder()
	r, _ := http.NewRequest("GET", "/path", nil)

	handler := loggingInfoForWorflow(client, template)
	err := handler(sirius.PermissionSet{}, w, r)

	assert.Nil(err)

	resp := w.Result()
	assert.Equal(http.StatusOK, resp.StatusCode)
	assert.Equal(getContext(r), client.lastCtx)

	assert.Equal(4, client.count)

	assert.Equal(1, template.count)
	assert.Equal("page", template.lastName)
	assert.Equal(workflowVars{
		Path:      "/path",
		ID:        123,
		Firstname: "John",
		Surname:   "Doe",
		TaskList: sirius.TaskList{
			WholeTaskList: []sirius.ApiTask{
				{
					ApiTaskAssignee: sirius.AssigneeDetails{
						AssigneeDisplayName: "Assignee Duke Clive Henry Hetley Junior Jones",
					},
					ApiTaskCaseItems: []sirius.CaseItemsDetails{
						{
							CaseItemClient: sirius.ClientDetails{
								ClientCaseRecNumber: "caseRecNumber",
								ClientFirstName:     "Client Alexander Zacchaeus",
								ClientId:            3333,
								ClientSupervisionCaseOwner: sirius.SupervisionCaseOwnerDetail{
									SupervisionCaseOwnerName: "Supervision - Team - Name",
								},
								ClientSurname: "Client Wolfeschlegelsteinhausenbergerdorff",
							},
						},
					},
					ApiTaskDueDate: "01/02/2021",
					ApiTaskType:    "Case work - General",
				},
			},
		},
		LoadTasks: sirius.TaskTypes{
			TaskTypeList: sirius.ApiTaskTypes{
				Handle:     "CDFC",
				Incomplete: "Correspondence - Review failed draft",
				Category:   "supervision",
				Complete:   "Correspondence - Reviewed draft failure",
				User:       true,
			},
		},
		TeamSelection: []sirius.TeamCollection{
			{
				Id: 13,
				Members: []sirius.TeamMembers{
					{
						TeamMembersId:   96,
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

	handler := loggingInfoForWorflow(client, template)
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

	handler := loggingInfoForWorflow(client, template)
	err := handler(sirius.PermissionSet{}, w, r)

	assert.Equal("err", err.Error())

	assert.Equal(0, template.count)
}

func TestPostWorkflow(t *testing.T) {
	assert := assert.New(t)
	template := &mockTemplates{}

	w := httptest.NewRecorder()
	r, _ := http.NewRequest("POST", "", nil)

	handler := loggingInfoForWorflow(nil, template)
	err := handler(sirius.PermissionSet{}, w, r)

	assert.Equal(StatusError(http.StatusMethodNotAllowed), err)

	assert.Equal(0, template.count)
}
